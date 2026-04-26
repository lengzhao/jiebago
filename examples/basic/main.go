// Basic example demonstrates fundamental jiebago usage.
// This example shows all four segmentation modes using the default segmenter.
package main

import (
	"fmt"

	"github.com/lengzhao/jiebago"
)

func main() {
	// 使用默认分词器（开箱即用，内置字典，无需手动加载）
	// Use the default segmenter (works out-of-the-box with embedded dictionary)

	sentence := "我来到北京清华大学"

	// 1. 精确模式 - Accurate mode (default)
	fmt.Println("【精确模式】")
	for word := range jiebago.Default.Cut(sentence, false) {
		fmt.Printf("%s / ", word)
	}
	fmt.Println()

	// 2. 全模式 - Full mode (gets all possible words)
	fmt.Println("\n【全模式】")
	for word := range jiebago.Default.CutAll(sentence) {
		fmt.Printf("%s / ", word)
	}
	fmt.Println()

	// 3. 搜索引擎模式 - Search engine mode
	fmt.Println("\n【搜索引擎模式】")
	sentence2 := "小明硕士毕业于中国科学院计算所"
	for word := range jiebago.Default.CutForSearch(sentence2, true) {
		fmt.Printf("%s / ", word)
	}
	fmt.Println()

	// 4. 使用HMM新词识别
	fmt.Println("\n【新词识别(HMM)】")
	sentence3 := "他来到了网易杭研大厦"
	for word := range jiebago.Default.Cut(sentence3, true) {
		fmt.Printf("%s / ", word)
	}
	fmt.Println()

	// 也可以使用自定义分词器（需要手动加载字典）
	// You can also use a custom segmenter (requires manual dictionary loading)
	fmt.Println("\n【使用自定义分词器】")
	var seg jiebago.Segmenter
	if err := seg.LoadDictionary("../../embed/dict.txt"); err != nil {
		panic(err)
	}
	for word := range seg.Cut(sentence, false) {
		fmt.Printf("%s / ", word)
	}
	fmt.Println()
}
