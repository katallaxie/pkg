package urn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	tests := []struct {
		desc string
		urn  *ResourceURN
		err  error
	}{
		{
			desc: "should not return an error, if fully populated",
			urn: &ResourceURN{
				Canonical:  "urn:machine:1234:ulysses",
				Namespace:  "urn",
				Collection: "machine",
				Identifier: "1234",
				Resource:   "ulysses",
			},
		},
		{
			desc: "should return error on missing canonical",
			urn: &ResourceURN{
				Canonical:  "",
				Namespace:  "urn",
				Collection: "machine",
				Identifier: "1234",
				Resource:   "ulysses",
			},
			err: ErrMissingCanonical,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			err := tt.urn.Validate()

			if tt.err == nil {
				assert.NoError(t, err)

				return
			}

			assert.Error(t, err)
			assert.Equal(t, tt.err, err)
		})
	}
}
