package gora

import (
	"testing"
)

var marshalTextCase = "Root item\n" +
	"=========\n" +
	"- [ ] Item1\n" +
	"  - [ ] Sub1\n" +
	"  - [ ] Sub2\n" +
	"  - [X] Sub3\n" +
	"- [ ] Item2\n" +
	"  - [X] Sub1\n" +
	"    - [X] SubSub1\n"

func TestMarshalText(t *testing.T) {
	l := NewList("Root item")
	li := l.Add("Item1")
	li.Add("Sub1")
	li.Add("Sub2")
	li.Add("Sub3").state = DONE
	li = l.Add("Item2")
	li = li.Add("Sub1")
	li.state = DONE
	li = li.Add("SubSub1")
	li.state = DONE
	b, _ := l.MarshalText()
	if got := string(b); got != marshalTextCase {
		t.Fatalf("Expected:\n>%v<\n\nGot:\n>%v<", marshalTextCase, got)
	}
}

func TestUnMarshalMarshal(t *testing.T) {
	unmarshalled, _ := Parse([]byte(marshalTextCase))
	marshalled, _ := unmarshalled.MarshalText()
	if string(marshalled) != marshalTextCase {
		t.Fatalf("Expected:\n%v\n\nGot:\n%v", marshalTextCase, string(marshalled))
	}
}

func TestUnmarshalText(t *testing.T) {
	var unmarshalCase = struct {
		input string
		want  List
	}{
		"A list\n" +
			"======\n" +
			"- [X] First item\n" +
			"  - [ ] First sub item\n" +
			"  - [ ] Second sub item\n" +
			"- [ ] Second item\n",
		List{
			&ListItem{
				0,
				TODO,
				"A list",
				[]*ListItem{
					&ListItem{
						1,
						DONE,
						"First item",
						[]*ListItem{
							&ListItem{
								2,
								TODO,
								"First sub item",
								nil,
							},
							&ListItem{
								3,
								TODO,
								"Second sub item",
								nil,
							},
						},
					},
					&ListItem{
						4,
						TODO,
						"Second item",
						nil,
					},
				},
			},
		},
	}

	var l List
	err := l.UnmarshalText([]byte(unmarshalCase.input))
	if err != nil {
		t.Fatal(err)
	}
	want := unmarshalCase.want
	if l.rootItem.id != want.rootItem.id {
		t.Fatalf("Expected root item id to be %d, was %d", want.rootItem.id, l.rootItem.id)
	}
	if l.rootItem.desc != want.rootItem.desc {
		t.Fatalf("Exptected root item desc to be '%v', was '%v'", want.rootItem.desc, l.rootItem.desc)
	}
	if len(l.rootItem.items) != len(want.rootItem.items) {
		t.Fatalf("Exptected root item to have %d children, had %d", len(want.rootItem.items), len(l.rootItem.items))
	}
}

func TestAddToList(t *testing.T) {
	l := NewList("A new list")
	li1 := l.Add("item 1")
	if li1.id != 0 {
		t.Fatalf("Expected id 0, got %d", li1.id)
	}
	li2 := l.Add("item 2")
	if li2.id != 1 {
		t.Fatalf("Expected id 1, got %d", li2.id)
	}
	if len(l.rootItem.items) != 2 {
		t.Errorf("%v", l.rootItem.items)
		t.Fatalf("Expected 2 list items, got %d", len(l.rootItem.items))
	}
}

func TestAddToListItem(t *testing.T) {
	l := NewList("A list")
	if l.rootItem.desc != "A list" {
		t.Fatalf("Expected root item description to be 'A list'. Got %v", l.rootItem.desc)
	}
	li1 := l.Add("First level item")
	if len(l.rootItem.items) != 1 {
		t.Fatalf("Expected length 1, got %d", len(l.rootItem.items))
	}
	li1.Add("Second level item")
	if len(li1.items) != 1 {
		t.Fatalf("Expected length 1, got %d", len(li1.items))
	}
	li2 := li1.Add("Another second level item")
	if li2.id != 1 {
		t.Fatalf("Expected id 1, got %d", li2.id)
	}
	if len(li1.items) != 2 {
		t.Fatalf("Expected length 2, got %d", len(li1.items))
	}
}

func TestSetState(t *testing.T) {
	l := NewList("A list")
	li1 := l.Add("First level item")
	li2 := li1.Add("Second level item")
	if li1.state != TODO {
		t.Fatalf("Set state: First test, expected TODO, got %v", li1.state)
	}
	if li2.state != TODO {
		t.Fatalf("Set state: Second test, expected TODO, got %v", li2.state)
	}
	li2.state = DONE
	if li2.state != DONE {
		t.Fatalf("Set state: Second test, expected DONE, got %v", li2.state)
	}
}

func TestMoveUp(t *testing.T) {
	l := NewList("A list")
	l.Add("Item 1")
	l.Add("Item 2")
	l.MoveUp(0)
	if l.rootItem.items[0].desc != "Item 1" || l.rootItem.items[0].id != 0 {
		t.Fatalf("MoveUp: Top item can't be moved up. Expected id=0, desc='Item 1', got id=%d, desc='%v'", l.rootItem.items[0].id, l.rootItem.items[0].desc)
	}
	l.MoveUp(1)
	if l.rootItem.items[0].desc != "Item 2" || l.rootItem.items[0].id != 0 {
		t.Fatalf("MoveUp: Expected id=0, desc='Item 2' to be on top, but got id=%d, desc='%v'", l.rootItem.items[0].id, l.rootItem.items[0].desc)
	}
}

func TestMoveDown(t *testing.T) {
	l := NewList("A list")
	l.Add("Item 1")
	l.Add("Item 2")
	l.MoveDown(1)
	if l.rootItem.items[1].desc != "Item 2" || l.rootItem.items[1].id != 1 {
		t.Fatalf("MoveDown: Top item can't be moved down. Expected id=1, desc='Item 2', got id=%d, desc='%v'", l.rootItem.items[1].id, l.rootItem.items[1].desc)
	}
	l.MoveDown(0)
	if l.rootItem.items[1].desc != "Item 1" || l.rootItem.items[1].id != 1 {
		t.Fatalf("MoveDown: Expected id=1, desc='Item 1' to be on bottom, but got id=%d, desc='%v'", l.rootItem.items[1].id, l.rootItem.items[1].desc)
	}
}

func TestDelete(t *testing.T) {
	l := NewList("This is a list")
	l.Add("item1")
	li := l.Add("item2")
	l.Delete(0)
	if len(l.rootItem.items) != 1 {
		t.Fatalf("Expected list length to be 1, was %d", len(l.rootItem.items))
	}
	if l.rootItem.items[0].desc != "item2" {
		t.Fatalf("Exptected item desc to be 'item2', was %v", l.rootItem.items[0])
	}
	li.Add("subItem1")
	li.Add("subItem2")
	li.Add("subItem3")
	li.Delete(1)
	if len(li.items) != 2 {
		t.Fatalf("Expected list length to be 2, was %d", len(li.items))
	}
}
