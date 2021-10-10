package util

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
