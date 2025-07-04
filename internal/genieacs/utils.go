package genieacs

import "encoding/json"

// Converts any struct/map to JSON string for use in query params
func toJSONString(v interface{}) string {
	bytes, _ := json.Marshal(v)
	return string(bytes)
}
