package cipher

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"fmt"
	"strings"
)

type AesCbc struct{}

func (a *AesCbc) key1() []byte {
	return []byte{0x55, 0x48, 0x5B, 0x7A, 0x7C, 0x6D, 0x3E, 0x2A, 0x6C, 0x56, 0x4D, 0x2D, 0x22, 0x67, 0x56, 0x4D}
}

func (a *AesCbc) key2() []byte {
	return []byte{0x4E, 0x25, 0x53, 0x71, 0x5F, 0x7A, 0x5A, 0x5C, 0x60, 0x45, 0x63, 0x48, 0x66, 0x24, 0x65, 0x50}
}

func (a *AesCbc) iv() []byte {
	return []byte{0x54, 0x67, 0x70, 0x75, 0x60, 0x73, 0x5A, 0x5C, 0x69, 0x40, 0x42, 0x66, 0x73, 0x5A, 0x7D, 0x5E}
}

func (a *AesCbc) Encrypt(data []byte) ([]byte, error) {
	padded1 := a.padZero(data)
	cipher1, err := a.encrypt(padded1, a.key1(), a.iv())
	if err != nil {
		return nil, err
	}
	r1 := append(a.iv(), cipher1...)

	cipher2, err := a.encrypt(r1, a.key2(), a.iv())
	if err != nil {
		return nil, err
	}
	final := append(a.iv(), cipher2...)

	return []byte(strings.ToUpper(hex.EncodeToString(final))), nil
}

func (a *AesCbc) Decrypt(data []byte) ([]byte, error) {
	data, err := hex.DecodeString(strings.ToLower(string(data)))
	if err != nil {
		return nil, err
	}
	if len(data) < 16 {
		return nil, fmt.Errorf("invalid encrypted data length")
	}

	decrypted1, err := a.decrypt(data[16:], a.key2(), a.iv())
	if err != nil {
		return nil, err
	}
	if len(decrypted1) < 16 {
		return nil, fmt.Errorf("invalid intermediate data length")
	}

	decrypted2, err := a.decrypt(decrypted1[16:], a.key1(), a.iv())
	if err != nil {
		return nil, err
	}

	return a.removePadding(decrypted2), nil
}

func (a *AesCbc) padZero(data []byte) []byte {
	blockSize := 16
	if len(data)%blockSize == 0 {
		return data
	}
	padded := make([]byte, (len(data)/blockSize+1)*blockSize)
	copy(padded, data)
	return padded
}

func (a *AesCbc) removePadding(data []byte) []byte {
	end := len(data)
	for end > 0 && data[end-1] == 0 {
		end--
	}
	return data[:end]
}

func (a *AesCbc) encrypt(plaintext, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(plaintext)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("plaintext length must be multiple of block size")
	}

	ciphertext := make([]byte, len(plaintext))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plaintext)
	return ciphertext, nil
}

func (a *AesCbc) decrypt(ciphertext, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("ciphertext length must be multiple of block size")
	}

	plaintext := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, ciphertext)
	return plaintext, nil
}
