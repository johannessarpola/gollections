package gollections

import (
	"encoding/json"
	"reflect"
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

	s2 := NewSet[string]()
	s2.AddAll("kk", "k", "äää")
	if !(s2.Contains("kk") && s2.Contains("k") && s2.Contains("äää")) {
		t.Errorf("expected %s,%s,%s to be contained after addAll", "kk", "k", "äää")
	}

	want := []string{"kk", "k", "äää"}
	if !reflect.DeepEqual(s2.Items(), want) {
		t.Errorf("expected %v to be contained after addAll", want)
	}

	var jsonStr = `
	[
		"abc", 
		"cba", 
		"bca"
	]`

	s3 := NewSet[string]()
	err := json.Unmarshal([]byte(jsonStr), &s3)

	if err != nil {
		t.Errorf("expected %s to unmarshal json", jsonStr)
	}

	if !(s3.Contains("abc") && s3.Contains("cba") && s3.Contains("bca")) {
		t.Errorf("expected strings to be contained")
	}

	if s3.Contains("") {
		t.Errorf("contained string that should not be there")
	}

	var ns3 Set[string]
	jb, err := json.Marshal(s3)
	if err != nil || len(jb) == 0 {
		t.Errorf("expected %s to marshal json", jsonStr)
	}

	if json.Unmarshal(jb, &ns3) != nil || !reflect.DeepEqual(ns3.Items(), s3.Items()) {
		t.Errorf("expected %v, got %v", s3.Items(), ns3.Items())
	}

	type contained struct {
		Field string      `json:"Field"`
		Set   Set[string] `json:"Set"`
	}

	var jsonStr2 = `
{
	"Field" : "value",
	"Set" : [
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
	if c.Set.Size() != 3 {
		t.Errorf("expected size to be %d ", 3)
	}

	if !(c.Set.Contains("abc") && s.Contains("cba") && s.Contains("bca")) {
		t.Errorf("expected strings to be contained")
	}

}
