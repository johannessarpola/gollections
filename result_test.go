package gollections

import (
	"errors"
	"testing"
)

func TestUnwrap(t *testing.T) {
	var rss []Result[int]

	rss = append(rss, Wrap(1, nil))
	rss = append(rss, Wrap(2, nil))
	rss = append(rss, Wrap(3, nil))
	rss = append(rss, Wrap(0, errors.New("ping pong computer is broke")))

	ecount := 0
	r := UnwrapResults(rss, func(err error) {
		ecount++
	})

	if len(r) == len(rss) {
		t.Errorf("error should have been filtered and not included in result values")
	}

	if ecount == 0 {
		t.Errorf("error should have been filtered and sent into the handler")
	}

}

func TestFanout(t *testing.T) {
	var rss []Result[int]

	rss = append(rss, Wrap(1, nil))
	rss = append(rss, Wrap(2, nil))
	rss = append(rss, Wrap(3, nil))
	rss = append(rss, Wrap(0, errors.New("ping pong computer is broke")))

	r, e := FanOut(rss)

	if len(r) == len(rss) {
		t.Errorf("error should have been filtered and not included in result values")
	}

	if len(e) == 0 {
		t.Errorf("error should have been filtered into err list")
	}

}
