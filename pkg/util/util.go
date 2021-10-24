package util

import (
	"encoding/binary"
	"strconv"

	"github.com/spaolacci/murmur3"
)

func IntOr(v *int, d int) int {
	if v != nil {
		return *v
	}
	return d
}

func Int64Or(v *int64, d int64) int64 {
	if v != nil {
		return *v
	}
	return d
}

func BoolOr(v *bool, d bool) bool {
	if v != nil {
		return *v
	}
	return d
}

func ContainsInt64(s []int64, v int64) bool {
	for _, e := range s {
		if v == e {
			return true
		}
	}
	return false
}

func FormatID(id int64) string {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(id))
	return strconv.FormatInt(int64(murmur3.Sum32(b)), 16)
}
