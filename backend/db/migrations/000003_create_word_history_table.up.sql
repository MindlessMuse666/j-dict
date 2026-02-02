CREATE TABLE word_history (
    id SERIAL PRIMARY KEY,
    word_id INTEGER NOT NULL REFERENCES words(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    action VARCHAR(50) NOT NULL,
    snapshot JSONB NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_word_history_word_id ON word_history(word_id);
CREATE INDEX idx_word_history_user_id ON word_history(user_id);
CREATE INDEX idx_word_history_created_at ON word_history(created_at DESC);
