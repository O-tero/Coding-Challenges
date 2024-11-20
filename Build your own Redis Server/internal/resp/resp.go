package resp

import (
    "bytes"
    "errors"
    "fmt"
    "strconv"
    "strings"
)

func Deserialize(data []byte) (interface{}, error) {
    if len(data) == 0 {
        return nil, errors.New("empty input")
    }

    switch data[0] {
    case '+': // Simple Strings
        return strings.TrimSuffix(string(data[1:]), "\r\n"), nil
    case '-': // Errors
        return errors.New(strings.TrimSuffix(string(data[1:]), "\r\n")), nil
    case ':': // Integers
        val, err := strconv.Atoi(strings.TrimSuffix(string(data[1:]), "\r\n"))
        if err != nil {
            return nil, errors.New("invalid integer format")
        }
        return val, nil
    case '$': // Bulk Strings
        parts := bytes.SplitN(data[1:], []byte("\r\n"), 2)
        if len(parts) < 2 {
            return nil, errors.New("invalid bulk string format")
        }
        length, err := strconv.Atoi(string(parts[0]))
        if err != nil || length < 0 {
            return nil, errors.New("invalid bulk string length")
        }
        return string(parts[1][:length]), nil
    case '*': // Arrays
        parts := bytes.SplitN(data[1:], []byte("\r\n"), 2)
        if len(parts) < 2 {
            return nil, errors.New("invalid array format")
        }
        length, err := strconv.Atoi(string(parts[0]))
        if err != nil || length < 0 {
            return nil, errors.New("invalid array length")
        }
        arr := []interface{}{}
        remaining := parts[1]
        for i := 0; i < length; i++ {
            element, err := Deserialize(remaining)
            if err != nil {
                return nil, err
            }
            arr = append(arr, element)

            // Update `remaining` to skip the current element
            serializedElement, _ := Serialize(element)
            remaining = remaining[len(serializedElement):]
        }
        return arr, nil
    default:
        return nil, errors.New("unknown RESP type")
    }
}

func Serialize(value interface{}) (string, error) {
    switch v := value.(type) {
    case nil: // Null Bulk String
        return "$-1\r\n", nil
    case string:
        if len(v) == 0 || strings.Contains(v, "\r\n") {
            // Treat as a Bulk String
            return fmt.Sprintf("$%d\r\n%s\r\n", len(v), v), nil
        }
        // Default to Simple String
        return fmt.Sprintf("+%s\r\n", v), nil
    case error: // Error Message
        return fmt.Sprintf("-%s\r\n", v.Error()), nil
    case int: // Integer
        return fmt.Sprintf(":%d\r\n", v), nil
    case []interface{}: // Arrays
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