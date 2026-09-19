-- +goose Up
CREATE TABLE votes (
	event_id UUID,
	poll_id UUID,
	user_id UUID,
	option_id UUID,
	voted_at DateTime64(3),
	ingested_at DateTime64(3) DEFAULT now64(3)
) ENGINE = ReplacingMergeTree(ingested_at)
ORDER BY (poll_id, user_id, event_id) PARTITION BY toYYYYMM(voted_at);
CREATE TABLE vote_counts (
	poll_id UUID,
	option_id UUID,
	votes AggregateFunction(count)
) ENGINE = AggregatingMergeTree()
ORDER BY (poll_id, option_id);
CREATE MATERIALIZED VIEW votes_mv TO vote_counts AS
SELECT poll_id,
	option_id,
	countState() AS votes
FROM votes
GROUP BY poll_id,
	option_id;
-- +goose Down
DROP VIEW IF EXISTS votes_mv;
DROP TABLE IF EXISTS vote_counts;
DROP TABLE IF EXISTS votes;
