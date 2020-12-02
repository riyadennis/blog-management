package eventsource

import (
	"context"
	"encoding/json"
	"time"
)

func (c *Config) Events(ctx context.Context) ([]*Article, error) {
	ctx, cancel := context.WithTimeout(ctx, TimeOut*time.Second)
	defer cancel()

	rows, err := c.Conn.QueryContext(
		ctx,
		"SELECT content FROM events_store GROUP BY resourceID;",
	)

	if err != nil {
		return nil, err
	}

	var articles []*Article

	defer rows.Close()
	for rows.Next() {
		var data []byte

		if err := rows.Scan(&data); err != nil {
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
