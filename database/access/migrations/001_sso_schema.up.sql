-- ============================================================================
-- Evero Access Module (SSO) - Complete Database Schema
-- ============================================================================
-- This migration contains all schema changes for the SSO/Access module
-- Adapted from the standalone SSO project to work within Evero architecture
-- Table names are prefixed with "sso_" to avoid conflicts
-- ============================================================================

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================================================
-- Core Tables
-- ============================================================================

-- Users table
CREATE TABLE IF NOT EXISTS sso_users (
    id VARCHAR(100) PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    company_id VARCHAR(100),
    role VARCHAR(50) DEFAULT 'user',
    is_active BOOLEAN DEFAULT true,
    is_verified BOOLEAN DEFAULT false,
    email_verified BOOLEAN DEFAULT false,
    failed_login_attempts INTEGER DEFAULT 0,
    last_failed_login TIMESTAMP WITH TIME ZONE,
    last_login_at TIMESTAMP WITH TIME ZONE,
    last_login_ip VARCHAR(45),
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL
);

-- Companies table
CREATE TABLE IF NOT EXISTS sso_companies (
    id VARCHAR(100) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255),
    industry VARCHAR(100),
    status VARCHAR(50) DEFAULT 'active',
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL
);

-- Add foreign key for company_id in users
ALTER TABLE sso_users ADD CONSTRAINT fk_sso_users_company_id 
    FOREIGN KEY (company_id) REFERENCES sso_companies(id) ON DELETE SET NULL;

-- User-Company relationship
CREATE TABLE IF NOT EXISTS sso_user_companies (
    id VARCHAR(100) PRIMARY KEY,
    user_id VARCHAR(100) REFERENCES sso_users(id) ON DELETE CASCADE,
    company_id VARCHAR(100) REFERENCES sso_companies(id) ON DELETE CASCADE,
    role VARCHAR(50) NOT NULL,
    is_primary BOOLEAN DEFAULT false,
    created_at BIGINT NOT NULL,
    UNIQUE(user_id, company_id)
);

-- Refresh Tokens
CREATE TABLE IF NOT EXISTS sso_refresh_tokens (
    id VARCHAR(100) PRIMARY KEY,
    user_id VARCHAR(100) REFERENCES sso_users(id) ON DELETE CASCADE,
    token VARCHAR(500) UNIQUE NOT NULL,
    client_id VARCHAR(255),
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at BIGINT NOT NULL,
    revoked BOOLEAN DEFAULT false
);

-- Session Management
CREATE TABLE IF NOT EXISTS sso_sessions (
    id VARCHAR(100) PRIMARY KEY,
    user_id VARCHAR(100) REFERENCES sso_users(id) ON DELETE CASCADE,
    session_token VARCHAR(500) UNIQUE NOT NULL,
    ip_address VARCHAR(45),
    user_agent TEXT,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at BIGINT NOT NULL
);

-- Audit Log
CREATE TABLE IF NOT EXISTS sso_audit_logs (
    id VARCHAR(100) PRIMARY KEY,
    user_id VARCHAR(100) REFERENCES sso_users(id),
    action VARCHAR(100) NOT NULL,
    resource VARCHAR(100),
    details JSONB,
    ip_address VARCHAR(45),
    created_at BIGINT NOT NULL
);

-- ============================================================================
-- Security Features
-- ============================================================================

-- Account Lockouts Table
CREATE TABLE IF NOT EXISTS sso_account_lockouts (
    id VARCHAR(100) PRIMARY KEY,
    user_id VARCHAR(100) NOT NULL REFERENCES sso_users(id) ON DELETE CASCADE,
    locked_at TIMESTAMP WITH TIME ZONE NOT NULL,
    locked_until TIMESTAMP WITH TIME ZONE NOT NULL,
    failed_attempts INTEGER DEFAULT 0,
    reason TEXT NOT NULL,
    unlocked_at TIMESTAMP WITH TIME ZONE,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL
);

-- ============================================================================
-- Two-Factor Authentication
-- ============================================================================

-- Two-Factor Authentication Table
CREATE TABLE IF NOT EXISTS sso_user_two_factors (
    id VARCHAR(100) PRIMARY KEY,
    user_id VARCHAR(100) NOT NULL REFERENCES sso_users(id) ON DELETE CASCADE UNIQUE,
    method VARCHAR(20) NOT NULL DEFAULT 'totp',
    secret TEXT NOT NULL,
    phone_number VARCHAR(20),
    status VARCHAR(20) NOT NULL DEFAULT 'disabled',
    backup_codes_count INTEGER DEFAULT 0,
    verified_at TIMESTAMP WITH TIME ZONE,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL
);

-- Backup Codes for 2FA
CREATE TABLE IF NOT EXISTS sso_backup_codes (
    id VARCHAR(100) PRIMARY KEY,
    user_id VARCHAR(100) NOT NULL REFERENCES sso_users(id) ON DELETE CASCADE,
    code VARCHAR(255) NOT NULL,
    used_at TIMESTAMP WITH TIME ZONE,
    created_at BIGINT NOT NULL
);

-- ============================================================================
-- OAuth2
-- ============================================================================

-- OAuth2 Clients
CREATE TABLE IF NOT EXISTS sso_oauth_clients (
    id VARCHAR(100) PRIMARY KEY,
    client_id VARCHAR(255) NOT NULL UNIQUE,
    client_secret TEXT NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    redirect_uris TEXT NOT NULL,
    grant_types TEXT NOT NULL,
    scopes TEXT NOT NULL,
    owner_id VARCHAR(100) NOT NULL REFERENCES sso_users(id) ON DELETE CASCADE,
    logo_url TEXT,
    active BOOLEAN NOT NULL DEFAULT true,
    created_at BIGINT NOT NULL,
    updated_at BIGINT NOT NULL
);

-- OAuth2 Authorization Codes
CREATE TABLE IF NOT EXISTS sso_oauth_authorization_codes (
    id VARCHAR(100) PRIMARY KEY,
    code VARCHAR(255) NOT NULL UNIQUE,
    client_id VARCHAR(255) NOT NULL,
    user_id VARCHAR(100) NOT NULL REFERENCES sso_users(id) ON DELETE CASCADE,
    redirect_uri TEXT NOT NULL,
    scopes TEXT NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    used_at TIMESTAMP WITH TIME ZONE,
    created_at BIGINT NOT NULL
);

-- OAuth2 Access Tokens
CREATE TABLE IF NOT EXISTS sso_oauth_tokens (
    id VARCHAR(100) PRIMARY KEY,
    access_token TEXT NOT NULL UNIQUE,
    refresh_token TEXT UNIQUE,
    client_id VARCHAR(255) NOT NULL,
    user_id VARCHAR(100) NOT NULL REFERENCES sso_users(id) ON DELETE CASCADE,
    scopes TEXT NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    revoked_at TIMESTAMP WITH TIME ZONE,
    created_at BIGINT NOT NULL
);

-- ============================================================================
-- Password Management
-- ============================================================================

-- Email verification tokens
CREATE TABLE IF NOT EXISTS sso_email_verification_tokens (
    id VARCHAR(100) PRIMARY KEY,
    user_id VARCHAR(100) NOT NULL REFERENCES sso_users(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL UNIQUE,
    verified BOOLEAN NOT NULL DEFAULT FALSE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at BIGINT NOT NULL
);

-- Password reset tokens
CREATE TABLE IF NOT EXISTS sso_password_reset_tokens (
    id VARCHAR(100) PRIMARY KEY,
    user_id VARCHAR(100) NOT NULL REFERENCES sso_users(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL UNIQUE,
    used BOOLEAN NOT NULL DEFAULT FALSE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at BIGINT NOT NULL
);

-- ============================================================================
-- INDEXES for Performance
-- ============================================================================

-- Core Tables Indexes
CREATE INDEX IF NOT EXISTS idx_sso_users_email ON sso_users(email);
CREATE INDEX IF NOT EXISTS idx_sso_users_company_id ON sso_users(company_id);
CREATE INDEX IF NOT EXISTS idx_sso_users_role ON sso_users(role);
CREATE INDEX IF NOT EXISTS idx_sso_users_is_active ON sso_users(is_active);
CREATE INDEX IF NOT EXISTS idx_sso_users_email_verified ON sso_users(email_verified);
CREATE INDEX IF NOT EXISTS idx_sso_users_created_at ON sso_users(created_at);

CREATE INDEX IF NOT EXISTS idx_sso_companies_name ON sso_companies(name);
CREATE INDEX IF NOT EXISTS idx_sso_companies_status ON sso_companies(status);

CREATE INDEX IF NOT EXISTS idx_sso_refresh_tokens_user_id ON sso_refresh_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_sso_refresh_tokens_token ON sso_refresh_tokens(token);
CREATE INDEX IF NOT EXISTS idx_sso_refresh_tokens_expires_at ON sso_refresh_tokens(expires_at);

CREATE INDEX IF NOT EXISTS idx_sso_sessions_user_id ON sso_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_sso_sessions_token ON sso_sessions(session_token);
CREATE INDEX IF NOT EXISTS idx_sso_sessions_expires_at ON sso_sessions(expires_at);

CREATE INDEX IF NOT EXISTS idx_sso_user_companies_user_id ON sso_user_companies(user_id);
CREATE INDEX IF NOT EXISTS idx_sso_user_companies_company_id ON sso_user_companies(company_id);

CREATE INDEX IF NOT EXISTS idx_sso_audit_logs_user_id ON sso_audit_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_sso_audit_logs_action ON sso_audit_logs(action);
CREATE INDEX IF NOT EXISTS idx_sso_audit_logs_created_at ON sso_audit_logs(created_at);

-- Security Features Indexes
CREATE INDEX IF NOT EXISTS idx_sso_account_lockouts_user_id ON sso_account_lockouts(user_id);
CREATE INDEX IF NOT EXISTS idx_sso_account_lockouts_locked_until ON sso_account_lockouts(locked_until);

-- Enhanced Authentication Indexes
CREATE INDEX IF NOT EXISTS idx_sso_two_factor_auth_user_id ON sso_user_two_factors(user_id);
CREATE INDEX IF NOT EXISTS idx_sso_backup_codes_user_id ON sso_backup_codes(user_id);
CREATE INDEX IF NOT EXISTS idx_sso_backup_codes_used_at ON sso_backup_codes(used_at);

CREATE INDEX IF NOT EXISTS idx_sso_oauth_clients_client_id ON sso_oauth_clients(client_id);
CREATE INDEX IF NOT EXISTS idx_sso_oauth_clients_owner_id ON sso_oauth_clients(owner_id);
CREATE INDEX IF NOT EXISTS idx_sso_oauth_clients_active ON sso_oauth_clients(active);

CREATE INDEX IF NOT EXISTS idx_sso_oauth_auth_codes_code ON sso_oauth_authorization_codes(code);
CREATE INDEX IF NOT EXISTS idx_sso_oauth_auth_codes_client_id ON sso_oauth_authorization_codes(client_id);
CREATE INDEX IF NOT EXISTS idx_sso_oauth_auth_codes_user_id ON sso_oauth_authorization_codes(user_id);
CREATE INDEX IF NOT EXISTS idx_sso_oauth_auth_codes_expires_at ON sso_oauth_authorization_codes(expires_at);

CREATE INDEX IF NOT EXISTS idx_sso_oauth_tokens_access_token ON sso_oauth_tokens(access_token);
CREATE INDEX IF NOT EXISTS idx_sso_oauth_tokens_refresh_token ON sso_oauth_tokens(refresh_token);
CREATE INDEX IF NOT EXISTS idx_sso_oauth_tokens_client_id ON sso_oauth_tokens(client_id);
CREATE INDEX IF NOT EXISTS idx_sso_oauth_tokens_user_id ON sso_oauth_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_sso_oauth_tokens_expires_at ON sso_oauth_tokens(expires_at);

CREATE INDEX IF NOT EXISTS idx_sso_email_verifications_user_id ON sso_email_verification_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_sso_email_verifications_token ON sso_email_verification_tokens(token);

CREATE INDEX IF NOT EXISTS idx_sso_password_resets_user_id ON sso_password_reset_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_sso_password_resets_token ON sso_password_reset_tokens(token);

-- ============================================================================
-- TABLE COMMENTS
-- ============================================================================

COMMENT ON TABLE sso_users IS 'Stores user authentication and profile information for SSO';
COMMENT ON TABLE sso_companies IS 'Stores company/organization information';
COMMENT ON TABLE sso_user_companies IS 'Links users to companies with role information';
COMMENT ON TABLE sso_refresh_tokens IS 'Stores refresh tokens for token rotation';
COMMENT ON TABLE sso_sessions IS 'Active user sessions for security tracking';
COMMENT ON TABLE sso_audit_logs IS 'Audit trail for security and compliance';
COMMENT ON TABLE sso_account_lockouts IS 'Tracks account lockouts due to failed login attempts';
COMMENT ON TABLE sso_user_two_factors IS 'Stores 2FA configuration for users';
COMMENT ON TABLE sso_backup_codes IS 'Stores hashed backup codes for 2FA recovery';
COMMENT ON TABLE sso_oauth_clients IS 'OAuth2 client applications';
COMMENT ON TABLE sso_oauth_authorization_codes IS 'OAuth2 authorization codes (short-lived)';
COMMENT ON TABLE sso_oauth_tokens IS 'OAuth2 access and refresh tokens';
COMMENT ON TABLE sso_email_verification_tokens IS 'Stores email verification tokens';
COMMENT ON TABLE sso_password_reset_tokens IS 'Stores password reset tokens';

-- ============================================================================
-- MIGRATION COMPLETE
-- ============================================================================
