package access

import (
	"testing"

	"github.com/katallaxie/pkg/urn"

	"github.com/stretchr/testify/assert"
)

func TestAction_String(t *testing.T) {
	a := Action("iam:changePassword")
	assert.Equal(t, "iam:changePassword", a.String())
}

func TestResource_String(t *testing.T) {
	r := Resource("urn:cloud:machine:eu-central-1:1234567890:ulysses")
	assert.Equal(t, "urn:cloud:machine:eu-central-1:1234567890:ulysses", r.String())
}

func TestResource_URN(t *testing.T) {
	r := Resource("urn:cloud:machine:eu-central-1:1234567890:ulysses")
	u, err := r.URN()
	assert.NoError(t, err)
	assert.Equal(t, urn.Match("urn"), u.Namespace)
	assert.Equal(t, urn.Match("cloud"), u.Partition)
	assert.Equal(t, urn.Match("machine"), u.Service)
	assert.Equal(t, urn.Match("eu-central-1"), u.Region)
	assert.Equal(t, urn.Match("1234567890"), u.Identifier)
	assert.Equal(t, urn.Match("ulysses"), u.Resource)
}

func TestMatcher(t *testing.T) {
	tests := []struct {
		desc    string
		l       *urn.URN
		r       *urn.URN
		matcher Matcher
		expect  bool
	}{
		{
			desc: "return false on mismatched namespace with identity matcher",
			l: &urn.URN{
				Namespace:  urn.Match("urn"),
				Partition:  urn.Match("cloud"),
				Service:    urn.Match("machine"),
				Region:     urn.Match("eu-central-1"),
				Identifier: urn.Match("1234567890"),
				Resource:   urn.Match("ulysses"),
			},
			r: &urn.URN{
				Namespace:  urn.Match("arn"),
				Partition:  urn.Match("cloud"),
				Service:    urn.Match("machine"),
				Region:     urn.Match("eu-central-1"),
				Identifier: urn.Match("1234567890"),
				Resource:   urn.Match("ulysses"),
			},
			matcher: IdentityBasedMatcher,
		},
		{
			desc: "return true on exact match with identity matcher",
			l: &urn.URN{
				Namespace:  urn.Match("urn"),
				Partition:  urn.Match("cloud"),
				Service:    urn.Match("machine"),
				Region:     urn.Match("eu-central-1"),
				Identifier: urn.Match("1234567890"),
				Resource:   urn.Match("ulysses"),
			},
			r: &urn.URN{
				Namespace:  urn.Match("urn"),
				Partition:  urn.Match("cloud"),
				Service:    urn.Match("machine"),
				Region:     urn.Match("eu-central-1"),
				Identifier: urn.Match("1234567890"),
				Resource:   urn.Match("ulysses"),
			},
			matcher: IdentityBasedMatcher,
			expect:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert.Equal(t, tt.expect, tt.matcher(tt.l, tt.r))
		})
	}
}
