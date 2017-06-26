package featureflow_go_sdk

import "strings"

func TestString(op string, a string, b []string) bool {
	switch op {
	case "equals":
		return equals(a, b[0])
	case "contains":
		return contains(a, b[0])
	case "startsWith":
		return contains(a, b[0])
	case "endsWith":
		return contains(a, b[0])
	default:
		return false
	}
}

func equals(a string, b string) bool{
	return a == b
}

func contains(a string, b string) bool{
	return strings.Contains(a, b)
}

func startsWith(a string, b string) bool{
	return strings.HasPrefix(a, b)
}

func endsWith(a string, b string) bool{
	return strings.HasSuffix(a, b)
}

