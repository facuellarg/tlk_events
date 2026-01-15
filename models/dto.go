package models

import (
	"challenge/utils"
)

// CreateEventRequest represents the JSON payload for creating an event
type CreateEventRequest struct {
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	StartTime   string  `json:"start_time"` // ISO 8601 format
	EndTime     string  `json:"end_time"`   // ISO 8601 format
}
type ValidationError struct {
	Message string `json:"message"`
}

const (
	MaxTitleLength = 100
)

var (
	TitleTooLong       = ValidationError{"title exceeds maximum length of 100 characters"}
	TitleEmpty         = ValidationError{"title should not be empty"}
	EndTimeBeforeStart = ValidationError{"end_time should be after start_time"}
	InvalidTimeFormat  = ValidationError{"invalid time format, expected ISO 8601 format"}
)

func (m *ValidationError) Error() string {
	return m.Message
}

func IsValid(event *CreateEventRequest) error {
	if event.Title == "" {
		return &TitleEmpty
	}

	if len(event.Title) > MaxTitleLength {
		return &TitleTooLong
	}
	startTime, err := utils.ParseTimestamp(event.StartTime)
	if err != nil {
		return &InvalidTimeFormat
	}

	endTime, err := utils.ParseTimestamp(event.EndTime)
	if err != nil {
		return &InvalidTimeFormat
	}

	if endTime.Before(startTime) {
		return &EndTimeBeforeStart
	}
	return nil
}
