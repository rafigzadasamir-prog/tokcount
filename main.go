// Command tokcount estimates the token count and USD cost of a prompt
// without calling any external API or tokenizer library.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
)

func main() {
	model := flag.String("model", defaultModel, "model name to price against")
	outputTokens := flag.Int("output", 0, "assumed output tokens, for cost estimation")
	listModels := flag.Bool("list-models", false, "print known model names and exit")
	flag.Parse()

	if *listModels {
		names := SupportedModels()
		sort.Strings(names)
		for _, n := range names {
			fmt.Println(n)
		}
		return
	}

	text, err := readInput(flag.Args())
	if err != nil {
		fmt.Fprintln(os.Stderr, "tokcount:", err)
		os.Exit(1)
	}

	inputTokens := EstimateTokens(text)
	pricing, resolvedModel := LookupPricing(*model)

	inputCost := EstimateCost(inputTokens, pricing.InputPer1M)
	outputCost := EstimateCost(*outputTokens, pricing.OutputPer1M)

	fmt.Printf("model:          %s\n", resolvedModel)
	fmt.Printf("input tokens:   %d (estimated)\n", inputTokens)
	fmt.Printf("input cost:     %s\n", formatUSD(inputCost))
	if *outputTokens > 0 {
		fmt.Printf("output tokens:  %d (assumed)\n", *outputTokens)
		fmt.Printf("output cost:    %s\n", formatUSD(outputCost))
		fmt.Printf("total cost:     %s\n", formatUSD(inputCost+outputCost))
	}
}

// readInput reads text from a file path given as the first positional
// argument, or from stdin when no argument is given.
func readInput(args []string) (string, error) {
	if len(args) > 0 {
		data, err := os.ReadFile(args[0])
		if err != nil {
			return "", fmt.Errorf("reading %s: %w", args[0], err)
		}
		return string(data), nil
	}

	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", fmt.Errorf("reading stdin: %w", err)
	}
	return string(data), nil
}
