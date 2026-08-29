package nera

import (
	"encoding/csv"
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestParseValidDocument(t *testing.T) {
	input := `Sender
Sender name

InvoiceNumber, Date
123, 12/12/2020

Line item, Quantity, Unit price, Total price
1, 10, 100, 1000
2, 20, 200, 4000
3, 30, 300, 9000
4, 40, 400, 16000
`
	doc, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(doc.Entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(doc.Entries))
	}

	want0 := Literal{Key: "Sender", Value: "Sender name"}
	if !reflect.DeepEqual(doc.Entries[0], want0) {
		t.Errorf("entry 0: got %#v, want %#v", doc.Entries[0], want0)
	}

	want1 := LiteralGroup{Keys: []string{"InvoiceNumber", "Date"}, Values: []string{"123", "12/12/2020"}}
	if !reflect.DeepEqual(doc.Entries[1], want1) {
		t.Errorf("entry 1: got %#v, want %#v", doc.Entries[1], want1)
	}

	e2, ok := doc.Entries[2].(LiteralGroupCollection)
	if !ok || len(e2.Values) != 4 {
		t.Errorf("entry 2: expected LiteralGroupCollection with 4 rows, got %#v", doc.Entries[2])
	}
}

func TestParseNoValueRow(t *testing.T) {
	input := `Sender
Sender name

InvoiceNumber, Date

123, 12/12/2020
`
	_, err := Parse(input)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if !errors.Is(err, ErrEmptyValueRow) {
		t.Errorf("expected ErrEmptyValueRow, got %v", err)
	}

	if !strings.Contains(err.Error(), "line 4") {
		t.Errorf("expected error to mention line 4, got %v", err)
	}
}

func TestParseValueRowShortColumn(t *testing.T) {
	input := `Sender
Sender name

Line item, Quantity, Unit price, Total price
1, 10, 100
`
	_, err := Parse(input)
	if !errors.Is(err, csv.ErrFieldCount) {
		t.Errorf("expected csv.ErrFieldCount, got %v", err)
	}
}

func TestParseKeyRowShortColumn(t *testing.T) {
	input := `Sender
Sender name

Line item, Quantity, Unit price
1, 10, 100, 1000
`
	_, err := Parse(input)
	if !errors.Is(err, csv.ErrFieldCount) {
		t.Errorf("expected csv.ErrFieldCount, got %v", err)
	}
}

func testDoc(t *testing.T) *Document {
	t.Helper()
	input := `Sender
Sender name

InvoiceNumber, Date
123, 12/12/2020

Line item, Quantity, Unit price, Total price
1, 10, 100, 1000
2, 20, 200, 4000
3, 30, 300, 9000
4, 40, 400, 16000
`
	doc, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error setting up fixture: %v", err)
	}
	return doc
}

func TestDocumentAt(t *testing.T) {
	doc := testDoc(t)

	got, err := doc.At(0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := Literal{Key: "Sender", Value: "Sender name"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %#v, want %#v", got, want)
	}
}

func TestDocumentAtGroup(t *testing.T) {
	doc := testDoc(t)

	got, err := doc.AtGroup(1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := LiteralGroup{Keys: []string{"InvoiceNumber", "Date"}, Values: []string{"123", "12/12/2020"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %#v, want %#v", got, want)
	}
}

func TestDocumentAtGroupCollection(t *testing.T) {
	doc := testDoc(t)

	got, err := doc.AtGroupCollection(2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got.Values) != 4 {
		t.Errorf("expected 4 rows, got %d", len(got.Values))
	}
	if got.Values[1][1] != "20" { // second line, Quantity column
		t.Errorf("Values[1][1]: got %q, want %q", got.Values[1][1], "20")
	}
}

func TestDocumentAtCollection(t *testing.T) {
	// testDoc has no LiteralCollection entry, so build a dedicated fixture.
	input := `Tags
first
second
third
`
	doc, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error setting up fixture: %v", err)
	}

	got, err := doc.AtCollection(0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := LiteralCollection{Key: "Tags", Values: []string{"first", "second", "third"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %#v, want %#v", got, want)
	}
}

func TestDocumentAccessorErrors(t *testing.T) {
	doc := testDoc(t)

	cases := []struct {
		name    string
		call    func() error
		wantErr error
	}{
		{
			name:    "At index out of range",
			call:    func() error { _, err := doc.At(99); return err },
			wantErr: ErrIndexOutOfRange,
		},
		{
			name:    "At negative index",
			call:    func() error { _, err := doc.At(-1); return err },
			wantErr: ErrIndexOutOfRange,
		},
		{
			name:    "At type mismatch — actually a LiteralGroup",
			call:    func() error { _, err := doc.At(1); return err },
			wantErr: ErrTypeMismatch,
		},
		{
			name:    "AtGroup type mismatch — actually a Literal",
			call:    func() error { _, err := doc.AtGroup(0); return err },
			wantErr: ErrTypeMismatch,
		},
		{
			name:    "AtCollection type mismatch — actually a LiteralGroupCollection",
			call:    func() error { _, err := doc.AtCollection(2); return err },
			wantErr: ErrTypeMismatch,
		},
		{
			name:    "AtGroupCollection type mismatch — actually a LiteralGroup",
			call:    func() error { _, err := doc.AtGroupCollection(1); return err },
			wantErr: ErrTypeMismatch,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.call()
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("got %v, want error wrapping %v", err, tc.wantErr)
			}
		})
	}
}
