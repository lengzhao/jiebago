// Quick-start example demonstrates the out-of-the-box experience with jiebago.
// This example uses the default segmenter which has built-in dictionary.
package main

import (
	"fmt"

	"github.com/lengzhao/jiebago"
)

func main() {
	fmt.Println("=== Jiebago Quick Start ===")
	fmt.Println("No dictionary loading needed - works out of the box!")
	fmt.Println()

	sentence := "我来到北京清华大学"

	// 1. 精确模式 - 最准确的分词结果
	fmt.Println("【精确模式】")
	fmt.Print("Result: ")
	for word := range jiebago.Default.Cut(sentence, false) {
		fmt.Printf("%s / ", word)
	}
	fmt.Println()
	fmt.Println()

	// 2. 全模式 - 获取所有可能的分词结果
	fmt.Println("【全模式】")
	fmt.Print("Result: ")
	for word := range jiebago.Default.CutAll(sentence) {
		fmt.Printf("%s / ", word)
	}
	fmt.Println()
	fmt.Println()

	// 3. 搜索引擎模式 - 适合建立搜索索引
	fmt.Println("【搜索引擎模式】")
	sentence2 := "小明硕士毕业于中国科学院计算所"
	fmt.Print("Result: ")
	for word := range jiebago.Default.CutForSearch(sentence2, true) {
		fmt.Printf("%s / ", word)
	}
	fmt.Println()
	fmt.Println()

	// 4. 新词识别（HMM）- 识别未登录词
	fmt.Println("【新词识别(HMM)】")
	sentence3 := "他来到了网易杭研大厦"
	fmt.Print("Result: ")
	for word := range jiebago.Default.Cut(sentence3, true) {
		fmt.Printf("%s / ", word)
	}
	fmt.Println()
	fmt.Println()

	// 5. 添加自定义词
	fmt.Println("【添加自定义词】")
	jiebago.Default.AddWord("超敏C反应蛋白", 1000.0)
	sentence4 := "超敏C反应蛋白是什么？"
	fmt.Print("Before adding: ")
	// Reset to show original behavior by creating a new segmenter
	var seg jiebago.Segmenter
	seg.LoadDictionary("../../embed/dict.txt")
	for word := range seg.Cut(sentence4, false) {
		fmt.Printf("%s / ", word)
	}
	fmt.Println()

	fmt.Print("After adding:  ")
	for word := range jiebago.Default.Cut(sentence4, false) {
		fmt.Printf("%s / ", word)
	}
	fmt.Println()
	fmt.Println()

	// 6. 查询词频
	fmt.Println("【查询词频】")
	freq, ok := jiebago.Default.Frequency("清华大学")
	fmt.Printf("'清华大学' frequency: %.0f (exists: %v)\n", freq, ok)
	freq2, ok2 := jiebago.Default.Frequency("nonexistentword123")
	fmt.Printf("'nonexistentword123' frequency: %.0f (exists: %v)\n", freq2, ok2)
	fmt.Println()

	fmt.Println("=== Done! ===")
}
