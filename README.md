# tokcount

A small CLI that estimates the token count and USD cost of a prompt,
without calling any external API or tokenizer library.

Token counts are approximate: it uses a chars-per-token / words heuristic
rather than a model-specific tokenizer, so treat the output as a ballpark
figure for sizing prompts and estimating spend, not a billing-accurate
count.

## Install

```
go install github.com/rafigzadasamir-prog/tokcount@latest
```

Or build from source:

```
go build -o tokcount .
```

## Usage

```
# from a file
tokcount prompt.txt

# from stdin
cat prompt.txt | tokcount

# pick a model's pricing table and include an assumed output size
tokcount --model gpt-4o --output 500 prompt.txt

# see which models have known pricing
tokcount --list-models
```

Example output:

```
model:          claude-sonnet-5
input tokens:   128 (estimated)
input cost:     $0.000384
output tokens:  500 (assumed)
output cost:    $0.007500
total cost:     $0.007884
```

## How the estimate works

`EstimateTokens` averages two cheap heuristics — non-whitespace character
count divided by 4, and word count times 1.3 — then floors the result at
the raw word count. It deliberately avoids depending on any
model-specific tokenizer (e.g. tiktoken/BPE tables), so it has zero
external dependencies and works offline for any language or model.

Pricing is a small hardcoded table in [estimate.go](estimate.go); rates
are approximate and will drift from providers' actual pricing over time.

## License

MIT
