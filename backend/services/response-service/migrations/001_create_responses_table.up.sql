DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'response_status') THEN
        CREATE TYPE response_status AS ENUM ('review', 'approved', 'discarded', 'invited');
    END IF;
END$$;

CREATE TABLE IF NOT EXISTS responses (
    id          SERIAL PRIMARY KEY,
    vacancy_id  INTEGER NOT NULL, -- FK убраны: сущности в других сервисах/БД
    student_id  INTEGER NOT NULL,
    status      response_status NOT NULL,
    message     TEXT,
    created_dt  TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (vacancy_id, student_id)
);

CREATE INDEX IF NOT EXISTS idx_responses_student_id ON responses (student_id);
CREATE INDEX IF NOT EXISTS idx_responses_vacancy_id ON responses (vacancy_id);
