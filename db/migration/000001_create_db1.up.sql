CREATE TABLE "users" (
  "username" varchar PRIMARY KEY,
  "name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "phone" varchar NOT NULL,
  "password" varchar NOT NULL,
  "profile_url" varchar,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "assets" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "price" bigserial NOT NULL,
  "detail" text NOT NULL,
  "status" boolean NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "asset_contacts" (
  "id" bigserial PRIMARY KEY,
  "asset_id" bigserial NOT NULL,
  "contact_name" varchar NOT NULL,
  "contact_detail" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "asset_images" (
  "id" bigserial PRIMARY KEY,
  "asset_id" bigserial NOT NULL,
  "image_url" varchar NOT NULL
);

ALTER TABLE "assets" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");

ALTER TABLE "asset_contacts" ADD FOREIGN KEY ("asset_id") REFERENCES "assets" ("id");

ALTER TABLE "asset_images" ADD FOREIGN KEY ("asset_id") REFERENCES "assets" ("id");
