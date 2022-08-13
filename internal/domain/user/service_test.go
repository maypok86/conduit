package user_test

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/maypok86/conduit/internal/domain/user"
	"github.com/maypok86/conduit/pkg/logger"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var (
	errHasher     = errors.New("hasher error")
	errRepository = errors.New("repository error")
)

func mockService(t *testing.T) (user.Service, *MockRepository, *MockPasswordHasher) {
	t.Helper()

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	repository := NewMockRepository(mockCtrl)
	passwordHasher := NewMockPasswordHasher(mockCtrl)
	service := user.NewService(repository, passwordHasher)

	return service, repository, passwordHasher
}

func TestService_Create(t *testing.T) {
	t.Parallel()

	ctx := logger.ContextWithLogger(context.Background(), zap.L())
	dto := user.CreateDTO{
		Username: faker.Username(),
		Email:    faker.Email(),
		Password: faker.Password(),
	}
	now := time.Now()
	userBio := faker.Sentence()
	userImage := faker.URL()
	validUser := user.User{
		ID:        uuid.New(),
		Username:  faker.Username(),
		Email:     faker.Email(),
		Password:  faker.Password(),
		Bio:       &userBio,
		Image:     &userImage,
		CreatedAt: now,
		UpdatedAt: now,
	}

	type args struct {
		dto user.CreateDTO
	}

	tests := []struct {
		name    string
		mock    func(*MockRepository, *MockPasswordHasher)
		args    args
		want    user.User
		wantErr bool
	}{
		{
			name: "creation user",
			mock: func(repository *MockRepository, hasher *MockPasswordHasher) {
				repository.EXPECT().Create(ctx, gomock.Any()).Return(validUser, nil)
				hasher.EXPECT().Hash(dto.Password).Return(validUser.Password, nil)
			},
			args: args{
				dto: dto,
			},
			want: validUser,
		},
		{
			name: "hasher error",
			mock: func(repository *MockRepository, hasher *MockPasswordHasher) {
				hasher.EXPECT().Hash(dto.Password).Return("", errHasher)
			},
			args: args{
				dto: dto,
			},
			want:    user.User{},
			wantErr: true,
		},
		{
			name: "repository error",
			mock: func(repository *MockRepository, hasher *MockPasswordHasher) {
				hasher.EXPECT().Hash(dto.Password).Return(validUser.Password, nil)
				repository.EXPECT().Create(ctx, gomock.Any()).Return(user.User{}, errRepository)
			},
			args: args{
				dto: dto,
			},
			want:    user.User{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service, repository, passwordHasher := mockService(t)

			tt.mock(repository, passwordHasher)

			got, err := service.Create(ctx, tt.args.dto)
			require.True(t, (err != nil) == tt.wantErr)
			require.True(t, reflect.DeepEqual(tt.want, got))
		})
	}
}

func TestService_GetByEmail(t *testing.T) {
	t.Parallel()

	ctx := logger.ContextWithLogger(context.Background(), zap.L())
	email := faker.Email()
	now := time.Now()
	userBio := faker.Sentence()
	userImage := faker.URL()
	validUser := user.User{
		ID:        uuid.New(),
		Username:  faker.Username(),
		Email:     faker.Email(),
		Password:  faker.Password(),
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
		want    user.User
		wantErr bool
	}{
		{
			name: "success get user by email",
			mock: func(repository *MockRepository) {
				repository.EXPECT().GetByEmail(ctx, email).Return(validUser, nil)
			},
			args: args{
				email: email,
			},
			want: validUser,
		},
		{
			name: "repository error",
			mock: func(repository *MockRepository) {
				repository.EXPECT().GetByEmail(ctx, email).Return(user.User{}, errRepository)
			},
			args: args{
				email: email,
			},
			want:    user.User{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service, repository, _ := mockService(t)

			tt.mock(repository)

			got, err := service.GetByEmail(ctx, tt.args.email)
			require.True(t, (err != nil) == tt.wantErr)
			require.True(t, reflect.DeepEqual(tt.want, got))
		})
	}
}

func TestService_Login(t *testing.T) {
	t.Parallel()

	ctx := logger.ContextWithLogger(context.Background(), zap.L())
	email := faker.Email()
	password := faker.Password()
	now := time.Now()
	userBio := faker.Sentence()
	userImage := faker.URL()
	validUser := user.User{
		ID:        uuid.New(),
		Username:  faker.Username(),
		Email:     faker.Email(),
		Password:  faker.Password(),
		Bio:       &userBio,
		Image:     &userImage,
		CreatedAt: now,
		UpdatedAt: now,
	}

	type args struct {
		email    string
		password string
	}

	tests := []struct {
		name    string
		mock    func(*MockRepository, *MockPasswordHasher)
		args    args
		want    user.User
		wantErr bool
	}{
		{
			name: "creation user",
			mock: func(repository *MockRepository, hasher *MockPasswordHasher) {
				repository.EXPECT().GetByEmail(ctx, email).Return(validUser, nil)
				hasher.EXPECT().Check(password, validUser.Password).Return(nil)
			},
			args: args{
				email:    email,
				password: password,
			},
			want: validUser,
		},
		{
			name: "hasher error",
			mock: func(repository *MockRepository, hasher *MockPasswordHasher) {
				hasher.EXPECT().Check(password, validUser.Password).Return(errHasher)
				repository.EXPECT().GetByEmail(ctx, email).Return(validUser, nil)
			},
			args: args{
				email:    email,
				password: password,
			},
			want:    user.User{},
			wantErr: true,
		},
		{
			name: "repository error",
			mock: func(repository *MockRepository, hasher *MockPasswordHasher) {
				repository.EXPECT().GetByEmail(ctx, email).Return(user.User{}, errRepository)
			},
			args: args{
				email:    email,
				password: password,
			},
			want:    user.User{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service, repository, passwordHasher := mockService(t)

			tt.mock(repository, passwordHasher)

			got, err := service.Login(ctx, tt.args.email, tt.args.password)
			require.True(t, (err != nil) == tt.wantErr)
			require.True(t, reflect.DeepEqual(tt.want, got))
		})
	}
}

func TestService_UpdateByEmail(t *testing.T) {
	t.Parallel()

	ctx := logger.ContextWithLogger(context.Background(), zap.L())
	email := faker.Email()
	dtoEmail := faker.Email()
	dtoUsername := faker.Username()
	dtoBio := faker.Sentence()
	dtoImage := faker.URL()
	dto := user.UpdateDTO{
		Email:    &dtoEmail,
		Username: &dtoUsername,
		Bio:      &dtoBio,
		Image:    &dtoImage,
	}
	now := time.Now()
	validUser := user.User{
		ID:        uuid.New(),
		Username:  dtoUsername,
		Email:     dtoEmail,
		Password:  faker.Password(),
		Bio:       &dtoBio,
		Image:     &dtoImage,
		CreatedAt: now,
		UpdatedAt: now,
	}

	type args struct {
		email string
		dto   user.UpdateDTO
	}

	tests := []struct {
		name    string
		mock    func(*MockRepository)
		args    args
		want    user.User
		wantErr bool
	}{
		{
			name: "success get user by email",
			mock: func(repository *MockRepository) {
				repository.EXPECT().UpdateByEmail(ctx, email, gomock.Any()).Return(validUser, nil)
			},
			args: args{
				email: email,
				dto:   dto,
			},
			want: validUser,
		},
		{
			name: "repository error",
			mock: func(repository *MockRepository) {
				repository.EXPECT().UpdateByEmail(ctx, email, gomock.Any()).Return(user.User{}, errRepository)
			},
			args: args{
				email: email,
				dto:   dto,
			},
			want:    user.User{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			service, repository, _ := mockService(t)

			tt.mock(repository)

			got, err := service.UpdateByEmail(ctx, tt.args.email, tt.args.dto)
			require.True(t, (err != nil) == tt.wantErr)
			require.True(t, reflect.DeepEqual(tt.want, got))
		})
	}
}
