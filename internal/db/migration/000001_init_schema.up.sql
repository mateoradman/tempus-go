CREATE TABLE "roles" (
  "id" bigserial PRIMARY KEY,
  "role" int UNIQUE NOT NULL,
  "name" varchar(255) UNIQUE NOT NULL
);

-- Default roles
INSERT INTO roles (role, name) VALUES (1, 'Superuser');
INSERT INTO roles (role, name) VALUES (2, 'Admin');
INSERT INTO roles (role, name) VALUES (3, 'Company admin');
INSERT INTO roles (role, name) VALUES (4, 'Team admin');
INSERT INTO roles (role, name) VALUES (5, 'Normal user');

CREATE TABLE "teams" (
  "id" bigserial PRIMARY KEY,
  "name" varchar(255) NOT NULL,
  "manager_id" bigint DEFAULT null,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp DEFAULT null
);

CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "role" int NOT NULL DEFAULT 5,
  "username" varchar(255) UNIQUE NOT NULL,
  "email" varchar(255) UNIQUE NOT NULL,
  "name" varchar(255) NOT NULL,
  "surname" varchar(255) NOT NULL,
  "company_id" bigint DEFAULT null,
  "password" varchar(255) NOT NULL,
  "gender" varchar(255) DEFAULT 'unknown',
  "birth_date" date NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp DEFAULT null,
  "language" varchar(2),
  "country" varchar(2),
  "timezone" varchar(255),
  "manager_id" bigint DEFAULT null,
  "team_id" bigint DEFAULT null
);

CREATE TABLE "absences" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "start_time" timestamp NOT NULL,
  "end_time" timestamp DEFAULT null,
  "reason" varchar(255) NOT NULL,
  "paid" boolean NOT NULL DEFAULT true,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp DEFAULT null,
  "approved_by_id" bigint DEFAULT null
);

CREATE TABLE "entries" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "start_time" timestamp NOT NULL,
  "end_time" timestamp DEFAULT null,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp DEFAULT null
);

CREATE TABLE "companies" (
  "id" BIGSERIAL PRIMARY KEY,
  "name" varchar(255) NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp DEFAULT null
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

-- Trigers for saving the updated_at time 

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