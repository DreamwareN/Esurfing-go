package cipher

import (
	"bytes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/emmansun/gmsm/sm4"
)

type Sm4Cbc struct{}

func (s *Sm4Cbc) key() []byte {
	return []byte{0x28, 0x2f, 0x29, 0x25, 0x6f, 0x3c, 0x75, 0x48, 0x6d, 0x4c, 0x2e, 0x51, 0x55, 0x27, 0x22, 0x2d}
}

func (s *Sm4Cbc) iv() []byte {
	return []byte{0x68, 0x3c, 0x42, 0x51, 0x5a, 0x46, 0x3a, 0x52, 0x67, 0x77, 0x7e, 0x6e, 0x69, 0x70, 0x48, 0x5e}
}

func (s *Sm4Cbc) Encrypt(data []byte) ([]byte, error) {
	block, err := sm4.NewCipher(s.key())
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()

	paddedData := s.padPKCS5(data, blockSize)

	if len(s.iv()) != blockSize {
		return nil, errors.New("IV length must equal block size")
	}

	mode := cipher.NewCBCEncrypter(block, s.iv())
	ciphertext := make([]byte, len(paddedData))
	mode.CryptBlocks(ciphertext, paddedData)
	return []byte(strings.ToUpper(hex.EncodeToString(ciphertext))), nil
}

func (s *Sm4Cbc) Decrypt(data []byte) ([]byte, error) {
	data, err := hex.DecodeString(strings.ToLower(string(data)))
	if err != nil {
		return nil, err
	}
	block, err := sm4.NewCipher(s.key())
	if err != nil {
		return nil, err
	}

	if len(s.iv()) != block.BlockSize() {
		return nil, errors.New("IV length must equal block size")
	}

	if len(data)%block.BlockSize() != 0 {
		return nil, errors.New("ciphertext length is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, s.iv())
	plaintext := make([]byte, len(data))
	mode.CryptBlocks(plaintext, data)
	return s.unPadPKCS5(plaintext, block.BlockSize())
}

func (s *Sm4Cbc) padPKCS5(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padText...)
}

func (s *Sm4Cbc) unPadPKCS5(src []byte, blockSize int) ([]byte, error) {
	length := len(src)
	if length == 0 {
		return nil, errors.New("invalid padding")
	}
	padding := int(src[length-1])
	if padding < 1 || padding > blockSize {
		return nil, errors.New("invalid padding")
	}
	if length < padding {
		return nil, errors.New("invalid padding")
	}
	return src[:length-padding], nil
}
