package main

import "testing"

func TestEstimateTokens(t *testing.T) {
	cases := []struct {
		name string
		in   string
		min  int
		max  int
	}{
		{"empty", "", 0, 0},
		{"single word", "hello", 1, 3},
		{"short sentence", "The quick brown fox jumps over the lazy dog.", 8, 16},
		{"whitespace only", "   \n\t  ", 0, 0},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := EstimateTokens(c.in)
			if got < c.min || got > c.max {
				t.Errorf("EstimateTokens(%q) = %d, want between %d and %d", c.in, got, c.min, c.max)
			}
		})
	}
}

func TestEstimateTokensMonotonic(t *testing.T) {
	short := EstimateTokens("hello world")
	long := EstimateTokens("hello world, this is a much longer sentence with many more words in it")
	if long <= short {
		t.Errorf("expected longer text to have a higher token estimate: short=%d long=%d", short, long)
	}
}

func TestLookupPricingKnownModel(t *testing.T) {
	p, resolved := LookupPricing("gpt-4o-mini")
	if resolved != "gpt-4o-mini" {
		t.Errorf("resolved = %q, want gpt-4o-mini", resolved)
	}
	if p.InputPer1M <= 0 || p.OutputPer1M <= 0 {
		t.Errorf("expected positive pricing, got %+v", p)
	}
}

func TestLookupPricingUnknownModelFallsBack(t *testing.T) {
	_, resolved := LookupPricing("not-a-real-model")
	if resolved != defaultModel {
		t.Errorf("resolved = %q, want fallback %q", resolved, defaultModel)
	}
}

func TestEstimateCost(t *testing.T) {
	cost := EstimateCost(1_000_000, 3.0)
	if cost != 3.0 {
		t.Errorf("EstimateCost(1_000_000, 3.0) = %v, want 3.0", cost)
	}

	zero := EstimateCost(0, 3.0)
	if zero != 0 {
		t.Errorf("EstimateCost(0, 3.0) = %v, want 0", zero)
	}
}

func TestSupportedModelsNonEmpty(t *testing.T) {
	names := SupportedModels()
	if len(names) == 0 {
		t.Fatal("expected at least one supported model")
	}
	found := false
	for _, n := range names {
		if n == defaultModel {
			found = true
		}
	}
	if !found {
		t.Errorf("expected defaultModel %q to be in supported list %v", defaultModel, names)
	}
}
