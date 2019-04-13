package osc

import "testing"

func _prepareArguments() *Arguments {
	args := NewArguments()
	args.AddString("play")
	args.AddInt32(42)
	args.AddInt64(42)
	args.AddASCIICharacter('g')
	args.AddASCIICharacter('o')
	args.AddFloat32(42.0)
	args.AddFloat64(42.0)
	args.AddBlob([]byte("musica"))
	return args
}

func _marshalArguments() ([]byte, error) {
	return _prepareArguments().MarshalBinary()
}

func BenchmarkArguments(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_marshalArguments()
	}
}
