package utils

import "encoding/json"

// ParseJSON parse JSON bytes to interface
func ParseJSON(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
