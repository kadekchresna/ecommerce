package model

import (
	"time"

	"github.com/google/uuid"
)

type InboxStatusType string

const (
	InboxStatusCreated    InboxStatusType = "created"
	InboxStatusInProgress InboxStatusType = "in-progress"
	InboxStatusFailed     InboxStatusType = "failed"
	InboxStatusSuccess    InboxStatusType = "success"
)

type InboxGeneralMetaResponse struct {
	MessageID   string
	ErrorReason string
}

type Inbox struct {
	UUID       uuid.UUID       `json:"uuid"`
	Metadata   string          `json:"metadata"`
	Response   string          `json:"response"`
	Status     InboxStatusType `json:"status"`
	Action     string          `json:"action"`
	Type       string          `json:"type"`
	Reference  string          `json:"refernce"`
	RetryCount int             `json:"retry_count"`
	CreatedAt  time.Time       `json:"created_at"`
	UpdatedAt  time.Time       `json:"updated_at"`
}
