package postgres

import (
	"database/sql"
	"github.com/yash/snippetbox/pkg/models"
)

type SnippetModel struct {
	Pool *sql.DB
}

func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	lastInsertid := 0
	sqlQuery := `INSERT INTO snippets (title, content, created, expires) 
					VALUES ($1,$2, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP + $3 * interval '1 Day') RETURNING id`
	sqlErr := m.Pool.QueryRow(sqlQuery, title, content, expires).Scan(&lastInsertid)
	if sqlErr != nil {
		return 0, sqlErr
	}
	return lastInsertid, nil
}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {

	sqlQuery := `SELECT id, title, content, created, expires FROM snippets
				WHERE expires > CURRENT_TIMESTAMP and id = $1`
	row := m.Pool.QueryRow(sqlQuery, id)

	s := &models.Snippet{}
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	/*
		var s models.Snippet defining the var like this can work too just return &s */

	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	return s, nil
}

// Latest this will return the most recently created snippets limit 10
func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	var snippets []*models.Snippet

	sqlQuery := `SELECT id, title, content, created, expires FROM snippets WHERE expires > CURRENT_TIMESTAMP 
				ORDER BY created DESC LIMIT 5 `

	rows, err := m.Pool.Query(sqlQuery)
	if err != nil {
		return nil, err
	}
	//defer rows.Close should be after err handling of the above err, If not there will be panic for closing the nil set
	defer rows.Close()

	for rows.Next() {
		var snippet models.Snippet
		err := rows.Scan(
			&snippet.ID,
			&snippet.Title,
			&snippet.Content,
			&snippet.Created,
			&snippet.Expires,
		)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, &snippet)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return snippets, nil
}
