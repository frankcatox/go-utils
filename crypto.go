package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"encoding/base64"
	"encoding/hex"
)

func PKCSPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCSUnPadding(srcBytes []byte) []byte {
	length := len(srcBytes)
	unpadding := int(srcBytes[length-1])
	return srcBytes[:(length - unpadding)]
}

func XorEnc(src any, seed int) []byte {
	srcBytes := ToBytes(src)
	encrypted := make([]byte, len(srcBytes))
	for i, v := range srcBytes {
		encrypted[i] = v ^ (uint8(seed))
	}
	return encrypted
}

func XorEncStr(src string, seed int) string {
	return base64.StdEncoding.EncodeToString(XorEnc(src, seed))
}

func XorEncHex(src string, seed int) string {
	return hex.EncodeToString(XorEnc(src, seed))
}

func XorDecStr(src string, seed int) string {
	enc, _ := base64.StdEncoding.DecodeString(src)
	return ToString(XorEnc(enc, seed))
}

func XorDecHex(src string, seed int) string {
	enc, _ := hex.DecodeString(src)
	return ToString(XorEnc(enc, seed))
}

func AesCbcEnc(src, key, iv any) ([]byte, error) {
	block, err := aes.NewCipher(ToBytes(key))
	if err != nil {
		return nil, err
	}
	srcBytes := ToBytes(src)
	srcBytes = PKCSPadding(srcBytes, block.BlockSize())
	encrypter := cipher.NewCBCEncrypter(block, ToBytes(iv))
	encrypted := make([]byte, len(srcBytes))
	encrypter.CryptBlocks(encrypted, srcBytes)
	return encrypted, nil
}

func AesCbcDec(src, key, iv any) ([]byte, error) {
	block, err := aes.NewCipher(ToBytes(key))
	if err != nil {
		return nil, err
	}
	srcBytes := ToBytes(src)
	decrypter := cipher.NewCBCDecrypter(block, ToBytes(iv))
	decrypted := make([]byte, len(srcBytes))
	decrypter.CryptBlocks(decrypted, srcBytes)
	decrypted = PKCSUnPadding(decrypted)
	return decrypted, nil
}

func AesCfbEnc(src, key, iv any) ([]byte, error) {
	block, err := aes.NewCipher(ToBytes(key))
	if err != nil {
		return nil, err
	}
	srcBytes := ToBytes(src)
	encrypter := cipher.NewCFBEncrypter(block, ToBytes(iv))
	encrypted := make([]byte, len(srcBytes))
	encrypter.XORKeyStream(encrypted, srcBytes)
	return encrypted, nil
}

func AesCfbDec(src, key, iv any) ([]byte, error) {
	block, err := aes.NewCipher(ToBytes(key))
	if err != nil {
		return nil, err
	}
	srcBytes := ToBytes(src)
	decrypter := cipher.NewCFBDecrypter(block, ToBytes(iv))
	decrypted := make([]byte, len(srcBytes))
	decrypter.XORKeyStream(decrypted, srcBytes)
	return decrypted, nil
}

func DesEnc(src, key, iv any) ([]byte, error) {
	block, err := des.NewCipher(ToBytes(key))
	if err != nil {
		return nil, err
	}
	srcBytes := ToBytes(src)
	srcBytes = PKCSPadding(srcBytes, block.BlockSize())
	encrypter := cipher.NewCBCEncrypter(block, ToBytes(iv))
	encrypted := make([]byte, len(srcBytes))
	encrypter.CryptBlocks(encrypted, srcBytes)
	return encrypted, nil
}

func DesDec(src, key, iv any) ([]byte, error) {
	block, err := des.NewCipher(ToBytes(key))
	if err != nil {
		return nil, err
	}
	srcBytes := ToBytes(src)
	decrypter := cipher.NewCBCDecrypter(block, ToBytes(iv))
	decrypted := make([]byte, len(srcBytes))
	decrypter.CryptBlocks(decrypted, srcBytes)
	decrypted = PKCSUnPadding(decrypted)
	return decrypted, nil
}

func AesEncStr(src, key string) string {
	enc, _ := AesCbcEnc([]byte(src), key, key)
	return base64.StdEncoding.EncodeToString(enc)
}

func AesEncHex(src, key string) string {
	enc, _ := AesCbcEnc([]byte(src), key, key)
	return hex.EncodeToString(enc)
}

func AesDecStr(src, key string) string {
	enc, _ := base64.StdEncoding.DecodeString(src)
	dec, _ := AesCbcDec(enc, key, key)
	return string(dec)
}

func AesDesHex(src, key string) string {
	enc, _ := hex.DecodeString(src)
	dec, _ := AesCbcDec(enc, key, key)
	return string(dec)
}

func DesEncStr(src, key string) string {
	enc, _ := DesEnc([]byte(src), key, key)
	return base64.StdEncoding.EncodeToString(enc)
}

func DesEncHex(src, key string) string {
	enc, _ := DesEnc([]byte(src), key, key)
	return hex.EncodeToString(enc)
}

func DesDecStr(src, key string) string {
	enc, _ := base64.StdEncoding.DecodeString(src)
	dec, _ := DesDec(enc, key, key)
	return string(dec)
}

func DesDecHex(src, key string) string {
	enc, _ := hex.DecodeString(src)
	dec, _ := DesDec(enc, key, key)
	return string(dec)
}
