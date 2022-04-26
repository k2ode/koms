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
		return i + by
	}
}

func MakeInc(max int) IntMod {
	return MakeIncBy(max, 1)
}

func MakeDescBy(max int, by int) IntMod {
	return func(i int) int {
		if i == 0 { return max }
		return i - by
	}
}

func MakeDesc(max int) IntMod {
	return MakeDescBy(max, 1)
}