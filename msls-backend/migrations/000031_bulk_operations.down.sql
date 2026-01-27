-- Migration: 000031_bulk_operations.down.sql
-- Drops bulk operations tables

DROP TABLE IF EXISTS bulk_operation_items;
DROP TABLE IF EXISTS bulk_operations;
