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
	zoom, err := NewZoomClient(webhook, token)
	if err != nil {
		t.Fatal(err)
	}
	msg := Message{
		Title: "this is title",
		Summary: "this is summary",
		Body: "this is body",
	}
	err = zoom.SendMessage(&msg)
	if err != nil {
		t.Fatal(err)
	}
}