package crypto

import (
	"testing"
)

func TestEncryptWithAESCBC(t *testing.T) {
	// Test case 1: Normal case with valid key and IV
	asciiIV1 := []byte("1234567890123456")
	asciiKey1 := []byte("12345678901234567890123456789012")
	msg1 := "Hello, World!"
	encrypted1, err1 := EncryptWithAESCBC(asciiIV1, asciiKey1, msg1)
	if err1 != nil {
		t.Errorf("Test case 1 failed: %s", err1)
	}
	t.Logf("Test case 1: Encrypted message: %s", encrypted1)

	// Test case 2: Invalid key size
	asciiIV2 := []byte("1234567890123456")
	asciiKey2 := []byte("12345678901678901")
	msg2 := "Hello, World!"
	_, err2 := EncryptWithAESCBC(asciiIV2, asciiKey2, msg2)
	if err2 == nil {
		t.Errorf("Test case 2 failed: Expected error but got nil")
	}
	t.Logf("Test case 2: Error message: %s", err2)

	// Test case 3: Empty message
	asciiIV3 := []byte("1234567890123456")
	asciiKey3 := []byte("12345678901234567890123456789012")
	msg3 := ""
	encrypted3, err3 := EncryptWithAESCBC(asciiIV3, asciiKey3, msg3)
	if err3 != nil {
		t.Errorf("Test case 3 failed: %s", err3)
	}
	t.Logf("Test case 3: Encrypted message: %s", encrypted3)

	// Test case 4: Message with special characters
	asciiIV4 := []byte("1234567890123456")
	asciiKey4 := []byte("12345678901234567890123456789012")
	msg4 := "Hello, World!@#$%^&*()"
	encrypted4, err4 := EncryptWithAESCBC(asciiIV4, asciiKey4, msg4)
	if err4 != nil {
		t.Errorf("Test case 4 failed: %s", err4)
	}
	t.Logf("Test case 4: Encrypted message: %s", encrypted4)

	// Test case 5: Message with non-ASCII characters
	asciiIV5 := []byte("1234567890123456")
	asciiKey5 := []byte("12345678901234567890123456789012")
	msg5 := "你好，世界！"
	encrypted5, err5 := EncryptWithAESCBC(asciiIV5, asciiKey5, msg5)
	if err5 != nil {
		t.Errorf("Test case 5 failed: %s", err5)
	}
	t.Logf("Test case 5: Encrypted message: %s", encrypted5)
}
