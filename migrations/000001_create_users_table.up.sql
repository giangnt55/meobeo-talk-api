-- Migration: 001_initial_schema_up.sql
-- Description: Initial database schema for Meobeo Talk MVP 1

-- Enable extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS unaccent;

-- =====================================================
-- 1. USERS TABLE
-- =====================================================
CREATE TABLE users (
  id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  username         VARCHAR(50) NOT NULL UNIQUE,
  email            VARCHAR(255) NOT NULL UNIQUE,
  password_hash    TEXT NOT NULL,
  display_name     VARCHAR(150),
  avatar_url       TEXT,
  bio              TEXT,
  
  -- Status
  is_active        BOOLEAN DEFAULT true,
  email_verified   BOOLEAN DEFAULT false,
  
  -- Counters (denormalized for performance)
  post_count       INTEGER DEFAULT 0,
  follower_count   INTEGER DEFAULT 0,
  following_count  INTEGER DEFAULT 0,
  
  -- Timestamps
  created_at       TIMESTAMP WITH TIME ZONE DEFAULT now(),
  updated_at       TIMESTAMP WITH TIME ZONE DEFAULT now(),
  last_seen_at     TIMESTAMP WITH TIME ZONE
);

-- Indexes for users
CREATE INDEX idx_users_username_lower ON users (lower(username));
CREATE INDEX idx_users_email_lower ON users (lower(email));
CREATE INDEX idx_users_created_at ON users (created_at DESC);
CREATE INDEX idx_users_last_seen ON users (last_seen_at DESC) WHERE is_active = true;

COMMENT ON TABLE users IS 'Platform users who can post and interact';
COMMENT ON COLUMN users.email_verified IS 'Email verification status';
COMMENT ON COLUMN users.post_count IS 'Cached count of user posts';

-- =====================================================
-- 2. POSTS TABLE
-- =====================================================
CREATE TABLE posts (
  id                UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  author_id         UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  
  -- Content
  title             TEXT,
  content           TEXT NOT NULL,
  content_preview   TEXT,
  
  -- Mood/Emotion (special for mental wellness platform)
  mood              VARCHAR(50),
  emotion_intensity INTEGER CHECK (emotion_intensity BETWEEN 1 AND 5),
  
  -- Visibility & Status
  visibility        VARCHAR(20) NOT NULL DEFAULT 'public',
  status            VARCHAR(20) NOT NULL DEFAULT 'published',
  
  -- Features
  allow_comments    BOOLEAN DEFAULT true,
  is_sensitive      BOOLEAN DEFAULT false,
  
  -- Counters
  comment_count     INTEGER DEFAULT 0,
  reaction_count    INTEGER DEFAULT 0,
  view_count        INTEGER DEFAULT 0,
  
  -- Timestamps
  created_at        TIMESTAMP WITH TIME ZONE DEFAULT now(),
  updated_at        TIMESTAMP WITH TIME ZONE DEFAULT now(),
  deleted_at        TIMESTAMP WITH TIME ZONE,
  
  -- Full-text search
  tsv               tsvector,
  
  CONSTRAINT chk_visibility CHECK (visibility IN ('public', 'followers', 'private', 'anonymous')),
  CONSTRAINT chk_status CHECK (status IN ('draft', 'published', 'archived'))
);

-- Indexes for posts
CREATE INDEX idx_posts_author_created ON posts (author_id, created_at DESC);
CREATE INDEX idx_posts_visibility_created ON posts (visibility, created_at DESC) 
  WHERE deleted_at IS NULL AND status = 'published';
CREATE INDEX idx_posts_mood ON posts (mood) WHERE mood IS NOT NULL;
CREATE INDEX idx_posts_created_at ON posts (created_at DESC) 
  WHERE deleted_at IS NULL;
CREATE INDEX idx_posts_tsv ON posts USING gin(tsv);

COMMENT ON TABLE posts IS 'User posts/thoughts/feelings';
COMMENT ON COLUMN posts.mood IS 'User mood when posting: happy, sad, anxious, grateful, etc';
COMMENT ON COLUMN posts.visibility IS 'public, followers, private, or anonymous';
COMMENT ON COLUMN posts.is_sensitive IS 'Content may contain sensitive topics';

-- =====================================================
-- 3. COMMENTS TABLE
-- =====================================================
CREATE TABLE comments (
  id             UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  post_id        UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
  author_id      UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  parent_id      UUID REFERENCES comments(id) ON DELETE CASCADE,
  
  content        TEXT NOT NULL,
  reaction_count INTEGER DEFAULT 0,
  status         VARCHAR(20) DEFAULT 'visible',
  
  created_at     TIMESTAMP WITH TIME ZONE DEFAULT now(),
  updated_at     TIMESTAMP WITH TIME ZONE DEFAULT now(),
  
  CONSTRAINT chk_comment_status CHECK (status IN ('visible', 'hidden', 'deleted'))
);

-- Indexes for comments
CREATE INDEX idx_comments_post_created ON comments (post_id, created_at DESC);
CREATE INDEX idx_comments_author ON comments (author_id, created_at DESC);
CREATE INDEX idx_comments_parent ON comments (parent_id) WHERE parent_id IS NOT NULL;

COMMENT ON TABLE comments IS 'Comments on posts (supports threading)';

-- =====================================================
-- 4. REACTIONS TABLE
-- =====================================================
CREATE TABLE reactions (
  id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  target_type VARCHAR(20) NOT NULL,
  target_id   UUID NOT NULL,
  reaction    VARCHAR(50) NOT NULL,
  created_at  TIMESTAMP WITH TIME ZONE DEFAULT now(),
  
  UNIQUE (user_id, target_type, target_id),
  CONSTRAINT chk_target_type CHECK (target_type IN ('post', 'comment')),
  CONSTRAINT chk_reaction CHECK (reaction IN ('support', 'hug', 'love', 'understand', 'relate', 'like'))
);

-- Indexes for reactions
CREATE INDEX idx_reactions_target ON reactions (target_type, target_id);
CREATE INDEX idx_reactions_user ON reactions (user_id, created_at DESC);

COMMENT ON TABLE reactions IS 'Emotional reactions to posts/comments';
COMMENT ON COLUMN reactions.reaction IS 'support, hug, love, understand, relate, like';

-- =====================================================
-- 5. FOLLOWS TABLE
-- =====================================================
CREATE TABLE follows (
  follower_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  followee_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at  TIMESTAMP WITH TIME ZONE DEFAULT now(),
  
  PRIMARY KEY (follower_id, followee_id),
  CONSTRAINT chk_no_self_follow CHECK (follower_id != followee_id)
);

-- Indexes for follows
CREATE INDEX idx_follows_followee ON follows (followee_id);
CREATE INDEX idx_follows_follower ON follows (follower_id);

COMMENT ON TABLE follows IS 'User follow relationships';

-- =====================================================
-- 6. NOTIFICATIONS TABLE
-- =====================================================
CREATE TABLE notifications (
  id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  actor_id   UUID REFERENCES users(id) ON DELETE SET NULL,
  type       VARCHAR(50) NOT NULL,
  payload    JSONB NOT NULL,
  is_read    BOOLEAN DEFAULT false,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
  
  CONSTRAINT chk_notification_type CHECK (type IN ('new_comment', 'new_reaction', 'new_follower', 'mention'))
);

-- Indexes for notifications
CREATE INDEX idx_notifications_user_unread ON notifications (user_id, is_read, created_at DESC);
CREATE INDEX idx_notifications_user_created ON notifications (user_id, created_at DESC);

COMMENT ON TABLE notifications IS 'User notifications for interactions';

-- =====================================================
-- TRIGGERS & FUNCTIONS
-- =====================================================

-- Function: Update post_count on users
CREATE OR REPLACE FUNCTION update_user_post_count()
RETURNS TRIGGER AS $$
BEGIN
  IF TG_OP = 'INSERT' THEN
    UPDATE users SET post_count = post_count + 1 WHERE id = NEW.author_id;
  ELSIF TG_OP = 'DELETE' THEN
    UPDATE users SET post_count = post_count - 1 WHERE id = OLD.author_id;
  END IF;
  RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_user_post_count
AFTER INSERT OR DELETE ON posts
FOR EACH ROW EXECUTE FUNCTION update_user_post_count();

-- Function: Update follow counts
CREATE OR REPLACE FUNCTION update_follow_counts()
RETURNS TRIGGER AS $$
BEGIN
  IF TG_OP = 'INSERT' THEN
    UPDATE users SET follower_count = follower_count + 1 WHERE id = NEW.followee_id;
    UPDATE users SET following_count = following_count + 1 WHERE id = NEW.follower_id;
  ELSIF TG_OP = 'DELETE' THEN
    UPDATE users SET follower_count = follower_count - 1 WHERE id = OLD.followee_id;
    UPDATE users SET following_count = following_count - 1 WHERE id = OLD.follower_id;
  END IF;
  RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_follow_counts
AFTER INSERT OR DELETE ON follows
FOR EACH ROW EXECUTE FUNCTION update_follow_counts();

-- Function: Update comment_count on posts
CREATE OR REPLACE FUNCTION update_post_comment_count()
RETURNS TRIGGER AS $$
BEGIN
  IF TG_OP = 'INSERT' THEN
    UPDATE posts SET comment_count = comment_count + 1 WHERE id = NEW.post_id;
  ELSIF TG_OP = 'DELETE' THEN
    UPDATE posts SET comment_count = comment_count - 1 WHERE id = OLD.post_id;
  END IF;
  RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_post_comment_count
AFTER INSERT OR DELETE ON comments
FOR EACH ROW EXECUTE FUNCTION update_post_comment_count();

-- Function: Update reaction_count (polymorphic)
CREATE OR REPLACE FUNCTION update_reaction_count()
RETURNS TRIGGER AS $$
BEGIN
  IF TG_OP = 'INSERT' THEN
    IF NEW.target_type = 'post' THEN
      UPDATE posts SET reaction_count = reaction_count + 1 WHERE id = NEW.target_id;
    ELSIF NEW.target_type = 'comment' THEN
      UPDATE comments SET reaction_count = reaction_count + 1 WHERE id = NEW.target_id;
    END IF;
  ELSIF TG_OP = 'DELETE' THEN
    IF OLD.target_type = 'post' THEN
      UPDATE posts SET reaction_count = reaction_count - 1 WHERE id = OLD.target_id;
    ELSIF OLD.target_type = 'comment' THEN
      UPDATE comments SET reaction_count = reaction_count - 1 WHERE id = OLD.target_id;
    END IF;
  END IF;
  RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_reaction_count
AFTER INSERT OR DELETE ON reactions
FOR EACH ROW EXECUTE FUNCTION update_reaction_count();

-- Function: Update posts.tsv for full-text search
CREATE OR REPLACE FUNCTION update_post_tsv()
RETURNS TRIGGER AS $$
BEGIN
  NEW.tsv := 
    setweight(to_tsvector('english', coalesce(NEW.title, '')), 'A') ||
    setweight(to_tsvector('english', coalesce(NEW.content, '')), 'B');
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_post_tsv
BEFORE INSERT OR UPDATE OF title, content ON posts
FOR EACH ROW EXECUTE FUNCTION update_post_tsv();

-- Function: Update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_users_updated_at
BEFORE UPDATE ON users
FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER trigger_posts_updated_at
BEFORE UPDATE ON posts
FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER trigger_comments_updated_at
BEFORE UPDATE ON comments
FOR EACH ROW EXECUTE FUNCTION update_updated_at();

-- =====================================================
-- INITIAL DATA / SEED
-- =====================================================

-- Insert system user (for notifications without actor)
INSERT INTO users (id, username, email, password_hash, display_name, is_active)
VALUES 
  ('00000000-0000-0000-0000-000000000001', 'system', 'system@meobeotalk.com', '', 'Meobeo Talk System', false);

-- Success message
DO $$
BEGIN
  RAISE NOTICE 'Migration 001: Initial schema created successfully!';
END $$;