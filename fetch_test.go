package jpath_test

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/text/jpath"
)

type Output struct {
	Value interface{}
	Error string
}

func TestFetchArray(t *testing.T) {
	for _, tc := range []struct {
		in           string
		path         string
		expectOutput []Output
		expectError  error
	}{
		{`[{"text": "foo"}, {"text": "bar"}]`, "[*].text", okOut("foo", "bar"), nil},
		{`[{"text": "foo"}, {"text": "bar"}]`, "[0].text", okOut("foo"), nil},
		{`[{"text": "foo"}, {"text": "bar"}]`, "[1].text", okOut("bar"), nil},
		{`[{"text": "foo"}, {"text": "bar"}, {"text": "baz"}]`, "[^1].text", okOut("baz"), nil},
		{`[{"text": "foo"}, {"text": "bar"}, {"text": "baz"}]`, "[^2].text", okOut("bar"), nil},
		{`[{"text": "foo"}, {"text": "bar"}]`, "[2].text", errOut(jpath.ErrOutOfRange), nil},
	} {
		r := strings.NewReader(tc.in)
		ch, err := jpath.Fetch(r, tc.path)
		if errorString(err) != errorString(tc.expectError) {
			t.Errorf("%q expected error to be %v but got %v", tc.path, tc.expectError, err)
		}
		actual := readAll(ch)
		if !reflect.DeepEqual(actual, tc.expectOutput) {
			t.Errorf("%q expected value to be %v but got %v", tc.path, tc.expectOutput, actual)
		}
	}
}

func TestFetchNested(t *testing.T) {
	for _, tc := range []struct {
		in           string
		path         string
		expectOutput []Output
		expectError  error
	}{
		{`{"arr":[{"text": "foo"}, {"text": "bar"}]}`, ".arr[0].text", okOut("foo"), nil},
	} {
		r := strings.NewReader(tc.in)
		ch, err := jpath.Fetch(r, tc.path)
		if errorString(err) != errorString(tc.expectError) {
			t.Errorf("%q expected error to be %v but got %v", tc.path, tc.expectError, err)
		}
		actual := readAll(ch)
		if !reflect.DeepEqual(actual, tc.expectOutput) {
			t.Errorf("%q expected value to be %v but got %v", tc.path, tc.expectOutput, actual)
		}
	}
}

func TestFetch(t *testing.T) {
	for _, tc := range []struct {
		in           string
		path         string
		expectOutput []Output
		expectError  error
	}{
		{`{"text": "foo"}`, ".text", okOut("foo"), nil},
		{`{"text": "bar"}`, ".text", okOut("bar"), nil},
		{`{"text": "bar"}`, ".name", errOut(jpath.ErrNotFound), nil},
		{`{"org": {"name": "foo"}}`, ".org", errOut(jpath.ErrFoundOther(make(map[string]interface{}))), nil},
		{`{"text": "bar"}`, ".text.name", errOut(jpath.ErrNotFound), nil},
	} {
		r := strings.NewReader(tc.in)
		ch, err := jpath.Fetch(r, tc.path)
		if errorString(err) != errorString(tc.expectError) {
			t.Errorf("%q expected error to be %v but got %v", tc.path, tc.expectError, err)
		}
		actual := readAll(ch)
		if !reflect.DeepEqual(actual, tc.expectOutput) {
			t.Errorf("%q expected value to be %v but got %v", tc.path, tc.expectOutput, actual)
		}
	}
}

func okOut(a ...string) []Output {
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
