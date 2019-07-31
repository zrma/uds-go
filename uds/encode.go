package uds

import (
	"encoding/base64"
)

func encode(chunk []byte) string {
	return base64.StdEncoding.EncodeToString(chunk)
}

func decode(chunk string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(chunk)
}
