package models

type Event struct {
	Type string `json:"type"`
}

type DbEvent struct {
	Type  string
	Count int
}
