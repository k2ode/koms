package main

import (
	"encoding/json"
	"fmt"
	"os"

	. "github.com/k2on/koms/types"
)

const (
	PATH_CACHE_CONVERSATIONS = "cache.json"
)

type providerIG struct {}

func NewProviderIG() (Provider, error) {
	return &providerIG{}, nil
}

func (provider *providerIG) GetId() string {
	return "instagram"
}

func (provider *providerIG) GetConversations() ([]ConversationRaw, error) {
	bytes, err := os.ReadFile(PATH_CACHE_CONVERSATIONS)
	if err != nil { return nil, err }

	var threads struct { Threads []Thread `json:"threads"` }
	err = json.Unmarshal(bytes, &threads)
	if err != nil { return nil, err }

	var conversations []ConversationRaw

	for _, thread := range threads.Threads {
		var label string
		if thread.Named { label = thread.Title }

		isGroupChat := len(thread.Users) != 1

		var userIds []string
		for _, user := range thread.Users {
			var userId string
			// userId = strconv.Itoa(int(user.Id))
			userId = user.Username
			userIds = append(userIds, userId)
		}

		conversation := ConversationRaw{
			Id: thread.Id,
			ParticipantIds: userIds,
			IsGroupChat: isGroupChat,
			Label: label,
			Provider: provider.GetId(),
		}
		conversations = append(conversations, conversation)
	}

	return conversations, nil
}

func (provider *providerIG) GetConversationMessages(id string) ([]MessageRaw, error) {
	// return []MessageRaw{{}}, nil
	instagram := GetInstagram()
	err := Login(instagram)
	if err != nil { return nil, err }

	items, err := FetchMessages(instagram, id, "")
	if err != nil { return nil, err }

	var messages []MessageRaw

	fmt.Println(items)

	// for _, item := range items {

	// }



	return messages, nil 
}

func (provider *providerIG) SendMessage(id string, body string) error {
	return nil
}

func (provider *providerIG) Sync() error {
	instagram := GetInstagram()
	err := Login(instagram)
	if err != nil { return err }
	inbox, err := FetchInbox(instagram)
	if err != nil { return err }

	path := PATH_CACHE_CONVERSATIONS
	str, _ := json.Marshal(inbox)
	os.WriteFile(path, str, 0644)


	fmt.Println("inbox")
	fmt.Println(inbox)

	return nil
}
type ResponseInbox struct {
	Inbox                Inbox  `json:"inbox"`
	SeqID                int64  `json:"seq_id"`
	PendingRequestsTotal int    `json:"pending_requests_total"`
	SnapshotAtMs         int64  `json:"snapshot_at_ms"`
	Status               string `json:"status"`
}

type Inbox struct {
	Threads []Thread `json:"threads"`
}

type Thread struct {
	Id             string      `json:"thread_id"`
	IdV2           string      `json:"thread_v2_id"`
	Named          bool        `json:"named"`
	Title          string      `json:"thread_title"`
	Users          []User      `json:"users"`
	LastActivityAt int64       `json:"last_activity_at"`
	Items          []InboxItem `json:"items"`
}

type User struct {
	Id       int64 `json:"pk"`
	Username string `json:"username"`
	Name     string `json:"full_name"`
}

type InboxItem struct {
	Id        string `json:"item_id"`
	UserId    int64  `json:"user_id"`
	Timestamp int64  `json:"timestamp"`
	ItemType  string `json:"item_type"`
	Text      string `json:"text"`
	Like      string `json:"like"`
	MediaType int    `json:"media_type"`
	Media struct {
		ID                   int64  `json:"id"`
		Images               Images `json:"image_versions2"`
		OriginalWidth        int    `json:"original_width"`
		OriginalHeight       int    `json:"original_height"`
		MediaType            int    `json:"media_type"`
		MediaID              int64  `json:"media_id"`
		PlaybackDurationSecs int    `json:"playback_duration_secs"`
		URLExpireAtSecs      int    `json:"url_expire_at_secs"`
		OrganicTrackingToken string `json:"organic_tracking_token"`
	} `json:"media"`
	MediaShare struct {
		Images Images `json:"image_versions2"`
		User   User   `json:"user"`
		PhotoOfUser bool `json:"photo_of_you"`
		DeletedReason int `json:"deleted_reason"`
		Caption struct {
			Text string `json:"text"`
		} `json:"caption"`
	} `json:"media_share"`
	StoryShare struct {
		Images Images `json:"image_versions"`
		User   User   `json:"user"`
	} `json:"story_share"`
}

type Images struct {
	Versions []ImageVersion `json:"candidates"`
}

type ImageVersion struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	URL    string `json:"url"`
}

type ResponseMessages struct {
	Thread Thread `json:"thread"`
	Status string `json:"string"`
}
