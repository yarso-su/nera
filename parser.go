package nera

import (
	"encoding/csv"
	"fmt"
	"strings"
)

// block pairs raw block text with the 1-based line number it started on,
// so parse errors can point back to the source file.
type block struct {
	text      string
	startLine int
}

func Parse(input string) (*Document, error) {
	blocks := splitBlocks(input)
	doc := &Document{Entries: make([]entry, 0, len(blocks))}

	for _, b := range blocks {
		e, err := parseBlock(b)
		if err != nil {
			return nil, err // parseBlock errors already carry line info
		}
		doc.Entries = append(doc.Entries, e)
	}
	return doc, nil
}

// splitBlocks separates the file into blank-line-delimited chunks,
// recording each block's starting line number (1-based) for error
// reporting. Empty blocks (consecutive blank lines) are discarded.
func splitBlocks(input string) []block {
	lines := strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n")

	var blocks []block
	var current []string
	start := 0

	flush := func() {
		if len(current) > 0 {
			blocks = append(blocks, block{
				text:      strings.Join(current, "\n"),
				startLine: start + 1, // convert 0-based index to 1-based line number
			})
			current = nil
		}
	}

	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			flush()
			continue
		}
		if len(current) == 0 {
			start = i
		}
		current = append(current, line)
	}
	flush()

	return blocks
}

func parseBlock(b block) (entry, error) {
	r := csv.NewReader(strings.NewReader(b.text))
	r.TrimLeadingSpace = true

	records, err := r.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("line %d: csv parse: %w", b.startLine, err)
	}
	if len(records) < 2 {
		return nil, fmt.Errorf("line %d: %w", b.startLine, ErrEmptyValueRow)
	}

	keys := records[0]
	values := records[1:]
	isGroup := len(keys) > 1
	isCollection := len(values) > 1

	switch {
	case !isGroup && !isCollection:
		return Literal{Key: keys[0], Value: values[0][0]}, nil
	case !isGroup && isCollection:
		vs := make([]string, len(values))
		for i, row := range values {
			vs[i] = row[0]
		}
		return LiteralCollection{Key: keys[0], Values: vs}, nil
	case isGroup && !isCollection:
		return LiteralGroup{Keys: keys, Values: values[0]}, nil
	default:
		return LiteralGroupCollection{Keys: keys, Values: values}, nil
	}
}
