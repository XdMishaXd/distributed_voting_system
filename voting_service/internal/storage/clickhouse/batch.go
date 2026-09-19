package clickhouse

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type VoteEvent struct {
	EventID  uuid.UUID
	PollID   uuid.UUID
	UserID   uuid.UUID
	OptionID uuid.UUID
	VotedAt  time.Time
}

// BatchWriter буферизует события и коммитит их пачками.
type BatchWriter struct {
	client *Client
	buf    []VoteEvent
	mu     sync.Mutex
}

func NewBatchWriter(c *Client) *BatchWriter {
	return &BatchWriter{
		client: c,
		buf:    make([]VoteEvent, 0, c.batchSize),
	}
}

// Add добавляет событие в буфер. Возвращает true, если буфер достиг
// порога и требует Flush со стороны вызывающего кода.
func (w *BatchWriter) Add(e VoteEvent) bool {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.buf = append(w.buf, e)
	return len(w.buf) >= int(w.client.batchSize)
}

// Flush коммитит накопленный буфер одной batch-вставкой.
// При ошибке буфер НЕ очищается — вызывающий код обязан не ack'ать
// соответствующие сообщения стрима и повторить Flush позже.
func (w *BatchWriter) Flush(ctx context.Context) error {
	w.mu.Lock()
	if len(w.buf) == 0 {
		w.mu.Unlock()
		return nil
	}
	pending := w.buf
	w.mu.Unlock()

	batch, err := w.client.conn.PrepareBatch(ctx, "INSERT INTO votes (event_id, poll_id, user_id, option_id, voted_at)")
	if err != nil {
		return fmt.Errorf("prepare batch: %w", err)
	}

	for i, e := range pending {
		if err := batch.Append(e.EventID, e.PollID, e.UserID, e.OptionID, e.VotedAt); err != nil {
			return fmt.Errorf("append event at index %d (event_id=%s): %w", i, e.EventID, err)
		}
	}

	if err := batch.Send(); err != nil {
		return fmt.Errorf("send batch of %d events: %w", len(pending), err)
	}

	// Очищаем буфер ТОЛЬКО после успешного Send — иначе при ошибке
	// потеряем события, которые никто не ack'нул в стриме.
	w.mu.Lock()
	w.buf = w.buf[:0]
	w.mu.Unlock()

	return nil
}
