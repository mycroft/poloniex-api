package poloniexapi

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"

	"net/url"
)

func getSha256(input []byte) []byte {
	sha := sha256.New()
	sha.Write(input)

	return sha.Sum(nil)
}

func getHMacSha512(message, secret []byte) []byte {
	mac := hmac.New(sha512.New, secret)
	mac.Write(message)

	return mac.Sum(nil)
}

func createPoloniexSignature(values url.Values, secret string) string {
	macsum := getHMacSha512([]byte(values.Encode()), []byte(secret))

	return hex.EncodeToString(macsum)
}
