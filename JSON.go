package gocord

import (
	"encoding/json"
	"fmt"
)

// DumpMap - Dumps the parsed JSON into a map.
func DumpMap(space string, m map[string]interface{}) {
	for k, v := range m {
		if mv, ok := v.(map[string]interface{}); ok {
			fmt.Printf("{ \"%v\": \n", k)
			DumpMap(space+"\t", mv)
			fmt.Printf("}\n")
		} else {
			fmt.Printf("%v %v : %v\n", space, k, v)
		}
	}
}

// ParseJSONObject - Parses a JSON object into a map;
func ParseJSONObject(rawJSON string) map[string]interface{} {
	jsonMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(rawJSON), &jsonMap)
	if err != nil {
		panic(err)
	}
	DumpMap("", jsonMap)
	return jsonMap
}
