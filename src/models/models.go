package models

import "time"

type Relationship struct {
	ID        uint      `json:"id"`
	Subject   string    `json:"subject"`
	Relation  string    `json:"relation"`
	Object    string    `json:"object"`
	CreatedAt time.Time `json:"created_at"`
}
type Object struct {
	ID        uint      `json:"id"`
	ObjectID  string    `json:"object_id"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

type NamespacePolicy struct {
	ID         uint      `json:"id"`
	ObjectType string    `json:"object_type"`
	Relation   string    `json:"relation"`
	Permission string    `json:"permission"`
	CreatedAt  time.Time `json:"created_at"`
}
