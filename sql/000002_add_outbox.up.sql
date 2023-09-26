CREATE TABLE outbox (
    id             UUID                     PRIMARY KEY,
    aggregate_type VARCHAR(255)             NOT NULL,
    aggregate_id   VARCHAR(255)             NOT NULL,
    event_type     VARCHAR(255)             NOT NULL,
    payload        BYTEA                    NOT NULL,
    created_at     TIMESTAMP with time zone NOT NULL     DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_outbox_created_at ON outbox (created_at);

CREATE TABLE outbox_lock (
    id           INTEGER                   PRIMARY KEY,
    locked       BOOLEAN                   NOT NULL DEFAULT false,
    locked_at    TIMESTAMP with time zone,
    locked_until TIMESTAMP with time zone,
    version      BIGINT                    NOT NULL
);

INSERT INTO outbox_lock (id, locked, locked_at, locked_until, version)
VALUES (1, false, null, null, 1);
