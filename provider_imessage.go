package main

import (
	"database/sql"
	"os"
	"path"
	"strconv"

	. "github.com/k2on/koms/types"

	// sqlite3 required to parse messages database
	_ "github.com/mattn/go-sqlite3"
)

type providerIMessage struct {
	db *sql.DB
	handles Handles
}

type Handles map[int]string

func NewProviderIMessage() (Provider, error) {
	db, err := sql.Open(
		"sqlite3",
		path.Join(os.Getenv("HOME"), "Library/Messages/chat.db?mode=ro"),
	)
	if err != nil { return nil, err }

	provider := &providerIMessage{ db, Handles{} }

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
		SELECT chat.ROWID, display_name, COALESCE(MAX(message.date),0) as last_activity, chat.style
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
		var id int
		var displayName *string
		var lastActivity int64
		var style int64

		err = rows.Scan(&id, &displayName, &lastActivity, &style)
		if err != nil { return []ConversationRaw{}, err }

		isGroupChat := style == 43

		label := ""
		if displayName != nil { label = *displayName }

		handles, err := provider.getConversationHandles(id)
		if err != nil { return []ConversationRaw{}, err }

		conversation := ConversationRaw{
			Id: strconv.Itoa(id),
			Label: label,
			IsGroupChat: isGroupChat,
			ParticipantIds: handles,
		}

		conversations = append(conversations, conversation)
	}

	return conversations, nil
}


func (provider *providerIMessage) GetConversationMessages(id string) ([]MessageRaw, error) {
	return []MessageRaw{}, nil
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

func (provider *providerIMessage) getConversationHandles(id int) ([]string, error) {
	rows, err := provider.runSQL(`SELECT handle_id FROM chat_handle_join WHERE chat_id == ?`, id)
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