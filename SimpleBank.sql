CREATE TYPE "entry_type" AS ENUM (
  'DEBIT',
  'CREDIT'
);

CREATE TABLE "users" (
  "id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "public_id" varchar UNIQUE NOT NULL,
  "is_blocked" bool DEFAULT false,
  "blocked_at" timestamptz,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "deleted_at" timestamptz,
  "firstname" varchar NOT NULL,
  "lastname" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password" varchar NOT NULL,
  "salt" varchar NOT NULL,
  "security_token" varchar(16),
  "email_confirmed" bool DEFAULT false,
  "security_token_requested_at" timestamptz
);

CREATE TABLE "accounts" (
  "id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "public_id" varchar UNIQUE NOT NULL,
  "is_blocked" bool DEFAULT false,
  "blocked_at" timestamptz,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "deleted_at" timestamptz,
  "user_id" int NOT NULL,
  "balance" bigint NOT NULL,
  "currency" varchar NOT NULL
);

CREATE TABLE "entries" (
  "id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "account_id" int NOT NULL,
  "public_id" varchar UNIQUE NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "deleted_at" timestamptz,
  "amount" bigint NOT NULL,
  "type" entry_type NOT NULL,
  "last_balance" bigint NOT NULL
);

CREATE TABLE "transfers" (
  "id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "from_account_id" int NOT NULL,
  "to_account_id" int NOT NULL,
  "public_id" varchar UNIQUE NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now()),
  "deleted_at" timestamptz,
  "amount" bigint NOT NULL
);

CREATE INDEX ON "entries" ("account_id");

CREATE INDEX ON "entries" ("type");

CREATE INDEX ON "transfers" ("from_account_id");

CREATE INDEX ON "transfers" ("to_account_id");

CREATE INDEX ON "transfers" ("from_account_id", "to_account_id");

COMMENT ON COLUMN "entries"."amount" IS 'can be negative or positive';

COMMENT ON COLUMN "transfers"."amount" IS 'must be positive';

ALTER TABLE "entries" ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers" ADD FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id");

ALTER TABLE "accounts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
