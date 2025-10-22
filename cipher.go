package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"encoding/binary"
	"encoding/hex"
	"errors"

	"github.com/emmansun/gmsm/sm4"
	"github.com/emmansun/gmsm/zuc"
)

type Cipher interface {
	Encrypt(data []byte) ([]byte, error)
	Decrypt(data []byte) ([]byte, error)
}

const (
	AlgoAesCbc    = "CAFBCBAD-B6E7-4CAB-8A67-14D39F00CE1E"
	AlgoAesEcb    = "A474B1C2-3DE0-4EA2-8C5F-7093409CE6C4"
	AlgoDesEdeCbc = "5BFBA864-BBA9-42DB-8EAD-49B5F412BD81"
	AlgoDesEdeEcb = "6E0B65FF-0B5B-459C-8FCE-EC7F2BEA9FF5"
	AlgoZUC       = "B809531F-0007-4B5B-923B-4BD560398113"
	AlgoSm4Cbc    = "F3974434-C0DD-4C20-9E87-DDB6814A1C48"
	AlgoSm4Ecb    = "ED382482-F72C-4C41-A76D-28EEA0F1F2AF"
	AlgoXTea      = "B3047D4E-67DF-4864-A6A5-DF9B9E525C79"
	AlgoXTeaIv    = "C32C68F9-CA81-4260-A329-BBAFD1A9CCD1"
)

var (
	ErrInvalidCiphertextLength = errors.New("ciphertext length must be a multiple of the block size")
	ErrInvalidPadding          = errors.New("invalid padding")
	ErrInvalidDataLength       = errors.New("invalid data length")
)

var cipherRegistry = map[string]func() Cipher{
	AlgoAesCbc:    func() Cipher { return new(AesCbc) },
	AlgoAesEcb:    func() Cipher { return new(AesEcb) },
	AlgoDesEdeCbc: func() Cipher { return new(DesEdeCbc) },
	AlgoDesEdeEcb: func() Cipher { return new(DesEdeEcb) },
	AlgoZUC:       func() Cipher { return new(Zuc) },
	AlgoSm4Cbc:    func() Cipher { return new(Sm4Cbc) },
	AlgoSm4Ecb:    func() Cipher { return new(Sm4Ecb) },
	AlgoXTea:      func() Cipher { return new(XTea) },
	AlgoXTeaIv:    func() Cipher { return new(XTeaIv) },
}

func NewCipher(algoID string) Cipher {
	if factory, ok := cipherRegistry[algoID]; ok {
		return factory()
	}
	return nil
}

func encodeHexUpper(data []byte) []byte {
	dst := make([]byte, hex.EncodedLen(len(data)))
	hex.Encode(dst, data)
	return bytes.ToUpper(dst)
}

func decodeHex(data []byte) ([]byte, error) {
	dst := make([]byte, hex.DecodedLen(len(data)))
	n, err := hex.Decode(dst, data)
	return dst[:n], err
}

func zeroPadding(data []byte, blockSize int) []byte {
	padding := (blockSize - len(data)%blockSize) % blockSize
	if padding == 0 {
		return data
	}
	return append(data, bytes.Repeat([]byte{0}, padding)...)
}

func zeroUnpadding(data []byte) []byte {
	return bytes.TrimRight(data, "\x00")
}

func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

func pkcs7Unpadding(data []byte, blockSize int) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, ErrInvalidPadding
	}
	padding := int(data[length-1])
	if padding < 1 || padding > blockSize || length < padding {
		return nil, ErrInvalidPadding
	}
	return data[:length-padding], nil
}

type ecb struct {
	b         cipher.Block
	blockSize int
}

func newECB(b cipher.Block) *ecb {
	return &ecb{b: b, blockSize: b.BlockSize()}
}

type ecbEncrypter ecb
type ecbDecrypter ecb

func NewECBEncrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbEncrypter)(newECB(b))
}

func NewECBDecrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbDecrypter)(newECB(b))
}

func (x *ecbEncrypter) BlockSize() int { return x.blockSize }
func (x *ecbEncrypter) CryptBlocks(dst, src []byte) {
	for i := 0; i < len(src); i += x.blockSize {
		x.b.Encrypt(dst[i:i+x.blockSize], src[i:i+x.blockSize])
	}
}

func (x *ecbDecrypter) BlockSize() int { return x.blockSize }
func (x *ecbDecrypter) CryptBlocks(dst, src []byte) {
	for i := 0; i < len(src); i += x.blockSize {
		x.b.Decrypt(dst[i:i+x.blockSize], src[i:i+x.blockSize])
	}
}

var (
	aesCbcKey1 = []byte{0x55, 0x48, 0x5B, 0x7A, 0x7C, 0x6D, 0x3E, 0x2A, 0x6C, 0x56, 0x4D, 0x2D, 0x22, 0x67, 0x56, 0x4D}
	aesCbcKey2 = []byte{0x4E, 0x25, 0x53, 0x71, 0x5F, 0x7A, 0x5A, 0x5C, 0x60, 0x45, 0x63, 0x48, 0x66, 0x24, 0x65, 0x50}
	aesCbcIv   = []byte{0x54, 0x67, 0x70, 0x75, 0x60, 0x73, 0x5A, 0x5C, 0x69, 0x40, 0x42, 0x66, 0x73, 0x5A, 0x7D, 0x5E}
)

type AesCbc struct{}

func (a *AesCbc) Encrypt(data []byte) ([]byte, error) {
	padded := zeroPadding(data, aes.BlockSize)

	block1, err := aes.NewCipher(aesCbcKey1)
	if err != nil {
		return nil, err
	}
	cipher1 := make([]byte, len(padded))
	cipher.NewCBCEncrypter(block1, aesCbcIv).CryptBlocks(cipher1, padded)

	r1 := append(aesCbcIv, cipher1...)

	block2, err := aes.NewCipher(aesCbcKey2)
	if err != nil {
		return nil, err
	}
	cipher2 := make([]byte, len(r1))
	cipher.NewCBCEncrypter(block2, aesCbcIv).CryptBlocks(cipher2, r1)

	final := append(aesCbcIv, cipher2...)
	return encodeHexUpper(final), nil
}

func (a *AesCbc) Decrypt(data []byte) ([]byte, error) {
	ciphertext, err := decodeHex(data)
	if err != nil {
		return nil, err
	}
	if len(ciphertext) < aes.BlockSize {
		return nil, ErrInvalidDataLength
	}

	block2, err := aes.NewCipher(aesCbcKey2)
	if err != nil {
		return nil, err
	}
	if len(ciphertext[aes.BlockSize:])%aes.BlockSize != 0 {
		return nil, ErrInvalidCiphertextLength
	}
	decrypted1 := make([]byte, len(ciphertext)-aes.BlockSize)
	cipher.NewCBCDecrypter(block2, aesCbcIv).CryptBlocks(decrypted1, ciphertext[aes.BlockSize:])

	if len(decrypted1) < aes.BlockSize {
		return nil, ErrInvalidDataLength
	}

	block1, err := aes.NewCipher(aesCbcKey1)
	if err != nil {
		return nil, err
	}
	if len(decrypted1[aes.BlockSize:])%aes.BlockSize != 0 {
		return nil, ErrInvalidCiphertextLength
	}
	decrypted2 := make([]byte, len(decrypted1)-aes.BlockSize)
	cipher.NewCBCDecrypter(block1, aesCbcIv).CryptBlocks(decrypted2, decrypted1[aes.BlockSize:])

	return zeroUnpadding(decrypted2), nil
}

var (
	aesEcbKey1 = []byte{0x3A, 0x71, 0x7C, 0x4C, 0x51, 0x4F, 0x3C, 0x6A, 0x2E, 0x43, 0x7A, 0x43, 0x3B, 0x56, 0x57, 0x59}
	aesEcbKey2 = []byte{0x72, 0x6E, 0x25, 0x41, 0x45, 0x2F, 0x41, 0x54, 0x27, 0x4B, 0x3B, 0x3B, 0x59, 0x25, 0x52, 0x24}
)

type AesEcb struct{}

func (a *AesEcb) Encrypt(data []byte) ([]byte, error) {
	padded := zeroPadding(data, aes.BlockSize)

	block1, err := aes.NewCipher(aesEcbKey1)
	if err != nil {
		return nil, err
	}
	encrypted1 := make([]byte, len(padded))
	NewECBEncrypter(block1).CryptBlocks(encrypted1, padded)

	block2, err := aes.NewCipher(aesEcbKey2)
	if err != nil {
		return nil, err
	}
	encrypted2 := make([]byte, len(encrypted1))
	NewECBEncrypter(block2).CryptBlocks(encrypted2, encrypted1)

	return encodeHexUpper(encrypted2), nil
}

func (a *AesEcb) Decrypt(data []byte) ([]byte, error) {
	ciphertext, err := decodeHex(data)
	if err != nil {
		return nil, err
	}
	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, ErrInvalidCiphertextLength
	}

	block2, err := aes.NewCipher(aesEcbKey2)
	if err != nil {
		return nil, err
	}
	decrypted1 := make([]byte, len(ciphertext))
	NewECBDecrypter(block2).CryptBlocks(decrypted1, ciphertext)

	block1, err := aes.NewCipher(aesEcbKey1)
	if err != nil {
		return nil, err
	}
	decrypted2 := make([]byte, len(decrypted1))
	NewECBDecrypter(block1).CryptBlocks(decrypted2, decrypted1)

	return zeroUnpadding(decrypted2), nil
}

var (
	desEdeCbcKey1 = []byte{0x5E, 0x67, 0x72, 0x79, 0x28, 0x50, 0x47, 0x75, 0x6D, 0x48, 0x63, 0x74, 0x5D, 0x29, 0x21, 0x3C, 0x7E, 0x6B, 0x56, 0x29, 0x4F, 0x21, 0x52, 0x40}
	desEdeCbcKey2 = []byte{0x63, 0x73, 0x63, 0x26, 0x72, 0x5C, 0x5E, 0x73, 0x6B, 0x60, 0x74, 0x51, 0x7B, 0x74, 0x76, 0x7D, 0x3F, 0x59, 0x2E, 0x6D, 0x6F, 0x64, 0x3E, 0x69}
	desEdeCbcIv   = []byte{0x77, 0x2D, 0x56, 0x51, 0x28, 0x49, 0x7E, 0x57}
)

type DesEdeCbc struct{}

func (d *DesEdeCbc) Encrypt(data []byte) ([]byte, error) {
	padded := zeroPadding(data, des.BlockSize)

	block1, err := des.NewTripleDESCipher(desEdeCbcKey1)
	if err != nil {
		return nil, err
	}
	encrypted1 := make([]byte, len(padded))
	cipher.NewCBCEncrypter(block1, desEdeCbcIv).CryptBlocks(encrypted1, padded)

	block2, err := des.NewTripleDESCipher(desEdeCbcKey2)
	if err != nil {
		return nil, err
	}
	encrypted2 := make([]byte, len(encrypted1))
	cipher.NewCBCEncrypter(block2, desEdeCbcIv).CryptBlocks(encrypted2, encrypted1)

	return encodeHexUpper(encrypted2), nil
}

func (d *DesEdeCbc) Decrypt(data []byte) ([]byte, error) {
	ciphertext, err := decodeHex(data)
	if err != nil {
		return nil, err
	}
	if len(ciphertext)%des.BlockSize != 0 {
		return nil, ErrInvalidCiphertextLength
	}

	block2, err := des.NewTripleDESCipher(desEdeCbcKey2)
	if err != nil {
		return nil, err
	}
	decrypted1 := make([]byte, len(ciphertext))
	cipher.NewCBCDecrypter(block2, desEdeCbcIv).CryptBlocks(decrypted1, ciphertext)

	block1, err := des.NewTripleDESCipher(desEdeCbcKey1)
	if err != nil {
		return nil, err
	}
	decrypted2 := make([]byte, len(decrypted1))
	cipher.NewCBCDecrypter(block1, desEdeCbcIv).CryptBlocks(decrypted2, decrypted1)

	return zeroUnpadding(decrypted2), nil
}

var (
	desEdeEcbKey1 = []byte{0x25, 0x6A, 0x63, 0x5A, 0x46, 0x3F, 0x26, 0x64, 0x53, 0x7A, 0x2E, 0x5B, 0x24, 0x4C, 0x62, 0x67, 0x2B, 0x2D, 0x67, 0x68, 0x43, 0x74, 0x69, 0x51}
	desEdeEcbKey2 = []byte{0x59, 0x28, 0x5B, 0x7E, 0x7D, 0x26, 0x74, 0x49, 0x48, 0x76, 0x59, 0x58, 0x62, 0x75, 0x51, 0x55, 0x26, 0x73, 0x55, 0x5C, 0x67, 0x52, 0x2E, 0x6C}
)

type DesEdeEcb struct{}

func (d *DesEdeEcb) Encrypt(data []byte) ([]byte, error) {
	padded := zeroPadding(data, des.BlockSize)

	block1, err := des.NewTripleDESCipher(desEdeEcbKey1)
	if err != nil {
		return nil, err
	}
	encrypted1 := make([]byte, len(padded))
	NewECBEncrypter(block1).CryptBlocks(encrypted1, padded)

	block2, err := des.NewTripleDESCipher(desEdeEcbKey2)
	if err != nil {
		return nil, err
	}
	encrypted2 := make([]byte, len(encrypted1))
	NewECBEncrypter(block2).CryptBlocks(encrypted2, encrypted1)

	return encodeHexUpper(encrypted2), nil
}

func (d *DesEdeEcb) Decrypt(data []byte) ([]byte, error) {
	ciphertext, err := decodeHex(data)
	if err != nil {
		return nil, err
	}
	if len(ciphertext)%des.BlockSize != 0 {
		return nil, ErrInvalidCiphertextLength
	}

	block2, err := des.NewTripleDESCipher(desEdeEcbKey2)
	if err != nil {
		return nil, err
	}
	decrypted1 := make([]byte, len(ciphertext))
	NewECBDecrypter(block2).CryptBlocks(decrypted1, ciphertext)

	block1, err := des.NewTripleDESCipher(desEdeEcbKey1)
	if err != nil {
		return nil, err
	}
	decrypted2 := make([]byte, len(decrypted1))
	NewECBDecrypter(block1).CryptBlocks(decrypted2, decrypted1)

	return zeroUnpadding(decrypted2), nil
}

var (
	sm4CbcKey = []byte{0x28, 0x2f, 0x29, 0x25, 0x6f, 0x3c, 0x75, 0x48, 0x6d, 0x4c, 0x2e, 0x51, 0x55, 0x27, 0x22, 0x2d}
	sm4CbcIv  = []byte{0x68, 0x3c, 0x42, 0x51, 0x5a, 0x46, 0x3a, 0x52, 0x67, 0x77, 0x7e, 0x6e, 0x69, 0x70, 0x48, 0x5e}
)

type Sm4Cbc struct{}

func (s *Sm4Cbc) Encrypt(data []byte) ([]byte, error) {
	block, err := sm4.NewCipher(sm4CbcKey)
	if err != nil {
		return nil, err
	}
	padded := pkcs7Padding(data, block.BlockSize())
	encrypted := make([]byte, len(padded))
	cipher.NewCBCEncrypter(block, sm4CbcIv).CryptBlocks(encrypted, padded)
	return encodeHexUpper(encrypted), nil
}

func (s *Sm4Cbc) Decrypt(data []byte) ([]byte, error) {
	ciphertext, err := decodeHex(data)
	if err != nil {
		return nil, err
	}
	block, err := sm4.NewCipher(sm4CbcKey)
	if err != nil {
		return nil, err
	}
	if len(ciphertext)%block.BlockSize() != 0 {
		return nil, ErrInvalidCiphertextLength
	}
	decrypted := make([]byte, len(ciphertext))
	cipher.NewCBCDecrypter(block, sm4CbcIv).CryptBlocks(decrypted, ciphertext)
	return pkcs7Unpadding(decrypted, block.BlockSize())
}

var sm4EcbKey = []byte{0x53, 0x2f, 0x79, 0x4a, 0x4e, 0x79, 0x74, 0x4d, 0x67, 0x66, 0x57, 0x5a, 0x2d, 0x44, 0x5c, 0x57}

type Sm4Ecb struct{}

func (s *Sm4Ecb) Encrypt(data []byte) ([]byte, error) {
	block, err := sm4.NewCipher(sm4EcbKey)
	if err != nil {
		return nil, err
	}
	padded := pkcs7Padding(data, block.BlockSize())
	encrypted := make([]byte, len(padded))
	NewECBEncrypter(block).CryptBlocks(encrypted, padded)
	return encodeHexUpper(encrypted), nil
}

func (s *Sm4Ecb) Decrypt(data []byte) ([]byte, error) {
	ciphertext, err := decodeHex(data)
	if err != nil {
		return nil, err
	}
	block, err := sm4.NewCipher(sm4EcbKey)
	if err != nil {
		return nil, err
	}
	if len(ciphertext)%block.BlockSize() != 0 {
		return nil, ErrInvalidCiphertextLength
	}
	decrypted := make([]byte, len(ciphertext))
	NewECBDecrypter(block).CryptBlocks(decrypted, ciphertext)
	return pkcs7Unpadding(decrypted, block.BlockSize())
}

var (
	xteaKey1 = []uint32{0x7a7a676a, 0x277e4a73, 0x3e43296c, 0x577d7d7a}
	xteaKey2 = []uint32{0x3d3c695f, 0x71797a74, 0x445f5763, 0x6f692765}
	xteaKey3 = []uint32{0x5b5a683d, 0x2e572a77, 0x4a474465, 0x663d7e5c}
)

type XTea struct{}

func (x *XTea) Encrypt(data []byte) ([]byte, error) {
	padded := zeroPadding(data, 8)
	encrypted := make([]byte, len(padded))
	for i := 0; i < len(padded); i += 8 {
		v0 := binary.BigEndian.Uint32(padded[i : i+4])
		v1 := binary.BigEndian.Uint32(padded[i+4 : i+8])
		v0, v1 = encryptXTeaBlock(v0, v1, xteaKey1)
		v0, v1 = encryptXTeaBlock(v0, v1, xteaKey2)
		v0, v1 = encryptXTeaBlock(v0, v1, xteaKey3)
		binary.BigEndian.PutUint32(encrypted[i:], v0)
		binary.BigEndian.PutUint32(encrypted[i+4:], v1)
	}
	return encodeHexUpper(encrypted), nil
}

func (x *XTea) Decrypt(data []byte) ([]byte, error) {
	ciphertext, err := decodeHex(data)
	if err != nil {
		return nil, err
	}
	if len(ciphertext)%8 != 0 {
		return nil, ErrInvalidCiphertextLength
	}
	decrypted := make([]byte, len(ciphertext))
	for i := 0; i < len(ciphertext); i += 8 {
		v0 := binary.BigEndian.Uint32(ciphertext[i : i+4])
		v1 := binary.BigEndian.Uint32(ciphertext[i+4 : i+8])
		v0, v1 = decryptXTeaBlock(v0, v1, xteaKey3)
		v0, v1 = decryptXTeaBlock(v0, v1, xteaKey2)
		v0, v1 = decryptXTeaBlock(v0, v1, xteaKey1)
		binary.BigEndian.PutUint32(decrypted[i:], v0)
		binary.BigEndian.PutUint32(decrypted[i+4:], v1)
	}
	return zeroUnpadding(decrypted), nil
}

var (
	xteaIvKey1 = []uint32{0x796d7855, 0x297b2355, 0x587d726e, 0x4d3d4423}
	xteaIvKey2 = []uint32{0x7c70525d, 0x5a585d3d, 0x413e4029, 0x28755d6a}
	xteaIvKey3 = []uint32{0x425e5f6e, 0x46754e24, 0x507b233d, 0x2d644641}
	xteaIv     = []uint32{0x544c2f3f, 0x6f485121}
)

type XTeaIv struct{}

func (x *XTeaIv) Encrypt(data []byte) ([]byte, error) {
	padded := zeroPadding(data, 8)
	encrypted := make([]byte, len(padded))
	prevV0, prevV1 := xteaIv[0], xteaIv[1]
	for i := 0; i < len(padded); i += 8 {
		v0 := binary.BigEndian.Uint32(padded[i:]) ^ prevV0
		v1 := binary.BigEndian.Uint32(padded[i+4:]) ^ prevV1
		v0, v1 = encryptXTeaBlock(v0, v1, xteaIvKey3)
		v0, v1 = encryptXTeaBlock(v0, v1, xteaIvKey2)
		v0, v1 = encryptXTeaBlock(v0, v1, xteaIvKey1)
		binary.BigEndian.PutUint32(encrypted[i:], v0)
		binary.BigEndian.PutUint32(encrypted[i+4:], v1)
		prevV0, prevV1 = v0, v1
	}
	return encodeHexUpper(encrypted), nil
}

func (x *XTeaIv) Decrypt(data []byte) ([]byte, error) {
	ciphertext, err := decodeHex(data)
	if err != nil {
		return nil, err
	}
	if len(ciphertext)%8 != 0 {
		return nil, ErrInvalidCiphertextLength
	}
	decrypted := make([]byte, len(ciphertext))
	prevV0, prevV1 := xteaIv[0], xteaIv[1]
	for i := 0; i < len(ciphertext); i += 8 {
		v0 := binary.BigEndian.Uint32(ciphertext[i:])
		v1 := binary.BigEndian.Uint32(ciphertext[i+4:])
		r0, r1 := decryptXTeaBlock(v0, v1, xteaIvKey1)
		r0, r1 = decryptXTeaBlock(r0, r1, xteaIvKey2)
		r0, r1 = decryptXTeaBlock(r0, r1, xteaIvKey3)
		binary.BigEndian.PutUint32(decrypted[i:], r0^prevV0)
		binary.BigEndian.PutUint32(decrypted[i+4:], r1^prevV1)
		prevV0, prevV1 = v0, v1
	}
	return zeroUnpadding(decrypted), nil
}

const xteaDelta uint32 = 0x9E3779B9
const xteaRounds = 32

func encryptXTeaBlock(v0, v1 uint32, key []uint32) (uint32, uint32) {
	var sum uint32
	for i := 0; i < xteaRounds; i++ {
		v0 += (((v1 << 4) ^ (v1 >> 5)) + v1) ^ (sum + key[sum&3])
		sum += xteaDelta
		v1 += (((v0 << 4) ^ (v0 >> 5)) + v0) ^ (sum + key[(sum>>11)&3])
	}
	return v0, v1
}

func decryptXTeaBlock(v0, v1 uint32, key []uint32) (uint32, uint32) {
	roundsAsVar := uint32(xteaRounds)
	sum := xteaDelta * roundsAsVar
	for i := 0; i < xteaRounds; i++ {
		v1 -= (((v0 << 4) ^ (v0 >> 5)) + v0) ^ (sum + key[(sum>>11)&3])
		sum -= xteaDelta
		v0 -= (((v1 << 4) ^ (v1 >> 5)) + v1) ^ (sum + key[sum&3])
	}
	return v0, v1
}

var (
	zucKey = []byte{0x4f, 0x3f, 0x25, 0x70, 0x53, 0x2b, 0x4b, 0x59, 0x3b, 0x5d, 0x5b, 0x21, 0x3a, 0x41, 0x7a, 0x48}
	zucIv  = []byte{0x41, 0x3c, 0x7a, 0x55, 0x4a, 0x21, 0x48, 0x3d, 0x5d, 0x2d, 0x24, 0x45, 0x45, 0x3c, 0x57, 0x79}
)

type Zuc struct{}

func (z *Zuc) Encrypt(data []byte) ([]byte, error) {
	padded := zeroPadding(data, 4) // ZUC works on 32-bit words
	c, err := zuc.NewCipher(zucKey, zucIv)
	if err != nil {
		return nil, err
	}
	c.XORKeyStream(padded, padded)
	return encodeHexUpper(padded), nil
}

func (z *Zuc) Decrypt(data []byte) ([]byte, error) {
	ciphertext, err := decodeHex(data)
	if err != nil {
		return nil, err
	}
	c, err := zuc.NewCipher(zucKey, zucIv)
	if err != nil {
		return nil, err
	}
	decrypted := make([]byte, len(ciphertext))
	c.XORKeyStream(decrypted, ciphertext)
	return zeroUnpadding(decrypted), nil
}
