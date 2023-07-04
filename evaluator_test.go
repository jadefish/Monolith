package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSplitPath(t *testing.T) {
	type test struct {
		input  string
		output []string
	}
	tests := map[string]test{
		"empty string":      {"", []string{""}},
		"whitespace string": {" ", []string{""}},
		"single key":        {"key", []string{"key"}},
		"nested key":        {"key1.key2", []string{"key1", "key2"}},
		"malformed":         {"....", []string{""}},
		"duplicate keys":    {"key1.key2.key2.key3", []string{"key1", "key2", "key2", "key3"}},
		"numeric key":       {"key1.2.key2", []string{"key1", "2", "key2"}},
		"extra whitespace":  {" key1 . 2 . key2 ", []string{"key1", "2", "key2"}},
		"extra separators":  {".key1.key2.key3.", []string{"key1", "key2", "key3"}},
		"keys with symbols": {"DeviceProperties.Add.PciRoot(0x0)/Pci(0x1f,0x3).layout-id", []string{"DeviceProperties", "Add", "PciRoot(0x0)/Pci(0x1f,0x3)", "layout-id"}},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			keys, last := splitPath(test.input)
			keys = append(keys, last)

			if diff := cmp.Diff(test.output, keys); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestIndex(t *testing.T) {
	type test struct {
		value any
		key   string
		want  any
	}
	tests := map[string]test{
		"map":              {map[string]any{"test": "value"}, "test", "value"},
		"map, int key":     {map[string]any{"4": "value"}, "4", "value"},
		"map, missing key": {map[string]any{"test": "value"}, "key", nil},
		"other structure":  {7, "0", 7},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got, _ := index(test.value, test.key)

			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestMutate(t *testing.T) {
	type test struct {
		s     any
		key   string
		value any
		want  any
	}
	tests := map[string]test{
		"map":   {map[string]any{"one": 1, "two": 2, "three": 3}, "two", "TWO!", map[string]any{"one": 1, "two": "TWO!", "three": 3}},
		"slice": {[]any{1, 2, 3, 4}, "0", 10, []any{10, 2, 3, 4}},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mutate(test.s, test.key, test.value)

			if diff := cmp.Diff(test.s, test.want); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestSet(t *testing.T) {
	type test struct {
		data  map[string]any
		path  string
		value any
		want  map[string]any
	}
	tests := map[string]test{
		"single key": {map[string]any{"key": "value"}, "key", 7, map[string]any{"key": 7}},
		"nested key": {map[string]any{"key": map[string]any{"nested": "value"}}, "key.nested", 7, map[string]any{"key": map[string]any{"nested": 7}}},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			evaluator := Evaluator{test.data}
			err := evaluator.set(test.path, test.value)

			if err != nil {
				t.Error(err)
			}

			if diff := cmp.Diff(test.want, evaluator.data); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestAppend(t *testing.T) {
	type test struct {
		data  map[string]any
		path  string
		value any
		want  map[string]any
	}
	tests := map[string]test{
		"single key": {map[string]any{"key": []any{1, 2, 3}}, "key", 7, map[string]any{"key": []any{1, 2, 3, 7}}},
		"nested key": {map[string]any{"key": map[string]any{"nested": []any{1, 2, 3}}}, "key.nested", 7, map[string]any{"key": map[string]any{"nested": []any{1, 2, 3, 7}}}},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			evaluator := Evaluator{test.data}
			err := evaluator.append(test.path, test.value)

			if err != nil {
				t.Error(err)
			}

			if diff := cmp.Diff(test.want, evaluator.data); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	type test struct {
		data map[string]any
		path string
		want map[string]any
	}
	tests := map[string]test{
		"single key": {map[string]any{"key": []any{1, 2, 3}}, "key", map[string]any{}},
		"nested key": {map[string]any{"key": map[string]any{"nested": []any{1, 2, 3}}}, "key.nested", map[string]any{"key": map[string]any{}}},
		"array":      {map[string]any{"key": []any{1, 2, 3}}, "key.1", map[string]any{"key": []any{1, 3}}},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			evaluator := Evaluator{test.data}
			err := evaluator.delete(test.path)

			if err != nil {
				t.Error(err)
			}

			if diff := cmp.Diff(test.want, evaluator.data); diff != "" {
				t.Error(diff)
			}
		})
	}
}
