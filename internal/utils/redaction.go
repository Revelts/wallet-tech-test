package utils

import (
	"encoding/json"
	"strings"
)

var sensitiveFields = []string{
	"password",
	"pin",
	"password_hash",
	"pin_hash",
	"authorization",
	"access_token",
	"refresh_token",
	"jwt",
	"token",
}

const redactedValue = "***REDACTED***"

func RedactSensitiveData(data []byte) (redacted []byte) {
	if len(data) == 0 {
		return data
	}

	var jsonData map[string]interface{}
	err := json.Unmarshal(data, &jsonData)
	if err != nil {
		return data
	}

	redactMap(jsonData)

	redacted, err = json.Marshal(jsonData)
	if err != nil {
		return data
	}

	return
}

func redactMap(data map[string]interface{}) {
	for key, value := range data {
		lowerKey := strings.ToLower(key)

		if isSensitiveField(lowerKey) {
			data[key] = redactedValue
			continue
		}

		switch v := value.(type) {
		case map[string]interface{}:
			redactMap(v)
		case []interface{}:
			redactSlice(v)
		}
	}
}

func redactSlice(data []interface{}) {
	for i, item := range data {
		switch v := item.(type) {
		case map[string]interface{}:
			redactMap(v)
		case []interface{}:
			redactSlice(v)
		default:
			data[i] = v
		}
	}
}

func isSensitiveField(fieldName string) (result bool) {
	for _, sensitive := range sensitiveFields {
		if strings.Contains(fieldName, sensitive) {
			result = true
			return
		}
	}
	result = false
	return
}
