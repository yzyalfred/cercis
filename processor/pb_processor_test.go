package processor

import (
	"testing"
)

func BenchmarkPBProcessor_Unmarshalarshal(b *testing.B) {
	processor := NewPBProcessor()
	processor.Register(uint32(1), &RpcSayHello{})
	msg := &RpcSayHello{Name: "yzy"}
	buf, _ := processor.Marshal(uint32(1), msg)
	for i := 0; i < b.N; i++ {
		processor.Unmarshal(buf)
	}
}

func BenchmarkPBProcessor_Marshalarshal(b *testing.B) {
	processor := NewPBProcessor()
	msg := &RpcSayHello{Name: "yzy"}
	for i := 0; i < b.N; i++ {
		processor.Marshal(uint32(1), msg)
	}
}