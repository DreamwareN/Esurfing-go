package cipher

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"strings"
)

type XTea struct{}

func (x *XTea) key1() []uint32 {
	return []uint32{0x7a7a676a, 0x277e4a73, 0x3e43296c, 0x577d7d7a}
}

func (x *XTea) key2() []uint32 {
	return []uint32{0x3d3c695f, 0x71797a74, 0x445f5763, 0x6f692765}
}

func (x *XTea) key3() []uint32 {
	return []uint32{0x5b5a683d, 0x2e572a77, 0x4a474465, 0x663d7e5c}
}

func (x *XTea) Encrypt(data []byte) ([]byte, error) {
	padded := x.padToMultipleOf8(data)
	encrypted := make([]byte, len(padded))
	copy(encrypted, padded)

	for i := 0; i < len(encrypted); i += 8 {
		v0 := binary.BigEndian.Uint32(encrypted[i : i+4])
		v1 := binary.BigEndian.Uint32(encrypted[i+4 : i+8])

		r0, r1 := x.encryptBlock(v0, v1, x.key1())
		r0, r1 = x.encryptBlock(r0, r1, x.key2())
		r0, r1 = x.encryptBlock(r0, r1, x.key3())

		binary.BigEndian.PutUint32(encrypted[i:], r0)
		binary.BigEndian.PutUint32(encrypted[i+4:], r1)
	}

	return []byte(strings.ToUpper(hex.EncodeToString(encrypted))), nil
}

func (x *XTea) Decrypt(data []byte) ([]byte, error) {
	data, err := hex.DecodeString(strings.ToLower(string(data)))
	if err != nil {
		return nil, err
	}
	if len(data)%8 != 0 {
		return nil, errors.New("data length is not a multiple of 8")
	}

	decrypted := make([]byte, len(data))
	copy(decrypted, data)

	for i := 0; i < len(decrypted); i += 8 {
		v0 := binary.BigEndian.Uint32(decrypted[i : i+4])
		v1 := binary.BigEndian.Uint32(decrypted[i+4 : i+8])

		r0, r1 := x.decryptBlock(v0, v1, x.key3())
		r0, r1 = x.decryptBlock(r0, r1, x.key2())
		r0, r1 = x.decryptBlock(r0, r1, x.key1())

		binary.BigEndian.PutUint32(decrypted[i:], r0)
		binary.BigEndian.PutUint32(decrypted[i+4:], r1)
	}

	decrypted = bytes.TrimRight(decrypted, "\x00")
	return decrypted, nil
}

func (x *XTea) padToMultipleOf8(data []byte) []byte {
	padding := (8 - len(data)%8) % 8
	if padding == 0 {
		return data
	}
	padded := make([]byte, len(data)+padding)
	copy(padded, data)
	return padded
}

func (x *XTea) encryptBlock(v0In, v1In uint32, key []uint32) (uint32, uint32) {
	v0, v1 := v0In, v1In
	var sum uint32 = 0
	delta := uint32(0x9E3779B9)

	for i := 0; i < 32; i++ {
		v0 += (v1 ^ sum) + key[sum&3] + ((v1 << 4) ^ (v1 >> 5))
		sum += delta
		v1 += key[(sum>>11)&3] + (v0 ^ sum) + ((v0 << 4) ^ (v0 >> 5))
	}

	return v0, v1
}

func (x *XTea) decryptBlock(v0In, v1In uint32, key []uint32) (uint32, uint32) {
	v0, v1 := v0In, v1In
	delta := uint32(0x9E3779B9)
	sum := delta * 32

	for i := 0; i < 32; i++ {
		v1 -= key[(sum>>11)&3] + (v0 ^ sum) + ((v0 << 4) ^ (v0 >> 5))
		sum -= delta
		v0 -= (v1 ^ sum) + key[sum&3] + ((v1 << 4) ^ (v1 >> 5))
	}

	return v0, v1
}
