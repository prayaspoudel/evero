-- Rollback migration for SSO Access Module
-- This will drop all tables created in the up migration

DROP TABLE IF EXISTS sso_password_reset_tokens CASCADE;
DROP TABLE IF EXISTS sso_email_verification_tokens CASCADE;
DROP TABLE IF EXISTS sso_oauth_tokens CASCADE;
DROP TABLE IF NOT EXISTS sso_oauth_authorization_codes CASCADE;
DROP TABLE IF EXISTS sso_oauth_clients CASCADE;
DROP TABLE IF EXISTS sso_backup_codes CASCADE;
DROP TABLE IF EXISTS sso_user_two_factors CASCADE;
DROP TABLE IF EXISTS sso_account_lockouts CASCADE;
DROP TABLE IF EXISTS sso_audit_logs CASCADE;
DROP TABLE IF EXISTS sso_sessions CASCADE;
DROP TABLE IF EXISTS sso_refresh_tokens CASCADE;
DROP TABLE IF EXISTS sso_user_companies CASCADE;
DROP TABLE IF EXISTS sso_users CASCADE;
DROP TABLE IF EXISTS sso_companies CASCADE;
