package main

import "strings"

func Lines(str string) []string {
	return strings.Split(str, "\n")
}

func Between(str string, start string, end string) string {
	return strings.Split(strings.Split(str, start)[1], end)[0]
}

type IntMod = func(int) int

func MakeIncBy(max int, by int) IntMod {
	return func(i int) int {
		if i == max { return 0 }
		if i + by > max { return max }
		return i + by
	}
}

func MakeInc(max int) IntMod {
	return MakeIncBy(max, 1)
}

func MakeDescBy(max int, by int) IntMod {
	return func(i int) int {
		if i == 0 { return max }
		if i - by < 0 { return 0 }
		return i - by
	}
}

func MakeDesc(max int) IntMod {
	return MakeDescBy(max, 1)
}

func Contains(haystack []string, needle string) bool { 
	return Find(haystack, needle) != -1
}

func Find(haystack []string, needle string) int {
	for index, item := range haystack {
		if item != needle { continue } 
		return index
	}
	return -1
}
