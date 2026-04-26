package jiebago_test

import (
	"fmt"

	"github.com/lengzhao/jiebago"
)

func Example() {
	// 使用默认分词器（开箱即用，无需加载字典）
	print := func(ch <-chan string) {
		for word := range ch {
			fmt.Printf(" %s /", word)
		}
		fmt.Println()
	}

	fmt.Print("【全模式】：")
	print(jiebago.Default.CutAll("我来到北京清华大学"))

	fmt.Print("【精确模式】：")
	print(jiebago.Default.Cut("我来到北京清华大学", false))

	fmt.Print("【新词识别】：")
	print(jiebago.Default.Cut("他来到了网易杭研大厦", true))

	fmt.Print("【搜索引擎模式】：")
	print(jiebago.Default.CutForSearch("小明硕士毕业于中国科学院计算所，后在日本京都大学深造", true))
	// Output:
	// 【全模式】： 我 / 来到 / 北京 / 清华 / 清华大学 / 华大 / 大学 /
	// 【精确模式】： 我 / 来到 / 北京 / 清华大学 /
	// 【新词识别】： 他 / 来到 / 了 / 网易 / 杭研 / 大厦 /
	// 【搜索引擎模式】： 小明 / 硕士 / 毕业 / 于 / 中国 / 科学 / 学院 / 科学院 / 中国科学院 / 计算 / 计算所 / ， / 后 / 在 / 日本 / 京都 / 大学 / 日本京都大学 / 深造 /
}

func Example_customDictionary() {
	// 使用自定义字典文件（传统方式）
	var seg jiebago.Segmenter
	seg.LoadDictionary("embed/dict.txt")

	print := func(ch <-chan string) {
		for word := range ch {
			fmt.Printf(" %s /", word)
		}
		fmt.Println()
	}

	fmt.Print("【精确模式】：")
	print(seg.Cut("我来到北京清华大学", false))
	// Output:
	// 【精确模式】： 我 / 来到 / 北京 / 清华大学 /
}

func Example_suggestFrequency() {
	// 使用默认分词器进行词频建议
	print := func(ch <-chan string) {
		for word := range ch {
			fmt.Printf(" %s /", word)
		}
		fmt.Println()
	}
	sentence := "超敏C反应蛋白是什么？"
	fmt.Print("Before:")
	print(jiebago.Default.Cut(sentence, false))
	word := "超敏C反应蛋白"
	oldFrequency, _ := jiebago.Default.Frequency(word)
	frequency := jiebago.Default.SuggestFrequency(word)
	fmt.Printf("%s current frequency: %f, suggest: %f.\n", word, oldFrequency, frequency)
	jiebago.Default.AddWord(word, frequency)
	fmt.Print("After:")
	print(jiebago.Default.Cut(sentence, false))
	// Output:
	// Before: 超敏 / C / 反应 / 蛋白 / 是 / 什么 / ？ /
	// 超敏C反应蛋白 current frequency: 0.000000, suggest: 1.000000.
	// After: 超敏C反应蛋白 / 是 / 什么 / ？ /
}

func Example_loadUserDictionary() {
	// 在默认分词器基础上加载用户自定义词典
	print := func(ch <-chan string) {
		for word := range ch {
			fmt.Printf(" %s /", word)
		}
		fmt.Println()
	}
	sentence := "李小福是创新办主任也是云计算方面的专家"
	fmt.Print("Before:")
	print(jiebago.Default.Cut(sentence, true))

	jiebago.Default.LoadUserDictionary("userdict.txt")

	fmt.Print("After:")
	print(jiebago.Default.Cut(sentence, true))
	// Output:
	// Before: 李小福 / 是 / 创新 / 办 / 主任 / 也 / 是 / 云 / 计算 / 方面 / 的 / 专家 /
	// After: 李小福 / 是 / 创新办 / 主任 / 也 / 是 / 云计算 / 方面 / 的 / 专家 /
}
