package clickhouse

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

type OptionCount struct {
	OptionID uuid.UUID
	Votes    uint64
}

func (c *Client) GetPollResults(ctx context.Context, pollID uuid.UUID) ([]OptionCount, error) {
	rows, err := c.conn.Query(ctx, `
		SELECT option_id, countMerge(votes) AS total
		FROM vote_counts
		WHERE poll_id = $1
		GROUP BY option_id
	`, pollID)
	if err != nil {
		return nil, fmt.Errorf("query poll results for poll_id=%s: %w", pollID, err)
	}
	defer rows.Close()

	var results []OptionCount
	for rows.Next() {
		var oc OptionCount
		if err := rows.Scan(&oc.OptionID, &oc.Votes); err != nil {
			return nil, fmt.Errorf("scan poll result row: %w", err)
		}
		results = append(results, oc)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate poll results: %w", err)
	}

	return results, nil
}
