package zoomus

import (
	"os"
	"testing"
)

func getEnvs() (string, string) {
	u := os.Getenv("ZOOM_WEBHOOK_URL")
	token := os.Getenv("ZOOM_TOKEN")
	return u, token
}

func TestSendMessage(t *testing.T) {
	webhook, token := getEnvs()
	zoom, err := NewClient(webhook, token)
	if err != nil {
		t.Fatal(err)
	}
	msg := Message{
		Title:   "this is title",
		Summary: "this is summary",
		Body:    "this is body",
	}
	err = zoom.SendMessage(&msg)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMakeJSONMassage(t *testing.T) {
	tests := []Message{
		Message{"this is title", "this is summary", "this is body", "send"},
		Message{"", "", "", ""},
		Message{Title: "this is title", Summary: "this is summary", Body: "this is body"},
	}
	for _, test := range tests {
		_, err := makeJSONMassage(&test)
		if err != nil {
			t.Fatal(err)
		}
	}
}
