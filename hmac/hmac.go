package hmac

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"strconv"
	"time"
)

const veracodeRequestVersionString = "vcode_request_version_1"

const dataFormat = "id=%s&host=%s&url=%s&method=%s"

const headerFormat = "%s id=%s,ts=%s,nonce=%X,sig=%X"

const veracodeHMACSHA256 = "VERACODE-HMAC-SHA-256"

func CalculateAuthorizationHeader(url *url.URL, httpMethod, apiKeyID, apiKeySecret string) (string, error) {
	nonce, err := createNonce(16)

	if err != nil {
		return "", err
	}

	secret, err := fromHexString(apiKeySecret)

	if err != nil {
		return "", err
	}

	timestampMilliseconds := strconv.FormatInt(time.Now().UnixNano()/int64(1000000), 10)
	data := fmt.Sprintf(dataFormat, apiKeyID, url.Hostname(), url.RequestURI(), httpMethod)
	dataSignature := calculateSignature(secret, nonce, []byte(timestampMilliseconds), []byte(data))
	return fmt.Sprintf(headerFormat, veracodeHMACSHA256, apiKeyID, timestampMilliseconds, nonce, dataSignature), nil
}

func createNonce(size int) ([]byte, error) {
	nonce := make([]byte, size)

	_, err := rand.Read(nonce)

	if err != nil {
		return nil, err
	}

	return nonce, nil
}

func fromHexString(input string) ([]byte, error) {
	decoded, err := hex.DecodeString(input)

	if err != nil {
		return nil, err
	}

	return decoded, nil
}

func calculateSignature(key, nonce, timestamp, data []byte) []byte {
	encryptedNonce := hmac256(nonce, key)
	encryptedTimestampMilliseconds := hmac256(timestamp, encryptedNonce)
	signingKey := hmac256([]byte(veracodeRequestVersionString), encryptedTimestampMilliseconds)
	return hmac256(data, signingKey)
}

func hmac256(message, key []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	return mac.Sum(nil)
}
