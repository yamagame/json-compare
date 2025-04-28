# JSON比較ツール

このツールは、特定のフィールドを無視しながら2つのJSONファイルを比較することができます。無視するフィールドは、JSONPath式を使用してコマンドラインオプションで直接指定するか、式を含むファイルを指定することで設定できます。

## 使用方法

```bash
json-compare [-i <ignore-file>] [-d <ignore-paths>] <file1.json> <file2.json>
```

### オプション

- `-i <ignore-file>`: 比較中に無視するJSONPath式を含むファイルへのパス。ファイル内の各行に1つのJSONPath式を記述します。
- `-d <ignore-paths>`: 比較中に無視するJSONPath式をカンマ区切りで指定します。このオプションは複数回指定できます。

### 使用例

#### フィールドを無視せずに2つのJSONファイルを比較
```bash
json-compare file1.json file2.json
```

#### `-d`オプションを使用して特定のフィールドを無視しながら2つのJSONファイルを比較
```bash
json-compare -d "$.ignoreField,$.nested.ignoreField" file1.json file2.json
```

#### ファイルで指定されたフィールドを無視しながら2つのJSONファイルを比較
```bash
json-compare -i ignore-file.txt file1.json file2.json
```

### 無視ファイルのフォーマット

無視ファイルには、1行に1つのJSONPath式を記述します。例:
```
$.ignoreField
$.nested.ignoreField
```

## 必要条件

- Go 1.16以降

## テストの実行

このツールのテストを実行するには、以下のコマンドを使用します:
```bash
go test ./...
```