package main

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

func TestCompareJSON(t *testing.T) {
	tests := []struct {
		name        string
		file1       string
		file2       string
		ignorePaths []string
		wantError   bool
	}{
		{
			name:        "Equal JSON without ignore paths",
			file1:       "testdata/equal1.json",
			file2:       "testdata/equal2.json",
			ignorePaths: nil,
			wantError:   false,
		},
		{
			name:        "Different JSON without ignore paths",
			file1:       "testdata/different1.json",
			file2:       "testdata/different2.json",
			ignorePaths: nil,
			wantError:   true,
		},
		{
			name:        "Equal JSON with ignore paths",
			file1:       "testdata/ignore1.json",
			file2:       "testdata/ignore2.json",
			ignorePaths: []string{"$.ignoreField", "$.nested.ignoreField"},
			wantError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data1, err := os.ReadFile(tt.file1)
			if err != nil {
				t.Fatalf("Failed to read file1: %v", err)
			}
			var json1 map[string]interface{}
			if err := json.Unmarshal(data1, &json1); err != nil {
				t.Fatalf("Failed to parse JSON from file1: %v", err)
			}

			data2, err := os.ReadFile(tt.file2)
			if err != nil {
				t.Fatalf("Failed to read file2: %v", err)
			}
			var json2 map[string]interface{}
			if err := json.Unmarshal(data2, &json2); err != nil {
				t.Fatalf("Failed to parse JSON from file2: %v", err)
			}

			err = compareJSONMaps(json1, json2, tt.ignorePaths)
			if (err != nil) != tt.wantError {
				t.Errorf("compareJSONMaps() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestIgnoreFileOption(t *testing.T) {
	// Create a temporary ignore file
	ignoreFilePath := "testdata/ignore_file.txt"
	ignoreFileContent := `
$.ignoreField
$.nested.ignoreField
`
	if err := os.WriteFile(ignoreFilePath, []byte(ignoreFileContent), 0644); err != nil {
		t.Fatalf("Failed to create ignore file: %v", err)
	}
	defer os.Remove(ignoreFilePath)

	tests := []struct {
		name       string
		file1      string
		file2      string
		ignoreFile string
		wantError  bool
	}{
		{
			name:       "Equal JSON with ignore file",
			file1:      "testdata/ignore1.json",
			file2:      "testdata/ignore2.json",
			ignoreFile: ignoreFilePath,
			wantError:  false,
		},
		{
			name:       "Different JSON with ignore file",
			file1:      "testdata/different1.json",
			file2:      "testdata/different2.json",
			ignoreFile: ignoreFilePath,
			wantError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ignorePaths, err := loadIgnorePathsFromFile(tt.ignoreFile)
			if err != nil {
				t.Fatalf("Failed to load ignore file: %v", err)
			}

			data1, err := os.ReadFile(tt.file1)
			if err != nil {
				t.Fatalf("Failed to read file1: %v", err)
			}
			var json1 map[string]interface{}
			if err := json.Unmarshal(data1, &json1); err != nil {
				t.Fatalf("Failed to parse JSON from file1: %v", err)
			}

			data2, err := os.ReadFile(tt.file2)
			if err != nil {
				t.Fatalf("Failed to read file2: %v", err)
			}
			var json2 map[string]interface{}
			if err := json.Unmarshal(data2, &json2); err != nil {
				t.Fatalf("Failed to parse JSON from file2: %v", err)
			}

			err = compareJSONMaps(json1, json2, ignorePaths)
			if (err != nil) != tt.wantError {
				t.Errorf("compareJSONMaps() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestRemoveFields(t *testing.T) {
	tests := []struct {
		name      string
		jsonData  map[string]interface{}
		paths     []string
		expected  map[string]interface{}
		wantError bool
	}{
		{
			name: "Remove single field",
			jsonData: map[string]interface{}{
				"key":         "value",
				"ignoreField": "ignore",
			},
			paths: []string{"$.ignoreField"},
			expected: map[string]interface{}{
				"key": "value",
			},
			wantError: false,
		},
		{
			name: "Remove nested field",
			jsonData: map[string]interface{}{
				"key": "value",
				"nested": map[string]interface{}{
					"ignoreField": "ignore",
				},
			},
			paths: []string{"$.nested.ignoreField"},
			expected: map[string]interface{}{
				"key":    "value",
				"nested": map[string]interface{}{},
			},
			wantError: false,
		},
		{
			name: "Field not found",
			jsonData: map[string]interface{}{
				"key": "value",
			},
			paths: []string{"$.nonExistentField"},
			expected: map[string]interface{}{
				"key": "value",
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData := tt.jsonData
			err := removeFields(jsonData, tt.paths)
			if (err != nil) != tt.wantError {
				t.Errorf("removeFields() error = %v, wantError %v", err, tt.wantError)
			}
			if !reflect.DeepEqual(jsonData, tt.expected) {
				t.Errorf("removeFields() = %v, expected %v", jsonData, tt.expected)
			}
		})
	}
}

func TestMain(m *testing.M) {
	// Setup test data directory
	if err := os.MkdirAll("testdata", 0755); err != nil {
		panic(err)
	}
	defer os.RemoveAll("testdata")

	// Create test files
	files := map[string]string{
		"testdata/equal1.json":     `{"key": "value"}`,
		"testdata/equal2.json":     `{"key": "value"}`,
		"testdata/different1.json": `{"key": "value1"}`,
		"testdata/different2.json": `{"key": "value2"}`,
		"testdata/ignore1.json":    `{"key": "value", "ignoreField": "ignore", "nested": {"ignoreField": "ignore"}}`,
		"testdata/ignore2.json":    `{"key": "value", "ignoreField": "different", "nested": {"ignoreField": "different"}}`,
	}

	for path, content := range files {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			panic(err)
		}
	}

	os.Exit(m.Run())
}
