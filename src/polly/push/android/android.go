package android

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var GCMSendApi = "https://android.googleapis.com/gcm/send"

type Client struct {
	key  string
	http *http.Client
}

// Create a client with key. Grab the key from [Google APIs Console](https://code.google.com/apis/console)
func New(key string) *Client {
	return &Client{
		key:  key,
		http: new(http.Client),
	}
}

func (c *Client) Send(message *Message) (*Response, error) {
	j, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("POST", GCMSendApi, bytes.NewBuffer(j))
	if err != nil {
		return nil, err
	}
	request.Header.Add("Authorization", fmt.Sprintf("key=%s", c.key))
	request.Header.Add("Content-Type", "application/json")

	resp, err := c.http.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return responseReply(resp)
}

func httpClientWithoutSecureVerify() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	return &http.Client{Transport: tr}
}

func responseReply(resp *http.Response) (*Response, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("%s %s", resp.Status, string(body))
	}

	ret := new(Response)
	err = json.Unmarshal(body, ret)
	return ret, err
}

// The field meaning explained at [GCM Architectural Overview](http://developer.android.com/guide/google/gcm/gcm.html#send-msg)
type Message struct {
	RegistrationIDs []string          `json:"registration_ids"`
	CollapseKey     string            `json:"collapse_key,omitempty"`
	DelayWhileIdle  bool              `json:"delay_while_idle,omitempty"`
	Data            map[string]string `json:"data,omitempty"`
	TimeToLive      int               `json:"time_to_live,omitempty"`
}

func NewMessage(ids ...string) *Message {
	return &Message{
		RegistrationIDs: ids,
		Data:            make(map[string]string),
	}
}

func (m *Message) AddRecipient(ids ...string) {
	m.RegistrationIDs = append(m.RegistrationIDs, ids...)
}

func (m *Message) SetPayload(key string, value string) {
	if m.Data == nil {
		m.Data = make(map[string]string)
	}
	m.Data[key] = value
}

// The field meaning explained at [GCM Architectural Overview](http://developer.android.com/guide/google/gcm/gcm.html#send-msg)
type Response struct {
	MulticastID  int64 `json:"multicast_id"`
	Success      int   `json:"success"`
	Failure      int   `json:"failure"`
	CanonicalIDs int   `json:"canonical_ids"`
	Results      []struct {
		MessageID      string `json:"message_id"`
		RegistrationID string `json:"registration_id"`
		Error          string `json:"error"`
	} `json:"results"`
}

// Return the indexes of succeed sent registration ids
func (r *Response) SuccessIndexes() []int {
	ret := make([]int, 0, r.Success)
	for i, result := range r.Results {
		if result.Error == "" {
			ret = append(ret, i)
		}
	}
	return ret
}

// Return the indexes of failed sent registration ids
func (r *Response) ErrorIndexes() []int {
	ret := make([]int, 0, r.Failure)
	for i, result := range r.Results {
		if result.Error != "" {
			ret = append(ret, i)
		}
	}
	return ret
}

// Return the indexes of registration ids which need update
func (r *Response) RefreshIndexes() []int {
	ret := make([]int, 0, r.CanonicalIDs)
	for i, result := range r.Results {
		if result.RegistrationID != "" {
			ret = append(ret, i)
		}
	}
	return ret
}
