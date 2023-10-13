package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

var (
	jsonDt, newJsonData []map[string]interface{}
)

func main() {
	// Open the JSON file
	file, err := os.Open("temp.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Read the JSON data from the file
	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	file.Close()

	var jsonData []map[string]interface{}

	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		fmt.Println(err)
		return
	}

	jsonDt = make([]map[string]interface{}, len(jsonData))
	for k, v := range jsonData {
		jsonDt[k] = v
	}

	for _, v := range jsonData {
		jsonObj := make(map[string]interface{})
		jsonObj["id"] = v["id"]
		jsonObj["uuid"] = v["uuid"]
		jsonObj["type"] = v["type"]
		jsonObj["post_code"] = v["post_code"]
		jsonObj["lng"] = v["lng"]
		jsonObj["lat"] = v["lat"]
		jsonObj["name"] = handleName(v["name"].(map[string]interface{}))

		if v["parent_id"] != nil {
			parent, ok := v["parent_id"].(float64)
			if !ok {
				continue // skip this object if parent is not an int
			}
			hromada, region := handleParent(parent)
			jsonObj["hromada"], jsonObj["region"] = hromada, region
		}

		newJsonData = append(newJsonData, jsonObj)
	}

	// Marshal the JSON data back into a string
	outputData, err := json.Marshal(newJsonData)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Write the JSON data to a file
	err = os.WriteFile("output.json", outputData, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func handleName(nameObj map[string]interface{}) any {
	enName := nameObj["en"].(string)

	if strings.Contains(enName, "’ї") {
		nameObj["en"] = strings.ReplaceAll(enName, "’ї", "yi")
	}
	if strings.Contains(enName, "'ї") {
		nameObj["en"] = strings.ReplaceAll(enName, "'ї", "yi")
	}
	if strings.Contains(enName, "’") {
		nameObj["en"] = strings.ReplaceAll(enName, "’", "'")
	}
	if strings.Contains(enName, "ї") {
		nameObj["en"] = strings.ReplaceAll(enName, "ї", "yi")
	}
	if strings.Contains(enName, "Ї") {
		nameObj["en"] = strings.ReplaceAll(enName, "Ї", "Yi")
	}
	if strings.Contains(enName, "І") {
		nameObj["en"] = strings.ReplaceAll(enName, "І", "I")
	}
	if strings.Contains(enName, "і") {
		nameObj["en"] = strings.ReplaceAll(enName, "і", "i")
	}
	if strings.Contains(enName, "є") {
		nameObj["en"] = strings.ReplaceAll(enName, "є", "ye")
	}
	if strings.Contains(enName, "Є") {
		nameObj["en"] = strings.ReplaceAll(enName, "Є", "Ye")
	}

	return nameObj
}

func handleParent(parent float64) (string, string) {

	for _, v := range jsonDt {
		if fmt.Sprintf("%.0f", v["id"]) != fmt.Sprintf("%.0f", parent) {
			continue
		}

		if v["type"] == "COMMUNITY" {
			parentID, ok := v["parent_id"].(float64)
			if !ok {
				continue // skip this object if parent is not a float64
			}
			_, region := handleParent(parentID)
			return v["uuid"].(string), region
		}

		if v["type"] == "DISTRICT" {
			return "", v["uuid"].(string)
		}

	}

	return "", ""
}
