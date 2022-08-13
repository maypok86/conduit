package user_test

import (
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/maypok86/conduit/internal/domain/user"
	"github.com/stretchr/testify/require"
)

func TestUser_GetBio(t *testing.T) {
	t.Parallel()

	bio := faker.Sentence()

	tests := []struct {
		name string
		user user.User
		want string
	}{
		{
			name: "without bio",
			user: user.User{},
			want: "",
		},
		{
			name: "with bio",
			user: user.User{Bio: &bio},
			want: bio,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tt.want, tt.user.GetBio())
		})
	}
}

func TestUser_GetImage(t *testing.T) {
	t.Parallel()

	image := faker.URL()

	tests := []struct {
		name string
		user user.User
		want string
	}{
		{
			name: "without image",
			user: user.User{},
			want: "",
		},
		{
			name: "with image",
			user: user.User{Image: &image},
			want: image,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tt.want, tt.user.GetImage())
		})
	}
}
