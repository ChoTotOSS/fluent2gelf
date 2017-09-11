package fluentd

import (
	"errors"

	"go.uber.org/zap"

	"github.com/ChoTotOSS/fluent2gelf/entry"
	"github.com/ChoTotOSS/fluent2gelf/quickmsgpack"
	"github.com/ChoTotOSS/fluent2gelf/quickmsgpack/family"
	"github.com/ChoTotOSS/fluent2gelf/quickmsgpack/format"
)

type Forward struct {
	n   int8 //numbers of segments
	r   *quickmsgpack.Reader
	tag []byte
}

func NewForwardReader(r *quickmsgpack.Reader) *Forward {
	return &Forward{0, r, nil}
}

func (fw *Forward) CountSegments() int8 {
	if fw.n == 0 {
		f, b := fw.r.NextFormat() //type should be fixarray

		if f != family.Array || !quickmsgpack.IsFixedFormat(b) {
			return 0
		}
		fw.n = quickmsgpack.FixedValueOf(b)
	}
	return fw.n
}

func (fw *Forward) ReadTag() ([]byte, bool) {
	if fw.CountSegments() == 0 {
		return nil, false
	}

	f, b := fw.r.NextFormat()
	if f != family.String {
		return nil, false
	}
	if quickmsgpack.IsFixedFormat(b) {
		return fw.r.NextBytes(uint(quickmsgpack.FixedValueOf(b))), true
	} else {
		return fw.r.NextBytes(fw.r.NextLength(quickmsgpack.ExtraOf(b))), true
	}
}

func (fw *Forward) CountEntry() uint {
	f, b := fw.r.NextFormat()
	if f != family.Array {
		logger.Warn("Not a array", zap.String("format", format.StringOf(b)))
		return 0
	}

	return fw.r.NextLengthOf(b)
}

func (fw *Forward) ReadEntry() (*entry.Entry, error) {
	if err := skipEntryStruct(fw.r); err != nil {
		return nil, err
	}

	e := new(entry.Entry)

	if err := e.Read(fw.r); err != nil {
		return nil, err
	}

	return e, nil
}

func skipEntryStruct(r *quickmsgpack.Reader) error {
	f, b := r.NextFormat()

	if f != family.Array || quickmsgpack.FixedValueOf(b) != 2 {
		return errors.New("Invalid entry struct")
	}

	f, b = r.NextFormat()
	switch f {
	case family.Integer:
		r.NextBytes(4)
	case family.Extension:
		r.NextBytes(r.NextLengthOf(b) + 1)
	default:
		return errors.New("Not a valid struct for timestamp")
	}
	return nil
}
