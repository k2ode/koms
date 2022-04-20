package main

import "time"

type Conversation struct {
	conversations []ConversationRaw
	contactIds    []string
	isGroupChat   bool
	label         string
}

type ConversationRaw struct {
	id             string
	participantIds []string
	isGroupChat    bool
	label          string
	provider       string
}

type MessageRaw struct {
	id        string
	from      string
	body      string
	timestamp time.Time
	reactions []Reaction
}

type Message struct {
	id        string
	from      Contact
	body      string
	provider  string
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