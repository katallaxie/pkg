package urlx

import (
	"maps"
	"net/url"

	"github.com/katallaxie/pkg/errorx"
)

// CopyValues is merging values in the query string.
func CopyValues(s string, values url.Values) (string, error) {
	u, err := url.Parse(s)
	if err != nil {
		return "", err
	}

	q := u.Query()
	maps.Copy(q, values)

	u.RawQuery = q.Encode()

	return u.String(), nil
}

// RemoveQueryValues is removing values from the query string.
func RemoveQueryValues(s string, keys ...string) (string, error) {
	u, err := url.Parse(s)
	if err != nil {
		return "", err
	}

	q := u.Query()
	for _, key := range keys {
		q.Del(key)
	}

	u.RawQuery = q.Encode()

	return u.String(), nil
}

// MustRemoveQueryValues is removing values from the query string.
func MustRemoveQueryValues(s string, keys ...string) string {
	u, err := RemoveQueryValues(s, keys...)
	errorx.Panic(err)

	return u
}

// MustCopyValues is merging values in the query string.
func MustCopyValues(s string, values url.Values) string {
	u, err := CopyValues(s, values)
	errorx.Panic(err)

	return u
}
