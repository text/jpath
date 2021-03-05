package jpath

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

var (
	ErrNotFound   = errors.New("not found")
	ErrFoundOther = func(x interface{}) error { return fmt.Errorf("found other type: %T", x) }
	ErrOutOfRange = errors.New("out of range")
)

func Fetch(r io.Reader, path string) (<-chan struct {
	Value interface{}
	Error error
}, error) {
	pre := strings.NewReader(`{"root":`)
	suf := strings.NewReader(`}`)
	dec := json.NewDecoder(io.MultiReader(pre, r, suf))
	m := make(map[string]interface{})
	if err := dec.Decode(&m); err != nil && err != io.EOF {
		return nil, err
	}
	omitFirstEmpty := func() []string {
		a := strings.Split(path, ".")
		if a[0] == "" {
			return a[1:]
		}
		return a
	}
	ch := make(chan struct {
		Value interface{}
		Error error
	})
	go func() {
		fetch(m["root"], omitFirstEmpty(), ch)
		close(ch)
	}()
	return ch, nil
}

func fetch(m interface{}, path []string, ch chan<- struct {
	Value interface{}
	Error error
}) {
	l := len(path)
	switch x := m.(type) {
	case map[string]interface{}:
		if l == 0 {
			ch <- struct {
				Value interface{}
				Error error
			}{nil, ErrFoundOther(x)}
			return
		}
		name, square := takeSquareBracket(path[0])
		n, ok := x[name]
		if !ok {
			ch <- struct {
				Value interface{}
				Error error
			}{nil, ErrNotFound}
			return
		}
		if square == "" {
			fetch(n, path[1:], ch)
			return
		}
		path[0] = square
		fetch(n, path, ch)
	case []interface{}:
		i, ok := index(x, path[0])
		if !ok {
			ch <- struct {
				Value interface{}
				Error error
			}{nil, ErrOutOfRange}
			return
		}
		if i > -1 {
			fetch(x[i], path[1:], ch)
			return
		}
		for _, n := range x {
			fetch(n, path[1:], ch)
		}
	case string:
		if l != 0 {
			ch <- struct {
				Value interface{}
				Error error
			}{nil, ErrNotFound}
			return
		}
		ch <- struct {
			Value interface{}
			Error error
		}{x, nil}
	default:
		panic(fmt.Sprintf("consider type %T", x))
	}
}

func takeSquareBracket(s string) (name, square string) {
	i := strings.Index(s, "[")
	if i < 0 {
		return s, ""
	}
	return s[:i], s[i:]
}

func index(a []interface{}, s string) (int, bool) {
	if s == "[*]" {
		return -1, true
	}
	l := len(s) - 1
	if s[0] != '[' && s[l] != ']' {
		return 0, false
	}
	fromEnd := false
	if s[1] == '^' {
		fromEnd = true
		s = s[2:l]
	} else {
		s = s[1:l]
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, false
	}
	al := len(a)
	if fromEnd {
		i = al - i
	}
	if i < 0 || i >= al {
		return 0, false
	}
	return i, true
}
