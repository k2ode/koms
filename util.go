package main

import "strings"

func Lines(str string) []string {
	return strings.Split(str, "\n")
}

func Between(str string, start string, end string) string {
	return strings.Split(strings.Split(str, start)[1], end)[0]
}