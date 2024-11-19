package resp

import (
    "bytes"
    "errors"
    "fmt"
    "strconv"
    "strings"
)

// Deserialize parses a RESP message into Go types
func Deserialize(data []byte) (interface{}, error) {
    if len(data) == 0 {
        return nil, errors.New("empty input")
    }

    switch data[0] {
    case '+': // Simple Strings
        return strings.TrimSuffix(string(data[1:]), "\r\n"), nil
    case '-': // Errors
        return strings.TrimSuffix(string(data[1:]), "\r\n"), nil
    case ':': // Integers
        val, err := strconv.Atoi(strings.TrimSuffix(string(data[1:]), "\r\n"))
        if err != nil {
            return nil, errors.New("invalid integer format")
        }
        return val, nil
    case '$': // Bulk Strings
        parts := bytes.SplitN(data[1:], []byte("\r\n"), 2)
        if string(parts[0]) == "-1" {
            return nil, nil // Null Bulk String
        }
        length, err := strconv.Atoi(string(parts[0]))
        if err != nil || length < 0 {
            return nil, errors.New("invalid bulk string length")
        }
        return string(parts[1][:length]), nil
    case '*': // Arrays
        parts := bytes.Split(data[1:], []byte("\r\n"))
        length, err := strconv.Atoi(string(parts[0]))
        if err != nil || length < 0 {
            return nil, errors.New("invalid array length")
        }
        var arr []interface{}
        remaining := parts[1:]
        for i := 0; i < length; i++ {
            value, err := Deserialize(remaining)
            if err != nil {
                return nil, err
            }
            arr = append(arr, value)
            // Adjust remaining based on value length
        }
        return arr, nil
    default:
        return nil, errors.New("unknown RESP type")
    }
}

// Serialize converts Go types into RESP format
func Serialize(value interface{}) (string, error) {
    switch v := value.(type) {
    case nil: // Null Bulk String
        return "$-1\r\n", nil
    case string:
        if strings.HasPrefix(v, "+") || strings.HasPrefix(v, "-") {
            return v + "\r\n", nil
        }
        return fmt.Sprintf("$%d\r\n%s\r\n", len(v), v), nil
    case int:
        return fmt.Sprintf(":%d\r\n", v), nil
    case []interface{}:
        var buf bytes.Buffer
        buf.WriteString(fmt.Sprintf("*%d\r\n", len(v)))
        for _, elem := range v {
            serialized, err := Serialize(elem)
            if err != nil {
                return "", err
            }
            buf.WriteString(serialized)
        }
        return buf.String(), nil
    default:
        return "", errors.New("unsupported type")
    }
}
