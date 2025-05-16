-- 1. Create the enum type
CREATE TYPE user_role AS ENUM ('user', 'admin');

-- 2. Add the column using that enum type
ALTER TABLE users ADD COLUMN roles user_role NOT NULL DEFAULT 'user';
