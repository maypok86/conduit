package psql_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/bxcodec/faker/v3"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgx/v4"
	"github.com/maypok86/conduit/internal/domain/profile"
	"github.com/maypok86/conduit/internal/repository/psql"
	mockPsql "github.com/maypok86/conduit/internal/repository/psql/mocks"
	"github.com/maypok86/conduit/pkg/logger"
	"github.com/maypok86/conduit/pkg/postgres"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var errProfileRepository = errors.New("profile repository error")

func mockProfileRepository(
	t *testing.T,
) (psql.ProfileRepository, *mockPsql.MockPgxPool, *mockPsql.MockRow) {
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

	profileRepository := psql.NewProfileRepository(db)

	return profileRepository, mockPgxPool, mockRow
}

func TestProfileRepository_GetByUsername(t *testing.T) {
	t.Parallel()

	ctx := logger.ContextWithLogger(context.Background(), zap.L())
	expectedSQL := "SELECT bio, image, created_at, updated_at FROM users WHERE username = $1 LIMIT 1"
	username := faker.Username()
	profileEntity := profile.Profile{
		Username: username,
	}

	type args struct {
		username string
	}

	tests := []struct {
		name    string
		mock    func(*mockPsql.MockRow, *mockPsql.MockPgxPool)
		args    args
		want    profile.Profile
		wantErr bool
	}{
		{
			name: "success get by username",
			mock: func(row *mockPsql.MockRow, pool *mockPsql.MockPgxPool) {
				row.EXPECT().Scan(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(nil).Times(1)
				pool.EXPECT().QueryRow(ctx, expectedSQL, username).Return(row).Times(1)
			},
			args: args{
				username: username,
			},
			want: profileEntity,
		},
		{
			name: "scan error",
			mock: func(row *mockPsql.MockRow, pool *mockPsql.MockPgxPool) {
				row.EXPECT().Scan(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(errProfileRepository).Times(1)
				pool.EXPECT().QueryRow(ctx, expectedSQL, username).Return(row).Times(1)
			},
			args: args{
				username: username,
			},
			want:    profile.Profile{},
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
				).Return(pgx.ErrNoRows).Times(1)
				pool.EXPECT().QueryRow(ctx, expectedSQL, username).Return(row).Times(1)
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

			profileRepository, mockPgxPool, mockRow := mockProfileRepository(t)

			tt.mock(mockRow, mockPgxPool)

			got, err := profileRepository.GetByUsername(ctx, tt.args.username)
			require.True(t, (err != nil) == tt.wantErr)
			require.True(t, reflect.DeepEqual(tt.want, got))
		})
	}
}
