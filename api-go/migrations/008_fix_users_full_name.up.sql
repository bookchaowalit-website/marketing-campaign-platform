-- Fix schema mismatch: Add full_name column and migrate data from first_name/last_name

-- Add full_name column to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS full_name VARCHAR(255);

-- Migrate existing data: combine first_name and last_name
UPDATE users 
SET full_name = TRIM(COALESCE(first_name, '') || ' ' || COALESCE(last_name, ''))
WHERE full_name IS NULL;

-- Update the default admin user's full_name
UPDATE users 
SET full_name = 'Admin User' 
WHERE email = 'admin@marketing.com' AND full_name IS NULL;
