package cipher

import (
	"crypto/cipher"
	"crypto/des"
	"encoding/hex"
	"errors"
	"strings"
)

type DesEdeCbc struct{}

func (d *DesEdeCbc) key1() []byte {
	return []byte{0x5E, 0x67, 0x72, 0x79, 0x28, 0x50, 0x47, 0x75, 0x6D, 0x48, 0x63, 0x74, 0x5D, 0x29, 0x21, 0x3C, 0x7E, 0x6B, 0x56, 0x29, 0x4F, 0x21, 0x52, 0x40}
}

func (d *DesEdeCbc) key2() []byte {
	return []byte{0x63, 0x73, 0x63, 0x26, 0x72, 0x5C, 0x5E, 0x73, 0x6B, 0x60, 0x74, 0x51, 0x7B, 0x74, 0x76, 0x7D, 0x3F, 0x59, 0x2E, 0x6D, 0x6F, 0x64, 0x3E, 0x69}
}

func (d *DesEdeCbc) iv() []byte {
	return []byte{0x77, 0x2D, 0x56, 0x51, 0x28, 0x49, 0x7E, 0x57}
}

func (d *DesEdeCbc) Encrypt(data []byte) ([]byte, error) {
	padded1 := d.padZero(data)
	encrypted1, err := d.encrypt(padded1, d.key1(), d.iv())
	if err != nil {
		return nil, err
	}

	encrypted2, err := d.encrypt(encrypted1, d.key2(), d.iv())
	if err != nil {
		return nil, err
	}

	return []byte(strings.ToUpper(hex.EncodeToString(encrypted2))), nil
}

func (d *DesEdeCbc) Decrypt(data []byte) ([]byte, error) {
	data, err := hex.DecodeString(strings.ToLower(string(data)))
	if err != nil {
		return nil, err
	}
	decrypted1, err := d.decrypt(data, d.key2(), d.iv())
	if err != nil {
		return nil, err
	}

	decrypted2, err := d.decrypt(decrypted1, d.key1(), d.iv())
	if err != nil {
		return nil, err
	}

	return d.removePadding(decrypted2), nil
}

func (d *DesEdeCbc) padZero(data []byte) []byte {
	blockSize := 16
	if len(data)%blockSize == 0 {
		return data
	}
	padded := make([]byte, (len(data)/blockSize+1)*blockSize)
	copy(padded, data)
	return padded
}

func (d *DesEdeCbc) removePadding(data []byte) []byte {
	end := len(data)
	for end > 0 && data[end-1] == 0 {
		end--
	}
	return data[:end]
}

func (d *DesEdeCbc) encrypt(plaintext, key, iv []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}

	if len(plaintext)%block.BlockSize() != 0 {
		return nil, errors.New("plaintext length must be multiple of block size")
	}

	ciphertext := make([]byte, len(plaintext))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plaintext)
	return ciphertext, nil
}

func (d *DesEdeCbc) decrypt(ciphertext, key, iv []byte) ([]byte, error) {
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext)%block.BlockSize() != 0 {
		return nil, errors.New("ciphertext length must be multiple of block size")
	}

	plaintext := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, ciphertext)
	return plaintext, nil
}
