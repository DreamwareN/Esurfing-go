package cipher

import (
	"crypto/aes"
	"encoding/hex"
	"errors"
	"strings"
)

type AesEcb struct{}

func (a *AesEcb) key1() []byte {
	return []byte{0x3A, 0x71, 0x7C, 0x4C, 0x51, 0x4F, 0x3C, 0x6A, 0x2E, 0x43, 0x7A, 0x43, 0x3B, 0x56, 0x57, 0x59}
}

func (a *AesEcb) key2() []byte {
	return []byte{0x72, 0x6E, 0x25, 0x41, 0x45, 0x2F, 0x41, 0x54, 0x27, 0x4B, 0x3B, 0x3B, 0x59, 0x25, 0x52, 0x24}
}

func (a *AesEcb) Encrypt(data []byte) ([]byte, error) {
	padded := a.padZero(data)
	encrypted1, err := a.encrypt(padded, a.key1())
	if err != nil {
		return nil, err
	}
	encrypted2, err := a.encrypt(encrypted1, a.key2())
	if err != nil {
		return nil, err
	}
	return []byte(strings.ToUpper(hex.EncodeToString(encrypted2))), nil
}

func (a *AesEcb) Decrypt(data []byte) ([]byte, error) {
	data, err := hex.DecodeString(strings.ToLower(string(data)))
	if err != nil {
		return nil, err
	}
	decrypted1, err := a.decrypt(data, a.key2())
	if err != nil {
		return nil, err
	}
	decrypted2, err := a.decrypt(decrypted1, a.key1())
	if err != nil {
		return nil, err
	}
	return a.removePadding(decrypted2), nil
}

func (a *AesEcb) padZero(data []byte) []byte {
	blockSize := 16
	if len(data)%blockSize == 0 {
		return data
	}
	padded := make([]byte, (len(data)/blockSize+1)*blockSize)
	copy(padded, data)
	return padded
}

func (a *AesEcb) removePadding(data []byte) []byte {
	end := len(data)
	for end > 0 && data[end-1] == 0 {
		end--
	}
	return data[:end]
}

func (a *AesEcb) encrypt(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(plaintext)%aes.BlockSize != 0 {
		return nil, errors.New("plaintext length must be multiple of block size")
	}

	ciphertext := make([]byte, len(plaintext))
	for i := 0; i < len(plaintext); i += aes.BlockSize {
		block.Encrypt(ciphertext[i:i+aes.BlockSize], plaintext[i:i+aes.BlockSize])
	}
	return ciphertext, nil
}

func (a *AesEcb) decrypt(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, errors.New("ciphertext length must be multiple of block size")
	}

	plaintext := make([]byte, len(ciphertext))
	for i := 0; i < len(ciphertext); i += aes.BlockSize {
		block.Decrypt(plaintext[i:i+aes.BlockSize], ciphertext[i:i+aes.BlockSize])
	}
	return plaintext, nil
}
