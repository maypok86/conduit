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
	errFollowRepository         = errors.New("follow repository error")
	errUnfollowRepository       = errors.New("unfollow repository error")
)

func mockService(t *testing.T) (profile.Service, *MockRepository) {
	t.Helper()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	repository := NewMockRepository(mockCtrl)
	service := profile.NewService(repository)

	return service, repository
}

func createProfile(t *testing.T, following bool) profile.Profile {
	t.Helper()

	id := uuid.New()
	username := faker.Username()
	bio := faker.Sentence()
	image := faker.URL()
	now := time.Now()

	return profile.Profile{
		ID:        id,
		Username:  username,
		Bio:       &bio,
		Image:     &image,
		CreatedAt: now,
		UpdatedAt: now,
		Following: following,
	}
}

func TestService_GetByUsername(t *testing.T) {
	t.Parallel()

	ctx := logger.ContextWithLogger(context.Background(), zap.L())
	validProfile := createProfile(t, false)

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
				repository.EXPECT().GetByUsername(ctx, validProfile.Username).Return(validProfile, nil)
			},
			args: args{
				username: validProfile.Username,
			},
			want: validProfile,
		},
		{
			name: "repository error",
			mock: func(repository *MockRepository) {
				repository.EXPECT().
					GetByUsername(ctx, validProfile.Username).
					Return(profile.Profile{}, errGetByUsernameRepository)
			},
			args: args{
				username: validProfile.Username,
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
	email := faker.Email()
	validProfile := createProfile(t, false)

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
	email := faker.Email()

	followee := createProfile(t, false)
	validProfile := createProfile(t, true)
	notFoundProfile := validProfile
	notFoundProfile.Following = false

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
				repository.EXPECT().GetByUsername(ctx, followee.Username).Return(followee, nil)
				repository.EXPECT().CheckFollowing(ctx, email, followee.ID).Return(validProfile, nil)
			},
			args: args{
				followeeUsername: followee.Username,
				email:            email,
			},
			want: validProfile,
		},
		{
			name: "get by username error",
			mock: func(repository *MockRepository) {
				repository.EXPECT().
					GetByUsername(ctx, followee.Username).
					Return(profile.Profile{}, errGetByUsernameRepository)
			},
			args: args{
				followeeUsername: followee.Username,
				email:            email,
			},
			want:    profile.Profile{},
			wantErr: true,
		},
		{
			name: "not found following error",
			mock: func(repository *MockRepository) {
				repository.EXPECT().GetByUsername(ctx, followee.Username).Return(followee, nil)
				repository.EXPECT().
					CheckFollowing(ctx, email, followee.ID).
					Return(notFoundProfile, errGetWithFollowRepository)
				repository.EXPECT().GetByEmail(ctx, email).Return(notFoundProfile, nil)
			},
			args: args{
				followeeUsername: followee.Username,
				email:            email,
			},
			want: notFoundProfile,
		},
		{
			name: "checking following error",
			mock: func(repository *MockRepository) {
				repository.EXPECT().GetByUsername(ctx, followee.Username).Return(followee, nil)
				repository.EXPECT().
					CheckFollowing(ctx, email, followee.ID).
					Return(profile.Profile{}, errCheckFollowingRepository)
			},
			args: args{
				followeeUsername: followee.Username,
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

func TestService_Follow(t *testing.T) {
	t.Parallel()

	ctx := logger.ContextWithLogger(context.Background(), zap.L())
	email := faker.Email()

	followee := createProfile(t, false)
	follower := createProfile(t, false)
	realFollower := follower
	realFollower.Following = true

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
			name: "success follow",
			mock: func(repository *MockRepository) {
				repository.EXPECT().GetByUsername(ctx, followee.Username).Return(followee, nil)
				repository.EXPECT().GetByEmail(ctx, email).Return(follower, nil)
				repository.EXPECT().Follow(ctx, followee.ID, follower.ID).Return(nil)
			},
			args: args{
				followeeUsername: followee.Username,
				email:            email,
			},
			want: realFollower,
		},
		{
			name: "get by username error",
			mock: func(repository *MockRepository) {
				repository.EXPECT().
					GetByUsername(ctx, followee.Username).
					Return(profile.Profile{}, errGetByUsernameRepository)
			},
			args: args{
				followeeUsername: followee.Username,
				email:            email,
			},
			want:    profile.Profile{},
			wantErr: true,
		},
		{
			name: "get by email error",
			mock: func(repository *MockRepository) {
				repository.EXPECT().GetByUsername(ctx, followee.Username).Return(followee, nil)
				repository.EXPECT().GetByEmail(ctx, email).Return(profile.Profile{}, errGetByEmailRepository)
			},
			args: args{
				followeeUsername: followee.Username,
				email:            email,
			},
			want:    profile.Profile{},
			wantErr: true,
		},
		{
			name: "follow error",
			mock: func(repository *MockRepository) {
				repository.EXPECT().GetByUsername(ctx, followee.Username).Return(followee, nil)
				repository.EXPECT().GetByEmail(ctx, email).Return(follower, nil)
				repository.EXPECT().Follow(ctx, followee.ID, follower.ID).Return(errFollowRepository)
			},
			args: args{
				followeeUsername: followee.Username,
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

			got, err := service.Follow(ctx, tt.args.email, tt.args.followeeUsername)
			require.True(t, (err != nil) == tt.wantErr)
			require.True(t, reflect.DeepEqual(tt.want, got))
		})
	}
}

func TestService_Unfollow(t *testing.T) {
	t.Parallel()

	ctx := logger.ContextWithLogger(context.Background(), zap.L())
	email := faker.Email()

	followee := createProfile(t, false)
	follower := createProfile(t, true)
	realFollower := follower
	realFollower.Following = false

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
			name: "success unfollow",
			mock: func(repository *MockRepository) {
				repository.EXPECT().GetByUsername(ctx, followee.Username).Return(followee, nil)
				repository.EXPECT().GetByEmail(ctx, email).Return(follower, nil)
				repository.EXPECT().Unfollow(ctx, followee.ID, follower.ID).Return(nil)
			},
			args: args{
				followeeUsername: followee.Username,
				email:            email,
			},
			want: realFollower,
		},
		{
			name: "get by username error",
			mock: func(repository *MockRepository) {
				repository.EXPECT().
					GetByUsername(ctx, followee.Username).
					Return(profile.Profile{}, errGetByUsernameRepository)
			},
			args: args{
				followeeUsername: followee.Username,
				email:            email,
			},
			want:    profile.Profile{},
			wantErr: true,
		},
		{
			name: "get by email error",
			mock: func(repository *MockRepository) {
				repository.EXPECT().GetByUsername(ctx, followee.Username).Return(followee, nil)
				repository.EXPECT().GetByEmail(ctx, email).Return(profile.Profile{}, errGetByEmailRepository)
			},
			args: args{
				followeeUsername: followee.Username,
				email:            email,
			},
			want:    profile.Profile{},
			wantErr: true,
		},
		{
			name: "unfollow error",
			mock: func(repository *MockRepository) {
				repository.EXPECT().GetByUsername(ctx, followee.Username).Return(followee, nil)
				repository.EXPECT().GetByEmail(ctx, email).Return(follower, nil)
				repository.EXPECT().Unfollow(ctx, followee.ID, follower.ID).Return(errUnfollowRepository)
			},
			args: args{
				followeeUsername: followee.Username,
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

			got, err := service.Unfollow(ctx, tt.args.email, tt.args.followeeUsername)
			require.True(t, (err != nil) == tt.wantErr)
			require.True(t, reflect.DeepEqual(tt.want, got))
		})
	}
}
