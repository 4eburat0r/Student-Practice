CREATE TABLE IF NOT EXISTS vacancies (
    id          BIGSERIAL PRIMARY KEY,
    employer_id BIGINT NOT NULL,
    title       VARCHAR(255) NOT NULL,
    salary      NUMERIC(12,2) NOT NULL DEFAULT 0,
    description TEXT NOT NULL,
    hours       INT NOT NULL CHECK (hours BETWEEN 2 AND 8),
    format      VARCHAR(20) NOT NULL CHECK (format IN ('onsite','remote','hybrid')),
    status      VARCHAR(20) NOT NULL DEFAULT 'inactive' CHECK (status IN ('active','inactive')),
    created_dt  TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_dt  TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_vacancies_employer_id ON vacancies (employer_id);
CREATE INDEX IF NOT EXISTS idx_vacancies_status ON vacancies (status);
CREATE INDEX IF NOT EXISTS idx_vacancies_created_dt ON vacancies (created_dt DESC);
