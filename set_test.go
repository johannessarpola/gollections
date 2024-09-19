package gollections

import "testing"

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
}
