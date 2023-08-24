package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
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

type MovieModel struct {
	DB *sql.DB
}

func (m MovieModel) Insert(movie *Movie) error {

	query := `INSERT INTO movies (title, year, runtime, genres)
			VALUES($1, $2, $3, $4)
			RETURNING id, created_at, version	
	`
	args := []any{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// Use QueryRowContext() and pass the context as the first argument.
	return m.DB.QueryRowContext(ctx, query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

func (m MovieModel) Get(id int64) (*Movie, error) {

	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `SELECT  id, created_at, title, year, runtime, genres, version
			from movies 
			where id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	var movie Movie

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&movie.ID,
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		pq.Array(&movie.Genres),
		&movie.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &movie, nil
}

func (m MovieModel) Update(movie *Movie) error {
	query := `UPDATE movies
			SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1
			WHERE id = $5 and version = $6
			RETURNING version`

	args := []any{
		movie.Title,
		movie.Year,
		movie.Runtime,
		pq.Array(movie.Genres),
		movie.ID,
		movie.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&movie.Version)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return err
}

func (m MovieModel) Delete(id int64) error {

	if id < 1 {
		return nil
	}

	query := `DELETE
			from movies 
			where id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result, err := m.DB.ExecContext(ctx, query, id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
func (m MovieModel) GetAll(title string, genres []string, filters Filters) ([]*Movie, MetaData, error) {

	query := fmt.Sprintf(`SELECT count(*) OVER(), id, created_at, title, year, runtime, genres, version
			from movies  
			WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
			AND (genres @> $2 OR $2 = '{}')
			ORDER BY %s %s, id ASC
			LIMIT %d OFFSET %d
			`, filters.sortColumn(), filters.sortDirection(), filters.limit(), filters.offset())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, title, pq.Array(genres))
	if err != nil {
		return nil, MetaData{}, err
	}

	defer rows.Close()
	totalRecords := 0
	movies := []*Movie{}

	for rows.Next() {
		var movie Movie

		err := rows.Scan(&totalRecords, &movie.ID, &movie.CreatedAt, &movie.Title, &movie.Year, &movie.Runtime, pq.Array(&movie.Genres), &movie.Version)

		if err != nil {
			return nil, MetaData{}, err
		}

		movies = append(movies, &movie)
	}

	metaData := CalculateMetaData(totalRecords, filters.Page, filters.PageSize)

	return movies, metaData, nil

}
