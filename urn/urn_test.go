package urn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		desc          string
		namespace     string
		collection    string
		identifier    string
		resource      string
		validateFuncs []ValidateFunc
		expected      *URN
		expectedErr   error
	}{
		{
			desc:       "returns a resource URN",
			namespace:  "urn",
			collection: "machine",
			identifier: "1234567890",
			resource:   "ulysses",
			expected: &URN{
				namespace:  "urn",
				collection: "machine",
				identifier: "1234567890",
				resource:   "ulysses",
			},
		},
		{
			desc:       "allows letters, numbers, dashes, underscores, dots and tildes in identifiers",
			namespace:  "urn",
			collection: "machine",
			identifier: "this_could-be-pr377y~much.anything",
			expected: &URN{
				namespace:  "urn",
				collection: "machine",
				identifier: "this_could-be-pr377y~much.anything",
			},
		},
		{
			desc:          "allows letters, numbers, dashes, underscores, dots and tildes in identifiers",
			namespace:     "",
			collection:    "machine",
			identifier:    "this_could-be-pr377y~much.anything",
			validateFuncs: []ValidateFunc{ValidateNamespace()},
			expectedErr:   ErrorInvalid,
			expected:      nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			urn, err := New(tc.namespace, tc.collection, tc.identifier, tc.resource, tc.validateFuncs...)

			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expected, urn)
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		desc          string
		urnStr        string
		expected      *URN
		validateFuncs []ValidateFunc
		expectedErr   error
	}{
		{
			desc:   "parses a machine URN",
			urnStr: "urn:machine:1234567890:ulysses",
			expected: &URN{
				namespace:  "urn",
				collection: "machine",
				identifier: "1234567890",
				resource:   "ulysses",
			},
		},
		{
			desc:   "allows letters, numbers, dashes, underscores, dots and tildes in namespaces and identifiers",
			urnStr: "n~4m3-5p4_3.:machine:this_could-be-pr377y~much.anything:ulysses",
			expected: &URN{
				namespace:  "n~4m3-5p4_3.",
				collection: "machine",
				identifier: "this_could-be-pr377y~much.anything",
				resource:   "ulysses",
			},
		},
		{
			desc:        "returns an ErrorInvalid when missing the namespace",
			urnStr:      ":machine:1234567890:ulysses",
			expectedErr: ErrorInvalid,
		},
		{
			desc:        "returns an ErrorInvalid when namespace contains invalid characters",
			urnStr:      "names#machine:1234567890:ulysses",
			expectedErr: ErrorInvalid,
		},
		{
			desc:        "returns an ErrorInvalid when namespace is longer than 256 chars",
			urnStr:      "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaas:machine:1234567890:ulysses",
			expectedErr: ErrorInvalid,
		},
		{
			desc:        "returns an ErrorInvalid when collection contains invalid characters",
			urnStr:      "urn:colle(tion:1234567890:ulysses",
			expectedErr: ErrorInvalid,
		},
		{
			desc:        "returns an ErrorInvalid when missing the identifier",
			urnStr:      "urn:collection:",
			expectedErr: ErrorInvalid,
		},
		{
			desc:        "returns an ErrorInvalid when identifier contains invalid characters",
			urnStr:      "urn:collection:1234$",
			expectedErr: ErrorInvalid,
		},
		{
			desc:        "returns an ErrorInvalid when resource contains invalid characters",
			urnStr:      "urn:collection:identifier:1234$",
			expectedErr: ErrorInvalid,
		},
		{
			desc:        "returns an ErrorInvalid when identifier is longer than 256 chars",
			urnStr:      "urn:collection:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaas:123456",
			expectedErr: ErrorInvalid,
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			urn, err := Parse(tc.urnStr)

			assert.Equal(t, tc.expectedErr, err)
			assert.Equal(t, tc.expected, urn)
		})
	}
}
