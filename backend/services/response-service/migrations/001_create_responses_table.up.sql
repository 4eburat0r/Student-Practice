-- enum для статуса отклика
CREATE TYPE response_status AS ENUM ('review', 'approved', 'discarded', 'invited');

CREATE TABLE IF NOT EXISTS responses (
    id          SERIAL PRIMARY KEY,
    vacancy_id  INTEGER NOT NULL REFERENCES vacancies(id) ON DELETE CASCADE,
    student_id  INTEGER NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    status      response_status NOT NULL,
    message     TEXT,
    created_dt  TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE (vacancy_id, student_id)
);

CREATE INDEX IF NOT EXISTS idx_responses_student_id ON responses (student_id);
CREATE INDEX IF NOT EXISTS idx_responses_vacancy_id ON responses (vacancy_id);
