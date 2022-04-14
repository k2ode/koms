package main

import "time"

type PersonOrGroupChat struct {
	conversations []Conversation
	contactIds    []string
	isGroupChat   bool
	label         string
}

type Conversation struct {
	proivder       string
	id             string
	participantIds []string
	isGroupChat    bool
	label          string
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