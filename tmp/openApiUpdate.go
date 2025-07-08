package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"os"
	"errors"
        "gopkg.in/yaml.v3"
	
)

/* the 'data' var MUST be an interface in order to handle changing OpenAPI file structure. New structure would require struct regeneration, which an interface avoids.
 */

var outFile = "output.OpenApi.yml"               
var data map[string]interface{} 

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !errors.Is(err, os.ErrNotExist)
}

func getApiFile() ([]byte, error) {
	//not ideal, I'll need to update this to be smarter
	url := "https://cdn.cribl.io/dl/4.12.2/cribl-apidocs-4.12.2-4b17c8d4.yml" 

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("\nError making GET request: %v\n", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("\nError: Received non-OK status code: %d\n", resp.StatusCode)
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("\nError reading response body: %v\n", err)
		return nil, err
	}

	return body, nil
}

func main() {
	fmt.Printf("Downloading current upstream API spec...")

	yamlData, err := getApiFile()
	if err != nil {
		return
	}

	fmt.Printf(" Done!\n")

	err = yaml.Unmarshal([]byte(yamlData), &data)
	if err != nil {
		fmt.Printf("\nError while parsing local API yml: %v", err)
		return
	}

	info := data["info"].(map[string]interface{})
	fmt.Printf("Upstream API File Version: %s\n", info["version"])

	fmt.Printf("Updating API Yaml...")
	data["info"].(map[string]interface{})["version"] = fmt.Sprintf("%s-TfProviderUpdated", info["version"])
	data["servers"] = serverData
        data["components"].(map[string]interface{})["schemas"].(map[string]interface{})["Pack"] = schemaPackData

	for path, value := range(pathSpeakeasyOperation) {
		pathParts := strings.Split(path, ".")
		switch len(pathParts) {
		case 1:
			data[pathParts[0]] = value 
		case 2:
			data[pathParts[0]].(map[string]interface{})[pathParts[1]] = value 
		case 3:
			data[pathParts[0]].(map[string]interface{})[pathParts[1]].(map[string]interface{})[pathParts[2]] = value 
		case 4:
			data[pathParts[0]].(map[string]interface{})[pathParts[1]].(map[string]interface{})[pathParts[2]].(map[string]interface{})[pathParts[3]] = value 
		}
	}

	updatedYaml, err := yaml.Marshal(&data)
	if err != nil {
		fmt.Printf("\nError while marshalling new API yml: %v", err)
		return
	}

	fmt.Printf(" Done!\n")

	if fileExists(outFile) {
		fmt.Printf("Detected existing outfile, removing... ")
		err := os.Remove(outFile)
		if err != nil {
			fmt.Printf("\nError removing file '%s': %v\n", outFile, err)
		} else {
			fmt.Printf(" Done!\n")
		}
        }

	fmt.Printf("Writing new API Yaml file...")

	f, err := os.OpenFile(outFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	//write our updated yaml to final file
	if _, err := f.Write([]byte(updatedYaml)); err != nil {
		fmt.Println(err)
	}

	fmt.Printf(" Done!\n")

}

