package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Evaluator struct {
	data map[string]any
}

// splitPath returns a (tokens, last) tuple containing the values of a
// dot-delimited string path.
//
// For example, an input path of "Misc.Security.AllowSetDefault" returns
// a tuple (["Misc", "Security"], "AllowSetDefault").
func splitPath(path string) ([]string, string) {
	tokens := strings.Split(path, ".")
	keys := make([]string, 0, len(tokens))

	for _, key := range tokens {
		strippedKey := strings.TrimSpace(key)

		if len(strippedKey) < 1 {
			continue
		}

		keys = append(keys, strippedKey)
	}

	if len(keys) < 1 {
		return []string{}, ""
	}

	n := len(keys) - 1

	return keys[0:n], keys[n]
}

// index retrieves the value of key within thing, if thing is a map.
// If no value for the provided key exists, nil is returned.
// If thing is not a map, thing is returned with no error.
func index(thing any, key string) (any, error) {
	switch typedThing := thing.(type) {
	case map[string]any:
		if value, ok := typedThing[key]; !ok {
			return nil, fmt.Errorf("key not found: \"%s\"", key)
		} else {
			return value, nil
		}
	}

	return thing, nil
}

func mutate(s any, key string, value any) {
	switch typedS := s.(type) {
	case map[string]any:
		typedS[key] = value
	case []any:
		if i, err := strconv.Atoi(key); err != nil {
			panic("found slice, but have string index")
		} else {
			typedS[i] = value
		}
	default:
		panic("refusing to handle this type")
	}
}

func (e *Evaluator) set(path string, value any) error {
	keys, last := splitPath(path)
	var f any = e.data
	var err error

	for _, key := range keys {
		f, err = index(f, key)

		if err != nil {
			return fmt.Errorf("mutate: %w", err)
		}
	}

	mutate(f, last, value)

	return nil
}

func (e *Evaluator) append(path string, value any) error {
	keys, last := splitPath(path)
	var f any = e.data
	var err error

	for _, key := range keys {
		f, err = index(f, key)

		if err != nil {
			return fmt.Errorf("append: %w", err)
		}
	}

	slice, err := index(f, last)

	if err != nil {
		return fmt.Errorf("append: %w", err)
	}

	slice2, ok := slice.([]any)

	if !ok {
		return fmt.Errorf("append: cannot append to non-array entry")
	}

	slice2 = append(slice2, value)
	mutate(f, last, slice2)

	return nil
}

func (e *Evaluator) delete(path string) error {
	keys, last := splitPath(path)
	var f any = e.data
	var err error

	for _, key := range keys {
		f, err = index(f, key)

		if err != nil {
			return fmt.Errorf("delete: %w", err)
		}
	}

	switch f.(type) {
	case map[string]any:
		m := f.(map[string]any)
		delete(m, last)
	case []any:
		s := f.([]any)
		i, err := strconv.Atoi(last)

		if err != nil {
			return fmt.Errorf(
				"delete: cannot index array with string key \"%s\"",
				last,
			)
		}

		if i >= len(s) {
			return fmt.Errorf("delete: index out of bounds")
		}

		tmp := make([]any, 0)
		tmp = append(tmp, s[:i]...)
		tmp = append(tmp, s[i+1:]...)

		return e.set(strings.Join(keys, "."), tmp)
	}

	return nil
}
