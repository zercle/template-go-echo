-- Initial schema for template-go-echo

-- Create users table as example domain
CREATE TABLE IF NOT EXISTS users (
    id CHAR(36) PRIMARY KEY COMMENT 'UUIDv7 unique identifier',
    email VARCHAR(255) NOT NULL UNIQUE COMMENT 'User email address',
    name VARCHAR(255) NOT NULL COMMENT 'User full name',
    password_hash VARCHAR(255) NOT NULL COMMENT 'Bcrypt hashed password',
    is_active BOOLEAN DEFAULT TRUE COMMENT 'Account status',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation timestamp',
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Last update timestamp',
    deleted_at TIMESTAMP NULL COMMENT 'Soft delete timestamp',

    INDEX idx_email (email),
    INDEX idx_created_at (created_at),
    INDEX idx_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='User accounts';

-- Create user sessions table for JWT refresh tokens
CREATE TABLE IF NOT EXISTS user_sessions (
    id CHAR(36) PRIMARY KEY COMMENT 'UUIDv7 unique identifier',
    user_id CHAR(36) NOT NULL COMMENT 'Foreign key to users',
    refresh_token_hash VARCHAR(255) NOT NULL UNIQUE COMMENT 'Hashed refresh token',
    ip_address VARCHAR(45) COMMENT 'Client IP address',
    user_agent VARCHAR(500) COMMENT 'Client user agent',
    expires_at TIMESTAMP NOT NULL COMMENT 'Token expiration time',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT 'Creation timestamp',

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_expires_at (expires_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='User session tokens';
