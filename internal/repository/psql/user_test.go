package psql_test

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/bxcodec/faker/v3"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"
	"github.com/maypok86/conduit/internal/domain/user"
	"github.com/maypok86/conduit/internal/repository/psql"
	mockPsql "github.com/maypok86/conduit/internal/repository/psql/mocks"
	"github.com/maypok86/conduit/pkg/logger"
	"github.com/maypok86/conduit/pkg/postgres"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var errRepository = errors.New("repository error")

func mockUserRepository(
	t *testing.T,
) (psql.UserRepository, *mockPsql.MockPgxPool, *mockPsql.MockRow) {
	t.Helper()

	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	mockPgxPool := mockPsql.NewMockPgxPool(mockCtl)
	mockRow := mockPsql.NewMockRow(mockCtl)

	queryBuilder := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	db := &postgres.Postgres{
		Builder: queryBuilder,
		Pool:    mockPgxPool,
	}

	userRepository := psql.NewUserRepository(db)

	return userRepository, mockPgxPool, mockRow
}

func TestUserRepository_Create(t *testing.T) {
	t.Parallel()

	ctx := logger.ContextWithLogger(context.Background(), zap.L())
	expectedSQL := "INSERT INTO users (email,username,password) VALUES ($1,$2,$3) RETURNING id"
	dto := user.User{
		Username: faker.Username(),
		Email:    faker.Email(),
		Password: faker.Password(),
	}

	type args struct {
		dto user.User
	}

	tests := []struct {
		name    string
		mock    func(*mockPsql.MockRow, *mockPsql.MockPgxPool)
		args    args
		want    user.User
		wantErr bool
	}{
		{
			name: "creation user",
			mock: func(row *mockPsql.MockRow, pool *mockPsql.MockPgxPool) {
				row.EXPECT().Scan(gomock.Any()).Return(nil).Times(1)
				pool.EXPECT().QueryRow(ctx, expectedSQL, dto.Email, dto.Username, dto.Password).Return(row).Times(1)
			},
			args: args{
				dto: dto,
			},
			want: dto,
		},
		{
			name: "scan error",
			mock: func(row *mockPsql.MockRow, pool *mockPsql.MockPgxPool) {
				row.EXPECT().Scan(gomock.Any()).Return(errRepository).Times(1)
				pool.EXPECT().QueryRow(ctx, expectedSQL, dto.Email, dto.Username, dto.Password).Return(row).Times(1)
			},
			args: args{
				dto: dto,
			},
			want:    user.User{},
			wantErr: true,
		},
		{
			name: "unique violation error",
			mock: func(row *mockPsql.MockRow, pool *mockPsql.MockPgxPool) {
				row.EXPECT().Scan(gomock.Any()).Return(&pgconn.PgError{Code: pgerrcode.UniqueViolation}).Times(1)
				pool.EXPECT().QueryRow(ctx, expectedSQL, dto.Email, dto.Username, dto.Password).Return(row).Times(1)
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

			userRepository, mockPgxPool, mockRow := mockUserRepository(t)

			tt.mock(mockRow, mockPgxPool)

			got, err := userRepository.Create(ctx, tt.args.dto)
			require.True(t, (err != nil) == tt.wantErr)
			require.True(t, reflect.DeepEqual(tt.want, got))
		})
	}
}

func TestUserRepository_GetByEmail(t *testing.T) {
	t.Parallel()

	ctx := logger.ContextWithLogger(context.Background(), zap.L())
	expectedSQL := "SELECT id, username, password, bio, image, created_at, updated_at FROM users WHERE email = $1 LIMIT 1"
	email := faker.Email()
	userEntity := user.User{
		Email: email,
	}

	type args struct {
		email string
	}

	tests := []struct {
		name    string
		mock    func(*mockPsql.MockRow, *mockPsql.MockPgxPool)
		args    args
		want    user.User
		wantErr bool
	}{
		{
			name: "success get by email",
			mock: func(row *mockPsql.MockRow, pool *mockPsql.MockPgxPool) {
				row.EXPECT().Scan(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(nil).Times(1)
				pool.EXPECT().QueryRow(ctx, expectedSQL, email).Return(row).Times(1)
			},
			args: args{
				email: email,
			},
			want: userEntity,
		},
		{
			name: "scan error",
			mock: func(row *mockPsql.MockRow, pool *mockPsql.MockPgxPool) {
				row.EXPECT().Scan(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(errRepository).Times(1)
				pool.EXPECT().QueryRow(ctx, expectedSQL, email).Return(row).Times(1)
			},
			args: args{
				email: email,
			},
			want:    user.User{},
			wantErr: true,
		},
		{
			name: "no rows error",
			mock: func(row *mockPsql.MockRow, pool *mockPsql.MockPgxPool) {
				row.EXPECT().Scan(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(pgx.ErrNoRows).Times(1)
				pool.EXPECT().QueryRow(ctx, expectedSQL, email).Return(row).Times(1)
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

			userRepository, mockPgxPool, mockRow := mockUserRepository(t)

			tt.mock(mockRow, mockPgxPool)

			got, err := userRepository.GetByEmail(ctx, tt.args.email)
			fmt.Println(err)
			require.True(t, (err != nil) == tt.wantErr)
			require.True(t, reflect.DeepEqual(tt.want, got))
		})
	}
}

func TestUserRepository_UpdateByEmail(t *testing.T) {
	t.Parallel()

	ctx := logger.ContextWithLogger(context.Background(), zap.L())
	expectedSQL := "UPDATE users SET username = $1, email = $2, bio = $3, image = $4, updated_at = $5 WHERE email = $6 RETURNING id, username, email, password, bio, image, created_at" //nolint:lll
	email := faker.Email()
	dtoEmail := faker.Email()
	dtoUsername := faker.Username()
	dtoBio := faker.Sentence()
	dtoImage := faker.URL()
	now := time.Now()
	dto := user.UpdateDTO{
		Email:     &dtoEmail,
		Username:  &dtoUsername,
		Bio:       &dtoBio,
		Image:     &dtoImage,
		UpdatedAt: now,
	}
	userEntity := user.User{
		UpdatedAt: now,
	}

	type args struct {
		email string
		dto   user.UpdateDTO
	}

	tests := []struct {
		name    string
		mock    func(*mockPsql.MockRow, *mockPsql.MockPgxPool)
		args    args
		want    user.User
		wantErr bool
	}{
		{
			name: "success update by email",
			mock: func(row *mockPsql.MockRow, pool *mockPsql.MockPgxPool) {
				row.EXPECT().Scan(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(nil).Times(1)
				pool.EXPECT().
					QueryRow(ctx, expectedSQL, *dto.Username, *dto.Email, *dto.Bio, *dto.Image, dto.UpdatedAt, email).
					Return(row).
					Times(1)
			},
			args: args{
				email: email,
				dto:   dto,
			},
			want: userEntity,
		},
		{
			name: "scan error",
			mock: func(row *mockPsql.MockRow, pool *mockPsql.MockPgxPool) {
				row.EXPECT().Scan(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(errRepository).Times(1)
				pool.EXPECT().
					QueryRow(ctx, expectedSQL, *dto.Username, *dto.Email, *dto.Bio, *dto.Image, dto.UpdatedAt, email).
					Return(row).
					Times(1)
			},
			args: args{
				email: email,
				dto:   dto,
			},
			want:    user.User{},
			wantErr: true,
		},
		{
			name: "no rows error",
			mock: func(row *mockPsql.MockRow, pool *mockPsql.MockPgxPool) {
				row.EXPECT().Scan(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(pgx.ErrNoRows).Times(1)
				pool.EXPECT().
					QueryRow(ctx, expectedSQL, *dto.Username, *dto.Email, *dto.Bio, *dto.Image, dto.UpdatedAt, email).
					Return(row).
					Times(1)
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

			userRepository, mockPgxPool, mockRow := mockUserRepository(t)

			tt.mock(mockRow, mockPgxPool)

			got, err := userRepository.UpdateByEmail(ctx, tt.args.email, tt.args.dto)
			fmt.Println(err)
			require.True(t, (err != nil) == tt.wantErr)
			require.True(t, reflect.DeepEqual(tt.want, got))
		})
	}
}
