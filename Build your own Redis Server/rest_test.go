package resp

import (
    "bytes"
    "testing"
)

func TestDeserialize(t *testing.T) {
    testCases := []struct {
        name     string
        input    string
        expected interface{}
        isError  bool
    }{
        {"Null Bulk String", "$-1\r\n", nil, false},
        {"Simple String", "+OK\r\n", "OK", false},
        {"Error Message", "-Error message\r\n", "Error message", false},
        {"Integer", ":100\r\n", 100, false},
        {"Empty Bulk String", "$0\r\n\r\n", "", false},
        {"Bulk String", "$11\r\nhello world\r\n", "hello world", false},
        {"Array with Bulk Strings", "*2\r\n$3\r\nget\r\n$3\r\nkey\r\n", []interface{}{"get", "key"}, false},
        {"Invalid Prefix", "@123\r\n", nil, true},
        {"Invalid Integer", ":abc\r\n", nil, true},
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result, err := Deserialize([]byte(tc.input))
            if tc.isError {
                if err == nil {
                    t.Errorf("expected error but got none")
                }
                return
            }
            if err != nil {
                t.Errorf("unexpected error: %v", err)
                return
            }
            if !compareValues(result, tc.expected) {
                t.Errorf("expected %v, got %v", tc.expected, result)
            }
        })
    }
}

func TestSerialize(t *testing.T) {
    testCases := []struct {
        name     string
        input    interface{}
        expected string
        isError  bool
    }{
        {"Null Bulk String", nil, "$-1\r\n", false},
        {"Simple String", "OK", "+OK\r\n", false},
        {"Error Message", "Error message", "-Error message\r\n", false},
        {"Integer", 100, ":100\r\n", false},
        {"Empty Bulk String", "", "$0\r\n\r\n", false},
        {"Bulk String", "hello world", "$11\r\nhello world\r\n", false},
        {"Array with Bulk Strings", []interface{}{"get", "key"}, "*2\r\n$3\r\nget\r\n$3\r\nkey\r\n", false},
        {"Invalid Type", map[string]string{"key": "value"}, "", true},
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result, err := Serialize(tc.input)
            if tc.isError {
                if err == nil {
                    t.Errorf("expected error but got none")
                }
                return
            }
            if err != nil {
                t.Errorf("unexpected error: %v", err)
                return
            }
            if result != tc.expected {
                t.Errorf("expected %v, got %v", tc.expected, result)
            }
        })
    }
}

func compareValues(a, b interface{}) bool {
    switch a := a.(type) {
    case []interface{}:
        b, ok := b.([]interface{})
        if !ok || len(a) != len(b) {
            return false
        }
        for i := range a {
            if !compareValues(a[i], b[i]) {
                return false
            }
        }
        return true
    default:
        return a == b
    }
}
