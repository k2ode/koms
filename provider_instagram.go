package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	cryptoRand "crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"time"
	"unsafe"

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
	return []ConversationRaw{}, nil
}

func (provider *providerIG) GetConversationMessages(id string) ([]MessageRaw, error) {
	return []MessageRaw{}, nil 
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

type Instagram struct {
	username string
	password string
	deviceId string
	uuid     string
	phoneId  string
	token    string
	client   *http.Client
}

func GetInstagram() *Instagram {
	cookieJar, _ := cookiejar.New(nil)
	username := INSTAGRAM_USERNAME
	password := INSTAGRAM_PASSWORD
	hash := generateMD5Hash(username + password)
	deviceId := generateDeviceID(hash)
	uuid := generateUUID()
	phoneId := generateUUID()
	token := ""
	client := &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
			},
			Jar: cookieJar,
		}
	instagram := &Instagram{
		username,
		password,
		deviceId,
		uuid,
		phoneId,
		token,
		client,
	}
	
	return instagram
}

const (
	volatileSeed         = "12345"
	goInstaSigKeyVersion = "4"
	goInstaIGSigKey      = "c36436a942ea1dbb40d7f2d7d45280a620d991ce8c62fb4ce600f0a048c32c11"
	URL_BASE             = "https://i.instagram.com/api/v1/"
	URL_LOGIN            = "accounts/login/"
	URL_INBOX            = "direct_v2/inbox/"
	USER_AGENT           = "Instagram 107.0.0.27.121 Android (24/7.0; 380dpi; 1080x1920; OnePlus; ONEPLUS A3010; OnePlus3T; qcom; en_US)"
	FB_ANALYTICS         = "567067343352427"
	IG_CAPABILITIES      = "3brTBw=="
	CONNECTION_TYPE      = "WIFI"
)

func generateMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func generateDeviceID(seed string) string {
	hash := generateMD5Hash(seed + volatileSeed)
	return "android-" + hash[:16]
}

func generateUUID() string {
	uuid, err := newUUID()
	if err != nil {
		return "cb479ee7-a50d-49e7-8b7b-60cc1a105e22" // default value when error occurred
	}
	return uuid
}

func newUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(cryptoRand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

func generateHMAC(text, key string) string {
	hasher := hmac.New(sha256.New, []byte(key))
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func generateSignature(data string) map[string]string {
	m := make(map[string]string)
	m["ig_sig_key_version"] = goInstaSigKeyVersion
	m["signed_body"] = fmt.Sprintf(
		"%s.%s", generateHMAC(data, goInstaIGSigKey), data,
	)
	return m
}

type Method string
const (
	MethodGet Method = "GET"
	MethodPost       = "POST"
)

type RequestOptions struct {
	Endpoint string
	Payload  map[string]string
	Method Method
	Login bool
}

func SendRequest(ig Instagram, options RequestOptions) ([]byte, error) {
	method := options.Method
	connection := "keep-alive"
	baseUrl := URL_BASE
	
	var err error
	
	requestUrl, err := url.Parse(baseUrl + options.Endpoint)
	if err != nil { return nil, err }

	values := url.Values{}
	payload := bytes.NewBuffer([]byte{})

	for key, value := range options.Payload {
		values.Add(key, value)
	}

	if options.Method == MethodPost {
		payload.WriteString(values.Encode())
	} else {
		for key, value := range requestUrl.Query() {
			values.Add(key, strings.Join(value, " "))
		}
		requestUrl.RawQuery = values.Encode()
	}

	var request *http.Request
	request, err = http.NewRequest(string(method), requestUrl.String(), payload)
	if err != nil { return nil, err }

	
	request.Header.Set("Connection", connection)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	request.Header.Set("Accept-Language", "en-US")
	request.Header.Set("User-Agent", USER_AGENT)
	request.Header.Set("X-IG-App-ID", FB_ANALYTICS)
	request.Header.Set("X-IG-Capabilities", IG_CAPABILITIES)
	request.Header.Set("X-IG-Connection-Type", CONNECTION_TYPE)
	request.Header.Set("X-IG-Connection-Speed", fmt.Sprintf("%dkbps", acquireRand(1000, 3700)))
	request.Header.Set("X-IG-Bandwidth-Speed-KBPS", "-1.000")
	request.Header.Set("X-IG-Bandwidth-TotalBytes-B", "0")
	request.Header.Set("X-IG-Bandwidth-TotalTime-MS", "0")

	response, err := ig.client.Do(request)
	if err != nil { return nil, err }
	defer response.Body.Close()

	cookieUrl, _ := url.Parse(URL_BASE)
	for _, value := range ig.client.Jar.Cookies(cookieUrl) {
		if strings.Contains(value.Name, "csrftoken") {
			ig.token = value.Value
		}
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil { return nil, err }

	return body, nil
}

func acquireRand(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func Login(ig *Instagram) error {
	var err error

	result, err := json.Marshal(
		map[string]interface{}{
			"guid":                ig.uuid,
			"login_attempt_count": 0,
			"_csrftoken":          ig.token,
			"device_id":           ig.deviceId,
			"adid":                "",
			"phone_id":            ig.phoneId,
			"username":            ig.username,
			"password":            ig.password,
			"google_tokens":       "[]",
		},
	)
	if err != nil { return err }
	payload := generateSignature(b2s(result))
	body, err := SendRequest(*ig, RequestOptions{
		Endpoint: URL_LOGIN,
		Payload: payload,
		Method: MethodPost,
		Login: true,
	})
	if err != nil { return err }
	ig.password = ""
	
	fmt.Println(string(body))

	return err

}

func b2s(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
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
	Type      string `json:"type"`
	Text      string `json:"text"`
	Like      string `json:"like"`
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
	}
}

type Images struct {
	Versions []ImageVersion `json:"candidates"`
}

type ImageVersion struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	URL    string `json:"url"`
}

func FetchInbox(ig *Instagram) (Inbox, error) {
	endpoint := URL_INBOX
	payload := map[string]string{
		"persistentBadging": "true",
		"use_unified_inbox": "true",
	}
	body, err := SendRequest(
		*ig,
		RequestOptions{
			Endpoint: endpoint,
			Payload: payload,
		},
	)
	if err != nil { return Inbox{}, err }

	response := ResponseInbox{}
	err = json.Unmarshal(body, &response)
	if err != nil { return Inbox{}, err }

	return response.Inbox, nil
}
