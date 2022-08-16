package profile_test

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/maypok86/conduit/internal/domain/profile"
	"github.com/maypok86/conduit/pkg/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var (
	errGetByUsernameRepository  = errors.New("get by username repository error")
	errGetByEmailRepository     = errors.New("get by email repository error")
	errGetWithFollowRepository  = fmt.Errorf("get with follow repository error: %w", profile.ErrNotFound)
	errCheckFollowingRepository = errors.New("check following repository error")
)

func mockService(t *testing.T) (profile.Service, *MockRepository) {
	t.Helper()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	repository := NewMockRepository(mockCtrl)
	service := profile.NewService(repository)

	return service, repository
}

func TestService_GetByUsername(t *testing.T) {
	t.Parallel()

	ctx := logger.ContextWithLogger(context.Background(), zap.L())
	userID := uuid.New()
	username := faker.Username()
	now := time.Now()
	userBio := faker.Sentence()
	userImage := faker.URL()
	validProfile := profile.Profile{
		ID:        userID,
		Username:  username,
		Bio:       &userBio,
		Image:     &userImage,
		CreatedAt: now,
		UpdatedAt: now,
	}

	type args struct {
		username string
	}

	tests := []struct {
		name    string
		mock    func(*MockRepository)
		args    args
		want    profile.Profile
		wantErr bool
	}{
		{
			name: "success get user by email",
			mock: func(repository *MockRepository) {
				repository.EXPECT().GetByUsername(ctx, username).Return(validProfile, nil)
			},
			args: args{
				username: username,
			},
			want: validProfile,
		},
		{
			name: "repository error",
			mock: func(repository *MockRepository) {
				repository.EXPECT().GetByUsername(ctx, username).Return(profile.Profile{}, errGetByUsernameRepository)
			},
			args: args{
				username: username,
			},
			want:    profile.Profile{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service, repository := mockService(t)

			tt.mock(repository)

			got, err := service.GetByUsername(ctx, tt.args.username)
			require.True(t, (err != nil) == tt.wantErr)
			require.True(t, reflect.DeepEqual(tt.want, got))
		})
	}
}

func TestService_GetByEmail(t *testing.T) {
	t.Parallel()

	ctx := logger.ContextWithLogger(context.Background(), zap.L())
	userID := uuid.New()
	email := faker.Email()
	username := faker.Username()
	now := time.Now()
	userBio := faker.Sentence()
	userImage := faker.URL()
	validProfile := profile.Profile{
		ID:        userID,
		Username:  username,
		Bio:       &userBio,
		Image:     &userImage,
		CreatedAt: now,
		UpdatedAt: now,
	}

	type args struct {
		email string
	}

	tests := []struct {
		name    string
		mock    func(*MockRepository)
		args    args
		want    profile.Profile
		wantErr bool
	}{
		{
			name: "success get user by email",
			mock: func(repository *MockRepository) {
				repository.EXPECT().GetByEmail(ctx, email).Return(validProfile, nil)
			},
			args: args{
				email: email,
			},
			want: validProfile,
		},
		{
			name: "repository error",
			mock: func(repository *MockRepository) {
				repository.EXPECT().GetByEmail(ctx, email).Return(profile.Profile{}, errGetByEmailRepository)
			},
			args: args{
				email: email,
			},
			want:    profile.Profile{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service, repository := mockService(t)

			tt.mock(repository)

			got, err := service.GetByEmail(ctx, tt.args.email)
			require.True(t, (err != nil) == tt.wantErr)
			require.True(t, reflect.DeepEqual(tt.want, got))
		})
	}
}

func TestService_GetWithFollow(t *testing.T) {
	t.Parallel()

	ctx := logger.ContextWithLogger(context.Background(), zap.L())
	followeeUsername := faker.Username()
	followeeID := uuid.New()
	followeeBio := faker.Sentence()
	followeeImage := faker.URL()
	now := time.Now()
	followee := profile.Profile{
		ID:        followeeID,
		Username:  followeeUsername,
		Bio:       &followeeBio,
		Image:     &followeeImage,
		CreatedAt: now,
		UpdatedAt: now,
	}

	userID := uuid.New()
	email := faker.Email()
	username := faker.Username()
	now = time.Now()
	userBio := faker.Sentence()
	userImage := faker.URL()
	validProfile := profile.Profile{
		ID:        userID,
		Username:  username,
		Bio:       &userBio,
		Image:     &userImage,
		CreatedAt: now,
		UpdatedAt: now,
		Following: true,
	}
	notFoundProfile := profile.Profile{
		ID:        userID,
		Username:  username,
		Bio:       &userBio,
		Image:     &userImage,
		CreatedAt: now,
		UpdatedAt: now,
	}

	type args struct {
		followeeUsername string
		email            string
	}

	tests := []struct {
		name    string
		mock    func(*MockRepository)
		args    args
		want    profile.Profile
		wantErr bool
	}{
		{
			name: "success get user by email",
			mock: func(repository *MockRepository) {
				repository.EXPECT().GetByUsername(ctx, followeeUsername).Return(followee, nil)
				repository.EXPECT().CheckFollowing(ctx, email, followeeID).Return(validProfile, nil)
			},
			args: args{
				followeeUsername: followeeUsername,
				email:            email,
			},
			want: validProfile,
		},
		{
			name: "get by username error",
			mock: func(repository *MockRepository) {
				repository.EXPECT().
					GetByUsername(ctx, followeeUsername).
					Return(profile.Profile{}, errGetByUsernameRepository)
			},
			args: args{
				followeeUsername: followeeUsername,
				email:            email,
			},
			want:    profile.Profile{},
			wantErr: true,
		},
		{
			name: "not found following error",
			mock: func(repository *MockRepository) {
				repository.EXPECT().GetByUsername(ctx, followeeUsername).Return(followee, nil)
				repository.EXPECT().
					CheckFollowing(ctx, email, followeeID).
					Return(notFoundProfile, errGetWithFollowRepository)
				repository.EXPECT().GetByEmail(ctx, email).Return(notFoundProfile, nil)
			},
			args: args{
				followeeUsername: followeeUsername,
				email:            email,
			},
			want: notFoundProfile,
		},
		{
			name: "checking following error",
			mock: func(repository *MockRepository) {
				repository.EXPECT().GetByUsername(ctx, followeeUsername).Return(followee, nil)
				repository.EXPECT().
					CheckFollowing(ctx, email, followeeID).
					Return(profile.Profile{}, errCheckFollowingRepository)
			},
			args: args{
				followeeUsername: followeeUsername,
				email:            email,
			},
			want:    profile.Profile{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service, repository := mockService(t)

			tt.mock(repository)

			got, err := service.GetWithFollow(ctx, tt.args.email, tt.args.followeeUsername)
			require.True(t, (err != nil) == tt.wantErr)
			require.True(t, reflect.DeepEqual(tt.want, got))
		})
	}
}
