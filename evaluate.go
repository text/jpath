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
	ErrOutOfRange = func(square string, length int) error {
		return fmt.Errorf("index out of range %s with length %d", square, length)
	}
	ErrParse = func(square string, err error) error {
		return fmt.Errorf("could not parse %s: %w", square, err)
	}
)

func Evaluate(expression string, r io.Reader) (<-chan struct {
	Value interface{}
	Error error
}, error) {
	pre := strings.NewReader(`{"root":`)
	suf := strings.NewReader(`}`)
	dec := json.NewDecoder(io.MultiReader(pre, r, suf))
	context := make(map[string]interface{})
	if err := dec.Decode(&context); err != nil && err != io.EOF {
		return nil, err
	}
	omitFirstEmpty := func() []string {
		a := strings.Split(expression, ".")
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
		evaluate(omitFirstEmpty(), context["root"], ch)
		close(ch)
	}()
	return ch, nil
}

func evaluate(expression []string, context interface{}, ch chan<- struct {
	Value interface{}
	Error error
}) {
	l := len(expression)
	if l == 0 {
		ch <- struct {
			Value interface{}
			Error error
		}{context, nil}
		return
	}
	switch x := context.(type) {
	case map[string]interface{}:
		name, square := takeSquareBracket(expression[0])
		nextContext, ok := x[name]
		if !ok {
			ch <- struct {
				Value interface{}
				Error error
			}{nil, ErrNotFound}
			return
		}
		if square == "" {
			evaluate(expression[1:], nextContext, ch)
			return
		}
		expression[0] = square
		evaluate(expression, nextContext, ch)
	case []interface{}:
		low, high, err := slice(x, expression[0])
		if err != nil {
			ch <- struct {
				Value interface{}
				Error error
			}{nil, err}
		}
		nextExpression := expression[1:]
		for _, nextContext := range x[low:high] {
			evaluate(nextExpression, nextContext, ch)
		}
	case string:
		if l != 0 {
			ch <- struct {
				Value interface{}
				Error error
			}{nil, ErrNotFound}
			return
		}
	}
}

func takeSquareBracket(s string) (name, square string) {
	i := strings.Index(s, "[")
	if i < 0 {
		return s, ""
	}
	return s[:i], s[i:]
}

func slice(a []interface{}, square string) (low, high int, err error) {
	if square == "[:]" {
		return 0, len(a), nil
	}
	s := square[1 : len(square)-1]
	t := strings.Split(s, ":")
	if len(t) == 2 {
		low, err := strconv.Atoi(t[0])
		if err != nil {
			return 0, 0, ErrParse(square, err)
		}
		high, err := strconv.Atoi(t[1])
		if err != nil {
			return 0, 0, ErrParse(square, err)
		}
		return low, high, nil
	}
	s = t[0]
	fromEnd := false
	if s[0] == '^' {
		fromEnd = true
		s = s[1:]
	}
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, 0, ErrParse(square, err)
	}
	al := len(a)
	if fromEnd {
		i = al - i
	}
	if i < 0 || i >= al {
		return 0, 0, ErrOutOfRange(square, al)
	}
	return i, i + 1, nil
}
