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
    case '+': // Simple String
        return string(data[1 : len(data)-2]), nil
    case '-': // Error Message
        return string(data[1 : len(data)-2]), nil
    case ':': // Integer
        return strconv.Atoi(string(data[1 : len(data)-2]))
    case '$': // Bulk String
        parts := bytes.SplitN(data[1:], []byte("\r\n"), 2)
        length, err := strconv.Atoi(string(parts[0]))
        if err != nil {
            return nil, err
        }
        if length == -1 {
            return nil, nil
        }
        return string(parts[1][:length]), nil
    case '*': // Array
        parts := bytes.SplitN(data[1:], []byte("\r\n"), 2)
        length, err := strconv.Atoi(string(parts[0]))
        if err != nil || length < 0 {
            return nil, errors.New("invalid array length")
        }
        var arr []interface{}
        remaining := parts[1]
        for i := 0; i < length; i++ {
            element, err := Deserialize(remaining)
            if err != nil {
                return nil, err
            }
            arr = append(arr, element)

            // Skip serialized element
            serialized, _ := Serialize(element)
            remaining = remaining[len(serialized):]
        }
        return arr, nil
    default:
        return nil, errors.New("unknown RESP type")
    }
}



func Serialize(value interface{}) (string, error) {
    switch v := value.(type) {
    case nil:
        return "$-1\r\n", nil
    case string:
        // Treat as Bulk String
        return fmt.Sprintf("$%d\r\n%s\r\n", len(v), v), nil
    case int:
        return fmt.Sprintf(":%d\r\n", v), nil
    case []interface{}:
        var sb strings.Builder
        sb.WriteString(fmt.Sprintf("*%d\r\n", len(v)))
        for _, elem := range v {
            serialized, err := Serialize(elem)
            if err != nil {
                return "", err
            }
            sb.WriteString(serialized)
        }
        return sb.String(), nil
    default:
        return "", fmt.Errorf("unsupported type: %T", value)
    }
}
