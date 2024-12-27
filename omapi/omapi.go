// Package omapi implements the OMAPI protocol of the ISC DHCP server,
// allowing it to query and modify objects as well as control the
// server itself.
package omapi

import (
	"encoding/binary"
	"math/rand"
	"sync"
	"time"
)

const DefaultPort = 7911

var (
	True  = []byte{0, 0, 0, 1}
	False = []byte{0, 0, 0, 0}
)

type syncRng struct {
	sync.Mutex
	*rand.Rand
}

var rng = syncRng{sync.Mutex{}, rand.New(rand.NewSource(time.Now().UTC().UnixNano()))}

func bytesToInt32(b []byte) int32 {
	if len(b) < 4 {
		return 0
	}

	return int32(binary.BigEndian.Uint32(b))
}

func int32ToBytes(i int32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(i))

	return b
}
