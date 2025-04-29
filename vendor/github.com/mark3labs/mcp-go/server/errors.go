package server

import (
	"errors"
)

var (
	// Common server errors
	ErrUnsupported      = errors.New("not supported")
	ErrResourceNotFound = errors.New("resource not found")
	ErrPromptNotFound   = errors.New("prompt not found")
	ErrToolNotFound     = errors.New("tool not found")

	// Session-related errors
	ErrSessionNotFound            = errors.New("session not found")
	ErrSessionExists              = errors.New("session already exists")
	ErrSessionNotInitialized      = errors.New("session not properly initialized")
	ErrSessionDoesNotSupportTools = errors.New("session does not support per-session tools")

	// Notification-related errors
	ErrNotificationNotInitialized = errors.New("notification channel not initialized")
	ErrNotificationChannelBlocked = errors.New("notification channel full or blocked")
)
