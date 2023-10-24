package snowflake

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	_, err := New(0)
	require.NoError(t, err)

	_, err = New(5000)
	require.Error(t, err)
}

func BenchmarkGenerate(b *testing.B) {
	node, err := New(1)
	require.NoError(b, err)

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
	require.NoError(b, err)

	b.ReportAllocs()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		id := node.Generate()
		id.ProtoMessage()
	}
}
