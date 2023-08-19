CREATE TABLE "users" (
   "id" uuid DEFAULT gen_random_uuid() NOT NULL,
   "username" VARCHAR(255) NOT NULL,
   "password_hash" VARCHAR(255) NOT NULL,
   PRIMARY KEY ("id"),
   CONSTRAINT "unique_users_username" UNIQUE("username")
);

-- Insert default user admim with the password 123456 
-- (it is true only with the salt of .env-axample)
INSERT INTO public.users
(username, password_hash)
VALUES('admin', 'eb75161d3d9b41c106051f755b41c79195c8c82112e93ebecba1121ce29df389');