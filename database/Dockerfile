# MySQL with custom configuration
FROM mysql:8.0

# Set environment variables
ENV MYSQL_DATABASE=fintrack
ENV MYSQL_USER=fintrack_user
ENV MYSQL_PASSWORD=fintrack_password
ENV MYSQL_ROOT_PASSWORD=root_password

# Copy initialization scripts in correct order
# Docker executes SQL files in alphabetical order
COPY ./docker-entrypoint-initdb.d/ /docker-entrypoint-initdb.d/
COPY ./migrations/ /docker-entrypoint-initdb.d/
COPY ./schemas/ /docker-entrypoint-initdb.d/
COPY ./seeds/ /docker-entrypoint-initdb.d/

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=30s --retries=3 \
  CMD mysqladmin ping -h localhost -u $MYSQL_USER -p$MYSQL_PASSWORD || exit 1

# Expose port
EXPOSE 3306

# Use default entrypoint and CMD from official image