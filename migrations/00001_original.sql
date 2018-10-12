-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE IF NOT EXISTS users (
    id UUID NOT NULL UNIQUE,
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    salt UUID NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    emailNotifications BOOL NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS poets (
    id UUID NOT NULL UNIQUE,
    designer UUID REFERENCES users NOT NULL,
    birthDate TIMESTAMP WITH TIME ZONE NOT NULL,
    deathDate TIMESTAMP WITH TIME ZONE NOT NULL,
    name VARCHAR(255) NOT NULL UNIQUE,
    description TEXT NOT NULL,
    language VARCHAR(255) NOT NULL,
    programFileName VARCHAR(255) NOT NULL,
    parameterFileName VARCHAR(255) NOT NULL,
    parameterFileIncluded BOOL NOT NULL,
    path VARCHAR(255) NOT NULL UNIQUE,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS issues (
    id UUID NOT NULL UNIQUE,
    volume SERIAL,
    date TIMESTAMP WITH TIME ZONE NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    upcoming BOOL NOT NULL,
    latest BOOL NOT NULL,
    PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS poems (
    id UUID NOT NULL UNIQUE,
    title VARCHAR(255) NOT NULL,
    date TIMESTAMP WITH TIME ZONE,
    author UUID REFERENCES poets NOT NULL,
    content TEXT NOT NULL,
    issue UUID REFERENCES issues NOT NULL,
    score NUMERIC(1) DEFAULT 0 CONSTRAINT normalized CHECK (score >= 0 and score <= 1),
    PRIMARY KEY (id)
);


CREATE TABLE IF NOT EXISTS issue_committee_membership (
    poet UUID REFERENCES poets NOT NULL,
    issue UUID REFERENCES issues NOT NULL,
    PRIMARY KEY (poet, issue)
);

CREATE TABLE IF NOT EXISTS issue_contributions (
    poet UUID REFERENCES poets NOT NULL,
    issue UUID REFERENCES issues NOT NULL,
    PRIMARY KEY (poet, issue)
);

CREATE TABLE IF NOT EXISTS user_poem_likes (
    usr UUID REFERENCES users NOT NULL,
    poem UUID REFERENCES poems NOT NULL,
    PRIMARY KEY (usr, poem)
);

CREATE TABLE IF NOT EXISTS user_issue_likes (
    usr UUID REFERENCES users NOT NULL,
    issue UUID REFERENCES issues NOT NULL,
    PRIMARY KEY (usr, issue)
);

CREATE TABLE IF NOT EXISTS user_poet_likes (
    usr UUID REFERENCES users NOT NULL,
    poet UUID REFERENCES poets NOT NULL,
    PRIMARY KEY (usr, poet)
);
-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE issue_committee_membership;
DROP TABLE issue_contributions;
DROP TABLE user_poem_likes;
DROP TABLE user_issue_likes;
DROP TABLE user_poet_likes;
DROP TABLE poems;
DROP TABLE poets;
DROP TABLE users;
DROP TABLE issues;

