package main

import "time"

type Conversation struct {
	id          string
	label       string
	isGroupChat bool
}

type Message struct {
	id        string
	from      string
	body      string
	timestamp time.Time
	reactions []Reaction
}

type Reaction struct {
	emoji string
	from  string
}

type Contact struct {
	id   string
	name string
	tags []string
}