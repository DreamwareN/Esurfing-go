package cipher

import (
	"bytes"
	"encoding/hex"
	"strings"

	"github.com/emmansun/gmsm/zuc"
)

type Zuc struct{}

func (z *Zuc) key() []byte {
	return []byte{0x4f, 0x3f, 0x25, 0x70, 0x53, 0x2b, 0x4b, 0x59, 0x3b, 0x5d, 0x5b, 0x21, 0x3a, 0x41, 0x7a, 0x48}
}

func (z *Zuc) iv() []byte {
	return []byte{0x41, 0x3c, 0x7a, 0x55, 0x4a, 0x21, 0x48, 0x3d, 0x5d, 0x2d, 0x24, 0x45, 0x45, 0x3c, 0x57, 0x79}
}

func (z *Zuc) padToMultipleOf4(data []byte) []byte {
	remainder := len(data) % 4
	if remainder == 0 {
		return data
	}
	padding := 4 - remainder
	padded := make([]byte, len(data)+padding)
	copy(padded, data)
	return padded
}

func (z *Zuc) removePadding(data []byte) []byte {
	return bytes.TrimRight(data, "\x00")
}

func (z *Zuc) Encrypt(data []byte) ([]byte, error) {
	padded := z.padToMultipleOf4(data)
	cipher, err := zuc.NewCipher(z.key(), z.iv())
	if err != nil {
		return nil, err
	}
	cipher.XORKeyStream(padded, padded)
	return []byte(strings.ToUpper(hex.EncodeToString(padded))), nil
}

func (z *Zuc) Decrypt(data []byte) ([]byte, error) {
	data, err := hex.DecodeString(strings.ToLower(string(data)))
	if err != nil {
		return nil, err
	}
	cipher, err := zuc.NewCipher(z.key(), z.iv())
	if err != nil {
		return nil, err
	}
	decrypted := make([]byte, len(data))
	cipher.XORKeyStream(decrypted, data)
	trimmed := z.removePadding(decrypted)
	return trimmed, nil
}
