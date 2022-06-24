BEGIN;

CREATE TABLE "users" (
    "id" bigserial PRIMARY KEY,
    "first_name" varchar NOT NULL,
    "last_name" varchar NOT NULL,
    "email" varchar NOT NULL UNIQUE,
    "gender" varchar NOT NULL,
    "date_of_birth" timestamptz NOT NULL
);

CREATE INDEX ON "users" ("email");

CREATE INDEX ON "users" ("date_of_birth");

COMMIT;