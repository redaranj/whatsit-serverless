package common

import (
	"errors"

	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

func Hash(number string) string {
	sha := sha256.New()
	sha.Write([]byte(number))
	shaString := hex.EncodeToString(sha.Sum(nil))

	return shaString
}

func ParseJSONBody(event map[string]interface{}) (map[string]interface{}, error) {
	var postParams map[string]interface{}

	body, bodyOk := event["body"].(string)
	if !bodyOk {
		return postParams, errors.New("invalid POST body")
	}

	if err := json.Unmarshal([]byte(body), &postParams); err != nil {
		return postParams, err
	}

	return postParams, nil
}
