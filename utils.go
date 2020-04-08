package main

// Contact contact string array to string by ","
func Contact(list []string) string {
	if len(list) == 0 {
		return ""
	}
	if len(list) == 1 {
		return list[0]
	}

	start := 0
	for start < len(list) && list[start] == "" {
		start++
	}

	if start >= len(list) {
		return ""
	}

	res := list[start]
	for i := start + 1; i < len(list); i++ {
		if list[i] != "" {
			res += "," + list[i]
		}
	}

	return res
}
