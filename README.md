# nera

A tiny Go library for parsing a positional, line-oriented micro-language into structured, ordered data — designed for text files that get turned into a downstream binary asset, where field *position* is the schema rather than field *names*.

## Why

`nera` exists for a specific shape of problem: you want a human-editable text file with dynamic, arbitrary labels, but the program consuming it only cares about *where* a value sits, not what its label says. Think of it as several small CSV blocks — literals, key/value pairs, and collections of either — stitched into one file with blank lines as separators.

The parser only produces `string` values. Type interpretation (int, float, bool, whatever) is left entirely to the consumer, since the consumer is the one who actually knows what shape each field should be.

## Install

```sh
go get github.com/yarso-su/nera
```

## Format

A `nera` file is a sequence of **blocks**, separated by one or more blank lines. Each block is, internally, a small CSV chunk:

- The **first line** of a block is one or more comma-separated **keys**.
- Every following line, up to the next blank line, is a comma-separated **row of values** — one value per key, in the same order.
- A block with **one key** and **one value row** is a `Literal`.
- A block with **one key** and **multiple value rows** is a `LiteralCollection`.
- A block with **multiple keys** and **one value row** is a `LiteralGroup`.
- A block with **multiple keys** and **multiple value rows** is a `LiteralGroupCollection`.

Leading whitespace around values is trimmed. Values containing a literal comma must be wrapped in double quotes: `"Smith, John & Co."` — the quotes are stripped and are not part of the parsed value. Unquoted commas are always treated as field separators.

A block's header and its first value row **must not** have a blank line between them — a blank line always closes the current block.

### Example

```text
Sender
Sender name

InvoiceNumber, Date
123, 12/12/2020

Line item, Quantity, Unit price, Total price
1, 10, 100, 1000
2, 20, 200, 4000
3, 30, 300, 9000
4, 40, 400, 16000
```

This parses into three entries, in file order:

1. `Literal{Key: "Sender", Value: "Sender name"}`
2. `LiteralGroup{Keys: ["InvoiceNumber", "Date"], Values: ["123", "12/12/2020"]}`
3. `LiteralGroupCollection{Keys: [...], Values: [][]string{...}}` — 4 rows

### File extension

`nera` does not read from the filesystem and has no opinion on file extensions — `Parse` takes a `string`, however you obtained it. As a **pure convention**, files in this format are suggested to use the `.nera` extension (e.g. `receipt-a2.nera`), but this is not enforced or checked anywhere in the library.

## Usage

```go
package main

import (
	"fmt"
	"log"

	"github.com/yarso-su/nera"
)

func main() {
	doc, err := nera.Parse(input)
	if err != nil {
		log.Fatal(err)
	}

	sender, err := doc.At(0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(sender.Value)

	invoice, err := doc.AtGroup(1)
	if err != nil {
		log.Fatal(err)
	}
	date, _ := invoice.Get("Date")
	fmt.Println(date)

	lines, err := doc.AtGroupCollection(2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(lines.Values), "line items")
}
```

`Document.Entries` preserves file order — position is the API, not key names. Use `At`, `AtCollection`, `AtGroup`, or `AtGroupCollection` when you already know the expected shape at a given index; each returns a typed error (`ErrIndexOutOfRange`, `ErrTypeMismatch`) via `errors.Is` if the assumption doesn't hold.

## License

MIT — see [LICENSE](LICENSE).
