package main

import (
	"fmt"
	"strings"
	"unicode"
)

// Pricing holds per-million-token USD rates for a model.
type Pricing struct {
	InputPer1M  float64
	OutputPer1M float64
}

// knownModels is a small, hardcoded pricing table. Rates are approximate
// and meant for rough estimation, not billing-accurate figures.
var knownModels = map[string]Pricing{
	"claude-sonnet-5":    {InputPer1M: 3.00, OutputPer1M: 15.00},
	"claude-haiku-4.5":   {InputPer1M: 0.80, OutputPer1M: 4.00},
	"claude-opus-5":      {InputPer1M: 15.00, OutputPer1M: 75.00},
	"gpt-4o":             {InputPer1M: 2.50, OutputPer1M: 10.00},
	"gpt-4o-mini":        {InputPer1M: 0.15, OutputPer1M: 0.60},
}

// defaultModel is used when the caller does not specify one.
const defaultModel = "claude-sonnet-5"

// EstimateTokens approximates a token count for text using a
// chars-per-token heuristic refined by word count, without depending on
// any model-specific tokenizer. It is intentionally conservative: real
// tokenizers vary per model and this is only meant for ballpark sizing.
func EstimateTokens(text string) int {
	text = strings.TrimSpace(text)
	if text == "" {
		return 0
	}

	runeCount := 0
	for _, r := range text {
		if !unicode.IsSpace(r) {
			runeCount++
		}
	}

	byCharRate := float64(runeCount) / 4.0

	words := strings.Fields(text)
	byWordRate := float64(len(words)) * 1.3

	estimate := (byCharRate + byWordRate) / 2
	if estimate < float64(len(words)) {
		estimate = float64(len(words))
	}

	return int(estimate + 0.5)
}

// LookupPricing returns the pricing for a model name, falling back to
// defaultModel when the name is unknown.
func LookupPricing(model string) (Pricing, string) {
	if p, ok := knownModels[model]; ok {
		return p, model
	}
	return knownModels[defaultModel], defaultModel
}

// EstimateCost returns the USD cost of a number of tokens at the given
// per-million-token rate.
func EstimateCost(tokens int, perMillion float64) float64 {
	return float64(tokens) / 1_000_000 * perMillion
}

// SupportedModels returns the list of model names with known pricing,
// used for the --list-models output.
func SupportedModels() []string {
	names := make([]string, 0, len(knownModels))
	for name := range knownModels {
		names = append(names, name)
	}
	return names
}

func formatUSD(amount float64) string {
	return fmt.Sprintf("$%.6f", amount)
}
