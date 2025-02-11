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
	for range s.All {
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
	want := []string{"kk", "k", "kek"}
	s2.AddAll(want...)

	for _, v := range want {
		if !s2.Contains(v) {
			t.Errorf("expected %s to be contained", v)
		}
	}

	jsonStr := `
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

	if json.Unmarshal(jb, &ns3) != nil || ns3.Size() != s3.Size() {
		t.Errorf("expected %v, got %v", s3.Items(), ns3.Items())
	}

	for _, v := range ns3.Items() {
		if !s3.Contains(v) {
			t.Errorf("expected %s to be contained", v)
		}
	}

	type contained struct {
		Field string      `json:"field"`
		Set   Set[string] `json:"set"`
	}

	jsonStr2 := `
{
	"field" : "value",
	"set" : [
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

	if !(c.Set.Contains("abc") && c.Set.Contains("cba") && c.Set.Contains("bca")) {
		t.Errorf("expected strings to be contained")
	}
}
