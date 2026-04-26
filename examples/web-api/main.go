// Web API example shows how to build a simple HTTP API for Chinese text segmentation.
package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lengzhao/jiebago"
)

var seg jiebago.Segmenter

func init() {
	if err := seg.LoadDictionary("../../embed/dict.txt"); err != nil {
		panic(err)
	}
}

type SegmentationRequest struct {
	Text   string `json:"text"`
	Mode   string `json:"mode"` // "accurate", "full", "search"
	UseHMM bool   `json:"use_hmm"`
}

type SegmentationResponse struct {
	Words []string `json:"words"`
}

func segmentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SegmentationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var words []string
	switch req.Mode {
	case "full":
		for word := range seg.CutAll(req.Text) {
			words = append(words, word)
		}
	case "search":
		for word := range seg.CutForSearch(req.Text, req.UseHMM) {
			words = append(words, word)
		}
	default: // accurate
		for word := range seg.Cut(req.Text, req.UseHMM) {
			words = append(words, word)
		}
	}

	resp := SegmentationResponse{Words: words}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	http.HandleFunc("/segment", segmentHandler)
	fmt.Println("Server starting on :8080")
	fmt.Println("Try: curl -X POST http://localhost:8080/segment -H 'Content-Type: application/json' -d '{\"text\":\"我爱北京天安门\",\"mode\":\"accurate\"}'")
	http.ListenAndServe(":8080", nil)
}
