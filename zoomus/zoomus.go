package zoomus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Client is zoom.us HTTP client.
type Client struct {
	// WebhookURL is URL of your zoom.us for webhook.
	WebhookURL *url.URL

	HTTPClient *http.Client

	// Header should have following key:
	// - "Content-Type": this is for HTTP header,
	//	 this should be "application/json".
	// - "X-Zoom-Token": this is necessary for request to zoom.us.
	Header map[string]string
}

// Message represents message which wiil be sent to your zoom room.
type Message struct {
	Title   string `json:"title"`
	Summary string `json:"summary"`
	Body    string `json:"body"`
	Action  string `json:"action"`
}

var (
	msgFormat    = "<p><label>%s</label></p>"
	actionFormat = map[string]string{
		"send": "<p><button onclick=\"sendMsg('1', %s)\">send</button></p>",
		"copy": "<p><button onclick=\"copyMsg('1', %s)\">copy</button></p>",
	}
)

// NewClient initialize Client with given webhook URL and Token.
func NewClient(webhook, token string) (*Client, error) {
	if len(token) == 0 {
		return nil, fmt.Errorf("token is missing")
	}

	p, err := url.Parse(webhook)
	if err != nil {
		return nil, err
	}
	c := &Client{
		WebhookURL: p,
		HTTPClient: &http.Client{},
		Header: map[string]string{
			"Content-Type": "application/json",
			"X-Zoom-Token": token,
		},
	}

	return c, nil
}

func makeJSONMassage(msg *Message) ([]byte, error) {
	newMsg := Message{
		Title:   fmt.Sprintf(msgFormat, msg.Title),
		Summary: fmt.Sprintf(msgFormat, msg.Title),
		Body:    fmt.Sprintf(msgFormat, msg.Body),
		Action:  fmt.Sprintf(actionFormat[msg.Action], msg.Action),
	}
	msgJSON, err := json.Marshal(&newMsg)
	if err != nil {
		return nil, err
	}
	return msgJSON, nil
}

func (c *Client) SendMessage(msg *Message) error {
	msgJSON, err := makeJSONMassage(msg)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", c.WebhookURL.String(), bytes.NewBuffer(msgJSON))
	for k, v := range c.Header {
		req.Header.Set(k, v)
	}
	res, err := c.HTTPClient.Do(req)
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("http request failed: status: %s: url=%s", res.Status, c.WebhookURL.String())
	}
	return nil
}
