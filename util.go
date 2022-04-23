package main

import "strings"

func Lines(str string) []string {
	return strings.Split(str, "\n")
}

func Between(str string, start string, end string) string {
	return strings.Split(strings.Split(str, start)[1], end)[0]
}

type IntMod = func(int) int

func MakeInc(max int) IntMod {
	return func(i int) int {
		if i == max { return 0 }
		return i + 1
	}
}

func MakeDesc(max int) IntMod {
	return func(i int) int {
		if i == 0 { return max }
		return i - 1
	}
}