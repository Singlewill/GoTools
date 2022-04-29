package netlink

import (
	"encoding/binary"
	"errors"
	"fmt"
	"syscall"
)

type Attribute struct {
	syscall.NlAttr
	Data []byte
}
type AttributeDecoder struct {
	// ByteOrder defines a specific byte order to use when processing integer
	// attributes.  ByteOrder should be set immediately after creating the
	// AttributeDecoder: before any attributes are parsed.
	//
	// If not set, the native byte order will be used.
	ByteOrder binary.ByteOrder

	// The current attribute being worked on.
	a Attribute

	// The slice of input bytes and its iterator index.
	b []byte
	i int

	length int

	// Any error encountered while decoding attributes.
	err error
}

// errInvalidAttribute specifies if an Attribute's length is incorrect.
var errInvalidAttribute = errors.New("invalid attribute; length too short or too large")

// NewAttributeDecoder creates an AttributeDecoder that unpacks Attributes
// from b and prepares the decoder for iteration.
func NewAttributeDecoder(b []byte, order binary.ByteOrder) (*AttributeDecoder, error) {
	ad := &AttributeDecoder{
		ByteOrder: order,
		b:         b,
	}

	var err error
	ad.length, err = ad.available()
	if err != nil {
		return nil, err
	}

	return ad, nil
}

// count scans the input slice to count the number of netlink attributes
// that could be decoded by Next().
func (ad *AttributeDecoder) available() (int, error) {
	var count int
	for i := 0; i < len(ad.b); {
		// Make sure there's at least a header's worth
		// of data to read on each iteration.
		if len(ad.b[i:]) < NlaHeaderLen {
			return 0, errInvalidAttribute
		}

		// Extract the length of the attribute.
		//l := int(nlenc.Uint16(ad.b[i : i+2]))
		l := int(ad.ByteOrder.Uint16(ad.b[i : i+2]))

		// Ignore zero-length attributes.
		if l != 0 {
			count++
		}

		// Advance by at least a header's worth of bytes.
		if l < NlaHeaderLen {
			l = NlaHeaderLen
		}

		i += NlaAlign(l)
	}

	return count, nil
}

// unmarshal unmarshals the contents of a byte slice into an Attribute.
func (a *Attribute) unmarshal(b []byte, order binary.ByteOrder) error {
	if len(b) < NlaHeaderLen {
		return errInvalidAttribute
	}

	a.Len = order.Uint16(b[0:2])
	a.Type = order.Uint16(b[2:4])

	if int(a.Len) > len(b) {
		return errInvalidAttribute
	}

	switch {
	// No length, no data
	case a.Len == 0:
		a.Data = make([]byte, 0)
	// Not enough length for any data
	case int(a.Len) < NlaHeaderLen:
		return errInvalidAttribute
	// Data present
	case int(a.Len) >= NlaHeaderLen:
		a.Data = make([]byte, len(b[NlaHeaderLen:a.Len]))
		copy(a.Data, b[NlaHeaderLen:a.Len])
	}

	return nil
}

// Next advances the decoder to the next netlink attribute.  It returns false
// when no more attributes are present, or an error was encountered.
func (ad *AttributeDecoder) Next() bool {
	if ad.err != nil {
		// Hit an error, stop iteration.
		fmt.Println("ad.err first")
		return false
	}

	// Exit if array pointer is at or beyond the end of the slice.
	if ad.i >= len(ad.b) {
		return false
	}

	if err := ad.a.unmarshal(ad.b[ad.i:], ad.ByteOrder); err != nil {
		ad.err = err
		fmt.Println(err)
		return false
	}

	// Advance the pointer by at least one header's length.
	if int(ad.a.Len) < NlaHeaderLen {
		ad.i += NlaHeaderLen
	} else {
		ad.i += NlaAlign(int(ad.a.Len))
	}

	return true
}

func (ad *AttributeDecoder) Type() uint16 {
	return ad.a.Type
}

func (ad *AttributeDecoder) AttrData() []byte {
	return ad.a.Data
}
