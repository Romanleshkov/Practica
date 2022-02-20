package AnsibleVault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"golang.org/x/crypto/pbkdf2"
	"log"
	"strings"
)

func Encrypt(body, password string) (result string, err error) {
	salt, err := GenerateRandomBytes(32)
	check(err)
	// salt_64 := "2262970e2309d5da757af6c473b0ed3034209cc0d48a3cc3d648c0b174c22fde"
	// salt,_ = hex.DecodeString(salt_64)
	key1, key2, iv := genKeyInitctr(password, salt)
	ciphertext := createCipherText(body, key1, iv)
	combined := combineParts(ciphertext, key2, salt)
	vaultText := hex.EncodeToString([]byte(combined))
	result = formatOutput(vaultText)
	return
}

func formatOutput(vaultText string) string {
	heading := "$ANSIBLE_VAULT"
	version := "1.1"
	cipherName := "AES256"

	headerElements := make([]string, 3)
	headerElements[0] = heading
	headerElements[1] = version
	headerElements[2] = cipherName
	header := strings.Join(headerElements, ";")

	elements := make([]string, 1)
	elements[0] = header
	for i := 0; i < len(vaultText); i += 80 {
		end := i + 80
		if end > len(vaultText) {
			end = len(vaultText)
		}
		elements = append(elements, vaultText[i:end])
	}
	elements = append(elements, "")

	whole := strings.Join(elements, "\n")
	return whole
}

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}
	return b, nil
}

func check(e error) {
	if e != nil {
		log.Fatalf("error: %v", e)
	}
}

func genKeyInitctr(password string, salt []byte) (key1, key2, iv []byte) {
	keylength := 32
	ivlength := 16
	key := pbkdf2.Key([]byte(password), salt, 10000, 2*keylength+ivlength, sha256.New)
	key1 = key[:keylength]
	key2 = key[keylength:(keylength * 2)]
	iv = key[(keylength * 2) : (keylength*2)+ivlength]
	return
}

func createCipherText(body string, key1, iv []byte) []byte {
	bs := aes.BlockSize
	padding := (bs - len(body)%bs)
	if padding == 0 {
		padding = bs
	}
	padChar := rune(padding)
	padArray := make([]byte, padding)
	for i := range padArray {
		padArray[i] = byte(padChar)
	}
	plaintext := []byte(body)
	plaintext = append(plaintext, padArray...)

	aesCipher, err := aes.NewCipher(key1)
	check(err)
	ciphertext := make([]byte, len(plaintext))

	aesBlock := cipher.NewCTR(aesCipher, iv)
	aesBlock.XORKeyStream(ciphertext, plaintext)
	return ciphertext
}

func combineParts(ciphertext, key2, salt []byte) string {
	hmacEncrypt := hmac.New(sha256.New, key2)
	_, err := hmacEncrypt.Write(ciphertext)
	check(err)
	hexSalt := hex.EncodeToString(salt)
	hexHmac := hmacEncrypt.Sum(nil)
	hexCipher := hex.EncodeToString(ciphertext)
	// nolint:unconvert
	combined := string(hexSalt) + "\n" + hex.EncodeToString([]byte(hexHmac)) + "\n" + string(hexCipher)
	return combined
}
