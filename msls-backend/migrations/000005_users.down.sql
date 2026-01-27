-- Migration: 000005_users.down.sql
-- Description: Remove users table

DROP TABLE IF EXISTS users CASCADE;
