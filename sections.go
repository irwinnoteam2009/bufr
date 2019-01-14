package bufr

import (
	"fmt"
	"io"

	"github.com/irwinnoteam2009/bitstream"
)

type decoder interface {
	Decode(io.Reader) error
}

// Section0 is an indicator section
type Section0 struct {
	MagicString [4]byte // octets 1-4. Value "BUFR"
	Len         int     // octets 5-7
	Version     byte    // octet 8
}

// Section1 is an identification section
type Section1 struct {
	Len                   int  // octets 1-3
	MasterTable           byte // octet 4
	Center                int  // octet 5-6. Code table 0 01 033
	UpdateSequenceNumber  byte // octet 7
	OptionalSectionExists bool // octet 8, bit 1. Bits 2->8 = 0
	DataCategory          byte // octet 9
	DataSubCategory       byte // octet 10
	MasterTableVersion    byte // octet 11
	LocalTableVersion     byte // octet 12
	Year                  int  // octet 13
	Month                 byte // octet 14
	Day                   byte // octet 15
	Hours                 byte // octet 16
	Minutes               byte // octet 17
}

// Section2 is an optional section
type Section2 struct {
	Len      int    // octets 1-3
	Reserved byte   // octet 4
	ADP      []byte // octet 5+
}

// Section3 is a data description section
type Section3 struct {
	Len            int    // octets 1-3
	Reserved       byte   // octet 4
	SubsetCount    int    // octets 5-6
	ObservedData   bool   // octet 7, bit 1
	CompressedData bool   // octet 7, bit 2. Bits 3-8 = 0
	Collection     []byte // octet 8+
}

// Section4 is a data section
type Section4 struct {
	Len      int    // octet 1-3
	Reserved byte   // octet 4
	Data     []byte //
}

// Section5 is an end section
type Section5 struct {
	END [4]byte // octets 1-4. Value 7777
}

// NewSection0 ...
func NewSection0(r io.Reader) (*Section0, error) {
	sect := new(Section0)
	if err := sect.Decode(r); err != nil {
		return nil, err
	}
	return sect, nil
}

// NewSection1 ...
func NewSection1(r io.Reader) (*Section1, error) {
	sect := new(Section1)
	if err := sect.Decode(r); err != nil {
		return nil, err
	}
	return sect, nil
}

// NewSection2 ...
func NewSection2(r io.Reader) (*Section2, error) {
	sect := new(Section2)
	if err := sect.Decode(r); err != nil {
		return nil, err
	}
	return sect, nil
}

// NewSection3 ...
func NewSection3(r io.Reader) (*Section3, error) {
	sect := new(Section3)
	if err := sect.Decode(r); err != nil {
		return nil, err
	}
	return sect, nil
}

// Decode is decoder implementation
func (s *Section0) Decode(r io.Reader) error {
	reader := bitstream.NewReader(r)
	var err error
	for i := 0; i < 4; i++ {
		if s.MagicString[i], err = reader.ReadByte(); err != nil {
			return err
		}
	}
	if s.Len, err = readByte3(reader); err != nil {
		return err
	}
	if s.Version, err = reader.ReadByte(); err != nil {
		return err
	}
	return nil
}

// Decode ...
func (s *Section1) Decode(r io.Reader) error {
	reader := bitstream.NewReader(r)
	var err error
	if s.Len, err = readByte3(reader); err != nil { // 1-3
		return err
	}
	if s.MasterTable, err = reader.ReadByte(); err != nil { // 4
		return err
	}
	if s.Center, err = readByte2(reader); err != nil { // 5-6
		return err
	}
	if s.UpdateSequenceNumber, err = reader.ReadByte(); err != nil { // 7
		return err
	}
	// 8
	b, err := reader.ReadByte()
	if err != nil {
		return err
	}
	s.OptionalSectionExists = (b & 0x80) != 0
	if s.DataCategory, err = reader.ReadByte(); err != nil { // 9
		return err
	}
	if s.DataSubCategory, err = reader.ReadByte(); err != nil { // 10
		return err
	}
	if s.MasterTableVersion, err = reader.ReadByte(); err != nil { // 11
		return err
	}
	if s.LocalTableVersion, err = reader.ReadByte(); err != nil { // 12
		return err
	}
	if s.Year, err = readByte(reader); err != nil { // 13
		return err
	}
	s.Year += 2000
	if s.Month, err = reader.ReadByte(); err != nil { // 14
		return err
	}
	if s.Day, err = reader.ReadByte(); err != nil { // 15
		return err
	}
	if s.Hours, err = reader.ReadByte(); err != nil { // 16
		return err
	}
	if s.Minutes, err = reader.ReadByte(); err != nil { // 17
		return err
	}

	return nil
}

// Decode ...
func (s *Section2) Decode(r io.Reader) error {
	return fmt.Errorf("not implemented yet")
}

// Decode ...
func (s *Section3) Decode(r io.Reader) error {

}
