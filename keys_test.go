package qrpc

import "testing"

func BenchmarkGenerateKey(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			topic := keySuffix()
			key := newKey(topic)
			array := transformKey(key)

			if keyPrefix(topic) != array[0] {
				b.Fatalf("prefix: %s, key: %s\narray: %v\n", topic, key, array)
			}
		}
	})
}
