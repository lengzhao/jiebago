# 结巴分词 Go 语言版：Jiebago

[![Go](https://github.com/lengzhao/jiebago/actions/workflows/go.yml/badge.svg)](https://github.com/lengzhao/jiebago/actions/workflows/go.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/lengzhao/jiebago.svg)](https://pkg.go.dev/github.com/lengzhao/jiebago)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

> **Fork 声明**：本项目 fork 自 [wangbin/jiebago](https://github.com/wangbin/jiebago)，原项目是 [结巴分词](https://github.com/fxsjy/jieba)（Python 版）的 Go 语言实现。
>
> 原仓库维护停滞，此 fork 进行了现代化更新：添加 Go Modules 支持、更新 CI/CD、修复兼容性问题和拼写错误。

[结巴分词](https://github.com/fxsjy/jieba) 是由 [@fxsjy](https://github.com/fxsjy) 使用 Python 编写的中文分词组件，Jiebago 是结巴分词的 Golang 语言实现。

## 特性

- **开箱即用**：内置字典数据，无需额外文件
- 支持多种分词模式：全模式、精确模式、搜索引擎模式
- 支持词性标注 (POS tagging)
- 支持关键词提取 (TF-IDF, TextRank)
- 支持并行分词
- 支持自定义词典
- 支持 Bleve 搜索引擎集成

## 安装

```bash
go get github.com/lengzhao/jiebago
```

## 快速开始（开箱即用）

使用默认分词器，无需加载字典文件：

```go
package main

import (
    "fmt"
    "github.com/lengzhao/jiebago"
)

func main() {
    sentence := "我来到北京清华大学"

    // 精确模式
    for word := range jiebago.Default.Cut(sentence, false) {
        fmt.Printf("%s / ", word)
    }
    // 输出: 我 / 来到 / 北京 / 清华大学 /
}
```

## 使用指南

### 默认分词器 vs 自定义分词器

| 方式 | 适用场景 | 示例 |
|------|----------|------|
| **默认分词器** | 快速开始、简单使用 | `jiebago.Default.Cut(text, true)` |
| **自定义分词器** | 需要自定义字典、精细控制 | `seg.LoadDictionary("dict.txt")` |

### 默认分词器（开箱即用）

```go
package main

import (
    "fmt"
    "github.com/lengzhao/jiebago"
)

func main() {
    sentence := "我来到北京清华大学"

    // 【全模式】获取所有可能的分词结果
    fmt.Println("【全模式】")
    for word := range jiebago.Default.CutAll(sentence) {
        fmt.Printf("%s / ", word)
    }

    // 【精确模式】最准确的分词结果（默认）
    fmt.Println("\n【精确模式】")
    for word := range jiebago.Default.Cut(sentence, false) {
        fmt.Printf("%s / ", word)
    }

    // 【新词识别】使用 HMM 模型识别未登录词
    fmt.Println("\n【新词识别】")
    for word := range jiebago.Default.Cut("他来到了网易杭研大厦", true) {
        fmt.Printf("%s / ", word)
    }

    // 【搜索引擎模式】适合建立搜索引擎索引
    fmt.Println("\n【搜索引擎模式】")
    for word := range jiebago.Default.CutForSearch("小明硕士毕业于中国科学院计算所", true) {
        fmt.Printf("%s / ", word)
    }

    // 添加自定义词
    jiebago.Default.AddWord("超敏C反应蛋白", 1000.0)

    // 加载用户词典（不覆盖内置词典）
    jiebago.Default.LoadUserDictionary("userdict.txt")
}
```

### 自定义分词器（传统方式）

如果需要使用自定义字典，可以创建 Segmenter 实例：

```go
package main

import (
    "fmt"
    "github.com/lengzhao/jiebago"
)

func main() {
    var seg jiebago.Segmenter

    // 从文件加载字典
    if err := seg.LoadDictionary("dict.txt"); err != nil {
        panic(err)
    }

    // 或者从字节加载（适合嵌入字典）
    // seg.LoadDictionaryFromBytes(dictData)

    sentence := "我来到北京清华大学"
    for word := range seg.Cut(sentence, false) {
        fmt.Printf("%s / ", word)
    }
}
```

## 使用指南：如何选择分词模式

| 场景 | 推荐方法 | 示例 | 优缺点 |
|------|----------|------|--------|
| **文本分析**<br>自然语言处理、情感分析 | `Cut(text, false)` | `seg.Cut("我来到北京清华大学", false)` | ✅ 结果精确，无冗余<br>❌ 可能遗漏复合词中的子词 |
| **全文检索**<br>搜索引擎索引 | `CutForSearch(text, true)` | `seg.CutForSearch(sentence, true)` | ✅ 召回率高，长短词兼顾<br>❌ 结果量大，需更多存储 |
| **歧义识别**<br>语义消歧、关键词提取 | `CutAll(text)` | `seg.CutAll("长春市长春药店")` | ✅ 展示所有可能性<br>❌ 大量冗余，需要后处理 |
| **新词发现**<br>处理网络用语、专有名词 | `Cut(text, true)` | `seg.Cut("他来到了网易杭研大厦", true)` | ✅ 自动识别未登录词<br>❌ 可能误分（如人名） |

### 详细说明

#### 1. 精确模式 (Accurate Mode) - `Cut(text, false)`
**适用场景**：文本分类、情感分析、机器翻译

```go
// 文本分析：需要最准确的分词结果
segments := seg.Cut("小明硕士毕业于中国科学院计算所", false)
// 输出: [小明 硕士 毕业 于 中国科学院 计算所]
```

- **优点**：消除歧义能力强，结果最准确
- **缺点**：不能识别词典中未收录的新词
- **性能**：~700KB/s

#### 2. 搜索引擎模式 (Search Mode) - `CutForSearch(text, hmm)`
**适用场景**：搜索引擎、全文检索、关键词索引

```go
// 搜索引擎索引：需要高召回率
segments := seg.CutForSearch("小明硕士毕业于中国科学院计算所", true)
// 输出: [小明 硕士 毕业 于 中国 科学 学院 科学院 中国科学院 计算 计算所]
```

- **优点**：在长词基础上再切分短词，提高召回率
- **缺点**：结果量大，索引体积增加
- **性能**：与精确模式相当

#### 3. 全模式 (Full Mode) - `CutAll(text)`
**适用场景**：歧义分析、语言学研究、展示所有可能

```go
// 歧义识别：展示所有可能的分词组合
segments := seg.CutAll("长春市长春药店")
// 输出: [长春 长春市 长春 药店, 长春 市长 春药店, ...]
```

- **优点**：速度最快，展示所有可能
- **缺点**：大量冗余结果，需要人工筛选
- **性能**：~2MB/s

#### 4. 新词识别 (HMM Mode) - `Cut(text, true)`
**适用场景**：处理网络新词、公司名、产品名

```go
// 新词识别：自动发现未登录词
segments := seg.Cut("他来到了网易杭研大厦", true)
// 输出: [他 来到 了 网易 杭研 大厦]
// "杭研" 是词典中没有的新词
```

- **优点**：基于 HMM 模型自动识别新词
- **缺点**：可能误判，计算量稍大
- **性能**：~500KB/s

## 高级用法

### 词性标注 (POS)

```go
import "github.com/lengzhao/jiebago/posseg"

var seg posseg.Segmenter
seg.LoadDictionary("embed/dict.txt")

for segment := range seg.Cut("我爱北京天安门", true) {
    fmt.Printf("%s %s\n", segment.Text(), segment.Pos())
}
```

### 关键词提取 (TF-IDF)

```go
import "github.com/lengzhao/jiebago/analyse"

var tagger analyse.TagExtracter
tagger.LoadDictionary("embed/dict.txt")
tagger.LoadIdf("idf.txt")  // 需要提供 IDF 字典文件

tags := tagger.ExtractTags(sentence, 10)
```

### TextRank 关键词提取

```go
var ranker analyse.TextRanker
ranker.LoadDictionary("embed/dict.txt")

result := ranker.TextRank(sentence, 10)
```

### 并行分词

```go
runtime.GOMAXPROCS(runtime.NumCPU())

for line := range lines {
    go func(l string) {
        for word := range seg.Cut(l, true) {
            // 处理分词结果
        }
    }(line)
}
```

## 示例代码

查看 [examples](examples/) 目录获取更多实用示例：

| 示例 | 说明 |
|------|------|
| `examples/quick-start/` | **快速开始（推荐）** - 开箱即用演示 |
| `examples/basic/` | 四种分词模式完整演示 |
| `examples/web-api/` | REST API 服务 |
| `examples/batch-processing/` | 并行批量处理 |
| `examples/keywords-extraction/` | TF-IDF / TextRank 关键词提取 |
| `examples/pos-tagging/` | 词性标注 |
| `examples/custom-dictionary/` | 自定义词典管理 |
| `examples/benchmark/` | 性能测试和对比 |

## 分词速度

在 Apple M4 (ARM64) 上测试：

| 模式 | 速度 | 单次耗时 | 内存分配 |
|------|------|---------|---------|
| 全模式 (CutAll) | ~78,000 次/秒 | 12.7 µs/op | 8.2 KB/op |
| 精确模式 (Cut) | ~63,000 次/秒 | 15.8 µs/op | 11.1 KB/op |
| 新词识别 (Cut+HMM) | ~47,000 次/秒 | 21.1 µs/op | 14.8 KB/op |
| 搜索引擎 (CutForSearch) | ~45,000 次/秒 | 22.0 µs/op | 14.2 KB/op |

> 详细性能报告见 [PERFORMANCE.md](PERFORMANCE.md)

历史参考数据（AMD Phenom II X6 1055T @ 2.8GHz）：
 - 2MB / Second in Full Mode
 - 700KB / Second in Default Mode

## 许可证

MIT - 详见 [LICENSE](LICENSE) 文件

## 来源与致谢

- **原项目**: [wangbin/jiebago](https://github.com/wangbin/jiebago) - 由 wangbin 开发的结巴分词 Go 语言实现
- **Python 原版**: [fxsjy/jieba](https://github.com/fxsjy/jieba) - 由 fxsjy 使用 Python 编写的中文分词库
- **Fork 更新**: 本 fork 由 lengzhao 维护，主要进行现代化改造和 bug 修复
