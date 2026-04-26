// Package analyse is the Golang implementation of Jieba's analyse module.
package analyse

import (
	"cmp"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/lengzhao/jiebago"
)

// Segment represents a word with weight.
type Segment struct {
	text   string
	weight float64
}

// Text returns the segment's text.
func (s Segment) Text() string {
	return s.text
}

// Weight returns the segment's weight.
func (s Segment) Weight() float64 {
	return s.weight
}

// Segments represents a slice of Segment.
type Segments []Segment

// TagExtracter is used to extract tags from sentence.
type TagExtracter struct {
	seg      *jiebago.Segmenter
	idf      *Idf
	stopWord *StopWord
}

// LoadDictionary reads the given filename and create a new dictionary.
func (t *TagExtracter) LoadDictionary(fileName string) error {
	t.stopWord = NewStopWord()
	t.seg = new(jiebago.Segmenter)
	return t.seg.LoadDictionary(fileName)
}

// LoadIdf reads the given file and create a new Idf dictionary.
func (t *TagExtracter) LoadIdf(fileName string) error {
	t.idf = NewIdf()
	return t.idf.loadDictionary(fileName)
}

// LoadStopWords reads the given file and create a new StopWord dictionary.
func (t *TagExtracter) LoadStopWords(fileName string) error {
	t.stopWord = NewStopWord()
	return t.stopWord.loadDictionary(fileName)
}

// ExtractTags extracts the topK key words from sentence.
func (t *TagExtracter) ExtractTags(sentence string, topK int) (tags Segments) {
	freqMap := make(map[string]float64, 64)

	for w := range t.seg.Cut(sentence, true) {
		w = strings.TrimSpace(w)
		if utf8.RuneCountInString(w) < 2 {
			continue
		}
		if t.stopWord.IsStopWord(w) {
			continue
		}
		freqMap[w] += 1.0
	}

	total := 0.0
	for _, freq := range freqMap {
		total += freq
	}
	for k, v := range freqMap {
		freqMap[k] = v / total
	}

	// Pre-allocate with capacity
	ws := make(Segments, 0, len(freqMap))
	for k, v := range freqMap {
		freq, ok := t.idf.Frequency(k)
		if !ok {
			freq = t.idf.median
		}
		ws = append(ws, Segment{text: k, weight: freq * v})
	}

	// Use slices.SortFunc instead of sort.Interface
	slices.SortFunc(ws, func(a, b Segment) int {
		if a.weight == b.weight {
			return cmp.Compare(b.text, a.text) // Reverse text order
		}
		return cmp.Compare(b.weight, a.weight) // Descending weight
	})

	if len(ws) > topK {
		return ws[:topK]
	}
	return ws
}
