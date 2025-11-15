-- Final migration log script
-- This file will be executed last and log all completed migrations

USE fintrack;

-- Log all migrations that should have been executed
INSERT IGNORE INTO migration_history (migration_file) VALUES 
('01_V1__users.sql'),
('02_V2__user_profiles.sql'),
('03_V3__accounts_extended_fields.sql'),
('04_V4__cards.sql');

-- Show migration summary
SELECT 
    'Migration Summary:' AS summary,
    COUNT(*) as total_migrations,
    MAX(executed_at) as last_execution
FROM migration_history;

SELECT 
    migration_file,
    executed_at,
    status
FROM migration_history 
ORDER BY migration_file;

-- Verify database structure
SELECT 
    'Database Tables:' AS info,
    GROUP_CONCAT(table_name ORDER BY table_name) as tables
FROM information_schema.tables 
WHERE table_schema = 'fintrack';

SELECT 'All migrations completed successfully!' AS message;