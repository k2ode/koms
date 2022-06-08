package main

import (
	"encoding/json"
	"errors"
	"os"
	"path"
	"strconv"
	"time"

	. "github.com/k2on/koms/types"
)

const (
	MEDIA_TYPE_PHOTO    = 1
	MEDIA_TYPE_VIDEO    = 2
	MEDIA_TYPE_CAROUSEL = 8
)

type InstagramResponseThreads struct {
	Viewer InstagramUser  `json:"viewer"`
	Inbox  InstagramInbox `json:"inbox"`
	Status string         `json:"status"`
}

type InstagramResponseThread struct {
	Status string          `json:"status"`
	Thread InstagramThread `json:"thread"`
}

type InstagramInbox struct {
	Threads  []InstagramThread `json:"threads"`
	HasOlder bool              `json:"has_older"`
}

type InstagramThread struct {
	HasOlder          bool                 `json:"has_older"`
	HasNewer          bool                 `json:"has_newer"`
	Pending           bool                 `json:"pending"`
	Items             []InstagramInboxItem `json:"items"`
	Canonical         bool                 `json:"canonical"`
	Id                string               `json:"thread_id"`
	IdV2              string               `json:"thread_v2_id"`
	Users             []InstagramUser      `json:"users"`
	LastActivityAt    int                  `json:"last_activity_at"`
	Muted             bool                 `json:"muted"`
	Archived          bool                 `json:"archived"`
	OldestCursor      string               `json:"oldest_cursor"`
	NewestCursor      string               `json:"newest_cursor"`
	Inviter           InstagramUser        `json:"inviter"`
	LastPermanentItem InstagramInboxItem   `json:"last_permanent_item"`
	Named             bool                 `json:"named"`
	NextCursor        string               `json:"next_cursor"`
	PrevCursor        string               `json:"prev_cursor"`
	Title             string               `json:"thread_title"`
	LeftUsers         []InstagramUser      `json:"left_users"`
	Spam              bool                 `json:"spam"`
	MentionsMuted     bool                 `json:"mentions_muted"`
	Type              string               `json:"thread_type"`
	ChatActivityMuted bool                 `json:"chat_activity_muted"`
	IsGroup           bool                 `json:"is_group"`
}

type InstagramInboxItem struct {
	Id               string                    `json:"item_id"`
	UserId           int                       `json:"user_id"`
	Timestamp        int                       `json:"timestamp"`
	Text             string                    `json:"text"`
	Reactions        InstagramReactions        `json:"reactions"`
	Type             InboxItemType             `json:"item_type"`
	StoryShare       InstagramStoryShare       `json:"story_share"`
	PlaceHolder      InstagramPlaceholder      `json:"placeholder"`
	DirectMediaShare InstagramDirectMediaShare `json:"direct_media_share"`
	MediaShare       InstagramMedia            `json:"media_share"`
	VisualMedia      InstagramVisualMedia      `json:"visual_media"`
	// FelixShare       InstagramFelixShare       `json:"felix_share"`
}

type InstagramReactions struct {
	Likes      []string         `json:"likes"`
	Emojis     []InstagramEmoji `json:"emojis"`
	LikesCount int              `json:"likes_count"`
}

type InstagramEmoji struct {
	Timestamp      int    `json:"timestamp"`
	SenderId       int    `json:"sender_id"`
	Emoji          string `json:"emoji"`
	// SuperReactType string `json:"super_react_type"`
}


type InboxItemType string
const (
	InboxItemTypePlaceholder InboxItemType = "placeholder"
	InboxItemTypeStoryShare                = "story_share"
	InboxItemTypeMediaShare                = "media_share"
	InboxItemTypeFelixShare                = "felix_share"
	InboxItemTypeRavenMedia                = "raven_media"
	InboxItemTypeText                      = "text"
)

type InstagramStoryShare struct {
	IsLinked bool   `json:"is_linked"`
	Title    string `json:"title"`
	Message  string `json:"message"`
	Reason   int    `json:"reason"`
	Text     string `json:"text"`
}

type InstagramPlaceholder struct {
	Message string `json:"message"`
}

type InstagramDirectMediaShare struct {
	Text           string         `json:"text"`
	MediaShareType string         `json:"media_share_type"`
	Media          InstagramMedia `json:"media"`
}

type InstagramMedia struct {
	TakenAt            int                      `json:"taken_at"`
	Pk                 int                      `json:"pk"`
	Id                 string                   `json:"id"`
	DeviceTimestamp    int                      `json:"device_timestamp"`
	MediaType          int                      `json:"media_type"`
	Code               string                   `json:"code"`
	IsPaidPartnership  bool                     `json:"is_paid_partnership"`
	CommentCount       int                      `json:"comment_count"`
	CarouselMediaCount int                      `json:"carousel_media_count"`
	CarouselMedia      []InstagramCarouselMedia `json:"carousel_media"`
	ImageVersions      InstagramImages          `json:"image_versions2"`
	OriginalWidth      int                      `json:"original_width"`
	OriginalHeight     int                      `json:"original_height"`
	User               InstagramUser            `json:"user"`
	LikeCount          int                      `json:"like_count"`
	HasLiked           bool                     `json:"has_liked"`
	PhotoOfYou         bool                     `json:"photo_of_you"`
	Usertags           InstagramUserTags        `json:"usertags"`
	Caption            InstagramCaption         `json:"caption"`
}

type InstagramCarouselMedia struct {
	Id             string            `json:"id"`
	MediaType      int               `json:"media_type"`
	ImageVersions  InstagramImages   `json:"image_versions2"`
	OriginalWidth  int               `json:"original_width"`
	OriginalHeight int               `json:"original_height"`
	Pk             int               `json:"pk"`
	UserTags       InstagramUserTags `json:"usertags"`
}

type InstagramImages struct {
	Candidates []InstagramCandidate `json:"candidates"`
}

type InstagramCandidate struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	URL    string `json:"url"`
}

type InstagramUserTags struct {
	In []InstagramUserTag `json:"in"`
}

type InstagramUserTag struct {
	User                  InstagramUser `json:"user"`
	Position              []float32     `json:"position"`
	StartTimeInVideoInSec float32       `json:"start_time_in_video_in_sec"`
	DurationInVideoInSec  float32       `json:"duration_in_video_in_sec"`

}

type InstagramUser struct {
	Id               int                       `json:"pk"`
	Username         string                    `json:"username"`
	FullName         string                    `json:"full_name"`
	IsPrivate        bool                      `json:"is_private"`
	ProfilePicURL    string                    `json:"profile_pic_url"`
	// FriendshipStatus InstagramFriendshipStatus `json:"friendship_status"`
}

// type InstagramFriendshipStatus struct {
// 	Following       bool `json:"following"`
// 	OutgoingRequest bool `json:"outgoing_request"`
// 	IsBestie        bool `json:"is_bestie"`
// 	IsRestricted    bool `json:"is_restricted"`
// 	IsFeedFavorite  bool `json:"is_feed_favorite"`
// }

type InstagramCaption struct {
	Text string `json:"text"`
	Type int    `json:"type"`
}

type InstagramVisualMedia struct {
	Media                InstagramMedia `json:"media"`
	UrlExpireAtSecs      int            `json:"url_expire_at_secs"`
	SeenUserIds          []int          `json:"seen_user_ids"`
	ViewMode             string         `json:"view_mode"`
	SeenCount            int            `json:"seen_count"`
	PlaybackDurationSecs int            `json:"playback_duration_secs "`
}

// type InstagramFelixShare struct {
// 	Video Instagram
// 	Text string `json:"text"`
// }

const (
	FILE_NAME_THREADS = "threads.json"
)

type providerIG struct {
	PathCache string
	CacheConversations map[string][]MessageRaw
}

func NewProviderIG(pathCache string) (Provider, error) {
	return &providerIG{ PathCache: pathCache, CacheConversations: make(map[string][]MessageRaw) }, nil
}

func (provider *providerIG) GetId() string {
	return "instagram"
}

func (provider *providerIG) GetConversations() ([]ConversationRaw, error) {
	pathThreads := path.Join(provider.PathCache, FILE_NAME_THREADS)
	bytes, err := os.ReadFile(pathThreads)
	if err != nil { return nil, err }

	var responseThreads InstagramResponseThreads
	err = json.Unmarshal(bytes, &responseThreads)
	if err != nil { return nil, err }

	status := responseThreads.Status
	if status != "ok" { return nil, errors.New("threads status: " + status) }

	var conversations []ConversationRaw
	threads := responseThreads.Inbox.Threads
	for _, thread := range threads {
		var label string
		if thread.Named { label = thread.Title }

		conversation := ConversationRaw{
			Id: thread.Id,
			ParticipantIds: GetUserIds(thread.Users),
			IsGroupChat: thread.IsGroup,
			Label: label,
			Provider: provider.GetId(),
		}

		provider.CacheConversations[thread.Id] = ParseMsgs(thread.Items)

		conversations = append(conversations, conversation)
	}

	return conversations, nil
}

func ParseMsgs(items []InstagramInboxItem) []MessageRaw {
		var msgs []MessageRaw
		for _, item := range items {
			msg := ParseInstagramInboxItem(item)
			msgs = append(msgs, msg)
		}
		return msgs
}

func GetUserIds(users []InstagramUser) []string {
	var userIds []string
	for _, user := range users {
		userId := user.Username
		userIds = append(userIds, userId)
	}
	return userIds
}

func (provider *providerIG) GetConversationMessages(id string) ([]MessageRaw, error) {
	threadPath := path.Join(provider.PathCache, id + ".json")
	bytes, err := os.ReadFile(threadPath)

	if err != nil { // THIS IS A HACK PLEASE FIX LATER
		return provider.CacheConversations[id], nil
	}

	var responseThread InstagramResponseThread
	err = json.Unmarshal(bytes, &responseThread)
	if err != nil { return nil, err }

	status := responseThread.Status
	if status != "ok" { return nil, errors.New("threads status: " + status) }

	return ParseMsgs(responseThread.Thread.Items), nil
}

func ParseInstagramInboxItem(item InstagramInboxItem) MessageRaw {
	var body string
	var url, bodyFull string
	var images []Image

	switch item.Type {
		case InboxItemTypeText:
			body = item.Text
			break
		case InboxItemTypeStoryShare:
			body = "[ " + item.StoryShare.Title + " ]"
			break
		case InboxItemTypeMediaShare:
			media := item.MediaShare
			if media.Pk == 0 { media = item.DirectMediaShare.Media }
			username := "@" + media.User.Username
			body = "[ " + username + "'s post" + " ]"

			url = "https://www.instagram.com/p/" + media.Code + "/"

			var imageURL string
			switch media.MediaType {
			case MEDIA_TYPE_PHOTO:
				imageURL = "/Users/koon/.koms/ig/media/" + media.Id + ".jpg"
				images = []Image{{
					URL: imageURL,
					Width: media.OriginalWidth,
					Height: media.OriginalHeight,
				}}
				break
			case MEDIA_TYPE_VIDEO:
				imageURL = "TODO"
			case MEDIA_TYPE_CAROUSEL:
				for _, carouselItem := range media.CarouselMedia {
					imageURL = "/Users/koon/.koms/ig/media/" + carouselItem.Id + ".jpg"
					images = append(images, Image{
						URL: imageURL,
						Width: media.OriginalWidth,
						Height: media.OriginalHeight,
					})

				}
				break
			}

			bodyFull = media.Caption.Text
			break
		case InboxItemTypePlaceholder:
			body = "[ TIKTOK ]"
			break
		default:
			body = "[ " + string(item.Type) + " ]"
			break
	}
	message := MessageRaw{
		Id: item.Id,
		From: strconv.Itoa(item.UserId),
		Body: body,
		Timestamp: time.Unix(int64(item.Timestamp), 0),
		URL: url,
		BodyFull: bodyFull,
		Images: images,
	}
	return message
}

func (provider *providerIG) SendMessage(id string, body string) error {
	return nil
}

func (provider *providerIG) Sync() error {
	return nil
}

func (img InstagramImages) GetBest() string {
	best := ""
	var mh, mw int
	for _, v := range img.Candidates {
		if v.Width > mw || v.Height > mh {
			best = v.URL
			mh, mw = v.Height, v.Width
		}
	}
	return best
}
