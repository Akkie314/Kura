CREATE EXTENSION IF NOT EXISTS "pgcrypto";


-- 許可されているDiscordID
CREATE TABLE allow_users (
    discord_id VARCHAR(20) PRIMARY KEY NOT NULL,
    memo TEXT
);


-- ユーザ情報
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    discord_id VARCHAR(20)  UNIQUE NOT NULL REFERENCES allow_users(discord_id) ON DELETE CASCADE,
    username VARCHAR(32) NOT NULL,
    role TEXT DEFAULT 'user',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);


-- ファイル情報
CREATE TABLE files (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    original_name VARCHAR(255) NOT NULL,
    mime_type VARCHAR(255) NOT NULL,
    has_thumbnail BOOLEAN DEFAULT false NOT NULL,
    size BIGINT NOT NULL,
    file_path TEXT NOT NULL,
    hash_sha256 CHAR(64) NOT NULL,

    width INTEGER DEFAULT -1 NOT NULL,
    height INTEGER DEFAULT -1 NOT NULL,

    delete_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);


-- アップロードセッション
CREATE TABLE upload_sessions (
    upload_id UUID PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE NOT NULL,
    original_name VARCHAR(255) NOT NULL,
    mime_type VARCHAR(255) NOT NULL,
    total_size BIGINT NOT NULL,
    stored_size BIGINT NOT NULL,
    temp_file_path TEXT NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL
);
