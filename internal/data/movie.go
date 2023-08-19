package data

import (
	"time"

	"github.com/rehmanm/greenlight/internal/validator"
)

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"` // Use the - directive, will never show the value
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`    // Add omitempty to hide in json response if empty
	Runtime   Runtime   `json:"runtime,omitempty"` // string directive will convert the response to string regardless of type
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version,omitempty"`
}

func ValiateMovie(v *validator.Validator, movie *Movie) {

	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1888, "year", "must be greater than 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")

	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than genre")

	v.Check(validator.UniqueValues(movie.Genres), "genres", "must not contain duplicate values")

}