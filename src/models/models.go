package models

import "time"

type Relationship struct {
	ID        uint      `json:"id"`
	Subject   string    `json:"subject"`
	Relation  string    `json:"relation"`
	Object    string    `json:"object"`
	CreatedAt time.Time `json:"created_at"`
}
