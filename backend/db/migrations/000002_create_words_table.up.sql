-- Создание таблицы слов
CREATE TABLE words (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    jp TEXT[] NOT NULL,
    ru TEXT[] NOT NULL,
    "on" TEXT[],
    kun TEXT[],
    ex_jp TEXT[],
    ex_ru TEXT[],
    tags TEXT[],
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    
    -- Ограничения на длину массивов
    CONSTRAINT chk_ex_jp_length CHECK (array_length(ex_jp, 1) <= 3),
    CONSTRAINT chk_ex_ru_length CHECK (array_length(ex_ru, 1) <= 3),
    CONSTRAINT chk_tags_length CHECK (array_length(tags, 1) <= 5)
);

-- Основные индексы для поиска
CREATE INDEX idx_words_user_id ON words(user_id);
CREATE INDEX idx_words_jp ON words USING GIN(jp);
CREATE INDEX idx_words_ru ON words USING GIN(ru);
CREATE INDEX idx_words_tags ON words USING GIN(tags);
CREATE INDEX idx_words_created_at ON words(created_at DESC);

-- Частичные индексы для часто используемых запросов
CREATE INDEX idx_words_onyomi ON words USING GIN("on") WHERE "on" IS NOT NULL;
CREATE INDEX idx_words_kunyomi ON words USING GIN(kun) WHERE kun IS NOT NULL;

-- Индекс для триграммного поиска (требует расширение pg_trgm)
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX idx_words_jp_trgm ON words USING GIN(jp gin_trgm_ops);
CREATE INDEX idx_words_ru_trgm ON words USING GIN(ru gin_trgm_ops);

-- Триггер для автоматического обновления updated_at
CREATE TRIGGER update_words_updated_at 
    BEFORE UPDATE ON words 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();