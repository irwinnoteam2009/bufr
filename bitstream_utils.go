package bufr

import (
	"github.com/irwinnoteam2009/bitstream"
)

func readByte3U64(bs *bitstream.Reader) (uint64, error) { return bs.ReadBits(24) }
func readByte3(bs *bitstream.Reader) (int, error) {
	b, err := readByte3U64(bs)
	if err != nil {
		return 0, err
	}
	return int(b), nil
}

func readByte2U64(bs *bitstream.Reader) (uint64, error) { return bs.ReadBits(16) }
func readByte2(bs *bitstream.Reader) (int, error) {
	b, err := readByte2U64(bs)
	if err != nil {
		return 0, err
	}
	return int(b), nil
}

func readByte(bs *bitstream.Reader) (int, error) {
	b, err := bs.ReadByte()
	if err != nil {
		return 0, err
	}
	return int(b), nil
}

func readBitBool(bs *bitstream.Reader) (bool, error) {
	b, err := bs.ReadBit()
	if err != nil {
		return false, err
	}
	return b != 0, err
}
