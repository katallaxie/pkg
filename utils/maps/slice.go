package maps

import "strings"

// FromSlice ...
func FromSlice(s []string) map[string]string {
	m := make(map[string]string)
	for _, v := range s {
		kv := strings.Split(v, "=")

		if len(kv) != 2 {
			m[kv[0]] = ""
			continue
		}

		m[kv[0]] = kv[1]
	}

	return m
}
