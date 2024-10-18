// Package pkcs7 implements PKCS#7 padding
//
// This is a standard way of encoding variable length buffers into
// buffers which are a multiple of an underlying crypto block size.
package pkcs7

import (
	"bytes"
	"errors"
	"fmt"
)

// Errors Unpad can return
var (
	ErrorPaddingNotFound      = errors.New("bad PKCS#7 padding - not padded")
	ErrorPaddingNotAMultiple  = errors.New("bad PKCS#7 padding - not a multiple of blocksize")
	ErrorPaddingTooLong       = errors.New("bad PKCS#7 padding - too long")
	ErrorPaddingTooShort      = errors.New("bad PKCS#7 padding - too short")
	ErrorPaddingNotAllTheSame = errors.New("bad PKCS#7 padding - not all the same")
)

// Pad buf using PKCS#7 to a multiple of n.
//
// Appends the padding to buf - make a copy of it first if you don't
// want it modified.

// Modified by IceCodeNew
func Pad(blockLen int, buf []byte) []byte {
	if blockLen <= 1 || blockLen >= 256 {
		err := fmt.Errorf("invalid block size: %d", blockLen)
		panic(err)
	}
	padding := blockLen - (len(buf) % blockLen)
	buf = append(buf,
		bytes.Repeat([]byte{byte(padding)}, padding)...)
	if (len(buf) % blockLen) != 0 {
		panic("padding failed")
	}
	return buf
}

// Unpad buf using PKCS#7 from a multiple of n returning a slice of
// buf or an error if malformed.
func Unpad(blockLen int, buf []byte) ([]byte, error) {
	if blockLen <= 1 || blockLen >= 256 {
		panic("bad multiple")
	}
	length := len(buf)
	if length == 0 {
		return nil, ErrorPaddingNotFound
	}
	if (length % blockLen) != 0 {
		return nil, ErrorPaddingNotAMultiple
	}
	padding := int(buf[length-1])
	if padding > blockLen {
		return nil, ErrorPaddingTooLong
	}
	if padding == 0 {
		return nil, ErrorPaddingTooShort
	}
	for i := 0; i < padding; i++ {
		if buf[length-1-i] != byte(padding) {
			return nil, ErrorPaddingNotAllTheSame
		}
	}
	return buf[:length-padding], nil
}
