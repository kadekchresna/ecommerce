package model

import (
	"time"

	"github.com/google/uuid"
)

type OutboxStatusType string

const (
	OutboxStatusCreated    OutboxStatusType = "created"
	OutboxStatusInProgress OutboxStatusType = "in-progress"
	OutboxStatusFailed     OutboxStatusType = "failed"
	OutboxStatusSuccess    OutboxStatusType = "success"
)

type OutboxOrderCreatedMetaRequest struct {
	Order       Order         `json:"order"`
	OrderDetail []OrderDetail `json:"order_detail"`
	MessageID   string        `json:"message_id"`
	Action      string        `json:"action"`
}

type OutboxOrderUpdateStatusMetaRequest struct {
	Order       Order         `json:"order"`
	OrderDetail []OrderDetail `json:"order_detail"`
	MessageID   string        `json:"message_id"`
	Action      string        `json:"action"`
}

type OutboxGeneralMetaResponse struct {
	MessageID   string
	ErrorReason string
}

type Outbox struct {
	UUID       uuid.UUID        `json:"uuid"`
	Metadata   string           `json:"metadata"`
	Response   string           `json:"response"`
	Status     OutboxStatusType `json:"status"`
	Action     string           `json:"action"`
	Type       string           `json:"type"`
	Reference  string           `json:"refernce"`
	RetryCount int              `json:"retry_count"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
}
