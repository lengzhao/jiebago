package jiebago

import (
	"testing"

	"github.com/lengzhao/jiebago/embed"
)

// BenchmarkDefaultSegmenter 测试默认分词器的性能
func BenchmarkDefaultSegmenter(b *testing.B) {
	sentence := "工信处女干事每月经过下属科室都要亲口交代24口交换机等技术性器件的安装工作"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range Default.Cut(sentence, true) {
		}
	}
}

// BenchmarkDefaultSegmenterNoHMM 测试默认分词器（无HMM）的性能
func BenchmarkDefaultSegmenterNoHMM(b *testing.B) {
	sentence := "工信处女干事每月经过下属科室都要亲口交代24口交换机等技术性器件的安装工作"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range Default.Cut(sentence, false) {
		}
	}
}

// BenchmarkDefaultSegmenterCutAll 测试默认分词器全模式
func BenchmarkDefaultSegmenterCutAll(b *testing.B) {
	sentence := "工信处女干事每月经过下属科室都要亲口交代24口交换机等技术性器件的安装工作"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range Default.CutAll(sentence) {
		}
	}
}

// BenchmarkCustomSegmenter 测试自定义分词器的性能
func BenchmarkCustomSegmenter(b *testing.B) {
	var seg Segmenter
	seg.LoadDictionaryFromBytes(embed.DictData)
	sentence := "工信处女干事每月经过下属科室都要亲口交代24口交换机等技术性器件的安装工作"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range seg.Cut(sentence, true) {
		}
	}
}

// BenchmarkCustomSegmenterNoHMM 测试自定义分词器（无HMM）的性能
func BenchmarkCustomSegmenterNoHMM(b *testing.B) {
	var seg Segmenter
	seg.LoadDictionaryFromBytes(embed.DictData)
	sentence := "工信处女干事每月经过下属科室都要亲口交代24口交换机等技术性器件的安装工作"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range seg.Cut(sentence, false) {
		}
	}
}

// BenchmarkLoadDictionaryFromBytes 测试从字节加载字典的性能
func BenchmarkLoadDictionaryFromBytes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var seg Segmenter
		seg.LoadDictionaryFromBytes(embed.DictData)
	}
}

// BenchmarkLoadDictionaryFromFile 测试从文件加载字典的性能
func BenchmarkLoadDictionaryFromFile(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var seg Segmenter
		seg.LoadDictionary("embed/dict.txt")
	}
}

// BenchmarkDefaultSegmenterInit 测试默认分词器延迟初始化的开销
func BenchmarkDefaultSegmenterInit(b *testing.B) {
	// 创建新的默认分词器实例来测试初始化开销
	for i := 0; i < b.N; i++ {
		d := &defaultSegmenter{}
		d.init()
	}
}

// BenchmarkSegmentLongText 测试长文本分词性能
func BenchmarkSegmentLongText(b *testing.B) {
	text := "自然语言处理是人工智能领域中的一个重要方向。它研究能实现人与计算机之间用自然语言进行有效通信的各种理论和方法。自然语言处理是一门融语言学、计算机科学、数学于一体的科学。因此，这一领域的研究将涉及自然语言，即人们日常使用的语言，所以它与语言学的研究有着密切的联系，但又有重要的区别。自然语言处理并不是一般地研究自然语言，而在于研制能有效地实现自然语言通信的计算机系统，特别是其中的软件系统。因而它是计算机科学的一部分。"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range Default.Cut(text, true) {
		}
	}
}

// BenchmarkSegmentShortText 测试短文本分词性能
func BenchmarkSegmentShortText(b *testing.B) {
	text := "我爱北京天安门"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range Default.Cut(text, true) {
		}
	}
}

// BenchmarkFrequencyLookup 测试词频查询性能
func BenchmarkFrequencyLookup(b *testing.B) {
	words := []string{"北京", "清华大学", "自然语言处理", "nonexistent"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, word := range words {
			Default.Frequency(word)
		}
	}
}

// BenchmarkCutWithNewWordRecognition 测试新词识别性能
func BenchmarkCutWithNewWordRecognition(b *testing.B) {
	// 包含新词的句子
	text := "他来到了网易杭研大厦"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range Default.Cut(text, true) {
		}
	}
}

// BenchmarkCutForSearchWithHMM 测试搜索引擎模式性能
func BenchmarkCutForSearchWithHMM(b *testing.B) {
	sentence := "小明硕士毕业于中国科学院计算所，后在日本京都大学深造"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range Default.CutForSearch(sentence, true) {
		}
	}
}

// BenchmarkParallelSegmentation 测试并行分词性能
func BenchmarkParallelSegmentation(b *testing.B) {
	sentences := []string{
		"我来到北京清华大学",
		"他来到了网易杭研大厦",
		"工信处女干事每月经过下属科室都要亲口交代24口交换机等技术性器件的安装工作",
		"自然语言处理是人工智能领域中的一个重要方向",
		"小明硕士毕业于中国科学院计算所",
	}
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			sentence := sentences[i%len(sentences)]
			for range Default.Cut(sentence, true) {
			}
			i++
		}
	})
}
