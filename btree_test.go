package gollections

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	type abc struct {
		a string
		b string
		c string
	}

	m, err := json.Marshal(abc{
		a: "aa",
		b: "bb",
		c: "cc",
	})
	if err != nil {
		t.Fatal(err)
	}
	hash := sha256.New()
	_, err = hash.Write(m)
	if err != nil {
		t.Fatal(err)
	}
	hsum := hex.EncodeToString(hash.Sum(nil))
	fmt.Println(hsum)
}
