package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strings"

	"github.com/yalp/jsonpath"
)

func removeFields(jsonData map[string]interface{}, paths []string) error {
	for _, path := range paths {
		// Use jsonpath to locate the parent object and key
		parentPath, key := splitJSONPath(path)
		parent, err := jsonpath.Read(jsonData, parentPath)
		if err != nil {
			// Skip if the parent path does not exist
			continue
		}

		// Ensure parent is a map and remove the key
		if parentMap, ok := parent.(map[string]interface{}); ok {
			delete(parentMap, key)
		} else {
			return fmt.Errorf("parent path %s is not a map", parentPath)
		}
	}
	return nil
}

func splitJSONPath(path string) (string, string) {
	// Split the JSONPath into parent path and key
	lastDot := strings.LastIndex(path, ".")
	if lastDot == -1 {
		return "", path
	}
	return path[:lastDot], path[lastDot+1:]
}

func compareJSON(file1, file2 string, ignorePaths []string) error {
	data1, err := ioutil.ReadFile(file1)
	if err != nil {
		return fmt.Errorf("failed to read %s: %v", file1, err)
	}

	data2, err := ioutil.ReadFile(file2)
	if err != nil {
		return fmt.Errorf("failed to read %s: %v", file2, err)
	}

	var json1, json2 map[string]interface{}
	if err := json.Unmarshal(data1, &json1); err != nil {
		return fmt.Errorf("failed to parse JSON from %s: %v", file1, err)
	}

	if err := json.Unmarshal(data2, &json2); err != nil {
		return fmt.Errorf("failed to parse JSON from %s: %v", file2, err)
	}

	if err := removeFields(json1, ignorePaths); err != nil {
		return fmt.Errorf("failed to remove fields from %s: %v", file1, err)
	}

	if err := removeFields(json2, ignorePaths); err != nil {
		return fmt.Errorf("failed to remove fields from %s: %v", file2, err)
	}

	if !reflect.DeepEqual(json1, json2) {
		return fmt.Errorf("JSON files are not equal")
	}

	return nil
}

func loadIgnorePathsFromFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open ignore file %s: %v", filePath, err)
	}
	defer file.Close()

	var paths []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			paths = append(paths, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read ignore file %s: %v", filePath, err)
	}

	return paths, nil
}

func main() {
	ignoreFile := flag.String("i", "", "Path to a file containing JSONPath expressions to ignore during comparison")
	ignore := flag.String("d", "", "Comma-separated JSONPath expressions to ignore during comparison. Can be specified multiple times.")
	flag.Parse()

	if len(flag.Args()) < 2 {
		fmt.Println("Usage: json-compare [-i <ignore-file>] [-d <ignore-paths>] <file1.json> <file2.json>")
		os.Exit(1)
	}

	file1 := flag.Arg(0)
	file2 := flag.Arg(1)

	ignorePaths := []string{}

	// Load ignore paths from file if specified
	if *ignoreFile != "" {
		filePaths, err := loadIgnorePathsFromFile(*ignoreFile)
		if err != nil {
			fmt.Printf("Error loading ignore file: %v\n", err)
			os.Exit(1)
		}
		ignorePaths = append(ignorePaths, filePaths...)
	}

	// Parse the -d flag values into ignorePaths
	if *ignore != "" {
		ignorePaths = append(ignorePaths, strings.Split(*ignore, ",")...)
	}

	fmt.Printf("Comparing %s and %s with ignore paths: %v\n", file1, file2, ignorePaths)

	if err := compareJSON(file1, file2, ignorePaths); err != nil {
		fmt.Printf("Comparison failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("JSON files are equal.")
}
