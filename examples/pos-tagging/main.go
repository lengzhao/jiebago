// POS (Part-of-Speech) tagging example demonstrates word segmentation with part-of-speech tags.
package main

import (
	"fmt"

	"github.com/lengzhao/jiebago/posseg"
)

func main() {
	var seg posseg.Segmenter

	// Load dictionary
	if err := seg.LoadDictionary("../../embed/dict.txt"); err != nil {
		panic(err)
	}

	// Example 1: Basic POS tagging
	sentence := "我爱北京天安门"
	fmt.Println("=== Basic POS Tagging ===")
	fmt.Printf("Text: %s\n\n", sentence)

	for segment := range seg.Cut(sentence, true) {
		fmt.Printf("Word: %-10s POS: %s\n", segment.Text(), segment.Pos())
	}

	// Example 2: Complex sentence
	fmt.Println("\n=== Complex Sentence ===")
	sentence2 := "我来到北京清华大学，看到了很多优秀学生。"
	fmt.Printf("Text: %s\n\n", sentence2)

	for segment := range seg.Cut(sentence2, true) {
		fmt.Printf("%-12s (%s)\n", segment.Text(), posName(segment.Pos()))
	}

	// Example 3: Names and places
	fmt.Println("\n=== Names and Places ===")
	sentence3 := "李小福和张三去上海市南京路购物"
	fmt.Printf("Text: %s\n\n", sentence3)

	for segment := range seg.Cut(sentence3, true) {
		fmt.Printf("%-10s (%s)\n", segment.Text(), segment.Pos())
	}

	// Example 4: Technical terms
	fmt.Println("\n=== Technical Terms ===")
	sentence4 := "自然语言处理和人工智能在云计算平台上的应用"
	fmt.Printf("Text: %s\n\n", sentence4)

	for segment := range seg.Cut(sentence4, true) {
		fmt.Printf("%-10s (%s)\n", segment.Text(), posName(segment.Pos()))
	}
}

// posName returns full name for POS tag
func posName(tag string) string {
	names := map[string]string{
		"n":   "名词",
		"v":   "动词",
		"a":   "形容词",
		"d":   "副词",
		"m":   "数词",
		"q":   "量词",
		"r":   "代词",
		"p":   "介词",
		"c":   "连词",
		"u":   "助词",
		"e":   "叹词",
		"y":   "语气词",
		"o":   "拟声词",
		"h":   "前缀",
		"k":   "后缀",
		"x":   "字符串",
		"w":   "标点符号",
		"nr":  "人名",
		"ns":  "地名",
		"nt":  "机构团体",
		"nz":  "其他专名",
		"t":   "时间词",
		"s":   "处所词",
		"f":   "方位词",
		"vn":  "名动词",
		"vd":  "副动词",
		"an":  "名形词",
		"ad":  "副形词",
	}

	if name, ok := names[tag]; ok {
		return name
	}
	return "未知"
}
