-- Enum for resume status
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'resume_status') THEN
        CREATE TYPE resume_status AS ENUM ('active', 'inactive');
    END IF;
END$$;

-- Main resumes table
CREATE TABLE IF NOT EXISTS resumes (
    id SERIAL PRIMARY KEY,
    student_id INTEGER NOT NULL,
    title VARCHAR(255) NOT NULL,
    summary TEXT,
    status resume_status NOT NULL DEFAULT 'inactive',
    created_dt TIMESTAMP DEFAULT NOW(),
    updated_dt TIMESTAMP DEFAULT NOW()
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_resumes_student_id ON resumes(student_id);
CREATE INDEX IF NOT EXISTS idx_resumes_status ON resumes(status);
