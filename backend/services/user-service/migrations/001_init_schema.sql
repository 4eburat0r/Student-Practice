DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_role') THEN
        CREATE TYPE user_role AS ENUM ('admin','student','employer');
    END IF;
END$$;

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role user_role NOT NULL,
    email_verified BOOLEAN DEFAULT FALSE,
    created_dt TIMESTAMP DEFAULT NOW(),
    updated_dt TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS students (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    middle_name VARCHAR(100),
    photo_path TEXT,
    facility VARCHAR(255),
    course SMALLINT CHECK (course >= 0 AND course <= 5),
    description TEXT,
    created_dt TIMESTAMP DEFAULT NOW(),
    updated_dt TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS employers (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    company_name VARCHAR(255) NOT NULL,
    description TEXT,
    photo_path TEXT,
    website_url VARCHAR(500),
    requisites TEXT,
    created_dt TIMESTAMP DEFAULT NOW(),
    updated_dt TIMESTAMP DEFAULT NOW()
);
