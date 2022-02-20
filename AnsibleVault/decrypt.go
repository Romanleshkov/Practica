package AnsibleVault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

func Decrypt(data, password string) (result string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("ERROR: %v", r)
		}
	}()
	data = replaceCarriageReturn(data)
	body := splitHeader([]byte(data))
	salt, cryptedHmac, ciphertext := decodeData(body)
	key1, key2, iv := genKeyInitctr(password, salt)
	checkDigest(key2, cryptedHmac, ciphertext)
	aesCipher, err := aes.NewCipher(key1)
	check(err)
	aesBlock := cipher.NewCTR(aesCipher, iv)
	plaintext := make([]byte, len(ciphertext))
	aesBlock.XORKeyStream(plaintext, ciphertext)
	padding := int(plaintext[len(plaintext)-1])
	result = string(plaintext[:len(plaintext)-padding])
	return
}

func replaceCarriageReturn(data string) string { // удаляет \r, которые ставятся в не-Unix системах
	return strings.ReplaceAll(data, "\r", "")
}

func splitHeader(data []byte) string { // проверка загаловка и возврат рабочей части
	contents := string(data)
	lines := strings.Split(contents, "\n")
	header := strings.Split(lines[0], ";")
	cipherName := strings.TrimSpace(header[2])
	if cipherName != "AES256" {
		panic("unsupported cipher: " + cipherName)
	}
	body := strings.Join(lines[1:], "")
	return body
}

func decodeData(body string) (salt, cryptedHmac, ciphertext []byte) {
	decoded, _ := hex.DecodeString(body)
	elements := strings.SplitN(string(decoded), "\n", 3)
	salt, err1 := hex.DecodeString(elements[0])
	if err1 != nil {
		panic(err1)
	}
	cryptedHmac, err2 := hex.DecodeString(elements[1])
	if err2 != nil {
		panic(err2)
	}
	ciphertext, err3 := hex.DecodeString(elements[2])
	if err3 != nil {
		panic(err3)
	}
	return
}

func checkDigest(key2, cryptedHmac, ciphertext []byte) {
	hmacDecrypt := hmac.New(sha256.New, key2)
	_, err := hmacDecrypt.Write(ciphertext)
	check(err)
	expectedMAC := hmacDecrypt.Sum(nil)
	if !hmac.Equal(cryptedHmac, expectedMAC) {
		panic("digests do not match - exiting")
	}
}