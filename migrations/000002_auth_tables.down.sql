-- Down migration for auth tables

-- Drop indexes
DROP INDEX IF EXISTS idx_email_verification_tokens_verified_at;
DROP INDEX IF EXISTS idx_email_verification_tokens_expires_at;
DROP INDEX IF EXISTS idx_email_verification_tokens_user_id;

DROP INDEX IF EXISTS idx_password_reset_tokens_used_at;
DROP INDEX IF EXISTS idx_password_reset_tokens_expires_at;
DROP INDEX IF EXISTS idx_password_reset_tokens_user_id;

DROP INDEX IF EXISTS idx_refresh_tokens_revoked_at;
DROP INDEX IF EXISTS idx_refresh_tokens_expires_at;
DROP INDEX IF EXISTS idx_refresh_tokens_user_id;

DROP INDEX IF EXISTS idx_users_email_verified;

-- Drop columns from users
ALTER TABLE users DROP COLUMN IF EXISTS last_login_at;
ALTER TABLE users DROP COLUMN IF EXISTS email_verified;

-- Drop tables
DROP TABLE IF EXISTS email_verification_tokens;
DROP TABLE IF EXISTS password_reset_tokens;
DROP TABLE IF EXISTS refresh_tokens;
