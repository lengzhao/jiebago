package tokenizers_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/blevesearch/bleve/v2"
	_ "github.com/lengzhao/jiebago/tokenizers"
)

// BleveSearchExample demonstrates using jiebago tokenizer with bleve search engine.
// Note: This test is skipped due to bleve v1/v2 API compatibility issues.
// To run this test, you need to use bleve v1.0.14 or update the code for bleve v2 API.
func Skip_BleveSearchExample(t *testing.T) {
	// open a new index
	indexMapping := bleve.NewIndexMapping()

	err := indexMapping.AddCustomTokenizer("jieba",
		map[string]interface{}{
			"file": "../dict.txt",
			"type": "jieba",
		})
	if err != nil {
		log.Fatal(err)
	}

	// create a custom analyzer
	err = indexMapping.AddCustomAnalyzer("jieba",
		map[string]interface{}{
			"type":      "custom",
			"tokenizer": "jieba",
			"token_filters": []string{
				"possessive_en",
				"to_lower",
				"stop_en",
			},
		})

	if err != nil {
		log.Fatal(err)
	}

	indexMapping.DefaultAnalyzer = "jieba"
	cacheDir := "jieba.bleve"
	os.RemoveAll(cacheDir)
	index, err := bleve.New(cacheDir, indexMapping)

	if err != nil {
		log.Fatal(err)
	}

	docs := []struct {
		Title string
		Name  string
	}{
		{
			Title: "Doc 1",
			Name:  "This is the first document we've added",
		},
		{
			Title: "Doc 2",
			Name:  "The second one 你 中文测试中文 is even more interesting! 吃水果",
		},
		{
			Title: "Doc 3",
			Name:  "买水果然后来世博园。",
		},
		{
			Title: "Doc 4",
			Name:  "工信处女干事每月经过下属科室都要亲口交代24口交换机等技术性器件的安装工作",
		},
		{
			Title: "Doc 5",
			Name:  "咱俩交换一下吧。",
		},
	}
	// index docs
	for _, doc := range docs {
		index.Index(doc.Title, doc)
	}

	// search for some text
	for _, keyword := range []string{"水果世博园", "你", "first", "中文", "交换机", "交换"} {
		query := bleve.NewQueryStringQuery(keyword)
		search := bleve.NewSearchRequest(query)
		search.Highlight = bleve.NewHighlight()
		searchResults, err := index.Search(search)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Result of \"%s\": %d matches:\n", keyword, searchResults.Total)
		for i, hit := range searchResults.Hits {
			rv := fmt.Sprintf("%d. %s, (%f)\n", i+searchResults.Request.From+1, hit.ID, hit.Score)
			for fragmentField, fragments := range hit.Fragments {
				rv += fmt.Sprintf("%s: ", fragmentField)
				for _, fragment := range fragments {
					rv += fmt.Sprintf("%s", fragment)
				}
			}
			fmt.Printf("%s\n", rv)
		}
	}
}
