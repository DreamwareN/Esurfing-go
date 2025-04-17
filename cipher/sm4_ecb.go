package cipher

import (
	"bytes"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/emmansun/gmsm/sm4"
)

type Sm4Ecb struct{}

func (s *Sm4Ecb) key() []byte {
	return []byte{0x53, 0x2f, 0x79, 0x4a, 0x4e, 0x79, 0x74, 0x4d, 0x67, 0x66, 0x57, 0x5a, 0x2d, 0x44, 0x5c, 0x57}
}

func (s *Sm4Ecb) Encrypt(data []byte) ([]byte, error) {
	block, err := sm4.NewCipher(s.key())
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	paddedData := s.padPKCS5(data, blockSize)
	encrypted := make([]byte, len(paddedData))
	for i := 0; i < len(paddedData); i += blockSize {
		block.Encrypt(encrypted[i:i+blockSize], paddedData[i:i+blockSize])
	}
	return []byte(strings.ToUpper(hex.EncodeToString(encrypted))), nil
}

func (s *Sm4Ecb) Decrypt(data []byte) ([]byte, error) {
	data, err := hex.DecodeString(strings.ToLower(string(data)))
	if err != nil {
		return nil, err
	}
	block, err := sm4.NewCipher(s.key())
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	if len(data)%blockSize != 0 {
		return nil, errors.New("ciphertext is not a multiple of the block size")
	}
	decrypted := make([]byte, len(data))
	for i := 0; i < len(data); i += blockSize {
		block.Decrypt(decrypted[i:i+blockSize], data[i:i+blockSize])
	}
	return s.unPadPKCS5(decrypted, blockSize)
}

func (s *Sm4Ecb) padPKCS5(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padText...)
}

func (s *Sm4Ecb) unPadPKCS5(src []byte, blockSize int) ([]byte, error) {
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
