package models

import "github.com/google/uuid"

type Website struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Token     string    `json:"token"`
	Pages     []Page    `json:"pages"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
}

type Page struct {
	ID           uuid.UUID `json:"id"`
	WebsiteID    uuid.UUID `json:"website_id"`
	Path         string    `json:"path"`
	VisitorCount int       `json:"visitor_count"`
	CreatedAt    string    `json:"created_at"`
	UpdatedAt    string    `json:"updated_at"`
}

type Visitor struct {
	IP        string    `json:"ip"`
	PageID    uuid.UUID `json:"page_id"`
	Agent     string    `json:"agent"`
	Timestamp string    `json:"timestamp"`
	Referrer  string    `json:"referrer"`
	UserID    string    `json:"user_id"`
}
