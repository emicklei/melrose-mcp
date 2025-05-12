package mcpserver

import (
	"errors"
	"strconv"
)

func toFloat64(input any) (float64, error) {
	switch v := input.(type) {
	case string:
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return 0, err
		}
		return f, nil
	case int:
		return float64(v), nil
	case float64:
		return v, nil
	default:
		return 0, errors.New("parameter must be a number")
	}
}

func toInt(input any) (int, error) {
	switch v := input.(type) {
	case string:
		i, err := strconv.Atoi(v)
		if err != nil {
			return 0, err
		}
		return i, nil
	case int:
		return v, nil
	case float64:
		return int(v), nil
	default:
		return 0, errors.New("parameter must be a number")
	}
}
