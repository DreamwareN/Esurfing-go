package cipher

import (
	"bytes"
	"crypto/cipher"
	"crypto/des"
	"encoding/hex"
	"errors"
	"strings"
)

type DesEdeEcb struct{}

func (d *DesEdeEcb) key1() []byte {
	return []byte{0x25, 0x6A, 0x63, 0x5A, 0x46, 0x3F, 0x26, 0x64, 0x53, 0x7A, 0x2E, 0x5B, 0x24, 0x4C, 0x62, 0x67, 0x2B, 0x2D, 0x67, 0x68, 0x43, 0x74, 0x69, 0x51}
}

func (d *DesEdeEcb) key2() []byte {
	return []byte{0x59, 0x28, 0x5B, 0x7E, 0x7D, 0x26, 0x74, 0x49, 0x48, 0x76, 0x59, 0x58, 0x62, 0x75, 0x51, 0x55, 0x26, 0x73, 0x55, 0x5C, 0x67, 0x52, 0x2E, 0x6C}
}

func (d *DesEdeEcb) Encrypt(data []byte) ([]byte, error) {
	padded := d.padZero(data)

	block1, err := des.NewTripleDESCipher(d.key1())
	if err != nil {
		return nil, err
	}
	encrypted1 := d.encryptECB(block1, padded)

	block2, err := des.NewTripleDESCipher(d.key2())
	if err != nil {
		return nil, err
	}
	encrypted2 := d.encryptECB(block2, encrypted1)

	return []byte(strings.ToUpper(hex.EncodeToString(encrypted2))), nil
}

func (d *DesEdeEcb) Decrypt(data []byte) ([]byte, error) {
	data, err := hex.DecodeString(strings.ToLower(string(data)))
	if err != nil {
		return nil, err
	}
	if len(data)%des.BlockSize != 0 {
		return nil, errors.New("ciphertext is not a multiple of the block size")
	}

	block2, err := des.NewTripleDESCipher(d.key2())
	if err != nil {
		return nil, err
	}
	decrypted2 := d.decryptECB(block2, data)

	block1, err := des.NewTripleDESCipher(d.key1())
	if err != nil {
		return nil, err
	}
	decrypted1 := d.decryptECB(block1, decrypted2)

	return bytes.TrimRight(decrypted1, "\x00"), nil
}

func (d *DesEdeEcb) padZero(data []byte) []byte {
	blockSize := 16
	if len(data)%blockSize == 0 {
		return data
	}
	padding := blockSize - len(data)%blockSize
	padded := make([]byte, len(data)+padding)
	copy(padded, data)
	return padded
}

func (d *DesEdeEcb) encryptECB(block cipher.Block, src []byte) []byte {
	dst := make([]byte, len(src))
	bs := block.BlockSize()
	for i := 0; i < len(src); i += bs {
		block.Encrypt(dst[i:i+bs], src[i:i+bs])
	}
	return dst
}

func (d *DesEdeEcb) decryptECB(block cipher.Block, src []byte) []byte {
	dst := make([]byte, len(src))
	bs := block.BlockSize()
	for i := 0; i < len(src); i += bs {
		block.Decrypt(dst[i:i+bs], src[i:i+bs])
	}
	return dst
}
