package profile_test

import (
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/maypok86/conduit/internal/domain/profile"
	"github.com/stretchr/testify/require"
)

func TestUser_GetBio(t *testing.T) {
	t.Parallel()

	bio := faker.Sentence()

	tests := []struct {
		name    string
		profile profile.Profile
		want    string
	}{
		{
			name:    "without bio",
			profile: profile.Profile{},
			want:    "",
		},
		{
			name:    "with bio",
			profile: profile.Profile{Bio: &bio},
			want:    bio,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tt.want, tt.profile.GetBio())
		})
	}
}

func TestUser_GetImage(t *testing.T) {
	t.Parallel()

	image := faker.URL()

	tests := []struct {
		name    string
		profile profile.Profile
		want    string
	}{
		{
			name:    "without image",
			profile: profile.Profile{},
			want:    "",
		},
		{
			name:    "with image",
			profile: profile.Profile{Image: &image},
			want:    image,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tt.want, tt.profile.GetImage())
		})
	}
}
