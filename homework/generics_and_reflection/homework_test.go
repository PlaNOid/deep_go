package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type Person struct {
	Name    string `properties:"name"`
	Address string `properties:"address,omitempty"`
	Age     int    `properties:"age"`
	Married bool   `properties:"married"`
}

func Serialize(v interface{}) (string, error) {
	val := reflect.ValueOf(v)

	if val.Kind() == reflect.Pointer {
		if val.IsNil() {
			return "", fmt.Errorf("nil pointer")
		}
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return "", fmt.Errorf("expected struct, got %s", val.Kind())
	}

	var lines []string
	t := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := t.Field(i)
		fieldVal := val.Field(i)

		tag := field.Tag.Get("properties")
		if tag == "" || tag == "-" {
			continue
		}

		name, omitempty := parseTag(tag)

		if omitempty && isEmpty(fieldVal) {
			continue
		}

		lines = append(lines, fmt.Sprintf("%s=%v", name, fieldVal.Interface()))
	}
	return strings.Join(lines, "\n"), nil
}

func parseTag(tag string) (string, bool) {
	omitempty := strings.Contains(tag, "omitempty")
	parts := strings.Split(tag, ",")
	var name string
	if len(parts) > 1 {
		for _, opt := range parts {
			if strings.TrimSpace(opt) != "omitempty" {
				name = opt
			}
		}
	} else {
		name = tag
	}

	return name, omitempty
}

func isEmpty(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Pointer, reflect.Interface:
		return v.IsNil()
	}
	return v.IsZero()
}

func TestSerialization(t *testing.T) {
	tests := map[string]struct {
		person Person
		result string
	}{
		"test case with empty fields": {
			result: "name=\nage=0\nmarried=false",
		},
		"test case with fields": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
			},
			result: "name=John Doe\nage=30\nmarried=true",
		},
		"test case with omitempty field": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
				Address: "Paris",
			},
			result: "name=John Doe\naddress=Paris\nage=30\nmarried=true",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result, err := Serialize(test.person)
			assert.Nil(t, err)
			assert.Equal(t, test.result, result)
		})
	}
}
