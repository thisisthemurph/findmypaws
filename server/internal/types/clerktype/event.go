package clerktype

import "encoding/json"

type HTTPRequestData struct {
	ClientIP  string `json:"client_ip"`
	UserAgent string `json:"user_agent"`
}

type EventAttributes struct {
	HTTPRequest HTTPRequestData `json:"http_request"`
}

type WebhookEvent struct {
	Data            json.RawMessage `json:"data"`
	EventAttributes EventAttributes `json:"event_attributes"`
	Object          string          `json:"object"`
	Timestamp       int64           `json:"timestamp"`
	Type            string          `json:"type"`
}

type UserDeletedWebhookEventData struct {
	Deleted bool   `json:"deleted"`
	UserId  string `json:"id"`
	Object  string `json:"object"`
}
