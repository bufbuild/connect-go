package connect

import (
	"github.com/bufconnect/connect/codec"
	"github.com/bufconnect/connect/codec/protobuf"
)

// readOnlyCodecs is a read-only interface to a map of named codecs.
type readOnlyCodecs interface {
	Get(string) codec.Codec
	Protobuf() codec.Codec
	Names() []string
}

type codecMap struct {
	codecs map[string]codec.Codec
}

func newReadOnlyCodecs(m map[string]codec.Codec) *codecMap {
	return &codecMap{m}
}

// Get the named codec.
func (m *codecMap) Get(name string) codec.Codec {
	return m.codecs[name]
}

// Get the user-supplied protobuf codec, falling back to the default
// implementation if necessary.
//
// This is helpful in the gRPC protocol, where the wire protocol requires
// marshaling protobuf structs to binary even if the RPC procedures were
// generated from a different IDL.
func (m *codecMap) Protobuf() codec.Codec {
	if pb, ok := m.codecs[protobuf.NameBinary]; ok {
		return pb
	}
	return protobuf.NewBinary()
}

// Names returns a copy of the registered codec names. The returned slice is
// safe for the caller to mutate.
func (m *codecMap) Names() []string {
	names := make([]string, 0, len(m.codecs))
	for name := range m.codecs {
		names = append(names, name)
	}
	return names
}
