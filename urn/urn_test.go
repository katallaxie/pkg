package urn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		desc        string
		namespace   string
		partition   string
		service     string
		region      string
		identifier  string
		resource    string
		expected    *URN
		expectedErr bool
	}{
		{
			desc:       "returns a resource URN",
			namespace:  "urn",
			partition:  "cloud",
			service:    "machine",
			region:     "eu-central-1",
			identifier: "1234567890",
			resource:   "ulysses",
			expected: &URN{
				Namespace:  "urn",
				Partition:  "cloud",
				Service:    "machine",
				Region:     "eu-central-1",
				Identifier: "1234567890",
				Resource:   "ulysses",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			urn, err := New(tc.namespace, tc.partition, tc.service, tc.region, tc.identifier, tc.resource)

			if tc.expectedErr {
				assert.Error(t, err)
			}
			assert.Equal(t, tc.expected, urn)
		})
	}
}

func TestMatch(t *testing.T) {
	tests := []struct {
		desc     string
		urn      *URN
		other    *URN
		expected bool
	}{
		{
			desc: "returns true when the URNs are equal",
			urn: &URN{
				Namespace:  "urn",
				Partition:  "cloud",
				Service:    "machine",
				Region:     "eu-central-1",
				Identifier: "1234567890",
				Resource:   "ulysses",
			},
			other: &URN{
				Namespace:  "urn",
				Partition:  "cloud",
				Service:    "machine",
				Region:     "eu-central-1",
				Identifier: "1234567890",
				Resource:   "ulysses",
			},
			expected: true,
		},
		{
			desc: "returns true when the URNs are equal and the other has a wildcard",
			urn: &URN{
				Namespace:  "urn",
				Partition:  "cloud",
				Service:    "machine",
				Region:     "eu-central-1",
				Identifier: "1234567890",
				Resource:   "*",
			},
			other: &URN{
				Namespace:  "urn",
				Partition:  "cloud",
				Service:    "machine",
				Region:     "eu-central-1",
				Identifier: "1234567890",
				Resource:   "ulysses",
			},
			expected: true,
		},
		{
			desc: "returns true when the URNs are equal and the other has a wildcard",
			urn: &URN{
				Namespace:  "urn",
				Partition:  "cloud",
				Service:    "machine",
				Region:     "eu-central-1",
				Identifier: "1234567890",
				Resource:   "",
			},
			other: &URN{
				Namespace:  "urn",
				Partition:  "cloud",
				Service:    "machine",
				Region:     "eu-central-1",
				Identifier: "1234567890",
				Resource:   "ulysses",
			},
			expected: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.urn.Match(tc.other))
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		desc        string
		urnStr      string
		expected    *URN
		expectedErr bool
	}{
		{
			desc:   "returns an ErrorInvalid when missing the namespace",
			urnStr: ":cloud:machine::1234567890:ulysses",
		},
		{
			desc:   "returns an ErrorInvalid when identifier is longer than 256 chars",
			urnStr: "urn:collection:::aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaas:123456",
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			_, err := Parse(tc.urnStr)

			assert.Error(t, err)
		})
	}
}
