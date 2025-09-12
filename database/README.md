# Database - FinTrack

## 📋 Descripción

Base de datos MySQL 8.0 para la plataforma FinTrack. Contiene toda la información financiera, usuarios, transacciones, cuentas y configuraciones del sistema. Diseñada con alta disponibilidad, seguridad y rendimiento optimizado.

## 🛠️ Tecnologías

- **Motor**: MySQL 8.0
- **Charset**: utf8mb4
- **Collation**: utf8mb4_unicode_ci
- **Storage Engine**: InnoDB
- **Backup**: mysqldump, Percona XtraBackup
- **Monitoring**: MySQL Performance Schema
- **Contenedor**: Docker con configuración personalizada

## 🏗️ Arquitectura

### Estructura del Proyecto

```
database/
├── Dockerfile               # Configuración Docker
├── docker-entrypoint.sh     # Script de inicialización
├── migrations/              # Scripts de migración
│   ├── 001_initial_schema.sql
│   ├── 002_add_indexes.sql
│   └── 003_add_constraints.sql
├── schemas/                 # Esquemas de base de datos
│   ├── users.sql
│   ├── accounts.sql
│   ├── transactions.sql
│   ├── wallets.sql
│   └── notifications.sql
├── seeds/                   # Datos de prueba
│   ├── users_seed.sql
│   ├── categories_seed.sql
│   └── currencies_seed.sql
└── README.md                # Este archivo
```

## 🚀 Configuración

### Variables de Entorno

```env
# Configuración de base de datos
MYSQL_DATABASE=fintrack
MYSQL_USER=fintrack_user
MYSQL_PASSWORD=fintrack_password
MYSQL_ROOT_PASSWORD=root_password

# Configuración de conexión
DB_HOST=localhost
DB_PORT=3306
DB_CHARSET=utf8mb4
DB_COLLATION=utf8mb4_unicode_ci

# Configuración de rendimiento
INNODB_BUFFER_POOL_SIZE=1G
INNODB_LOG_FILE_SIZE=256M
MAX_CONNECTIONS=200
QUERY_CACHE_SIZE=64M

# Configuración de seguridad
SQL_MODE=STRICT_TRANS_TABLES,NO_ZERO_DATE,NO_ZERO_IN_DATE,ERROR_FOR_DIVISION_BY_ZERO
LOCAL_INFILE=0
SECURE_FILE_PRIV=/var/lib/mysql-files/
```

### Configuración MySQL (my.cnf)

```ini
[mysqld]
# Basic Settings
user = mysql
port = 3306
basedir = /usr
datadir = /var/lib/mysql
tmpdir = /tmp
lc-messages-dir = /usr/share/mysql

# Character Set
character-set-server = utf8mb4
collation-server = utf8mb4_unicode_ci

# InnoDB Settings
innodb_buffer_pool_size = 1G
innodb_log_file_size = 256M
innodb_log_buffer_size = 16M
innodb_flush_log_at_trx_commit = 1
innodb_file_per_table = 1
innodb_open_files = 400

# Connection Settings
max_connections = 200
max_connect_errors = 1000
wait_timeout = 28800
interactive_timeout = 28800

# Query Cache
query_cache_type = 1
query_cache_size = 64M
query_cache_limit = 2M

# Logging
general_log = 1
general_log_file = /var/log/mysql/general.log
log_error = /var/log/mysql/error.log
slow_query_log = 1
slow_query_log_file = /var/log/mysql/slow.log
long_query_time = 2

# Security
sql_mode = STRICT_TRANS_TABLES,NO_ZERO_DATE,NO_ZERO_IN_DATE,ERROR_FOR_DIVISION_BY_ZERO
local_infile = 0
secure_file_priv = /var/lib/mysql-files/
```

## 🐳 Docker

### Construcción y Ejecución

```bash
# Build de la imagen
docker build -t fintrack-database .

# Ejecutar contenedor
docker run -d \
  --name fintrack-db \
  -p 3306:3306 \
  -e MYSQL_ROOT_PASSWORD=root_password \
  -e MYSQL_DATABASE=fintrack \
  -e MYSQL_USER=fintrack_user \
  -e MYSQL_PASSWORD=fintrack_password \
  -v fintrack_db_data:/var/lib/mysql \
  fintrack-database

# Docker Compose
docker-compose up mysql

# Conectar a la base de datos
docker exec -it fintrack-db mysql -u fintrack_user -p fintrack
```

### Volúmenes

```yaml
volumes:
  fintrack_db_data:
    driver: local
  fintrack_db_logs:
    driver: local
  fintrack_db_backups:
    driver: local
```

## 📊 Esquema de Base de Datos

### Tablas Principales

#### Users (Usuarios)
```sql
CREATE TABLE users (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    uuid CHAR(36) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone VARCHAR(20),
    date_of_birth DATE,
    country VARCHAR(2),
    timezone VARCHAR(50) DEFAULT 'UTC',
    language VARCHAR(5) DEFAULT 'en',
    currency VARCHAR(3) DEFAULT 'USD',
    email_verified BOOLEAN DEFAULT FALSE,
    phone_verified BOOLEAN DEFAULT FALSE,
    two_factor_enabled BOOLEAN DEFAULT FALSE,
    status ENUM('active', 'inactive', 'suspended') DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    last_login_at TIMESTAMP NULL,
    
    INDEX idx_email (email),
    INDEX idx_uuid (uuid),
    INDEX idx_status (status),
    INDEX idx_created_at (created_at)
);
```

#### Accounts (Cuentas)
```sql
CREATE TABLE accounts (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    uuid CHAR(36) UNIQUE NOT NULL,
    user_id BIGINT NOT NULL,
    name VARCHAR(100) NOT NULL,
    account_type ENUM('checking', 'savings', 'credit', 'investment', 'loan') NOT NULL,
    account_number VARCHAR(50),
    bank_name VARCHAR(100),
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    balance DECIMAL(15,2) DEFAULT 0.00,
    available_balance DECIMAL(15,2) DEFAULT 0.00,
    credit_limit DECIMAL(15,2) DEFAULT 0.00,
    interest_rate DECIMAL(5,4) DEFAULT 0.0000,
    is_active BOOLEAN DEFAULT TRUE,
    is_primary BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_uuid (uuid),
    INDEX idx_account_type (account_type),
    INDEX idx_is_active (is_active)
);
```

#### Transactions (Transacciones)
```sql
CREATE TABLE transactions (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    uuid CHAR(36) UNIQUE NOT NULL,
    user_id BIGINT NOT NULL,
    account_id BIGINT NOT NULL,
    category_id BIGINT,
    transaction_type ENUM('income', 'expense', 'transfer') NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    currency VARCHAR(3) NOT NULL,
    exchange_rate DECIMAL(10,6) DEFAULT 1.000000,
    amount_usd DECIMAL(15,2) NOT NULL,
    description TEXT,
    reference_number VARCHAR(100),
    merchant_name VARCHAR(100),
    location VARCHAR(255),
    transaction_date TIMESTAMP NOT NULL,
    processed_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status ENUM('pending', 'completed', 'failed', 'cancelled') DEFAULT 'completed',
    tags JSON,
    metadata JSON,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (account_id) REFERENCES accounts(id) ON DELETE CASCADE,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE SET NULL,
    INDEX idx_user_id (user_id),
    INDEX idx_account_id (account_id),
    INDEX idx_transaction_date (transaction_date),
    INDEX idx_transaction_type (transaction_type),
    INDEX idx_status (status),
    INDEX idx_amount (amount),
    INDEX idx_currency (currency)
);
```

#### Categories (Categorías)
```sql
CREATE TABLE categories (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    uuid CHAR(36) UNIQUE NOT NULL,
    user_id BIGINT,
    parent_id BIGINT,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    color VARCHAR(7),
    icon VARCHAR(50),
    category_type ENUM('income', 'expense') NOT NULL,
    is_system BOOLEAN DEFAULT FALSE,
    is_active BOOLEAN DEFAULT TRUE,
    sort_order INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (parent_id) REFERENCES categories(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_parent_id (parent_id),
    INDEX idx_category_type (category_type),
    INDEX idx_is_system (is_system)
);
```

#### Wallets (Billeteras Crypto)
```sql
CREATE TABLE wallets (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    uuid CHAR(36) UNIQUE NOT NULL,
    user_id BIGINT NOT NULL,
    name VARCHAR(100) NOT NULL,
    wallet_type ENUM('bitcoin', 'ethereum', 'other') NOT NULL,
    address VARCHAR(255) NOT NULL,
    public_key TEXT,
    encrypted_private_key TEXT,
    balance DECIMAL(20,8) DEFAULT 0.00000000,
    balance_usd DECIMAL(15,2) DEFAULT 0.00,
    is_active BOOLEAN DEFAULT TRUE,
    last_sync_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE KEY unique_user_address (user_id, address),
    INDEX idx_user_id (user_id),
    INDEX idx_wallet_type (wallet_type),
    INDEX idx_address (address)
);
```

#### Notifications (Notificaciones)
```sql
CREATE TABLE notifications (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    uuid CHAR(36) UNIQUE NOT NULL,
    user_id BIGINT NOT NULL,
    type ENUM('email', 'sms', 'push', 'in_app') NOT NULL,
    category ENUM('transaction', 'security', 'system', 'marketing') NOT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    data JSON,
    status ENUM('pending', 'sent', 'delivered', 'failed') DEFAULT 'pending',
    priority ENUM('low', 'medium', 'high', 'urgent') DEFAULT 'medium',
    scheduled_at TIMESTAMP NULL,
    sent_at TIMESTAMP NULL,
    read_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    INDEX idx_user_id (user_id),
    INDEX idx_type (type),
    INDEX idx_status (status),
    INDEX idx_scheduled_at (scheduled_at),
    INDEX idx_created_at (created_at)
);
```

## 🔧 Migraciones

### Estructura de Migraciones

```sql
-- migrations/001_initial_schema.sql
-- Crear tablas principales
CREATE TABLE IF NOT EXISTS schema_migrations (
    version VARCHAR(255) PRIMARY KEY,
    applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- migrations/002_add_indexes.sql
-- Agregar índices para optimización
CREATE INDEX idx_transactions_date_amount ON transactions(transaction_date, amount);
CREATE INDEX idx_accounts_user_type ON accounts(user_id, account_type);

-- migrations/003_add_constraints.sql
-- Agregar restricciones adicionales
ALTER TABLE accounts ADD CONSTRAINT chk_balance_positive 
    CHECK (account_type = 'credit' OR balance >= 0);
```

### Comandos de Migración

```bash
# Aplicar migraciones
mysql -u fintrack_user -p fintrack < migrations/001_initial_schema.sql

# Verificar estado de migraciones
mysql -u fintrack_user -p fintrack -e "SELECT * FROM schema_migrations;"

# Rollback (manual)
mysql -u fintrack_user -p fintrack < rollbacks/001_rollback.sql
```

## 🌱 Seeds (Datos de Prueba)

### Categorías por Defecto

```sql
-- seeds/categories_seed.sql
INSERT INTO categories (uuid, name, category_type, color, icon, is_system) VALUES
(UUID(), 'Salary', 'income', '#4CAF50', 'work', TRUE),
(UUID(), 'Freelance', 'income', '#8BC34A', 'laptop', TRUE),
(UUID(), 'Investment', 'income', '#2196F3', 'trending_up', TRUE),
(UUID(), 'Food & Dining', 'expense', '#FF5722', 'restaurant', TRUE),
(UUID(), 'Transportation', 'expense', '#9C27B0', 'directions_car', TRUE),
(UUID(), 'Shopping', 'expense', '#E91E63', 'shopping_cart', TRUE),
(UUID(), 'Utilities', 'expense', '#607D8B', 'home', TRUE),
(UUID(), 'Healthcare', 'expense', '#F44336', 'local_hospital', TRUE),
(UUID(), 'Entertainment', 'expense', '#FF9800', 'movie', TRUE);
```

### Monedas Soportadas

```sql
-- seeds/currencies_seed.sql
CREATE TABLE IF NOT EXISTS currencies (
    code VARCHAR(3) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    symbol VARCHAR(10) NOT NULL,
    decimal_places INT DEFAULT 2,
    is_active BOOLEAN DEFAULT TRUE
);

INSERT INTO currencies (code, name, symbol, decimal_places) VALUES
('USD', 'US Dollar', '$', 2),
('EUR', 'Euro', '€', 2),
('GBP', 'British Pound', '£', 2),
('JPY', 'Japanese Yen', '¥', 0),
('CAD', 'Canadian Dollar', 'C$', 2),
('AUD', 'Australian Dollar', 'A$', 2),
('CHF', 'Swiss Franc', 'CHF', 2),
('CNY', 'Chinese Yuan', '¥', 2),
('BTC', 'Bitcoin', '₿', 8),
('ETH', 'Ethereum', 'Ξ', 8);
```

## 🔐 Seguridad

### Medidas Implementadas

- **Encriptación**: Datos sensibles encriptados en reposo
- **Hashing**: Contraseñas con bcrypt
- **SSL/TLS**: Conexiones encriptadas
- **Backup Encryption**: Backups encriptados
- **Access Control**: Control de acceso granular
- **Audit Logging**: Registro de todas las operaciones
- **Data Masking**: Enmascaramiento de datos sensibles

### Configuración de Seguridad

```sql
-- Crear usuario de solo lectura
CREATE USER 'fintrack_readonly'@'%' IDENTIFIED BY 'readonly_password';
GRANT SELECT ON fintrack.* TO 'fintrack_readonly'@'%';

-- Crear usuario para backups
CREATE USER 'fintrack_backup'@'localhost' IDENTIFIED BY 'backup_password';
GRANT SELECT, LOCK TABLES, SHOW VIEW, EVENT, TRIGGER ON fintrack.* TO 'fintrack_backup'@'localhost';

-- Configurar SSL
ALTER USER 'fintrack_user'@'%' REQUIRE SSL;
```

## 📊 Monitoreo y Performance

### Métricas Importantes

```sql
-- Consultas lentas
SELECT * FROM mysql.slow_log 
WHERE start_time > DATE_SUB(NOW(), INTERVAL 1 HOUR)
ORDER BY query_time DESC;

-- Uso de índices
SELECT 
    TABLE_SCHEMA,
    TABLE_NAME,
    INDEX_NAME,
    CARDINALITY
FROM information_schema.STATISTICS 
WHERE TABLE_SCHEMA = 'fintrack'
ORDER BY CARDINALITY DESC;

-- Tamaño de tablas
SELECT 
    TABLE_NAME,
    ROUND(((DATA_LENGTH + INDEX_LENGTH) / 1024 / 1024), 2) AS 'Size (MB)'
FROM information_schema.TABLES 
WHERE TABLE_SCHEMA = 'fintrack'
ORDER BY (DATA_LENGTH + INDEX_LENGTH) DESC;
```

### Health Checks

```bash
# Verificar estado de la base de datos
mysqladmin ping -h localhost -u fintrack_user -p

# Verificar conexiones activas
mysql -u fintrack_user -p -e "SHOW PROCESSLIST;"

# Verificar estado de replicación (si aplica)
mysql -u fintrack_user -p -e "SHOW SLAVE STATUS\G"
```

## 💾 Backup y Recuperación

### Estrategia de Backup

```bash
# Backup completo diario
mysqldump --single-transaction --routines --triggers \
  -u fintrack_backup -p fintrack > backup_$(date +%Y%m%d).sql

# Backup incremental con binlog
mysqlbinlog --start-datetime="2024-01-01 00:00:00" \
  --stop-datetime="2024-01-01 23:59:59" \
  /var/lib/mysql/mysql-bin.000001 > incremental_backup.sql

# Backup con compresión
mysqldump --single-transaction -u fintrack_backup -p fintrack | \
  gzip > backup_$(date +%Y%m%d).sql.gz
```

### Restauración

```bash
# Restaurar desde backup completo
mysql -u fintrack_user -p fintrack < backup_20240101.sql

# Restaurar desde backup comprimido
gunzip < backup_20240101.sql.gz | mysql -u fintrack_user -p fintrack

# Restaurar punto en el tiempo
mysql -u fintrack_user -p fintrack < backup_20240101.sql
mysql -u fintrack_user -p fintrack < incremental_backup.sql
```

## 🚀 Optimización

### Configuración de Performance

```sql
-- Optimizar tablas
OPTIMIZE TABLE transactions, accounts, users;

-- Analizar tablas
ANALYZE TABLE transactions, accounts, users;

-- Verificar fragmentación
SELECT 
    TABLE_NAME,
    ROUND(DATA_FREE / 1024 / 1024, 2) AS 'Fragmentation (MB)'
FROM information_schema.TABLES 
WHERE TABLE_SCHEMA = 'fintrack' AND DATA_FREE > 0;
```

### Índices Recomendados

```sql
-- Índices compuestos para consultas frecuentes
CREATE INDEX idx_transactions_user_date ON transactions(user_id, transaction_date);
CREATE INDEX idx_transactions_account_status ON transactions(account_id, status);
CREATE INDEX idx_accounts_user_active ON accounts(user_id, is_active);

-- Índices para reportes
CREATE INDEX idx_transactions_date_type_amount ON transactions(transaction_date, transaction_type, amount);
CREATE INDEX idx_transactions_user_category_date ON transactions(user_id, category_id, transaction_date);
```

## 🔧 Mantenimiento

### Tareas de Mantenimiento

```bash
# Script de mantenimiento diario
#!/bin/bash

# Limpiar logs antiguos
mysql -u root -p -e "PURGE BINARY LOGS BEFORE DATE_SUB(NOW(), INTERVAL 7 DAY);"

# Optimizar tablas
mysql -u fintrack_user -p fintrack -e "OPTIMIZE TABLE transactions;"

# Verificar integridad
mysqlcheck --check --all-databases -u root -p

# Backup automático
mysqldump --single-transaction -u fintrack_backup -p fintrack | \
  gzip > /backups/fintrack_$(date +%Y%m%d_%H%M%S).sql.gz
```

### Limpieza de Datos

```sql
-- Limpiar notificaciones antiguas
DELETE FROM notifications 
WHERE created_at < DATE_SUB(NOW(), INTERVAL 90 DAY) 
AND status = 'delivered';

-- Archivar transacciones antiguas
CREATE TABLE transactions_archive LIKE transactions;
INSERT INTO transactions_archive 
SELECT * FROM transactions 
WHERE transaction_date < DATE_SUB(NOW(), INTERVAL 2 YEAR);
```

---

**Database** - El corazón de datos de FinTrack 💾🔒