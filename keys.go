package qrpc

import (
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"io"
	"time"
)

const (
	keyPrefixLen = 32
)

func newKey(topic string) string {
	return fmt.Sprintf("%s-%s", keyPrefix(topic), keySuffix())
}

func keyPrefix(topic string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(topic)))
}

func keySuffix() string {
	rnd := make([]byte, 2)
	io.ReadFull(rand.Reader, rnd)

	return fmt.Sprintf("%016x-%04x", time.Now().UnixNano(), rnd)
}

func transformKey(key string) []string {
	if len(key) > keyPrefixLen {
		return []string{key[0:keyPrefixLen]}
	}
	return []string{}
}
