CREATE TABLE "images" (
   "id" uuid DEFAULT gen_random_uuid() NOT NULL,
   "user_id" uuid NOT NULL,
   "image_path" VARCHAR(252) NOT NULL,
   "image_url" VARCHAR(255) NOT NULL,
   PRIMARY KEY ("id"),
   FOREIGN KEY ("user_id") REFERENCES "users" ("id")
);
