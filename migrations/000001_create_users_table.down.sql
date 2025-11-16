-- Migration: 001_initial_schema_down.sql
-- Description: Rollback initial database schema

-- Drop triggers first
DROP TRIGGER IF EXISTS trigger_users_updated_at ON users;
DROP TRIGGER IF EXISTS trigger_posts_updated_at ON posts;
DROP TRIGGER IF EXISTS trigger_comments_updated_at ON comments;
DROP TRIGGER IF EXISTS trigger_update_post_tsv ON posts;
DROP TRIGGER IF EXISTS trigger_update_reaction_count ON reactions;
DROP TRIGGER IF EXISTS trigger_update_post_comment_count ON comments;
DROP TRIGGER IF EXISTS trigger_update_follow_counts ON follows;
DROP TRIGGER IF EXISTS trigger_update_user_post_count ON posts;

-- Drop functions
DROP FUNCTION IF EXISTS update_updated_at();
DROP FUNCTION IF EXISTS update_post_tsv();
DROP FUNCTION IF EXISTS update_reaction_count();
DROP FUNCTION IF EXISTS update_post_comment_count();
DROP FUNCTION IF EXISTS update_follow_counts();
DROP FUNCTION IF EXISTS update_user_post_count();

-- Drop tables (reverse order due to foreign keys)
DROP TABLE IF EXISTS notifications CASCADE;
DROP TABLE IF EXISTS follows CASCADE;
DROP TABLE IF EXISTS reactions CASCADE;
DROP TABLE IF EXISTS comments CASCADE;
DROP TABLE IF EXISTS posts CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- Success message
DO $$
BEGIN
  RAISE NOTICE 'Migration 001: Schema rolled back successfully!';
END $$;