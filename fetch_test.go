package jpath_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/text/jpath"
)

type Output struct {
	Value interface{}
	Error string
}

func TestErrOutOfRange(t *testing.T) {
	for _, tc := range []struct {
		in           string
		query        string
		expectOutput []Output
		expectError  error
	}{
		{`["foo"]`, "[1]", errOut(jpath.ErrOutOfRange("[1]", 1)), nil},
	} {
		r := strings.NewReader(tc.in)
		ch, err := jpath.Fetch(r, tc.query)
		if errorString(err) != errorString(tc.expectError) {
			t.Errorf("%q expected error to be %v but got %v", tc.query, tc.expectError, err)
		}
		actual := readAll(ch)
		if !reflect.DeepEqual(actual, tc.expectOutput) {
			t.Errorf("%q expected value to be %v but got %v", tc.query, tc.expectOutput, actual)
		}
	}
}

func TestErrParse(t *testing.T) {
	parseAtoiError := func(s string) error {
		_, err := strconv.Atoi(s)
		return err
	}
	parseDecodeError := func(s string) error {
		dec := json.NewDecoder(strings.NewReader(s))
		m := map[string]interface{}{}
		return dec.Decode(&m)
	}
	for _, tc := range []struct {
		in           string
		query        string
		expectOutput []Output
		expectError  error
	}{
		{`["foo"]`, "[x]", errOut(jpath.ErrParse("[x]", parseAtoiError("x"))), nil},
		{`["foo"]`, "[y:0]", errOut(jpath.ErrParse("[y:0]", parseAtoiError("y"))), nil},
		{`["foo"]`, "[0:z]", errOut(jpath.ErrParse("[0:z]", parseAtoiError("z"))), nil},
		{"}", "[1]", nil, parseDecodeError("}")},
	} {
		r := strings.NewReader(tc.in)
		ch, err := jpath.Fetch(r, tc.query)
		if errorString(err) != errorString(tc.expectError) {
			t.Errorf("%q expected error to be %v but got %v", tc.query, tc.expectError, err)
		}
		actual := readAll(ch)
		if !reflect.DeepEqual(actual, tc.expectOutput) {
			t.Errorf("%q expected value to be %v but got %v", tc.query, tc.expectOutput, actual)
		}
	}
}

func TestFetchNested(t *testing.T) {
	for _, tc := range []struct {
		in           string
		query        string
		expectOutput []Output
		expectError  error
	}{
		{`{"arr":[{"text": "foo"}, {"text": "bar"}]}`, ".arr[0].text", okOut("foo"), nil},
	} {
		r := strings.NewReader(tc.in)
		ch, err := jpath.Fetch(r, tc.query)
		if errorString(err) != errorString(tc.expectError) {
			t.Errorf("%q expected error to be %v but got %v", tc.query, tc.expectError, err)
		}
		actual := readAll(ch)
		if !reflect.DeepEqual(actual, tc.expectOutput) {
			t.Errorf("%q expected value to be %v but got %v", tc.query, tc.expectOutput, actual)
		}
	}
}

func TestFetch(t *testing.T) {
	t.Run("array", func(t *testing.T) {
		for _, tc := range []struct {
			in           string
			query        string
			expectOutput []Output
			expectError  error
		}{
			{`[{"text": "foo"}, {"text": "bar"}]`, "[:].text", okOut("foo", "bar"), nil},
			{`[{"text": "foo"}, {"text": "bar"}]`, "[0].text", okOut("foo"), nil},
			{`[{"text": "foo"}, {"text": "bar"}]`, "[1].text", okOut("bar"), nil},
			{`[{"text": "foo"}, {"text": "bar"}, {"text": "baz"}]`, "[^1].text", okOut("baz"), nil},
			{`[{"text": "foo"}, {"text": "bar"}, {"text": "baz"}]`, "[^2].text", okOut("bar"), nil},
		} {
			r := strings.NewReader(tc.in)
			ch, err := jpath.Fetch(r, tc.query)
			if errorString(err) != errorString(tc.expectError) {
				t.Errorf("%q expected error to be %v but got %v", tc.query, tc.expectError, err)
			}
			actual := readAll(ch)
			if !reflect.DeepEqual(actual, tc.expectOutput) {
				t.Errorf("%q expected value to be %v but got %v", tc.query, tc.expectOutput, actual)
			}
		}
	})
	t.Run("slice", func(t *testing.T) {
		for _, tc := range []struct {
			in           string
			query        string
			expectOutput []Output
			expectError  error
		}{
			{`["foo", "bar", "baz", "pub"]`, "[1:3]", okOut("bar", "baz"), nil},
		} {
			r := strings.NewReader(tc.in)
			ch, err := jpath.Fetch(r, tc.query)
			if errorString(err) != errorString(tc.expectError) {
				t.Errorf("%q expected error to be %v but got %v", tc.query, tc.expectError, err)
			}
			actual := readAll(ch)
			if !reflect.DeepEqual(actual, tc.expectOutput) {
				t.Errorf("%q expected value to be %v but got %v", tc.query, tc.expectOutput, actual)
			}
		}
	})
	t.Run("root", func(t *testing.T) {
		for _, tc := range []struct {
			in           string
			query        string
			expectOutput []Output
			expectError  error
		}{
			{`{"text": "foo"}`, ".text", okOut("foo"), nil},
			{`{"text": "bar"}`, ".text", okOut("bar"), nil},
			{`{"text": "bar"}`, ".name", errOut(jpath.ErrNotFound), nil},
			{`{"org": {"name": "foo"}}`, ".org", okOut(map[string]interface{}{"name": "foo"}), nil},
			{`{"text": "bar"}`, ".text.name", errOut(jpath.ErrNotFound), nil},
		} {
			r := strings.NewReader(tc.in)
			ch, err := jpath.Fetch(r, tc.query)
			if errorString(err) != errorString(tc.expectError) {
				t.Errorf("%q expected error to be %v but got %v", tc.query, tc.expectError, err)
			}
			actual := readAll(ch)
			if !reflect.DeepEqual(actual, tc.expectOutput) {
				t.Errorf("%q expected value to be %v but got %v", tc.query, tc.expectOutput, actual)
			}
		}
	})
}

func TestErrNotFound(t *testing.T) {
	for _, tc := range []struct {
		in           string
		query        string
		expectOutput []Output
		expectError  error
	}{
		{`{"text": "bar"}`, ".text.name", errOut(jpath.ErrNotFound), nil},
	} {
		r := strings.NewReader(tc.in)
		ch, err := jpath.Fetch(r, tc.query)
		if errorString(err) != errorString(tc.expectError) {
			t.Errorf("%q expected error to be %v but got %v", tc.query, tc.expectError, err)
		}
		actual := readAll(ch)
		if !reflect.DeepEqual(actual, tc.expectOutput) {
			t.Errorf("%q expected value to be %v but got %v", tc.query, tc.expectOutput, actual)
		}
	}
}

func okOut(a ...interface{}) []Output {
	s := make([]Output, 0)
	for _, v := range a {
		s = append(s, Output{v, errorString(nil)})
	}
	return s
}

func errOut(a ...error) []Output {
	s := make([]Output, 0)
	for _, v := range a {
		s = append(s, Output{nil, errorString(v)})
	}
	return s
}

func readAll(ch <-chan struct {
	Value interface{}
	Error error
}) []Output {
	if ch == nil {
		return nil
	}
	s := make([]Output, 0)
	for v := range ch {
		s = append(s, Output{
			Value: v.Value,
			Error: errorString(v.Error),
		})
	}
	return s
}

func errorString(err error) string { return fmt.Sprintf("%v", err) }
