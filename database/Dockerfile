# MySQL with custom configuration
FROM mysql:8.0

# Install additional tools
RUN apt-get update && apt-get install -y \
    mysql-client \
    && rm -rf /var/lib/apt/lists/*

# Set environment variables
ENV MYSQL_DATABASE=fintrack
ENV MYSQL_USER=fintrack_user
ENV MYSQL_PASSWORD=fintrack_password
ENV MYSQL_ROOT_PASSWORD=root_password

# Copy initialization scripts
COPY ./migrations/*.sql /docker-entrypoint-initdb.d/
COPY ./seeds/*.sql /docker-entrypoint-initdb.d/
COPY ./schemas/*.sql /docker-entrypoint-initdb.d/

# Create custom entrypoint script
COPY docker-entrypoint.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/docker-entrypoint.sh

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=30s --retries=3 \
  CMD mysqladmin ping -h localhost -u $MYSQL_USER -p$MYSQL_PASSWORD || exit 1

# Expose port
EXPOSE 3306

# Use custom entrypoint
ENTRYPOINT ["/usr/local/bin/docker-entrypoint.sh"]
CMD ["mysqld"]