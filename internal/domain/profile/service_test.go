package profile_test

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/golang/mock/gomock"
	"github.com/maypok86/conduit/internal/domain/profile"
	"github.com/maypok86/conduit/pkg/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var errRepository = errors.New("repository error")

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
	username := faker.Username()
	now := time.Now()
	userBio := faker.Sentence()
	userImage := faker.URL()
	validProfile := profile.Profile{
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
				repository.EXPECT().GetByUsername(ctx, username).Return(profile.Profile{}, errRepository)
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
