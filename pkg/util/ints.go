package util

import "slices"

func AddIntToSet(s []int, i int) []int {
	if !slices.Contains(s, i) {
		s = append(s, i)
	}
	return s
}

func AddIntsToSet(s []int, i []int) []int {
	for _, v := range i {
		s = AddIntToSet(s, v)
	}
	return s
}

func RemoveIntFromSet(s []int, i int) []int {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func RemoveIntsFromSet(s []int, i []int) []int {
	for _, v := range i {
		s = RemoveIntFromSet(s, v)
	}
	return s
}
