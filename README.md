# JSON Compare Tool

This tool allows you to compare two JSON files while ignoring specific fields. You can specify fields to ignore using JSONPath expressions either directly via a command-line option or by providing a file containing the expressions.

## Usage

```bash
json-compare [-i <ignore-file>] [-d <ignore-paths>] <file1.json> <file2.json>
```

### Options

- `-i <ignore-file>`: Path to a file containing JSONPath expressions to ignore during comparison. Each line in the file should contain one JSONPath expression.
- `-d <ignore-paths>`: Comma-separated JSONPath expressions to ignore during comparison. This option can be specified multiple times.

### Examples

#### Compare two JSON files without ignoring any fields
```bash
json-compare file1.json file2.json
```

#### Compare two JSON files while ignoring specific fields using the `-d` option
```bash
json-compare -d "$.ignoreField,$.nested.ignoreField" file1.json file2.json
```

#### Compare two JSON files while ignoring fields specified in a file
```bash
json-compare -i ignore-file.txt file1.json file2.json
```

### Ignore File Format

The ignore file should contain one JSONPath expression per line. For example:
```
$.ignoreField
$.nested.ignoreField
```

## Requirements

- Go 1.16 or later

## Running Tests

To run the tests for this tool, use the following command:
```bash
go test ./...
```