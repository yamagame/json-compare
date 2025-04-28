package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
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

func compareJSONMaps(json1, json2 map[string]interface{}, ignorePaths []string) error {
	if err := removeFields(json1, ignorePaths); err != nil {
		return fmt.Errorf("failed to remove fields from first JSON: %v", err)
	}

	if err := removeFields(json2, ignorePaths); err != nil {
		return fmt.Errorf("failed to remove fields from second JSON: %v", err)
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

func readJSONFromStdin() (map[string]interface{}, error) {
	reader := bufio.NewReader(os.Stdin)
	input, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read from stdin: %v", err)
	}

	var jsonData map[string]interface{}
	if err := json.Unmarshal(input, &jsonData); err != nil {
		return nil, fmt.Errorf("failed to parse JSON from stdin: %v", err)
	}

	return jsonData, nil
}

func main() {
	// Add a version flag
	version := flag.Bool("v", false, "Print the version of the tool")
	ignoreFile := flag.String("i", "", "Path to a file containing JSONPath expressions to ignore during comparison")
	ignore := flag.String("d", "", "Comma-separated JSONPath expressions to ignore during comparison. Can be specified multiple times.")
	useStdin := flag.Bool("p", false, "Read one of the JSON files from standard input")
	flag.Parse()

	if *version {
		fmt.Println("json-compare version 0.1.0")
		os.Exit(0)
	}

	if len(flag.Args()) < 2 && !*useStdin {
		fmt.Println("Usage: json-compare [-i <ignore-file>] [-d <ignore-paths>] [-p] <file1.json> <file2.json>")
		os.Exit(1)
	}

	var json1, json2 map[string]interface{}
	var err error

	if *useStdin {
		json1, err = readJSONFromStdin()
		if err != nil {
			fmt.Printf("Error reading JSON from stdin: %v\n", err)
			os.Exit(1)
		}
		file2 := flag.Arg(0)
		data2, err := os.ReadFile(file2)
		if err != nil {
			fmt.Printf("Failed to read %s: %v\n", file2, err)
			os.Exit(1)
		}
		if err := json.Unmarshal(data2, &json2); err != nil {
			fmt.Printf("Failed to parse JSON from %s: %v\n", file2, err)
			os.Exit(1)
		}
	} else {
		file1 := flag.Arg(0)
		file2 := flag.Arg(1)
		data1, err := os.ReadFile(file1)
		if err != nil {
			fmt.Printf("Failed to read %s: %v\n", file1, err)
			os.Exit(1)
		}
		if err := json.Unmarshal(data1, &json1); err != nil {
			fmt.Printf("Failed to parse JSON from %s: %v\n", file1, err)
			os.Exit(1)
		}
		data2, err := os.ReadFile(file2)
		if err != nil {
			fmt.Printf("Failed to read %s: %v\n", file2, err)
			os.Exit(1)
		}
		if err := json.Unmarshal(data2, &json2); err != nil {
			fmt.Printf("Failed to parse JSON from %s: %v\n", file2, err)
			os.Exit(1)
		}
	}

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

	fmt.Printf("Comparing JSON files with ignore paths: %v\n", ignorePaths)

	if err := compareJSONMaps(json1, json2, ignorePaths); err != nil {
		fmt.Printf("Comparison failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("JSON files are equal.")
}
