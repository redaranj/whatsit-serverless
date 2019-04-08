package common

import (
	"testing"
)

func TestHash(t *testing.T) {
	val := "testing123"
	hashVal := Hash(val)
	expectedVal := "b822f1cd2dcfc685b47e83e3980289fd5d8e3ff3a82def24d7d1d68bb272eb32"
	if hashVal != expectedVal {
		t.Error("hash invalid")
	}
}

func TestParseJSONBody(t *testing.T) {
	t.Run("Valid POST", func(t *testing.T) {
		event := map[string]interface{}{
			"body": `{"number": "15555555555", "message": "This is a message"}`,
		}
		_, err := ParseJSONBody(event)
		if err != nil {
			t.Error(err.Error())
		}
	})

	t.Run("Invalid POST", func(t *testing.T) {
		event := map[string]interface{}{
			"body": `{"message": A message"}`,
		}
		_, err := ParseJSONBody(event)
		if err.Error() != "invalid character 'A' looking for beginning of value" {
			t.Error("expected invalid JSON string")
		}
	})
}
