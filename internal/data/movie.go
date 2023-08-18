package data

import "time"

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"` // Use the - directive, will never show the value
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`    // Add omitempty to hide in json response if empty
	Runtime   Runtime   `json:"runtime,omitempty"` // string directive will convert the response to string regardless of type
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version,omitempty"`
}
