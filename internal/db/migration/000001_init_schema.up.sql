CREATE TABLE "teams" (
  "id" bigserial PRIMARY KEY,
  "name" varchar(255) NOT NULL,
  "manager_id" bigint DEFAULT null,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp DEFAULT null
);

CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar(255) UNIQUE NOT NULL,
  "email" varchar(255) UNIQUE NOT NULL,
  "name" varchar(255) NOT NULL,
  "surname" varchar(255) NOT NULL,
  "company_id" bigint DEFAULT null,
  "password" varchar(255) NOT NULL,
  "gender" varchar(255) NOT NULL,
  "birth_date" date NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp DEFAULT null,
  "language" varchar(2) NOT NULL DEFAULT 'en',
  "country" varchar(2) NOT NULL,
  "timezone" varchar(255),
  "manager_id" bigint DEFAULT null,
  "team_id" bigint DEFAULT null
);

CREATE TABLE "roles" (
  "id" bigserial PRIMARY KEY,
  "name" varchar(255) UNIQUE NOT NULL,
  "description" varchar NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp DEFAULT null
);

CREATE TABLE "permissions" (
  "id" bigserial PRIMARY KEY,
  "name" varchar(255) UNIQUE NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp DEFAULT null
);

CREATE TABLE "absences" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "reason" varchar(255) NOT NULL,
  "paid" boolean NOT NULL DEFAULT true,
  "date" date NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp DEFAULT null,
  "approved_by_id" bigint DEFAULT null,
  "length" float4 NOT NULL
);

CREATE TABLE "entries" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "start_time" timestamp NOT NULL DEFAULT (now()),
  "end_time" timestamp DEFAULT null,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp DEFAULT null,
  "date" date NOT NULL DEFAULT (now())
);

CREATE TABLE "companies" (
  "id" BIGSERIAL PRIMARY KEY,
  "name" varchar(255) NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp DEFAULT null
);

CREATE INDEX ON "teams" ("id");

CREATE INDEX ON "teams" ("manager_id");

CREATE INDEX ON "users" ("email");

CREATE INDEX ON "users" ("username");

CREATE INDEX ON "users" ("manager_id");

CREATE INDEX ON "users" ("team_id");

CREATE INDEX ON "roles" ("name");

CREATE INDEX ON "permissions" ("name");

CREATE INDEX ON "absences" ("user_id");

CREATE INDEX ON "entries" ("user_id");

COMMENT ON COLUMN "users"."language" IS 'ISO-2 language code';

COMMENT ON COLUMN "users"."country" IS 'ISO-2 Country code';

ALTER TABLE "teams" ADD FOREIGN KEY ("manager_id") REFERENCES "users" ("id");

ALTER TABLE "users" ADD FOREIGN KEY ("manager_id") REFERENCES "users" ("id");

ALTER TABLE "users" ADD FOREIGN KEY ("team_id") REFERENCES "teams" ("id");

ALTER TABLE "absences" ADD FOREIGN KEY ("approved_by_id") REFERENCES "users" ("id");

ALTER TABLE "users" ADD FOREIGN KEY ("company_id") REFERENCES "companies" ("id");

ALTER TABLE "entries" ADD CONSTRAINT "user_entries" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

ALTER TABLE "absences" ADD CONSTRAINT "user_absences" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;

CREATE TABLE "users_roles" (
  "users_id" bigserial NOT NULL,
  "roles_id" bigserial NOT NULL,
  PRIMARY KEY ("users_id", "roles_id")
);

ALTER TABLE "users_roles" ADD FOREIGN KEY ("users_id") REFERENCES "users" ("id");

ALTER TABLE "users_roles" ADD FOREIGN KEY ("roles_id") REFERENCES "roles" ("id");


CREATE TABLE "roles_permissions" (
  "roles_id" bigserial NOT NULL,
  "permissions_id" bigserial NOT NULL,
  PRIMARY KEY ("roles_id", "permissions_id")
);

ALTER TABLE "roles_permissions" ADD FOREIGN KEY ("roles_id") REFERENCES "roles" ("id");

ALTER TABLE "roles_permissions" ADD FOREIGN KEY ("permissions_id") REFERENCES "permissions" ("id");


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
BEFORE UPDATE ON permissions
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON roles
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();

CREATE TRIGGER set_timestamp
BEFORE UPDATE ON teams
FOR EACH ROW
EXECUTE FUNCTION trigger_set_timestamp();