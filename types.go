package nera

import "fmt"

type Literal struct {
	Key   string
	Value string
}

type LiteralCollection struct {
	Key    string
	Values []string
}

type LiteralGroup struct {
	Keys   []string
	Values []string
}

type LiteralGroupCollection struct {
	Keys   []string
	Values [][]string
}

type entry interface{ isEntry() }

func (Literal) isEntry()                {}
func (LiteralCollection) isEntry()      {}
func (LiteralGroup) isEntry()           {}
func (LiteralGroupCollection) isEntry() {}

type Document struct {
	Entries []entry // order preserved == file order == the API
}

// At returns the entry at index i as a Literal, or an error if the
// index is out of range or the entry at that position isn't a Literal.
func (d *Document) At(i int) (Literal, error) {
	e, err := d.entryAt(i)
	if err != nil {
		return Literal{}, err
	}
	v, ok := e.(Literal)
	if !ok {
		return Literal{}, fmt.Errorf("entry %d: expected Literal, got %T: %w", i, e, ErrTypeMismatch)
	}
	return v, nil
}

func (d *Document) AtCollection(i int) (LiteralCollection, error) {
	e, err := d.entryAt(i)
	if err != nil {
		return LiteralCollection{}, err
	}
	v, ok := e.(LiteralCollection)
	if !ok {
		return LiteralCollection{}, fmt.Errorf("entry %d: expected LiteralCollection, got %T: %w", i, e, ErrTypeMismatch)
	}
	return v, nil
}

func (d *Document) AtGroup(i int) (LiteralGroup, error) {
	e, err := d.entryAt(i)
	if err != nil {
		return LiteralGroup{}, err
	}
	v, ok := e.(LiteralGroup)
	if !ok {
		return LiteralGroup{}, fmt.Errorf("entry %d: expected LiteralGroup, got %T: %w", i, e, ErrTypeMismatch)
	}
	return v, nil
}

func (d *Document) AtGroupCollection(i int) (LiteralGroupCollection, error) {
	e, err := d.entryAt(i)
	if err != nil {
		return LiteralGroupCollection{}, err
	}
	v, ok := e.(LiteralGroupCollection)
	if !ok {
		return LiteralGroupCollection{}, fmt.Errorf("entry %d: expected LiteralGroupCollection, got %T: %w", i, e, ErrTypeMismatch)
	}
	return v, nil
}

func (d *Document) entryAt(i int) (entry, error) {
	if i < 0 || i >= len(d.Entries) {
		return nil, fmt.Errorf("index %d out of range (%d entries): %w", i, len(d.Entries), ErrIndexOutOfRange)
	}
	return d.Entries[i], nil
}
