package uuid

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"time"
)

type UUID [16]byte

var Nil UUID = UUID{}

func New4() (UUID, error) {
	var uuid UUID
	if _, err := rand.Read(uuid[:]); err != nil {
		return Nil, err
	}

	uuid[6] = (uuid[6] & 0x0f) | 0x40 // Version 4
	uuid[8] = (uuid[8] & 0x3f) | 0x80

	return uuid, nil
}

func New5(namespace UUID, name string) UUID {
	h := sha1.New()
	h.Write(namespace[:])
	h.Write([]byte(name))
	sum := h.Sum(nil)

	var uuid UUID
	copy(uuid[:], sum[:16])

	uuid[6] = (uuid[6] & 0x0f) | 0x50 // Version 5
	uuid[8] = (uuid[8] & 0x3f) | 0x80

	return uuid
}

func New7() (UUID, error) {
	var uuid UUID

	timestamp := uint64(time.Now().UnixMilli())
	binary.BigEndian.PutUint64(uuid[0:8], (timestamp << 16))
	if _, err := rand.Read(uuid[8:]); err != nil {
		return Nil, err
	}

	uuid[6] = (uuid[6] & 0x0f) | 0x70 // Version 7
	uuid[8] = (uuid[8] & 0x3f) | 0x80

	return uuid, nil
}

func (this *UUID) Bytes() []byte { return this[:] }

func (this *UUID) String() string {
	return fmt.Sprintf(
		"%08x-%04x-%04x-%04x-%012x",
		this[0:4], this[4:6], this[6:8], this[8:10], this[10:16],
	)
}
