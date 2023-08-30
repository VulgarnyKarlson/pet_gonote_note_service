-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS notes_outbox (
    id SERIAL PRIMARY KEY,
    event_id UUID,
    note_id UUID,
    user_id UUID,
    action varchar(10),
    sent BOOLEAN DEFAULT FALSE
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS notes_outbox;
-- +goose StatementEnd
