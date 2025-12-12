-- -------------------------------------
-- ENUM TYPES
-- -------------------------------------
CREATE TYPE user_role AS ENUM ('student', 'employer', 'admin');
CREATE TYPE resume_status AS ENUM ('active', 'inactive');
CREATE TYPE vacancy_status AS ENUM ('active', 'inactive');
CREATE TYPE response_status AS ENUM ('review', 'approved', 'discarded');
CREATE TYPE practice_format AS ENUM ('online', 'offline', 'hybrid');

-- =========================
-- TABLE: users
-- =========================
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role user_role NOT NULL,
    email_verified BOOLEAN DEFAULT FALSE,
    created_dt TIMESTAMP DEFAULT NOW(),
    updated_dt TIMESTAMP DEFAULT NOW()
);

-- =========================
-- TABLE: students
-- =========================
CREATE TABLE students (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    middle_name VARCHAR(100),
    birth_date DATE NOT NULL,
    photo_path TEXT,
    facility VARCHAR(255),
    speciality VARCHAR(255),
    course SMALLINT CHECK (course >= 1 AND course <= 5),
    description TEXT,
    created_dt TIMESTAMP DEFAULT NOW(),
    updated_dt TIMESTAMP DEFAULT NOW()
);

-- =========================
-- TABLE: employers
-- =========================
CREATE TABLE employers (
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

-- =========================
-- TABLE: resumes
-- =========================
CREATE TABLE resumes (
    id SERIAL PRIMARY KEY,
    student_id INTEGER NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    summary TEXT,
    status resume_status NOT NULL,
    created_dt TIMESTAMP DEFAULT NOW(),
    updated_dt TIMESTAMP DEFAULT NOW()
);

-- =========================
-- TABLE: vacancies
-- =========================
CREATE TABLE vacancies (
    id SERIAL PRIMARY KEY,
    employer_id INTEGER NOT NULL REFERENCES employers(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    hours SMALLINT CHECK (hours >= 2 AND hours <= 8),
    format practice_format NOT NULL,
    status vacancy_status NOT NULL,
    created_dt TIMESTAMP DEFAULT NOW(),
    updated_dt TIMESTAMP DEFAULT NOW()
);

-- =========================
-- TABLE: responses
-- =========================
CREATE TABLE responses (
    id SERIAL PRIMARY KEY,
    vacancy_id INTEGER NOT NULL REFERENCES vacancies(id) ON DELETE CASCADE,
    student_id INTEGER NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    status response_status NOT NULL,
    message TEXT,
    created_dt TIMESTAMP DEFAULT NOW(),
    UNIQUE (vacancy_id, student_id) -- один отклик на вакансию
);

-- =========================
-- TABLE: reports
-- =========================
CREATE TABLE reports (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    created_dt TIMESTAMP DEFAULT NOW()
);

-- =========================
-- INDEXES
-- =========================
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_students_user_id ON students(user_id);
CREATE INDEX idx_employers_user_id ON employers(user_id);
CREATE INDEX idx_resumes_student_id ON resumes(student_id);
CREATE INDEX idx_resumes_status ON resumes(status);
CREATE INDEX idx_vacancies_employer_id ON vacancies(employer_id);
CREATE INDEX idx_vacancies_status ON vacancies(status);
CREATE INDEX idx_responses_vacancy_id ON responses(vacancy_id);
CREATE INDEX idx_responses_student_id ON responses(student_id);
CREATE INDEX idx_responses_status ON responses(status);
