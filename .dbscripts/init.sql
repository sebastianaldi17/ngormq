CREATE TABLE messages (
    id INT GENERATED ALWAYS AS IDENTITY,
    hostname VARCHAR NOT NULL,
    message VARCHAR NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);