// Custom dictionary example shows how to add custom words and use user dictionaries.
package main

import (
	"fmt"

	"github.com/lengzhao/jiebago"
)

func main() {
	var seg jiebago.Segmenter

	// Load main dictionary (or use jiebago.Default for out-of-the-box experience)
	if err := seg.LoadDictionary("../../embed/dict.txt"); err != nil {
		panic(err)
	}

	sentence := "李小福是创新办主任也是云计算方面的专家"

	// Before adding custom words
	fmt.Println("=== Before Custom Dictionary ===")
	fmt.Println("Text:", sentence)
	fmt.Println("\nSegmentation:")
	for word := range seg.Cut(sentence, true) {
		fmt.Printf("%s / ", word)
	}
	fmt.Println()

	// Load user dictionary
	fmt.Println("\n=== Loading User Dictionary ===")
	if err := seg.LoadUserDictionary("../../userdict.txt"); err != nil {
		fmt.Printf("Note: userdict.txt not found, using AddWord instead\n")
	}

	// Or add words programmatically
	seg.AddWord("李小福", 1000.0)  // Add name
	seg.AddWord("创新办", 500.0)   // Add organization
	seg.AddWord("云计算", 800.0)   // Add tech term

	// After adding custom words
	fmt.Println("\n=== After Adding Custom Words ===")
	for word := range seg.Cut(sentence, true) {
		fmt.Printf("%s / ", word)
	}
	fmt.Println()

	// Check word frequency
	freq, ok := seg.Frequency("李小福")
	fmt.Printf("\nFrequency of '李小福': %.0f (exists: %v)\n", freq, ok)

	// Delete a word
	seg.DeleteWord("李小福")
	fmt.Println("\n=== After Deleting '李小福' ===")
	for word := range seg.Cut(sentence, true) {
		fmt.Printf("%s / ", word)
	}
	fmt.Println()

	// Suggest frequency example
	fmt.Println("\n=== Suggest Frequency ===")
	// Suggest frequency for a new word based on its components
	suggestedFreq := seg.SuggestFrequency("创新办")
	fmt.Printf("Suggested frequency for '创新办': %.0f\n", suggestedFreq)

	// Suggest based on component words
	suggestedFreq2 := seg.SuggestFrequency("创新", "办")
	fmt.Printf("Suggested frequency based on components: %.0f\n", suggestedFreq2)
}
