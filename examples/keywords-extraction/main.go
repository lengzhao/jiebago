// Keywords extraction example demonstrates TF-IDF and TextRank for extracting keywords.
package main

import (
	"fmt"

	"github.com/lengzhao/jiebago/analyse"
)

func main() {
	// TF-IDF Example
	fmt.Println("=== TF-IDF Keywords Extraction ===")

	var tagger analyse.TagExtracter
	if err := tagger.LoadDictionary("../../embed/dict.txt"); err != nil {
		panic(err)
	}
	if err := tagger.LoadIdf("idf.txt"); err != nil {
		panic(err)
	}

	sentence := "此外，公司拟对全资子公司吉林欧亚置业有限公司增资4.3亿元，增资后，吉林欧亚置业注册资本由7000万元增加到5亿元。吉林欧亚置业主要经营范围为房地产开发及百货零售等业务。"

	keywords := tagger.ExtractTags(sentence, 10)
	fmt.Println("\nOriginal text:")
	fmt.Println(sentence)
	fmt.Println("\nTop 10 keywords (TF-IDF):")
	for i, kw := range keywords {
		fmt.Printf("%d. %s (weight: %.4f)\n", i+1, kw.Text(), kw.Weight())
	}

	// TextRank Example
	fmt.Println("\n\n=== TextRank Keywords Extraction ===")

	var ranker analyse.TextRanker
	if err := ranker.LoadDictionary("../../embed/dict.txt"); err != nil {
		panic(err)
	}

	keywords2 := ranker.TextRank(sentence, 10)
	fmt.Println("Top 10 keywords (TextRank):")
	for i, kw := range keywords2 {
		fmt.Printf("%d. %s (weight: %.4f)\n", i+1, kw.Text(), kw.Weight())
	}

	// Another example with custom text
	fmt.Println("\n\n=== Another Example ===")
	text2 := "自然语言处理是人工智能领域中的一个重要方向。它研究能实现人与计算机之间用自然语言进行有效通信的各种理论和方法。"

	keywords3 := tagger.ExtractTags(text2, 5)
	fmt.Println("Text:", text2)
	fmt.Println("\nTop 5 keywords:")
	for i, kw := range keywords3 {
		fmt.Printf("%d. %s\n", i+1, kw.Text())
	}
}
