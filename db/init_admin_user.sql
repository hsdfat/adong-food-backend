-- Initialize admin user
-- Username: admin
-- Password: admin@adong
-- Email: admin@adong.com
-- Role: Admin

BEGIN;

-- Clear all existing users and related data
DELETE FROM auth_user_sessions;
DELETE FROM auth_token_pairs;
DELETE FROM master_users;

-- Insert admin user
-- Password hash for 'admin@adong' generated with bcrypt cost 10
INSERT INTO master_users (
    user_id,
    user_name,
    password,
    plain_password,
    full_name,
    role,
    email,
    phone,
    active,
    created_date,
    modified_date
) VALUES (
    '1',  -- Simple admin user ID
    'admin',
    '$2a$10$K0EKWc0uOdtRnUY2jpcEUe4mBZgAeRYgQliXuEplk48x43YmwTTiu',  -- bcrypt hash of 'admin@adong'
    'admin@adong',  -- Store plain password for fallback
    'Administrator',
    'Admin',
    'admin@adong.com',
    '',
    true,
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP
);

COMMIT;
