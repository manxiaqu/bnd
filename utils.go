package main

// Contact contact string array to string by "," 
func Contact(list []string) string {
	if len(list) == 0 {
		return ""
	}
	if len(list) == 1 {
		return list[0] 
	}

	res := list[0]
	for i := 1; i < len(list); i++ {
		res += "," + list[i]
	}

	return res
}