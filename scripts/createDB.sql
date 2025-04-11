CREATE TYPE STATUS AS ENUM ('draft', 'saved', 'checked', 'appeal', 'appealed');

CREATE TABLE "user" (
    id SERIAL PRIMARY KEY,
    mail VARCHAR(100) NOT NULL UNIQUE,
    nickname VARCHAR(100) NOT NULL,
    password VARCHAR(250) NOT NULL,
    is_moderator BOOLEAN DEFAULT FALSE,
    count_checks INTEGER DEFAULT 2
);

CREATE TABLE variant (
    id SERIAL PRIMARY KEY,
    variant_title TEXT,
    variant_text TEXT,
    author_position TEXT,
    is_public BOOLEAN DEFAULT FALSE
);

CREATE TABLE essay (
    id SERIAL PRIMARY KEY,
    essay_text TEXT,
    completed_at TIMESTAMP,
    status STATUS,
    is_published BOOLEAN DEFAULT FALSE,
    user_id INTEGER,
    variant_id INTEGER,
    FOREIGN KEY (user_id) REFERENCES "user"(id),
    FOREIGN KEY (variant_id) REFERENCES variant(id)
);

CREATE TABLE comment (
    id SERIAL PRIMARY KEY,
    user_id INTEGER,
    essay_id INTEGER,
    comment_text TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES "user"(id),
    FOREIGN KEY (essay_id) REFERENCES essay(id)
);

CREATE TABLE "like" (
    user_id INTEGER,
    essay_id INTEGER,
    PRIMARY KEY (user_id, essay_id),
    FOREIGN KEY (user_id) REFERENCES "user"(id),
    FOREIGN KEY (essay_id) REFERENCES essay(id)
);

CREATE TABLE result (
    id SERIAL PRIMARY KEY,
    sum_score INTEGER,
    appeal_text TEXT,
    essay_id INTEGER,
    FOREIGN KEY (essay_id) REFERENCES essay(id)
);

CREATE TABLE criteria (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL
    max_score INTEGER NOT NULL,
);

CREATE TABLE result_criteria (
    result_id INTEGER,
    criteria_id INTEGER,
    score INTEGER,
    explanation TEXT,
    PRIMARY KEY (result_id, criteria_id),
    FOREIGN KEY (result_id) REFERENCES result(id),
    FOREIGN KEY (criteria_id) REFERENCES criteria(id)
);