# Jiebago Examples

This directory contains practical examples demonstrating various use cases of jiebago (结巴分词 Go 语言版).

## Quick Start

Each example is a standalone Go program. To run any example:

```bash
cd <example-folder>
go run main.go
```

## Examples Overview

| Example | Description | Key Features |
|---------|-------------|--------------|
| `basic/` | Fundamental segmentation modes | 精确模式, 全模式, 搜索引擎模式, HMM新词识别 |
| `web-api/` | HTTP API for text segmentation | REST API, JSON request/response |
| `batch-processing/` | Parallel processing of large datasets | goroutines, channels, worker pool |
| `keywords-extraction/` | TF-IDF and TextRank keyword extraction | analyse package |
| `custom-dictionary/` | Custom words and user dictionaries | AddWord, LoadUserDictionary, SuggestFrequency |

## Directory Structure

```
examples/
├── README.md
├── basic/
│   └── main.go          # Basic usage demo
├── web-api/
│   └── main.go          # REST API server
├── batch-processing/
│   └── main.go          # Parallel batch processing
├── keywords-extraction/
│   └── main.go          # TF-IDF and TextRank
└── custom-dictionary/
    └── main.go          # Custom dictionary management
```

## Prerequisites

- Go 1.25 or later
- Dictionary file `dict.txt` in project root
- For keywords extraction: `idf.txt` in examples folder

## Example Output

### Basic Example
```
【精确模式】
我 / 来到 / 北京 / 清华大学 /

【全模式】
我 / 来到 / 北京 / 清华 / 清华大学 / 华大 / 大学 /
```

### Web API Example
```bash
curl -X POST http://localhost:8080/segment \
  -H 'Content-Type: application/json' \
  -d '{"text":"我爱北京天安门","mode":"accurate"}'

# Response: {"words":["我","爱","北京","天安门"]}
```

### Keywords Extraction Example
```
Top 5 keywords (TF-IDF):
1. 欧亚 (weight: 0.8781)
2. 置业 (weight: 0.5620)
3. 吉林 (weight: 1.0000)
4. 增资 (weight: 0.3606)
5. 子公司 (weight: 0.3531)
```

## More Examples Coming

- `pos-tagging/` - Part-of-speech tagging
- `search-engine/` - Full-text search integration with Bleve
- `streaming/` - Real-time text processing
- `performance-benchmark/` - Performance comparison

## Contributing

Feel free to add more examples! Each example should:
1. Have its own directory
2. Include a `main.go` that can run independently
3. Demonstrate specific functionality
4. Include comments explaining the code
