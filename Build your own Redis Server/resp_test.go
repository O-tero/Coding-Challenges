package resp

import (
	"errors"
	"reflect"
	"testing"
)

func TestSerialize(t *testing.T) {
    tests := []struct {
        name     string
        input    interface{}
        expected string
    }{
        {"Simple String", "OK", "+OK\r\n"},
        {"Error Message", errors.New("Error message"), "-Error message\r\n"},
        {"Empty Bulk String", "", "$0\r\n\r\n"},
        {"Bulk String", "hello world", "$11\r\nhello world\r\n"},
        {
            "Array with Bulk Strings",
            []interface{}{"get", "key"},
            "*2\r\n$3\r\nget\r\n$3\r\nkey\r\n",
        },
    }

    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            result, err := Serialize(test.input)
            if err != nil {
                t.Fatalf("unexpected error: %v", err)
            }
            if result != test.expected {
                t.Errorf("expected %q, got %q", test.expected, result)
            }
        })
    }
}

func TestDeserialize(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        expected interface{}
    }{
        {"Simple String", "+OK\r\n", "OK"},
        {"Error Message", "-Error message\r\n", errors.New("Error message")},
        {"Empty Bulk String", "$0\r\n\r\n", ""},
        {"Bulk String", "$11\r\nhello world\r\n", "hello world"},
        {
            "Array with Bulk Strings",
            "*2\r\n$3\r\nget\r\n$3\r\nkey\r\n",
            []interface{}{"get", "key"},
        },
    }

    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            result, err := Deserialize([]byte(test.input))
            if err != nil {
                t.Fatalf("unexpected error: %v", err)
            }
            if !reflect.DeepEqual(result, test.expected) {
                t.Errorf("expected %v, got %v", test.expected, result)
            }
        })
    }
}
