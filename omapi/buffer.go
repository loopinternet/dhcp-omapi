package omapi

import (
	"bytes"
	"encoding/binary"
	"sort"

	log "github.com/sirupsen/logrus"
)

type buffer struct {
	buffer *bytes.Buffer
}

func newBuffer() *buffer {
	return &buffer{new(bytes.Buffer)}
}

func (b *buffer) addBytes(data []byte) {
	b.buffer.Write(data)
}

func (b *buffer) add(data interface{}) {
	if err := binary.Write(b.buffer, binary.BigEndian, data); err != nil {
		log.Fatal(err)
	}
}

func (b *buffer) addMap(data map[string][]byte) {
	// We need to add the map in a deterministic order for signing to
	// work, so we first sort the keys in alphabetical order, then use
	// that order to access the map entries.

	keys := make(sort.StringSlice, 0, len(data))

	for key := range data {
		keys = append(keys, key)
	}

	sort.Sort(keys)

	for _, key := range keys {
		value := data[key]

		b.add(int16(len(key)))
		b.add([]byte(key))

		b.add(int32(len(value)))
		if len(value) > 0 {
			b.add(value)
		}
	}

	b.add([]byte("\x00\x00"))
}

func (b *buffer) bytes() []byte {
	return b.buffer.Bytes()
}
