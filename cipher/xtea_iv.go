package cipher

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"strings"
)

type XTeaIv struct{}

func (x *XTeaIv) key1() []uint32 {
	return []uint32{0x796d7855, 0x297b2355, 0x587d726e, 0x4d3d4423}
}

func (x *XTeaIv) key2() []uint32 {
	return []uint32{0x7c70525d, 0x5a585d3d, 0x413e4029, 0x28755d6a}
}

func (x *XTeaIv) key3() []uint32 {
	return []uint32{0x425e5f6e, 0x46754e24, 0x507b233d, 0x2d644641}
}

func (x *XTeaIv) iv() []uint32 {
	return []uint32{0x544c2f3f, 0x6f485121}
}

func (x *XTeaIv) Encrypt(data []byte) ([]byte, error) {
	padded := x.padToMultipleOf8(data)
	blocks := make([]byte, len(padded))
	copy(blocks, padded)

	previous := make([]uint32, 2)
	copy(previous, x.iv())

	for i := 0; i < len(blocks); i += 8 {
		v0 := binary.BigEndian.Uint32(blocks[i:])
		v1 := binary.BigEndian.Uint32(blocks[i+4:])

		XORedV0 := v0 ^ previous[0]
		XORedV1 := v1 ^ previous[1]

		r1v0, r1v1 := x.encryptBlock(XORedV0, XORedV1, x.key3())
		r2v0, r2v1 := x.encryptBlock(r1v0, r1v1, x.key2())
		r3v0, r3v1 := x.encryptBlock(r2v0, r2v1, x.key1())

		binary.BigEndian.PutUint32(blocks[i:], r3v0)
		binary.BigEndian.PutUint32(blocks[i+4:], r3v1)

		previous[0], previous[1] = r3v0, r3v1
	}

	return []byte(strings.ToUpper(hex.EncodeToString(blocks))), nil
}

func (x *XTeaIv) Decrypt(data []byte) ([]byte, error) {
	data, err := hex.DecodeString(strings.ToLower(string(data)))
	if err != nil {
		return nil, err
	}
	if len(data)%8 != 0 {
		return nil, errors.New("data length is not a multiple of block size")
	}

	blocks := make([]byte, len(data))
	copy(blocks, data)

	previous := make([]uint32, 2)
	copy(previous, x.iv())

	for i := 0; i < len(blocks); i += 8 {
		v0 := binary.BigEndian.Uint32(blocks[i:])
		v1 := binary.BigEndian.Uint32(blocks[i+4:])

		r1v0, r1v1 := x.decryptBlock(v0, v1, x.key1())
		r2v0, r2v1 := x.decryptBlock(r1v0, r1v1, x.key2())
		r3v0, r3v1 := x.decryptBlock(r2v0, r2v1, x.key3())

		XORedV0 := r3v0 ^ previous[0]
		XORedV1 := r3v1 ^ previous[1]

		binary.BigEndian.PutUint32(blocks[i:], XORedV0)
		binary.BigEndian.PutUint32(blocks[i+4:], XORedV1)

		previous[0], previous[1] = v0, v1
	}

	return x.unPad(blocks), nil
}

func (x *XTeaIv) padToMultipleOf8(data []byte) []byte {
	padding := (8 - len(data)%8) % 8
	if padding == 0 {
		return data
	}
	padded := make([]byte, len(data)+padding)
	copy(padded, data)
	return padded
}

func (x *XTeaIv) unPad(data []byte) []byte {
	i := len(data) - 1
	for ; i >= 0 && data[i] == 0; i-- {
	}
	if i < 0 {
		return []byte{}
	}
	return data[:i+1]
}

func (x *XTeaIv) encryptBlock(v0, v1 uint32, key []uint32) (uint32, uint32) {
	var sum int32 = 0
	delta := uint32(0x9E3779B9)
	numRounds := 32

	for i := 0; i < numRounds; i++ {
		sumU := uint32(sum)
		term1 := v1 ^ sumU
		term2 := key[sumU&3]
		term3 := (v1 << 4) ^ (v1 >> 5)
		v0 += term1 + term2 + term3

		sum += int32(delta)

		sumU = uint32(sum)
		term1 = v0 ^ sumU
		term2 = key[(sumU>>11)&3]
		term3 = (v0 << 4) ^ (v0 >> 5)
		v1 += term1 + term2 + term3
	}
	return v0, v1
}

func (x *XTeaIv) decryptBlock(v0, v1 uint32, key []uint32) (uint32, uint32) {
	delta := uint32(0x9E3779B9)
	numRounds := 32
	sum := int32(delta) * int32(numRounds)

	for i := 0; i < numRounds; i++ {
		sumU := uint32(sum)
		term1 := key[(sumU>>11)&3]
		term2 := v0 ^ sumU
		term3 := (v0 << 4) ^ (v0 >> 5)
		v1 -= term1 + term2 + term3

		sum -= int32(delta)

		sumU = uint32(sum)
		term1 = v1 ^ sumU
		term2 = key[sumU&3]
		term3 = (v1 << 4) ^ (v1 >> 5)
		v0 -= term1 + term2 + term3
	}
	return v0, v1
}
