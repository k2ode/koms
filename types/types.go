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
	Images    []Image
	URL       string
	BodyFull  string
}

type Message struct {
	Raw       MessageRaw
	From      Contact
	FromUser  bool
	Provider  string
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

type Image struct {
	URL string
	Width int
	Height int
}
