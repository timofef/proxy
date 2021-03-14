CREATE TABLE requests
(
    id      BIGSERIAL PRIMARY KEY,
    host    VARCHAR NOT NULL,
    request VARCHAR NOT NULL
);