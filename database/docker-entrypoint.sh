#!/bin/bash
set -e

# Custom MySQL entrypoint for FinTrack

echo "Starting FinTrack MySQL initialization..."

# Set default values if not provided
export MYSQL_DATABASE=${MYSQL_DATABASE:-fintrack}
export MYSQL_USER=${MYSQL_USER:-fintrack_user}
export MYSQL_PASSWORD=${MYSQL_PASSWORD:-fintrack_password}
export MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD:-root_password}

# Initialize database if data directory is empty
if [ ! -d "/var/lib/mysql/mysql" ]; then
    echo "Initializing database..."
    
    # Initialize MySQL
    mysqld --initialize-insecure --user=mysql --datadir=/var/lib/mysql
    
    # Start MySQL temporarily for setup
    mysqld --user=mysql --datadir=/var/lib/mysql --skip-networking &
    MYSQL_PID=$!
    
    # Wait for MySQL to start
    until mysqladmin ping --silent; do
        sleep 1
    done
    
    # Create database and user
    mysql -u root <<-EOSQL
        CREATE DATABASE IF NOT EXISTS $MYSQL_DATABASE;
        CREATE USER IF NOT EXISTS '$MYSQL_USER'@'%' IDENTIFIED BY '$MYSQL_PASSWORD';
        GRANT ALL PRIVILEGES ON $MYSQL_DATABASE.* TO '$MYSQL_USER'@'%';
        FLUSH PRIVILEGES;
EOSQL
    
    echo "Running initialization scripts..."
    
    # Run initialization scripts in order
    for f in /docker-entrypoint-initdb.d/*; do
        case "$f" in
            *.sql)    echo "$0: running $f"; mysql -u "$MYSQL_USER" -p"$MYSQL_PASSWORD" "$MYSQL_DATABASE" < "$f"; echo ;;
            *.sql.gz) echo "$0: running $f"; gunzip -c "$f" | mysql -u "$MYSQL_USER" -p"$MYSQL_PASSWORD" "$MYSQL_DATABASE"; echo ;;
            *)        echo "$0: ignoring $f" ;;
        esac
        echo
    done
    
    # Stop temporary MySQL
    kill $MYSQL_PID
    wait $MYSQL_PID
    
    echo "Database initialization completed."
fi

# Start MySQL
echo "Starting MySQL server..."
exec "$@"