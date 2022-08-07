package snowflake

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {

	_, err := New(0)
	assert.NoError(t, err)

	_, err = New(5000)
	assert.Error(t, err)
}

func BenchmarkGenerate(b *testing.B) {
	node, err := New(1)
	assert.NoError(b, err)

	b.ReportAllocs()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = node.Generate()
	}
}

func BenchmarkGenerateMaxSequence(b *testing.B) {
	NodeBits = 1
	StepBits = 21
	node, _ := New(1)

	b.ReportAllocs()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = node.Generate()
	}
}

func BenchmarkGenerateProtoMessage(b *testing.B) {
	node, err := New(1)
	assert.NoError(b, err)

	b.ReportAllocs()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		id := node.Generate()
		id.ProtoMessage()
	}
}
