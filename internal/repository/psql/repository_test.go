package psql_test

import (
	"reflect"
	"testing"

	sq "github.com/Masterminds/squirrel"
	"github.com/maypok86/conduit/internal/repository/psql"
	"github.com/maypok86/conduit/pkg/postgres"
	"github.com/stretchr/testify/require"
)

func TestNewRepositories(t *testing.T) {
	t.Parallel()

	db := &postgres.Postgres{
		Builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}

	want := psql.Repositories{
		User:    psql.NewUserRepository(db),
		Profile: psql.NewProfileRepository(db),
	}

	got := psql.NewRepositories(db)

	require.True(t, reflect.DeepEqual(want, got))
}
