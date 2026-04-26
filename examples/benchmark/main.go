// Benchmark example demonstrates the performance of jiebago.
// Run with: go run main.go
package main

import (
	"fmt"
	"time"

	"github.com/lengzhao/jiebago"
)

func main() {
	fmt.Println("=== Jiebago Performance Benchmark ===")
	fmt.Println()

	// 预热默认分词器
	fmt.Println("Warming up default segmenter...")
	start := time.Now()
	for range jiebago.Default.Cut("预热", true) {
	}
	initTime := time.Since(start)
	fmt.Printf("Initialization time: %v\n", initTime)
	fmt.Println()

	// 测试用例
	testCases := []struct {
		name    string
		text    string
		mode    string
		hmm     bool
		iter    int
	}{
		{
			name: "Short text (10 chars)",
			text: "我爱北京天安门",
			mode: "cut",
			hmm:  true,
			iter: 10000,
		},
		{
			name: "Medium text (50 chars)",
			text: "工信处女干事每月经过下属科室都要亲口交代24口交换机等技术性器件",
			mode: "cut",
			hmm:  true,
			iter: 5000,
		},
		{
			name: "Long text (200 chars)",
			text: "自然语言处理是人工智能领域中的一个重要方向。它研究能实现人与计算机之间用自然语言进行有效通信的各种理论和方法。自然语言处理是一门融语言学、计算机科学、数学于一体的科学。",
			mode: "cut",
			hmm:  true,
			iter: 1000,
		},
		{
			name: "New word recognition",
			text: "他来到了网易杭研大厦",
			mode: "cut",
			hmm:  true,
			iter: 5000,
		},
	}

	// 运行基准测试
	fmt.Println("Benchmark Results:")
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Printf("%-25s %12s %12s %12s\n", "Test Case", "Iterations", "Total Time", "Per Op")
	fmt.Println("--------------------------------------------------------------------------------")

	for _, tc := range testCases {
		start := time.Now()
		for i := 0; i < tc.iter; i++ {
			switch tc.mode {
			case "cut":
				for range jiebago.Default.Cut(tc.text, tc.hmm) {
				}
			case "cutall":
				for range jiebago.Default.CutAll(tc.text) {
				}
			}
		}
		duration := time.Since(start)
		perOp := duration / time.Duration(tc.iter)

		fmt.Printf("%-25s %12d %12v %12v\n", tc.name, tc.iter, duration, perOp)
	}

	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println()

	// 不同分词模式对比
	fmt.Println("Segmentation Mode Comparison:")
	fmt.Println("--------------------------------------------------------------------------------")
	sentence := "小明硕士毕业于中国科学院计算所，后在日本京都大学深造"
	iterations := 3000

	modes := []struct {
		name string
		fn   func() int
	}{
		{
			name: "CutAll (全模式)",
			fn: func() int {
				count := 0
				for range jiebago.Default.CutAll(sentence) {
					count++
				}
				return count
			},
		},
		{
			name: "Cut no HMM (精确模式)",
			fn: func() int {
				count := 0
				for range jiebago.Default.Cut(sentence, false) {
					count++
				}
				return count
			},
		},
		{
			name: "Cut with HMM (新词识别)",
			fn: func() int {
				count := 0
				for range jiebago.Default.Cut(sentence, true) {
					count++
				}
				return count
			},
		},
		{
			name: "CutForSearch (搜索引擎)",
			fn: func() int {
				count := 0
				for range jiebago.Default.CutForSearch(sentence, true) {
					count++
				}
				return count
			},
		},
	}

	fmt.Printf("%-25s %12s %12s %12s\n", "Mode", "Words", "Total Time", "Per Op")
	fmt.Println("--------------------------------------------------------------------------------")

	for _, mode := range modes {
		// 预热
		mode.fn()

		start := time.Now()
		var wordCount int
		for i := 0; i < iterations; i++ {
			wordCount = mode.fn()
		}
		duration := time.Since(start)
		perOp := duration / time.Duration(iterations)

		fmt.Printf("%-25s %12d %12v %12v\n", mode.name, wordCount, duration, perOp)
	}

	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println()

	// 词频查询测试
	fmt.Println("Frequency Lookup Performance:")
	words := []string{"北京", "清华大学", "人工智能", "自然语言处理", "nonexistent"}
	lookupIter := 1000000

	start = time.Now()
	for i := 0; i < lookupIter; i++ {
		for _, word := range words {
			jiebago.Default.Frequency(word)
		}
	}
	lookupDuration := time.Since(start)
	lookupPerOp := lookupDuration / time.Duration(lookupIter*len(words))

	fmt.Printf("  Words: %v\n", words)
	fmt.Printf("  Iterations: %d per word\n", lookupIter)
	fmt.Printf("  Total lookups: %d\n", lookupIter*len(words))
	fmt.Printf("  Total time: %v\n", lookupDuration)
	fmt.Printf("  Per lookup: %v\n", lookupPerOp)
	fmt.Println()

	// 并发性能测试
	fmt.Println("Concurrent Segmentation Test:")
	concurrentIter := 10000
	numGoroutines := 10

	texts := []string{
		"我来到北京清华大学",
		"他来到了网易杭研大厦",
		"工信处女干事每月经过下属科室都要亲口交代24口交换机等技术性器件的安装工作",
		"自然语言处理是人工智能领域中的一个重要方向",
		"小明硕士毕业于中国科学院计算所",
	}

	start = time.Now()
	done := make(chan bool, numGoroutines)

	for g := 0; g < numGoroutines; g++ {
		go func(id int) {
			for i := 0; i < concurrentIter; i++ {
				text := texts[i%len(texts)]
				for range jiebago.Default.Cut(text, true) {
				}
			}
			done <- true
		}(g)
	}

	for g := 0; g < numGoroutines; g++ {
		<-done
	}

	concurrentDuration := time.Since(start)
	totalOps := numGoroutines * concurrentIter

	fmt.Printf("  Goroutines: %d\n", numGoroutines)
	fmt.Printf("  Operations per goroutine: %d\n", concurrentIter)
	fmt.Printf("  Total operations: %d\n", totalOps)
	fmt.Printf("  Total time: %v\n", concurrentDuration)
	fmt.Printf("  Ops per second: %.0f\n", float64(totalOps)/concurrentDuration.Seconds())
	fmt.Println()

	// 内存使用估算
	fmt.Println("Memory Usage:")
	fmt.Println("  Note: Dictionary loaded once (~89 MB)")
	fmt.Println("  Per-segment memory depends on text length")
	fmt.Println()

	fmt.Println("=== Benchmark Complete ===")
}
