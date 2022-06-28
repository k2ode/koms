package main

import (
	"database/sql"
	"os"
	"path"
	"strings"
	"time"

	. "github.com/k2on/koms/types"

	// sqlite3 required to parse messages database
	_ "github.com/mattn/go-sqlite3"
)

const cocoaUnixEpocDiff int64 = 978285600
const nanosecondsInSecond int64 = 1000000000

const (
	REACTION_TYPE_LOVE      = 2000
	REACTION_TYPE_LIKE      = 2001
	REACTION_TYPE_DISLIKE   = 2002
	REACTION_TYPE_LAUGH     = 2003
	REACTION_TYPE_EMPHASIZE = 2004
	REACTION_TYPE_QUESTION  = 2005

	REACTION_EMOJI_LOVE      = "❤️"
	REACTION_EMOJI_LIKE      = "👍"
	REACTION_EMOJI_DISLIKE   = "👎"
	REACTION_EMOJI_LAUGH     = "😂"
	REACTION_EMOJI_EMPHASIZE = "‼️"
	REACTION_EMOJI_QUESTION  = "❓"
)

const sendMessageApplescript = `
on run {msgText, handleId, serviceId}
	tell application "Messages"
		send msgText to buddy handleId of service id serviceId
	end tell
end run
`

type providerIMessage struct {
	db *sql.DB
	handles Handles
	rowIds map[string]int
}

type Handles map[int]string

func NewProviderIMessage() (Provider, error) {
	db, err := sql.Open(
		"sqlite3",
		path.Join(os.Getenv("HOME"), "Library/Messages/chat.db?mode=ro"),
	)
	if err != nil { return nil, err }

	provider := &providerIMessage{ db, Handles{}, make(map[string]int) }

	handles, err := provider.getHandles()
	if err != nil { return nil, err }
	provider.handles = handles

	return provider, nil
}

func (provider *providerIMessage) GetId() string {
	return "imessage"
}

func (provider *providerIMessage) GetConversations() ([]ConversationRaw, error) {
	rows, err := provider.runSQL(`
		SELECT chat.guid, chat.ROWID, display_name, COALESCE(MAX(message.date),0) as last_activity, chat.style
		FROM chat
		LEFT JOIN chat_message_join ON chat_message_join.chat_id = chat.ROWID
		LEFT JOIN message ON chat_message_join.message_id = message.ROWID
		GROUP BY chat.ROWID
		ORDER BY last_activity DESC
	`)
	if err != nil { return []ConversationRaw{}, err }
	defer rows.Close()


	conversations := []ConversationRaw{}

	for rows.Next() {
		var guid string
		var rowId int
		var displayName *string
		var lastActivity int64
		var style int64

		err = rows.Scan(&guid, &rowId, &displayName, &lastActivity, &style)
		if err != nil { return []ConversationRaw{}, err }

		isGroupChat := style == 43

		label := ""
		if displayName != nil { label = *displayName }

		handles, err := provider.getConversationHandles(rowId)
		if err != nil { return []ConversationRaw{}, err }

		provider.rowIds[guid] = rowId
		conversation := ConversationRaw{
			Id: guid,
			Label: label,
			IsGroupChat: isGroupChat,
			ParticipantIds: handles,
			Provider: provider.GetId(),
		}

		conversations = append(conversations, conversation)
	}

	return conversations, nil
}


func (provider *providerIMessage) GetConversationMessages(id string) ([]MessageRaw, error) {
	rowId, exists := provider.rowIds[id]
	if !exists { panic("rowId does not exist") }

	rows, err := provider.runSQL(`SELECT message.guid, message.date, message.text, message.handle_id, message.is_from_me, message.associated_message_guid, message.associated_message_type
	FROM message
	LEFT JOIN "chat_message_join" ON message.ROWID = "chat_message_join"."message_id"
	WHERE "chat_message_join"."chat_id" = ?
	ORDER BY date DESC
	LIMIT 20`, rowId)
	if err != nil { return nil, err }
	defer rows.Close()

	messages := []MessageRaw{}

	var msgPos int
	idMap := make(map[string]int)
	type MessageMeta struct { reactions []Reaction }
	metaMessages := make(map[string]MessageMeta)


	for rows.Next() {
		var id string
		var timestamp int64
		var text *string
		var handle_id int
		var from_me bool
		var associated_message_id *string
		var associated_message_type int

		err = rows.Scan(&id, &timestamp, &text, &handle_id, &from_me, &associated_message_id, &associated_message_type)
		if err != nil { return nil, err }

		var from string
		if from_me { from = "me" } else
		{ from = provider.handles[handle_id] }

		if associated_message_type >= REACTION_TYPE_LOVE {
			if associated_message_id == nil { panic("associated message type is not 0 but message guid is null") }

			messageId := extractMessageId(*associated_message_id)
			messageMeta, exists := metaMessages[messageId]
			if !exists {
				messageMeta = MessageMeta{[]Reaction{}}
			}
			messageMeta.reactions = append(messageMeta.reactions, Reaction{
				Emoji: getEmojiFromReactionType(associated_message_type),
				From: from,
			})
			metaMessages[messageId] = messageMeta

			continue
		}

		idMap[id] = msgPos
		msgPos++

		body := ""
		if text != nil { body = *text }



		message := MessageRaw{
			Id: id,
			Body: body,
			Timestamp: cocoaTimestampToTime(timestamp),
			From: from,
			Reactions: []Reaction{},
		}

		messages = append(messages, message)
	}


	for messageId, meta := range metaMessages {
		pos := idMap[messageId]
		msg := messages[pos]
		for _, reaction := range meta.reactions {
			msg.Reactions = append(msg.Reactions, reaction)
		}
		messages[pos] = msg
	}

	return messages, nil
}

func (provider *providerIMessage) SendMessage(id string, body string) error {
	
	return nil
}

func (provider *providerIMessage) Sync() error {
	return nil
}

func (provider *providerIMessage) runSQL(query string, args ...interface{}) (*sql.Rows, error) {
	if len(args) == 0 { return provider.db.Query(query) }

	stmt, err := provider.db.Prepare(query)
	if err != nil { return nil, err }
	defer stmt.Close()
	
	return stmt.Query(args...)
}

func (provider *providerIMessage) getHandles() (Handles, error) {
	rows, err := provider.runSQL(`SELECT ROWID, id FROM handle`)
	if err != nil { return nil, err }
	defer rows.Close()

	handles := Handles{}

	for rows.Next() {
		var rowId int
		var id string

		err = rows.Scan(&rowId, &id)
		if err != nil { return nil, err }

		handles[rowId] = id
	}

	return handles, nil
}

func (provider *providerIMessage) getConversationHandles(rowId int) ([]string, error) {
	rows, err := provider.runSQL(`SELECT handle_id FROM chat_handle_join WHERE chat_id == ?`, rowId)
	if err != nil { return nil, err }
	defer rows.Close()

	handles := []string{}

	for rows.Next() {
		var id int

		err = rows.Scan(&id)
		if err != nil { return nil, err }

		handle, exists := provider.handles[id]
		if !exists { handle = "unknown" }
		handles = append(handles, handle)
	}

	return handles, nil
}


func cocoaTimestampToTime(timestamp int64) time.Time {
	if timestamp > 1000000000000 {
		// If timestamp is bigger than 1000000000000 we can safely assume it's in nanoseconds
		// Older versions of macos use seconds, newer use nanoseconds
		timestamp = timestamp / nanosecondsInSecond
	}

	return time.Unix(timestamp+cocoaUnixEpocDiff, 0)
}

func getEmojiFromReactionType(reaction int) string {
	reactionEmojiMap := map[int]string{
		REACTION_TYPE_LOVE: REACTION_EMOJI_LOVE,
		REACTION_TYPE_LIKE: REACTION_EMOJI_LIKE,
		REACTION_TYPE_DISLIKE: REACTION_EMOJI_DISLIKE,
		REACTION_TYPE_LAUGH: REACTION_EMOJI_LAUGH,
		REACTION_TYPE_EMPHASIZE: REACTION_EMOJI_EMPHASIZE,
		REACTION_TYPE_QUESTION: REACTION_EMOJI_QUESTION,
	}
	emoji, exists := reactionEmojiMap[reaction]
	if !exists { return "?" }
	return emoji
}

func extractMessageId(id string) string {
	messageId := strings.Split(id, "/")
	if len(messageId) == 2 { return messageId[1] }
	messageId = strings.Split(id, ":")
	if len(messageId) == 2 { return messageId[1] }
	return id
}
