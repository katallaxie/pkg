package access

import (
	"encoding/json"
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

func TestDefaultServices(t *testing.T) {
	s := Services{
		"access": true,
	}

	assert.Equal(t, s, DefaultServices)

	s.Add("k8s")
	DefaultServices.Add("k8s")

	assert.Equal(t, s, DefaultServices)
}

func TestPolicyUnmarshalJSON(t *testing.T) {
	tests := []struct {
		desc string
		json json.RawMessage
		pol  *Policy
		err  error
	}{
		{
			desc: "suceessfully unmarshal policy",
			json: json.RawMessage(`{"version":"2023-03-28", "rules": []}`),
			pol: &Policy{
				Version: DefaultVersion,
				Rules:   Rules{},
			},
		},
		{
			desc: "suceessfully unmarshal rules",
			json: json.RawMessage(`{
				"version":"2023-03-28",
				"rules": [
					{
						"id": "1",
						"effect": "allow",
						"resources": [
							"ionos:k8s:de-fra:12345678910:/cluster/12345678910"
						],
						"actions": [
							"k8s:cluster-create"
						]
					}
				]
			}`),
			pol: &Policy{
				Version: DefaultVersion,
				Rules: Rules{
					{
						ID:     "1",
						Effect: Allow,
						Resources: Resources{
							"ionos:k8s:de-fra:12345678910:/cluster/12345678910",
						},
						Actions: Actions{
							"k8s:cluster-create",
						},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			p := DefaultPolicy()
			err := p.UnmarshalJSON(test.json)

			if test.err != nil {
				assert.EqualError(t, err, test.err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, test.pol.Version, p.Version)
			assert.Equal(t, test.pol.Rules, p.Rules)
		})
	}
}

func BenchmarkPolicyUnmarshal(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		var p Policy
		_ = json.Unmarshal([]byte(`{
			"name": "changePassword",
			"version": "2023-03-28",
			"description": "Allow users to change their password",
			"rules": [
				{
					"actions": [
						"iam:changePassword"
					],
					"effect": "Allow",
					"resources": [
						"urn:cloud:machine:eu-central-1:1234567890:root"
					]
				}
			]
		}`), &p)
	}
}

func TestIs(t *testing.T) {
	tests := []struct {
		desc     string
		urn      *urn.URN
		resource ResourceIdentifier
		expect   bool
	}{
		{
			desc: "return true on matching user",
			urn: &urn.URN{
				Namespace:  urn.Match("urn"),
				Partition:  urn.Match("cloud"),
				Service:    urn.Match("access"),
				Region:     urn.Match("eu-central-1"),
				Identifier: urn.Match("1234567890"),
				Resource:   urn.Match("users/ulysses"),
			},
			resource: UserResourceIdentifier,
			expect:   true,
		},
		{
			desc: "return false on non matching service",
			urn: &urn.URN{
				Namespace:  urn.Match("urn"),
				Partition:  urn.Match("cloud"),
				Service:    urn.Match("k8s"),
				Region:     urn.Match("eu-central-1"),
				Identifier: urn.Match("1234567890"),
				Resource:   urn.Match("users/ulysses"),
			},
			resource: UserResourceIdentifier,
			expect:   false,
		},
		{
			desc: "return false on non matching user",
			urn: &urn.URN{
				Namespace:  urn.Match("urn"),
				Partition:  urn.Match("cloud"),
				Service:    urn.Match("access"),
				Region:     urn.Match("eu-central-1"),
				Identifier: urn.Match("1234567890"),
				Resource:   urn.Match("groups/ulysses"),
			},
			resource: UserResourceIdentifier,
			expect:   false,
		},

		{
			desc: "return true on matching group",
			urn: &urn.URN{
				Namespace:  urn.Match("urn"),
				Partition:  urn.Match("cloud"),
				Service:    urn.Match("access"),
				Region:     urn.Match("eu-central-1"),
				Identifier: urn.Match("1234567890"),
				Resource:   urn.Match("groups/ulysses"),
			},
			resource: GroupResourceIdentifier,
			expect:   true,
		},
		{
			desc: "return false on non matching group",
			urn: &urn.URN{
				Namespace:  urn.Match("urn"),
				Partition:  urn.Match("cloud"),
				Service:    urn.Match("access"),
				Region:     urn.Match("eu-central-1"),
				Identifier: urn.Match("1234567890"),
				Resource:   urn.Match("users/ulysses"),
			},
			resource: GroupResourceIdentifier,
			expect:   false,
		},
		{
			desc: "return true on matching role",
			urn: &urn.URN{
				Namespace: urn.Match("urn"),

				Partition:  urn.Match("cloud"),
				Service:    urn.Match("access"),
				Region:     urn.Match("eu-central-1"),
				Identifier: urn.Match("1234567890"),
				Resource:   urn.Match("roles/ulysses"),
			},
			resource: RoleResourceIdentifier,
			expect:   true,
		},
		{
			desc: "return false on non matching role",
			urn: &urn.URN{
				Namespace: urn.Match("urn"),
				Partition: urn.Match("cloud"),
				Service:   urn.Match("access"),

				Region:     urn.Match("eu-central-1"),
				Identifier: urn.Match("1234567890"),
				Resource:   urn.Match("users/ulysses"),
			},
			resource: RoleResourceIdentifier,
			expect:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			assert.Equal(t, test.expect, Is(test.urn, test.resource))
		})
	}
}

func BenchmarkIs_UserResouceIdentifier(b *testing.B) {
	b.ReportAllocs()

	urn := &urn.URN{
		Namespace:  urn.Match("urn"),
		Partition:  urn.Match("cloud"),
		Service:    urn.Match("access"),
		Region:     urn.Match("eu-central-1"),
		Identifier: urn.Match("1234567890"),
		Resource:   urn.Match("users/ulysses"),
	}

	for i := 0; i < b.N; i++ {
		Is(urn, UserResourceIdentifier)
	}
}

func BenchmarkIs_RoleResouceIdentifier(b *testing.B) {
	b.ReportAllocs()

	urn := &urn.URN{
		Namespace:  urn.Match("urn"),
		Partition:  urn.Match("cloud"),
		Service:    urn.Match("access"),
		Region:     urn.Match("eu-central-1"),
		Identifier: urn.Match("1234567890"),
		Resource:   urn.Match("roles/ulysses"),
	}

	for i := 0; i < b.N; i++ {
		Is(urn, RoleResourceIdentifier)
	}
}

func BenchmarkIs_GroupResouceIdentifier(b *testing.B) {
	b.ReportAllocs()

	urn := &urn.URN{
		Namespace:  urn.Match("urn"),
		Partition:  urn.Match("cloud"),
		Service:    urn.Match("access"),
		Region:     urn.Match("eu-central-1"),
		Identifier: urn.Match("1234567890"),
		Resource:   urn.Match("groups/ulysses"),
	}

	for i := 0; i < b.N; i++ {
		Is(urn, GroupResourceIdentifier)
	}
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

func BenchmarkMatcher_IdentityBased(b *testing.B) {
	l := &urn.URN{
		Namespace:  urn.Match("urn"),
		Partition:  urn.Match("cloud"),
		Service:    urn.Match("machine"),
		Region:     urn.Match("eu-central-1"),
		Identifier: urn.Match("1234567890"),
		Resource:   urn.Match("ulysses"),
	}
	r := &urn.URN{
		Namespace:  urn.Match("urn"),
		Partition:  urn.Match("cloud"),
		Service:    urn.Match("machine"),
		Region:     urn.Match("eu-central-1"),
		Identifier: urn.Match("1234567890"),
		Resource:   urn.Match("ulysses"),
	}

	for i := 0; i < b.N; i++ {
		IdentityBasedMatcher(l, r)
	}
}
