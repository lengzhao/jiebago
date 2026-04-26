package analyse

import (
	"cmp"
	"slices"

	"github.com/lengzhao/jiebago/posseg"
)

const dampingFactor = 0.85

var (
	defaultAllowPOS = []string{"ns", "n", "vn", "v"}
)

type edge struct {
	start  string
	end    string
	weight float64
}

type undirectWeightedGraph struct {
	graph map[string][]edge
	keys  []string
}

func newUndirectWeightedGraph() *undirectWeightedGraph {
	return &undirectWeightedGraph{
		graph: make(map[string][]edge),
		keys:  make([]string, 0),
	}
}

func (u *undirectWeightedGraph) addEdge(start, end string, weight float64) {
	if _, ok := u.graph[start]; !ok {
		u.keys = append(u.keys, start)
		u.graph[start] = []edge{{start: start, end: end, weight: weight}}
	} else {
		u.graph[start] = append(u.graph[start], edge{start: start, end: end, weight: weight})
	}

	if _, ok := u.graph[end]; !ok {
		u.keys = append(u.keys, end)
		u.graph[end] = []edge{{start: end, end: start, weight: weight}}
	} else {
		u.graph[end] = append(u.graph[end], edge{start: end, end: start, weight: weight})
	}
}

func (u *undirectWeightedGraph) rank() Segments {
	slices.Sort(u.keys)

	ws := make(map[string]float64, len(u.graph))
	outSum := make(map[string]float64, len(u.graph))

	wsdef := 1.0
	if len(u.graph) > 0 {
		wsdef /= float64(len(u.graph))
	}
	for n, out := range u.graph {
		ws[n] = wsdef
		sum := 0.0
		for _, e := range out {
			sum += e.weight
		}
		outSum[n] = sum
	}

	for x := 0; x < 10; x++ {
		for _, n := range u.keys {
			s := 0.0
			inedges := u.graph[n]
			for _, e := range inedges {
				s += e.weight / outSum[e.end] * ws[e.end]
			}
			ws[n] = (1 - dampingFactor) + dampingFactor*s
		}
	}

	// Find min and max using Go 1.21+ built-in functions
	minRank := 1.0
	maxRank := 0.0
	for _, w := range ws {
		minRank = min(minRank, w)
		maxRank = max(maxRank, w)
	}

	result := make(Segments, 0, len(ws))
	for n, w := range ws {
		result = append(result, Segment{text: n, weight: (w - minRank/10.0) / (maxRank - minRank/10.0)})
	}

	// Use slices.SortFunc instead of sort.Interface
	slices.SortFunc(result, func(a, b Segment) int {
		return cmp.Compare(b.weight, a.weight) // Reverse order (descending)
	})
	return result
}

// TextRankWithPOS extracts keywords from sentence using TextRank algorithm.
// Parameter allowPOS allows a customized pos list.
func (t *TextRanker) TextRankWithPOS(sentence string, topK int, allowPOS []string) Segments {
	// Use map[string]struct{} as set (zero memory overhead)
	posFilt := make(map[string]struct{}, len(allowPOS))
	for _, pos := range allowPOS {
		posFilt[pos] = struct{}{}
	}

	g := newUndirectWeightedGraph()
	cm := make(map[[2]string]float64)
	span := 5

	// Pre-allocate slice with estimated capacity
	pairs := make([]posseg.Segment, 0, 64)
	for pair := range t.seg.Cut(sentence, true) {
		pairs = append(pairs, pair)
	}

	for i := range pairs {
		if _, ok := posFilt[pairs[i].Pos()]; ok {
			for j := i + 1; j < i+span && j < len(pairs); j++ {
				if _, ok := posFilt[pairs[j].Pos()]; !ok {
					continue
				}
				key := [2]string{pairs[i].Text(), pairs[j].Text()}
				cm[key] += 1.0
			}
		}
	}

	for startEnd, weight := range cm {
		g.addEdge(startEnd[0], startEnd[1], weight)
	}

	tags := g.rank()
	if topK > 0 && len(tags) > topK {
		tags = tags[:topK]
	}
	return tags
}

// TextRank extract keywords from sentence using TextRank algorithm.
// Parameter topK specify how many top keywords to be returned at most.
func (t *TextRanker) TextRank(sentence string, topK int) Segments {
	return t.TextRankWithPOS(sentence, topK, defaultAllowPOS)
}

// TextRanker is used to extract tags from sentence.
type TextRanker struct {
	seg *posseg.Segmenter
}

// LoadDictionary reads a given file and create a new dictionary file for Textranker.
func (t *TextRanker) LoadDictionary(fileName string) error {
	t.seg = new(posseg.Segmenter)
	return t.seg.LoadDictionary(fileName)
}
