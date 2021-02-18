package eventsource

import (
	"context"
	"database/sql"
	"encoding/json"
)

// Events extracts articles from each stored event
func (c *Config) Events(ctx context.Context) ([]*Article, error) {
	ctx, cancel := context.WithTimeout(ctx, TimeOut)
	defer cancel()

	ids, err := resourceIDS(ctx, c.Conn)
	if err != nil {
		return nil, err
	}

	articles := make([]*Article, 1)

	for _, id := range ids {
		var data []byte
		row := c.Conn.QueryRowContext(
			ctx,
			"SELECT content FROM events_store where resourceID = ? ORDER BY version desc LIMIT 1",
			id,
		)

		if err := row.Scan(&data); err != nil {
			return nil, err
		}

		ar := &Article{}
		err := json.Unmarshal(data, ar)
		if err != nil {
			return nil, err
		}

		articles = append(articles, ar)
	}

	return articles, nil
}

func resourceIDS(ctx context.Context, conn *sql.DB) ([]string, error) {
	rows, err := conn.QueryContext(
		ctx,
		"SELECT DISTINCT resourceID FROM events_store;",
	)

	if err != nil {
		return nil, err
	}

	var ids []string

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}

		ids = append(ids, id)
	}

	return ids, nil
}
