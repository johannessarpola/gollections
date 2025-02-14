package linkedlist

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestLinkedList1(t *testing.T) {
	l := NewLinkedList[int]()

	if !l.IsEmpty() {
		t.Errorf("list is not empty")
	}

	l.Append(1)
	l.Append(2)
	l.Append(77)

	l.Prepend(99)

	if l.IsEmpty() {
		t.Errorf("list is empty")
	}

	if v, ok := l.GetLast(); v != 77 || !ok {
		t.Errorf("expected last to be %d and %t but got %d and %t", 3, true, v, ok)
	}

	if v, ok := l.GetFirst(); v != 99 || !ok {
		t.Errorf("expected head to be %d and %t but got %d and %t", 1, true, v, ok)
	}

	if l.Size() != 4 {
		t.Errorf("expected size to be %d got %d", 4, l.Size())
	}

	if v, err := l.GetAt(0); v != 99 || err != nil {
		t.Errorf("expected element at 0 to be %d but got %d with error %v", 99, v, err)
	}

	if _, err := l.GetAt(10); err == nil {
		t.Error("expected error but got nil")
	}

	e, err := l.RemoveAt(0)
	if err != nil || e == 0 {
		t.Errorf("could not remove at idx %d", 0)
	}

	c := l.Contains(99)
	if c == true {
		t.Errorf("expected removed element to not be contained but got %t", c)
	}

	e2, err := l.RemoveAt(10)
	if err == nil || e2 != 0 {
		t.Errorf("expected error but got %s", err)
	}

	count := 0
	size := l.Size()
	for range l.All {
		count++
	}

	if count != size {
		t.Errorf("expected %d elements but got %d", size, count)
	}

	c2 := l.Contains(-99)
	if c2 == true {
		t.Errorf("expected element to be not contained but got %t", c2)
	}

	err = l.InsertAt(0, -66)
	if err != nil {
		t.Error("expected no error but got ", err)
	}

	if v, ok := l.GetFirst(); v != -66 || !ok {
		t.Errorf("expected head element to be %d and %t but got %d and %t", -66, true, v, ok)
	}

	c3 := l.Contains(-66)
	if c3 != true {
		t.Errorf("expected element to be contained but got %t", c3)
	}

	if v, ok := l.GetFirst(); v != -66 || !ok {
		t.Errorf("expected head to be %d but got %d and %t", -66, v, ok)
	}

	err = l.InsertAt(3, 55)
	if err != nil {
		t.Error("expected no error but got ", err)
	}

	if v, err := l.GetAt(3); v != 55 || err != nil {
		t.Errorf("expected %d to be %d but got %d and %s", 3, 55, v, err)
	}

	err = l.InsertAt(99, 33)
	if err == nil {
		t.Error("expected error but got nil")
	}

	l.Prepend(22)
	if v, b := l.RemoveFirst(); v != 22 || !b {
		t.Errorf("expected element to be %d but got %d and %t", 22, v, b)
	}

	if v, b := l.GetFirst(); v == 22 || !b {
		t.Errorf("expected element to not be %d but got %d and %t", 22, v, b)
	}

	l.Append(42)
	if v, b := l.RemoveLast(); v != 42 || !b {
		t.Errorf("expected element to be %d but got %d and %t", 42, v, b)
	}
	if v, _ := l.GetLast(); v == 42 {
		t.Errorf("expected element to not be %d but got %d", 42, v)
	}
	if l.Contains(42) {
		t.Errorf("expected element to be not contained after RemoveLast()")
	}

	if err = l.InsertAt(1, 31); err != nil {
		t.Error("expected no error but got ", err)
	}

	if i := l.IndexOf(31); i != 1 {
		t.Errorf("expected index of 31 to be %d but got %d", 1, i)
	}

	l.Prepend(67)
	if i := l.IndexOf(67); i != 0 {
		t.Errorf("expected index of 67 to be %d but got %d", 0, i)
	}

	if i := l.IndexOf(-999); i != -1 {
		t.Errorf("expected index of -999 to be -1 but got %d", i)
	}

	if ok := l.Remove(67); !ok {
		t.Errorf("expected element to be removed but got %t", ok)
	}
	if contained := l.Contains(67); contained {
		t.Errorf("expected element to be not contained but got %t", contained)
	}

	if err = l.InsertAt(2, 45); err != nil {
		t.Error("expected no error but got ", err)
	}
	if ok := l.Remove(45); !ok {
		t.Errorf("expected element to be removed but got %t", ok)
	}
	if contained := l.Contains(45); contained {
		t.Errorf("expected element to be not contained but got %t", contained)
	}

	fmt.Printf("List: %s", l.String())

	l.Clear()
	if !l.IsEmpty() {
		t.Errorf("expected list to be empty but got %t", l.IsEmpty())
	}
	if _, b := l.GetFirst(); b != false {
		t.Errorf("expected list to be empty but got %t", b)
	}

	l2 := NewLinkedList[string]()
	l2.Append("abc")
	sz2 := l2.Size()

	additions := []string{"xxx", "yyy", "zzz"}
	l2.AddAll(additions...)
	if l2.Size() != sz2+len(additions) {
		t.Errorf("expected len to be %d but got %d", sz2+len(additions), l2.Size())
	}
	for _, v := range additions {
		if !l2.Contains(v) {
			t.Errorf("expected element to be contained but got %t", l2.Contains(v))
		}
	}

	if !l2.Contains("abc") {
		t.Errorf("expected original element to be contained but got %t", l2.Contains("abc"))
	}
}

func TestLinkedList_UnmarshalJSON(t *testing.T) {
	jsonStr := `
	[
		"abc", 
		"cba", 
		"bca"
	]`

	ll := NewLinkedList[string]()
	err := json.Unmarshal([]byte(jsonStr), &ll)
	if err != nil {
		t.Errorf("expected %s to unmarshal json", jsonStr)
	}

	if !(ll.Contains("abc") && ll.Contains("cba") && ll.Contains("bca")) {
		t.Errorf("expected strings to be contained")
	}

	if ll.Contains("") {
		t.Errorf("contained string that should not be there")
	}

	if ll.Size() != 3 {
		t.Errorf("expected len to be 3 but got %d", ll.Size())
	}

	var ll2 LinkedList[string]
	jb, err := json.Marshal(&ll)

	if err != nil || len(jb) == 0 {
		t.Errorf("expected %s to marshal json", jsonStr)
	}

	if json.Unmarshal(jb, &ll2) != nil || !reflect.DeepEqual(ll2.Items(), ll.Items()) {
		t.Errorf("expected %v, got %v", ll.Items(), ll2.Items())
	}

	type contained struct {
		Field string             `json:"field"`
		List  LinkedList[string] `json:"list"`
	}

	jsonStr2 := `
{
	"field" : "value",
	"list" : [
		"abc", 
		"cba", 
		"bca"
	]
}`

	var c contained
	err = json.Unmarshal([]byte(jsonStr2), &c)
	if err != nil {
		t.Errorf("expected %s to unmarshal json", jsonStr2)
	}
	if c.Field != "value" {
		t.Errorf("expected Field to be %s but got %s", c.Field, "value")
	}
	if c.List.Size() != 3 {
		t.Errorf("expected size to be %d ", 3)
	}

	if !(c.List.Contains("abc") && c.List.Contains("cba") && c.List.Contains("bca")) {
		t.Errorf("expected strings to be contained")
	}
}
