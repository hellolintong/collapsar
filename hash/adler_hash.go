package hash

import (
	"hash/adler32"
	"unsafe"
)

type AdlerHash struct {
}

func NewAdlerHash() *AdlerHash {
	return &AdlerHash{}
}

func (d *AdlerHash) Hash(key string) uint64 {
	ret := adler32.Checksum(*(*[]byte)(unsafe.Pointer(&key)))
	return uint64(ret)
}