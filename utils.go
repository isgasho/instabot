package instabot

import (
	"bufio"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/winterssy/sreq"
)

const (
	volatileSeed    = "12345"
	igSigKey        = "5f3e50f435583c9ae626302a71f7340044087a7e2c60adacfc254205a993e305"
	igSigKeyVersion = 4
)

func readUserInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return readUserInput(prompt)
	}
	return input
}

func timeOffset() int {
	_, offset := time.Now().Zone()
	return offset
}

func GenerateUUID() string {
	uuid := make([]byte, 16)
	_, _ = rand.Read(uuid)
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	uuid[8] = (uuid[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

func generateMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func generateDeviceId(seed string) string {
	hash := generateMD5Hash(seed + volatileSeed)
	return "android-" + hash[:16]
}

func generateHMAC(text string, key string) string {
	hasher := hmac.New(sha256.New, []byte(key))
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func GenerateSignedForm(form sreq.Form) sreq.Form {
	signedForm := make(sreq.Form, 2)
	signedForm.Set("ig_sig_key_version", igSigKeyVersion)
	body := form.Marshal()
	signedForm.Set("signed_body", fmt.Sprintf("%s.%s", generateHMAC(body, igSigKey), body))
	return signedForm
}
