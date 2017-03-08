package gora

import (
	"fmt"
	"strings"
)

// List represents a TODO list.
type List struct {
	rootItem *ListItem
}

// ListItem represents TODO list items.
type ListItem struct {
	id    int
	state State
	desc  string
	items []*ListItem
}

// State represents the current state of a list item.
type State int

const (
	// TODO is the state of an item still waiting to be done.
	TODO State = iota
	// DONE is the state of a finished item.
	DONE
)

// NewList creates and initializes a new List.
func NewList(s string) *List {
	return &List{
		&ListItem{
			0,
			TODO,
			s,
			[]*ListItem{},
		},
	}
}

// Parse is a convenience function that constructs a list from bytes
func Parse(bytes []byte) (*List, error) {
	l := NewList("temp")
	err := l.UnmarshalText(bytes)
	if err != nil {
		return nil, fmt.Errorf("Unable to read list from text")
	}
	return l, nil
}

// Add creates a new list item and adds it to the list.
func (l *List) Add(s string) *ListItem {

	li := l.rootItem
	ni := li.Add(s)
	return ni
}

// Add adds a child item to a parent list item.
func (li *ListItem) Add(s string) *ListItem {
	ni := ListItem{
		len(li.items),
		TODO,
		s,
		make([]*ListItem, 0),
	}
	li.items = append(li.items, &ni)
	return &ni
}

// MoveUp is a wrapper for List.rootItem.ModeUp()
func (l *List) MoveUp(i int) {
	li := l.rootItem
	li.MoveUp(i)
}

// MoveUp moves the ListItem at index i up one slot
func (li *ListItem) MoveUp(i int) {
	if i > 0 && i < len(li.items) {
		a, b := li.items[i-1], li.items[i]
		a.id, b.id = b.id, a.id
		li.items[i-1] = b
		li.items[i] = a
	}
}

// MoveDown is a wrapper for List.rootItem.MoveDown()
func (l *List) MoveDown(i int) {
	li := l.rootItem
	li.MoveDown(i)
}

// MoveDown moves the ListItem at index i down one slot
func (li *ListItem) MoveDown(i int) {
	if i >= 0 && i < len(li.items)-1 {
		a, b := li.items[i], li.items[i+1]
		a.id, b.id = b.id, a.id
		li.items[i] = b
		li.items[i+1] = a
	}
}

// Delete is a wrapper for List.rootItem.Delete()
func (l *List) Delete(i int) {
	li := l.rootItem
	li.Delete(i)
}

// Delete deletes the item at index i
func (li *ListItem) Delete(i int) {
	if i >= 0 && i < len(li.items) {
		li.items = append(li.items[:i], li.items[i+1:]...)
	}
}

func recursiveMarshal(s string, li *ListItem, i int) string {
	s += strings.Repeat("  ", i)
	s += "- ["
	switch li.state {
	case DONE:
		s += "X"
	case TODO:
		fallthrough
	default:
		s += " "
	}
	s += "] " + strings.TrimSpace(li.desc) + "\n"
	for _, li := range li.items {
		s = recursiveMarshal(s, li, i+1)
	}
	return s
}

// MarshalText is used to build a textual representation from a list object
func (l *List) MarshalText() (text []byte, err error) {
	heading := strings.TrimSpace(l.rootItem.desc)
	s := heading + "\n"
	s += strings.Repeat("=", len(heading)) + "\n"
	for _, li := range l.rootItem.items {
		s = recursiveMarshal(s, li, 0)
	}
	return []byte(s), nil
}

// UnmarshalText is used to build a list object from its textual representation
func (l *List) UnmarshalText(text []byte) error {
	indentString := "  "
	lines := strings.Split(string(text), "\n")
	if lines[1] != strings.Repeat("=", len(lines[0])) {
		return fmt.Errorf("At line %d: The heading and its underline do not have the same length", 2)
	}
	id := 0
	l.rootItem = &ListItem{id, TODO, lines[0], nil}
	parentStack := []*ListItem{l.rootItem}
	itemLines := lines[2:]
	for i, line := range itemLines {
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}
		lineNbr := i + 3
		itemLevel := len(parentStack)
		expectedIndent := strings.Repeat(indentString, itemLevel-1)
		actualIndent := line[0:strings.Index(line, "-")]
		if actualIndent != expectedIndent {
			return fmt.Errorf("At line %d: Expected indent width %d, but found %d", lineNbr, len(expectedIndent), len(actualIndent))
		}
		lineSlice := line
		if len(actualIndent) > 0 {
			lineSlice = lineSlice[len(actualIndent):]
		}
		if lineSlice[2:3] != "[" || lineSlice[4:5] != "]" {
			return fmt.Errorf("At line %d: Expected matching brackets at position 2 and 4, but found '%v' and '%v'", lineNbr, lineSlice[2:3], lineSlice[4:5])
		}
		stateRune := lineSlice[3]
		state := TODO
		switch stateRune {
		case ' ':
			state = TODO
		case 'X':
			state = DONE
		default:
			return fmt.Errorf("At line  %d: Expected task state to be one of (' ', 'X'). Got '%v'", lineNbr, lineSlice[1])
		}
		desc := strings.TrimSpace(lineSlice[6:])
		id++
		li := &ListItem{
			id,
			state,
			desc,
			nil,
		}
		parent := parentStack[len(parentStack)-1]
		parent.items = append(parent.items, li)
		if len(itemLines) > i+1 {
			indentIndex := strings.Index(itemLines[i+1], "-")
			if indentIndex < 0 {
				indentIndex = 0
			}
			nextIndent := itemLines[i+1][:indentIndex]
			switch {
			case len(nextIndent) < len(actualIndent):
				parentStack = parentStack[:len(parentStack)-1]
			case len(nextIndent) > len(actualIndent):
				parentStack = append(parentStack, li)
			}
		}
	}

	return nil
}
