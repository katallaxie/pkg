package store

import (
	"bytes"
	"encoding/gob"
	"errors"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

// Store represents the key/value store.
// Use the Open() method to open one, and Close() to close one.
type Store struct {
	db *leveldb.DB
}

var (
	// ErrKeyNotExist is returned when the supplied key to Get() or Delete()
	// does not exists in the database.
	ErrKeyNotExist = errors.New("store: key does not exist")

	// ErrBadValue is returned when the value supplied to Put() in the database
	// could not be encoded or is nil.
	ErrBadValue = errors.New("store: bad value")
)

// Open a key/value store. "path" is the full path to LevelDB database file.
// The DB will be created if not exist, unless ErrorIfMissing is true.
// You can pass in options to LevelDB via "o".
func Open(path string, o *opt.Options) (*Store, error) {
	var err error

	db, err := leveldb.OpenFile(path, o)
	if err != nil {
		return nil, err
	}

	store := new(Store)
	store.db = db

	return store, err
}

// Get a value from the store. "value" must be a pointer to a type,
// because of underlying gob. If the key is not present in the store,
// an ErrKeyNotExist is returned.
//
//		type MyStruct struct {
//	    	Numbers []int
//		}
//		var val MyStruct
//		if err := s.Get(store.Byte("key42"), &val); err == store.ErrKeyNotExist {
//	    	// "key42" not found
//		}
//		if err != nil {
//			// some other error occurred
//		}
//
// The value passed to Get() can be nil. Which gives you the same result as
// Exist().
//
//		if err := s.Get(store.Byte("key42"), nil); err == nil {
//	    	fmt.Println("entry is present")
//		}
func (s *Store) Get(key []byte, value interface{}) error {
	var err error

	v, err := s.db.Get(key, nil)
	if err != nil || value == nil {
		return ErrKeyNotExist
	}

	d := gob.NewDecoder(bytes.NewReader(v))
	return d.Decode(value)
}

// Put is creating an entry in the store. It overwrites any previous existing
// value for the given key. The passed value is gob-encoded and stored.
// The value cannot be nil and Put() returns an ErrBadValue.
//
//	err := s.Put(store.Byte("key42"), 42)
//	err := s.Put(store.Byte("key42"), "the question to live and the universe.")
//	m := map[string]int{
//		"foo": 0,
//		"bar": 1
//	}
//	err := s.Put(store.Byte("key42"), m)
func (s *Store) Put(key []byte, value interface{}) error {
	var err error

	if value == nil {
		return ErrBadValue
	}

	var b bytes.Buffer
	if err = gob.NewEncoder(&b).Encode(value); err != nil {
		return err
	}

	return s.db.Put(key, b.Bytes(), nil)
}

// Delete an entry of a given key from the store.
// If no such key is present in the store it will not return an error.
//
//	s.Delete(store.Byte("key42"))
func (s *Store) Delete(key []byte) error {
	return s.db.Delete(key, nil)
}

// Close closes the key/value database file.
func (s *Store) Close() error {
	return s.db.Close()
}

// Byte is converting a string to []byte stream.
func Byte(b string) []byte {
	return []byte(b)
}
