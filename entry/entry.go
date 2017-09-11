package entry

import (
	"bytes"
	"errors"

	"github.com/ChoTotOSS/fluent2gelf/quickmsgpack"
	"github.com/ChoTotOSS/fluent2gelf/quickmsgpack/family"
)

type Entry struct {
	Log    []byte
	Source []byte
	Root   map[string]interface{}
	K8s    map[string]interface{}
}

func New(log []byte, source []byte) *Entry {
	return &Entry{log, source, nil, nil}
}

func (e *Entry) Read(r *quickmsgpack.Reader) error {
	f, b := r.NextFormat()

	if f != family.Map {
		return errors.New("Entry is not map")
	}

	for count := r.NextLengthOf(b); count > 0; count-- {
		//Read for key
		key, err := readString(r)
		if err != nil {
			return err
		}
		if bytes.Compare([]byte("log"), key) == 0 {
			log, err := readString(r)
			if err != nil {
				return err
			}
			e.Log = log
		} else {
			// deep read to entry
		}
	}

	return nil
}

func readString(r *quickmsgpack.Reader) ([]byte, error) {
	f, b := r.NextFormat()
	if f != family.String {
		return nil, errors.New("Value is not a string")
	}
	return r.NextBytes(r.NextLengthOf(b)), nil
}
