package t1

import (
	"encoding/json"
)

// JSONString returns a JSON encoded string of v.
func JSONString(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}
