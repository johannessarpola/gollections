package gollections

import (
	"encoding/json"
	"testing"
)

func TestSet(t *testing.T) {

	s := NewSet[string]()
	s.Add("abc")
	s.Add("cba")
	s.Add("bca")

	if !s.Contains("abc") {
		t.Errorf("expected %s to be contained", "abc")
	}

	if !s.Contains("cba") {
		t.Errorf("expected %s to be contained", "cba")
	}

	if !s.Contains("bca") {
		t.Errorf("expected %s to be contained", "bca")
	}

	if s.Contains("non") {
		t.Errorf("expected %s to be not contained", "non")
	}

	if s.Size() != 3 {
		t.Errorf("expected %d got %d", 3, s.Size())
	}

	cnt := 0
	for _, _ = range s.All {
		cnt++
	}

	if cnt != s.Size() || cnt != 3 {
		t.Errorf("expected %d got %d", 3, s.Size())
	}

	s.Remove("abc")
	if s.Contains("abc") {
		t.Errorf("expected %s to be not contained", "abc")
	}

	if s.Remove("non") != false {
		t.Errorf("expected %s to be not contained", "non")
	}

	s.AddAll("kk", "k", "äää")
	if !(s.Contains("kk") && s.Contains("k") && s.Contains("äää")) {
		t.Errorf("expected %s,%s,%s to be contained after addAll", "kk", "k", "äää")
	}

	var jsonStr = `
	[
		"abc", 
		"cba", 
		"bca"
	]`

	s2 := NewSet[string]()
	err := json.Unmarshal([]byte(jsonStr), &s2)

	if err != nil {
		t.Errorf("expected %s to unmarshal json", jsonStr)
	}

	if !(s2.Contains("abc") && s.Contains("cba") && s.Contains("bca")) {
		t.Errorf("expected strings to be contained")
	}

	if s2.Contains("") {
		t.Errorf("contained string that should not be there")
	}

	type contained struct {
		Field string      `json:"Field"`
		Ids   Set[string] `json:"Ids"`
	}

	var jsonStr2 = `
{
	"Field" : "value",
	"Ids" : [
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
	if c.Ids.Size() != 3 {
		t.Errorf("expected size to be %d ", 3)
	}

	if !(c.Ids.Contains("abc") && s.Contains("cba") && s.Contains("bca")) {
		t.Errorf("expected strings to be contained")
	}
}
