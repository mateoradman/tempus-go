CREATE TABLE "teams" (
  "id" BIGSERIAL PRIMARY KEY,
  "name" varchar(255) NOT NULL,
  "manager_id" bigint DEFAULT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp DEFAULT NULL
);

CREATE TABLE "users" (
  "id" BIGSERIAL PRIMARY KEY,
  "username" varchar(255) UNIQUE NOT NULL,
  "email" varchar(255) UNIQUE NOT NULL,
  "name" varchar(255) NOT NULL,
  "surname" varchar(255) NOT NULL,
  "company_id" bigint DEFAULT NULL,
  "password" varchar(255) NOT NULL,
  "gender" varchar(255) NOT NULL DEFAULT ('unknown'),
  "birth_date" date NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp DEFAULT NULL,
  "language" varchar(2) NOT NULL DEFAULT ('en'),
  "country" varchar(2) DEFAULT NULL,
  "timezone" varchar(64) NOT NULL DEFAULT ('UTC'),
  "manager_id" bigint DEFAULT NULL,
  "team_id" bigint DEFAULT NULL
);

CREATE TABLE "absences" (
  "id" BIGSERIAL PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "start_time" timestamp NOT NULL,
  "end_time" timestamp DEFAULT NULL,
  "reason" varchar(255) NOT NULL,
  "paid" boolean NOT NULL DEFAULT true,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp DEFAULT NULL,
  "approved_by_id" bigint DEFAULT NULL
);

CREATE TABLE "entries" (
  "id" BIGSERIAL PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "start_time" timestamp NOT NULL,
  "end_time" timestamp DEFAULT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp DEFAULT NULL
);

CREATE TABLE "companies" (
  "id" BIGSERIAL PRIMARY KEY,
  "name" varchar(255) NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp DEFAULT NULL
);

CREATE TABLE "sessions" (
  "id" uuid PRIMARY KEY,
  "username" varchar NOT NULL,
  "refresh_token" varchar NOT NULL,
  "user_agent" varchar NOT NULL,
  "client_ip" varchar NOT NULL,
  "is_blocked" boolean NOT NULL,
  "expires_at" timestamp NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now())
);

CREATE INDEX ON "teams" ("id");

CREATE INDEX ON "teams" ("manager_id");

CREATE INDEX ON "users" ("email");

CREATE INDEX ON "users" ("username");

CREATE INDEX ON "users" ("manager_id");

CREATE INDEX ON "users" ("team_id");

CREATE INDEX ON "absences" ("user_id");

CREATE INDEX ON "entries" ("user_id");

COMMENT ON COLUMN "users"."language" IS 'ISO-2 language code';

COMMENT ON COLUMN "users"."country" IS 'ISO-2 Country code';

COMMENT ON COLUMN "users"."timezone" IS 'Timezone name';

ALTER TABLE "teams" ADD FOREIGN KEY ("manager_id") REFERENCES "users" ("id");

ALTER TABLE "users" ADD FOREIGN KEY ("manager_id") REFERENCES "users" ("id");

ALTER TABLE "users" ADD FOREIGN KEY ("team_id") REFERENCES "teams" ("id");

ALTER TABLE "absences" ADD FOREIGN KEY ("approved_by_id") REFERENCES "users" ("id");

ALTER TABLE "users" ADD FOREIGN KEY ("company_id") REFERENCES "companies" ("id");

ALTER TABLE "entries" ADD CONSTRAINT "user_entries" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "absences" ADD CONSTRAINT "user_absences" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "sessions" ADD CONSTRAINT "user_sessions" FOREIGN KEY ("username") REFERENCES "users" ("username") ON DELETE CASCADE;

-- Triggers for saving the updated_at time 

CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON absences
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON companies
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON entries
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON teams
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();
