-- Создание таблицы words
CREATE TABLE words (
    "id" SERIAL PRIMARY KEY,
    "ru" TEXT[] NOT NULL,      -- Переводы на русский (массив)
    "jp" TEXT[] NOT NULL,      -- Написания на японском (массив)  
    "on" TEXT[],               -- Онъёми чтения (массив)
    "kun" TEXT[],              -- Кунъёми чтения (массив)
    "ex_jp" TEXT[],            -- Примеры на японском (массив)
    "ex_ru" TEXT[],            -- Переводы примеров (массив)
    "tags" TEXT[],             -- Теги для категоризации (массив)
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Основные индексы
CREATE INDEX idx_words_ru ON words USING GIN("ru");
CREATE INDEX idx_words_jp ON words USING GIN("jp");
CREATE INDEX idx_words_tags ON words USING GIN("tags");
CREATE INDEX idx_words_created_at ON words("created_at" DESC);

-- Частичные индексы для часто используемых запросов
CREATE INDEX idx_words_onyomi ON words USING GIN("on") WHERE "on" IS NOT NULL;
CREATE INDEX idx_words_kunyomi ON words USING GIN("kun") WHERE "kun" IS NOT NULL;

-- Триггер для автообновления updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_words_updated_at 
    BEFORE UPDATE ON words 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();