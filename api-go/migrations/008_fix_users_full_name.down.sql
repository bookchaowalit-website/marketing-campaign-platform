-- Revert: Remove full_name column
ALTER TABLE users DROP COLUMN IF EXISTS full_name;
