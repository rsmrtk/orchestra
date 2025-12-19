-- ============================================
-- ENUM Types
-- ============================================

CREATE TYPE user_status AS ENUM ('active', 'inactive');

CREATE TYPE user_role AS ENUM (
    'admin',
    'super_admin',
    'student',
    'user',
    'teacher',
    'sales_manager',
    'parent',
    'developer'
);

CREATE TYPE user_school_type AS ENUM ('kanapka', 'po_slovensky');

CREATE TYPE course_status AS ENUM ('valid', 'expired');

-- ============================================
-- Table: users
-- ============================================

CREATE TABLE users
(
    user_id       UUID PRIMARY KEY          DEFAULT gen_random_uuid(),
    first_name    VARCHAR          NOT NULL,
    last_name     VARCHAR          NOT NULL,
    email         VARCHAR UNIQUE,
    phone         VARCHAR,
    status        user_status      NOT NULL,
    role          user_role        NOT NULL,
    school        user_school_type NOT NULL,
    password_hash VARCHAR          NOT NULL,
    created_at    TIMESTAMPTZ      NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ      NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_users_status_role_school ON users (status, role, school);
CREATE INDEX idx_users_email ON users (email);

-- ============================================
-- Table: schools
-- ============================================

CREATE TABLE schools
(
    school_id   UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    school_name VARCHAR     NOT NULL UNIQUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ============================================
-- Table: user_schools (Many-to-Many relationship)
-- ============================================

CREATE TABLE user_schools
(
    user_id    UUID        NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
    school_id  UUID        NOT NULL REFERENCES schools (school_id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (user_id, school_id)
);

CREATE INDEX idx_user_schools_school_id ON user_schools (school_id);

-- ============================================
-- Table: courses
-- ============================================

CREATE TABLE courses
(
    course_id   UUID PRIMARY KEY     DEFAULT gen_random_uuid(),
    school_id   UUID        NOT NULL REFERENCES schools (school_id) ON DELETE CASCADE,
    course_name VARCHAR     NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_courses_school_id ON courses (school_id);

-- ============================================
-- Table: user_courses
-- ============================================

CREATE TABLE user_courses
(
    user_id    UUID          NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
    school_id   UUID          NOT NULL REFERENCES schools (school_id) ON DELETE CASCADE,
    course_id  UUID          NOT NULL REFERENCES courses (course_id) ON DELETE CASCADE,
    status     course_status NOT NULL,
    created_at TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ   NOT NULL DEFAULT NOW(),

    PRIMARY KEY (user_id, course_id)
);

CREATE INDEX idx_user_courses_status ON user_courses (status);
CREATE INDEX idx_user_courses_course_id ON user_courses (course_id);
