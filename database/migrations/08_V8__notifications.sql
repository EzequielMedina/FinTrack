-- Migration: 08_V8__notifications.sql
-- Description: Add notification service tables for email notifications and job tracking

-- Create job_runs table for tracking notification job executions
CREATE TABLE IF NOT EXISTS job_runs (
    id VARCHAR(36) PRIMARY KEY,
    started_at TIMESTAMP NOT NULL,
    completed_at TIMESTAMP NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'running',
    cards_found INT NOT NULL DEFAULT 0,
    emails_sent INT NOT NULL DEFAULT 0,
    errors INT NOT NULL DEFAULT 0,
    error_message TEXT NULL,
    INDEX idx_job_runs_started_at (started_at),
    INDEX idx_job_runs_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Create notification_logs table for tracking individual email notifications
CREATE TABLE IF NOT EXISTS notification_logs (
    id VARCHAR(36) PRIMARY KEY,
    job_run_id VARCHAR(36) NOT NULL,
    card_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    email VARCHAR(255) NOT NULL,
    status VARCHAR(20) NOT NULL,
    error_message TEXT NULL,
    sent_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_notification_logs_job_run_id (job_run_id),
    INDEX idx_notification_logs_card_id (card_id),
    INDEX idx_notification_logs_user_id (user_id),
    INDEX idx_notification_logs_status (status),
    INDEX idx_notification_logs_created_at (created_at),
    FOREIGN KEY (job_run_id) REFERENCES job_runs(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;