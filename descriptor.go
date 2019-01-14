package bufr

import (
	"fmt"
	"io"

	"github.com/irwinnoteam2009/bitstream"
)

// Descriptor represents BUFR descriptor (FXY). Length = 2 octets (16 bits)
type Descriptor struct {
	// type of descriptor.
	// F=0 - element descriptor(Table B)
	// F=1 - replication opeator
	// F=2 - operator descriptor (Table C)
	// F=3 - sequence descriptor (Table D)
	F byte // 2 bits
	X byte // 6 bits
	Y byte // 8 bits
}

// NewDescriptor ...
func NewDescriptor(r io.Reader) (*Descriptor, error) {
	d := new(Descriptor)
	if err := d.Decode(r); err != nil {
		return nil, err
	}
	return d, nil
}

// Decode ...
func (d *Descriptor) Decode(r io.Reader) error {
	reader := bitstream.NewReader(r)
	// 1-2
	b, err := reader.ReadBits(2)
	if err != nil {
		return err
	}
	d.F = byte(b)
	// 3-8
	b, err = reader.ReadBits(6)
	if err != nil {
		return err
	}
	d.X = byte(b)
	// 9-16
	d.Y, err = reader.ReadByte()
	return err
}

// String ...
func (d *Descriptor) String() string {
	return fmt.Sprintf("%d%.2d%.3d", d.F, d.X, d.Y)
}
