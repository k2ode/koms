package types

import "time"

type Conversation struct {
	Conversations []ConversationRaw
	ContactIds    []string
	IsGroupChat   bool
	Label         string
}

type ConversationRaw struct {
	Id             string
	ParticipantIds []string
	IsGroupChat    bool
	Label          string
	Provider       string
}

type MessageRaw struct {
	Id        string
	From      string
	Body      string
	Timestamp time.Time
	Reactions []Reaction
}

type Message struct {
	Id        string
	From      Contact
	FromUser  bool
	Body      string
	Provider  string
	Timestamp time.Time
	Reactions []Reaction
}

type Reaction struct {
	Emoji string
	From  string
}

type Contact struct {
	Id   string
	Name string
	Tags []string
}