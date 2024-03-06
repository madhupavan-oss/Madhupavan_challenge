package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	schemaFlag := flag.String("config", "schema.json", "Used to read the json file")
	mapOut, err := parseSchema(*schemaFlag)
	if err != nil {
		fmt.Println("error :", err)
		return
	}
	output, err := JsonTransform(mapOut)
	if err != nil {
		fmt.Println("error :", err)
		return
	}
	out, err := json.Marshal(output)
	if err != nil {
		fmt.Println("error :", err)
		return
	}
	fmt.Println(string(out))
}

func JsonTransform(inputMap map[string]interface{}) (map[string]interface{}, error) {
	output := make(map[string]interface{}, 0)

	for key, value := range inputMap {
		// Sanitize key
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}

		// Initialize output map
		outMap := make(map[string]interface{})

		switch val := value.(type) {
		case map[string]interface{}:
			// Process map
			for k, v := range val {
				k = strings.TrimSpace(k)
				switch k {
				case "S":
					outMap[key] = formatString(v)
				case "N":
					// Numeric type
					if vNum, err := formatNum(v); err == nil {
						outMap[key] = vNum
					}

				case "BOOL":
					// Boolean type
					outMap[key] = formatBool(v)

				case "NULL":
					// Null type
					nullStr := strings.ToLower(strings.TrimSpace(v.(string)))
					if nullStr == "1" || nullStr == "t" || nullStr == "true" {
						outMap[key] = nil
					}
				case "M":
					// Map type
					submap := v.(map[string]interface{})
					if len(submap) > 0 {
						subOutput, err := JsonTransform(submap)
						if err != nil {
							return nil, err
						}
						outMap[key] = subOutput
					}
				case "L":
					// List type
					listValue, ok := v.([]interface{})
					if !ok {
						// If the value is not a slice, continue to the next key
						continue
					}
					if len(listValue) == 0 {
						// If the list is empty, omit this field
						continue
					}
					outList := make([]interface{}, 0)
					for _, listItem := range listValue {
						switch listItem := listItem.(type) {
						case map[string]interface{}:
							for itemKey, item := range listItem {
								itemKey = strings.TrimSpace(itemKey)
								switch itemKey {
								case "S":
									itemStr := formatString(item)
									if itemStr != "" {
										outList = append(outList, itemStr)
									}
								case "N":
									// Numeric type
									if num, err := formatNum(item); err == nil {
										outList = append(outList, num)
									}
								case "BOOL":
									outList = append(outList, formatBool(item))
								}
							}
						}
					}
					if len(outList) > 0 {
						outMap[key] = outList
					}

				}
			}
		}

		if len(outMap) > 0 {
			for k, v := range outMap {
				output[k] = v
			}
		}
	}

	return output, nil
}

func formatString(v interface{}) interface{} {
	strVal := strings.TrimSpace(v.(string))
	if strVal != "" {
		// Transform RFC3339 formatted strings to Unix Epoch
		if t, err := time.Parse(time.RFC3339, strVal); err == nil {
			return t.Unix()
		}
		return strVal
	}
	return ""
}

func formatNum(v interface{}) (float64, error) {
	numStr := strings.TrimSpace(v.(string))
	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0, errors.New("Invalid number format")
	}
	return num, nil
}

func formatBool(v interface{}) bool {
	boolStr := strings.ToLower(strings.TrimSpace(v.(string)))
	switch boolStr {
	case "1", "t", "true":
		return true
	case "0", "f", "false":
		return false
	}
	return false
}

func parseSchema(fileName string) (map[string]interface{}, error) {
	output := make(map[string]interface{}, 0)
	if strings.Contains(fileName, ".json") {
		fileBytes, err := os.ReadFile(fileName)
		if err != nil {
			return nil, err
		}
		if err = json.Unmarshal(fileBytes, &output); err != nil {
			return nil, err
		}
		return output, nil
	}
	return nil, errors.New("config file is not json file")
}
