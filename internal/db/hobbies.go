package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type HobbyResult struct {
	Name  string
	Count int
}

func AddHobby(ctx context.Context, pool *pgxpool.Pool, chatID int64, hobby string) (bool, error) {
	hobbyName := strings.ToLower(strings.TrimSpace(hobby))
	if hobbyName == "" {
		return false, fmt.Errorf("hobby name cannot be empty")
	}

	query := `
		INSERT INTO hobbies (chat_id, name)
		VALUES ($1, $2)
		ON CONFLICT (chat_id, name) DO NOTHING;
	`

	cmdTag, err := pool.Exec(ctx, query, chatID, hobbyName)
	if err != nil {
		return false, err
	}

	return cmdTag.RowsAffected() > 0, nil
}

// GetWeeklyResults counts all completions in the last 7 days for every active hobby.
func GetWeeklyResults(ctx context.Context, pool *pgxpool.Pool, chatID int64) ([]HobbyResult, error) {

	query := `
		SELECT h.name, COUNT(l.id) AS total
		FROM hobbies h
		LEFT JOIN hobby_logs l 
			ON h.chat_id = l.chat_id 
			AND LOWER(h.name) = LOWER(l.hobby_name) 
			AND l.logged_at >= NOW() - INTERVAL '7 days'
		WHERE h.chat_id = $1
		GROUP BY h.name
		ORDER BY h.name ASC;
	`

	rows, err := pool.Query(ctx, query, chatID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []HobbyResult
	for rows.Next() {
		var res HobbyResult
		if err := rows.Scan(&res.Name, &res.Count); err != nil {
			return nil, err
		}
		results = append(results, res)
	}

	return results, rows.Err()
}

// LogHobbyCompletion logs standard text entries (e.g., "reading") into the database
func LogHobbyCompletion(ctx context.Context, pool *pgxpool.Pool, chatID int64, hobby string) error {
	hobbyName := strings.ToLower(strings.TrimSpace(hobby))

	query := `
		INSERT INTO hobby_logs (chat_id, hobby_name)
		VALUES ($1, $2);
	`
	_, err := pool.Exec(ctx, query, chatID, hobbyName)
	return err
}

func CheckHobby(ctx context.Context, pool *pgxpool.Pool, hobby string) bool {

	hobbyName := strings.ToLower(strings.TrimSpace(hobby))
	// see if hobby
	query := `SELECT EXISTS( SELECT 1 
	FROM hobbies
	WHERE name = $1)`

	inDb, err := pool.Exec(ctx, query, hobbyName)
	fmt.Println(inDb)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Print(inDb)
	return inDb.RowsAffected() > 0

}
