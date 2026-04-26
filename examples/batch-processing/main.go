// Batch processing example demonstrates parallel text segmentation for large datasets.
package main

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/lengzhao/jiebago"
)

var seg jiebago.Segmenter

func init() {
	if err := seg.LoadDictionary("../../embed/dict.txt"); err != nil {
		panic(err)
	}
}

type Task struct {
	ID   int
	Text string
}

type Result struct {
	ID    int
	Words []string
	Err   error
}

func worker(id int, tasks <-chan Task, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for task := range tasks {
		var words []string
		for word := range seg.Cut(task.Text, true) {
			words = append(words, word)
		}
		results <- Result{ID: task.ID, Words: words}
	}
}

func processBatch(lines []string, numWorkers int) []Result {
	tasks := make(chan Task, len(lines))
	results := make(chan Result, len(lines))

	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(i, tasks, results, &wg)
	}

	// Send tasks
	for i, line := range lines {
		tasks <- Task{ID: i, Text: line}
	}
	close(tasks)

	// Wait and collect results
	wg.Wait()
	close(results)

	var allResults []Result
	for r := range results {
		allResults = append(allResults, r)
	}

	return allResults
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Simulate batch data
	lines := []string{
		"这是一个伸手不见五指的黑夜",
		"小明硕士毕业于中国科学院计算所",
		"他来到了网易杭研大厦",
		"长春市长春节讲话",
		"南京市长江大桥",
		"乒乓球拍卖完了",
		"中华人民共和国",
		"中国人民银行",
	}

	fmt.Printf("Processing %d lines with %d workers...\n\n", len(lines), runtime.NumCPU())

	start := time.Now()
	results := processBatch(lines, runtime.NumCPU())
	duration := time.Since(start)

	// Output results
	for _, r := range results {
		fmt.Printf("Line %d: %s\n", r.ID, strings.Join(r.Words, " / "))
	}

	fmt.Printf("\nTotal time: %v\n", duration)
	fmt.Printf("Speed: %.2f lines/second\n", float64(len(lines))/duration.Seconds())
}
