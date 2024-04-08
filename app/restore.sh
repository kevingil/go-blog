#!/bin/bash

set -e

# Execute MySQL commands and capture output
mysql_exec() {
    local output
    if output=$(mysql -h db -u "$MYSQL_USER" -p"$MYSQL_PASSWORD" -e "$1" 2>&1); then
        return 0  
    else
        # Return error
        echo "Error: $output"
        return 1  
    fi
}

# Check if the database exists
if ! mysql_exec "USE $MYSQL_DATABASE;"; then
    echo "Database does not exists, you must create with Docker compose"
    exit 1
else 
    echo "Database present"
fi

# Check if data already exists
if mysql_exec "SELECT 1 FROM $MYSQL_DATABASE.users LIMIT 1;" >/dev/null; then
    echo "Users table present, database is in use"
    exit 1
fi

# Download the backup file
curl -L "$RESTORE_URL" -o /tmp/pscale_dump.zip

# Unzip the backup file
unzip /tmp/pscale_dump.zip -d /tmp/pscale_dump/

# Restore the schema files
for schema_file in /tmp/pscale_dump/*-schema.sql; do
    if ! mysql "$MYSQL_DATABASE" < "$schema_file"; then
        echo "Error: Failed to restore schema from $schema_file"
        exit 1
    fi
done

# Restore the table data files
for file in /tmp/pscale_dump/*.sql; do
    if [[ "$file" != *-schema.sql ]]; then
        if ! mysql "$MYSQL_DATABASE" < "$file"; then
            echo "Error: Failed to restore data from $file"
            exit 1
        fi
    fi
done
