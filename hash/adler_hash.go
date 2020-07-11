package hash

import "hash/adler32"

type AdlerHash struct {
}

func NewAdlerHash() *AdlerHash {
	return &AdlerHash{}
}

func (d *AdlerHash) Hash(key string) uint32 {
	ret := adler32.Checksum([]byte(key))
	return ret
}
