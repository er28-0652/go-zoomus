package zoomus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
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
	// Title will be shown in top of message.
	Title string `json:"title"`

	// Summary will be shown in middle of message.
	Summary string `json:"summary"`

	// Body will be shown in main message.
	Body string `json:"body"`

	// Action should be "send" or "copy".
	// if it's "send", message will be showen with send button.
	// if it's "copy", copy button will be showen.
	Action string `json:"action"`
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
	if len(webhook) == 0 {
		return nil, fmt.Errorf("webhook url is missing")
	}
	if len(token) == 0 {
		return nil, fmt.Errorf("token is missing")
	}

	p, err := url.Parse(webhook)
	if err != nil {
		return nil, errors.Wrap(err, "fail to parse url")
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

func makeJSONMassage(msg Message) ([]byte, error) {
	newMsg := Message{
		Title:   fmt.Sprintf(msgFormat, msg.Title),
		Summary: fmt.Sprintf(msgFormat, msg.Summary),
		Body:    fmt.Sprintf(msgFormat, msg.Body),
	}
	if len(msg.Action) != 0 {
		value, ok := actionFormat[msg.Action]
		if !ok {
			return nil, fmt.Errorf("invalid action is given")
		}
		newMsg.Action = fmt.Sprintf(value, msg.Action)
	} else {
		newMsg.Action = ""
	}

	msgJSON, err := json.Marshal(&newMsg)
	if err != nil {
		return nil, errors.Wrap(err, "fail to marshalize")
	}
	return msgJSON, nil
}

// SendMessage sends given message to your zoom room.
// if it's successful, nill will be returned.
func (c *Client) SendMessage(msg Message) error {
	// convert Message to JSON bytes
	msgJSON, err := makeJSONMassage(msg)
	if err != nil {
		return errors.Wrap(err, "fail to make json bytes")
	}

	// create new Request with given params
	req, err := http.NewRequest("POST", c.WebhookURL.String(), bytes.NewBuffer(msgJSON))
	if err != nil {
		return errors.Wrap(err, "fail to create new request")
	}

	// add necessary headers
	for k, v := range c.Header {
		req.Header.Add(k, v)
	}

	log.Printf("send request")
	// send http request
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		log.Printf("error after send request")
		return errors.Wrap(err, "fail to send request")
	}
	defer res.Body.Close()

	log.Printf("check status code")
	// check status code
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("http request failed: status: %s: url=%s", res.Status, c.WebhookURL.String())
	}
	return nil
}
