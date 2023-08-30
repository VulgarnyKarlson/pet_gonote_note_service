-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS notes (
    id UUID PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    user_id UUID,
    content TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE OR REPLACE TRIGGER set_updated_at
BEFORE UPDATE ON notes
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_updated_at();

CREATE INDEX idx_notes_user_id ON notes(id, user_id);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS notes;
-- +goose StatementEnd
