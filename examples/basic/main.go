// Basic example demonstrates fundamental jiebago usage.
// This example shows all four segmentation modes.
package main

import (
	"fmt"

	"github.com/lengzhao/jiebago"
)

func main() {
	var seg jiebago.Segmenter

	// Load dictionary
	if err := seg.LoadDictionary("../../dict.txt"); err != nil {
		panic(err)
	}

	sentence := "我来到北京清华大学"

	// 1. 精确模式 - Accurate mode (default)
	fmt.Println("【精确模式】")
	for word := range seg.Cut(sentence, false) {
		fmt.Printf("%s / ", word)
	}
	fmt.Println()

	// 2. 全模式 - Full mode (gets all possible words)
	fmt.Println("\n【全模式】")
	for word := range seg.CutAll(sentence) {
		fmt.Printf("%s / ", word)
	}
	fmt.Println()

	// 3. 搜索引擎模式 - Search engine mode
	fmt.Println("\n【搜索引擎模式】")
	sentence2 := "小明硕士毕业于中国科学院计算所"
	for word := range seg.CutForSearch(sentence2, true) {
		fmt.Printf("%s / ", word)
	}
	fmt.Println()

	// 4. 使用HMM新词识别
	fmt.Println("\n【新词识别(HMM)】")
	sentence3 := "他来到了网易杭研大厦"
	for word := range seg.Cut(sentence3, true) {
		fmt.Printf("%s / ", word)
	}
	fmt.Println()
}
