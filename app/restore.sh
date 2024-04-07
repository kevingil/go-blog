#!/bin/bash

set -e

# Check if the database exists
if ! mysql -e "USE $MYSQL_DATABASE;" &> /dev/null; then
  echo "Database $MYSQL_DATABASE does not exist"
  exit 1
fi

# Check if data already exists
if mysql -e "SELECT 1 FROM $MYSQL_DATABASE.users LIMIT 1;" &> /dev/null; then
  echo "Database alraedy has data"
  exit 1
fi

# Download the backup file
curl -L $RESTORE_URL -o /tmp/pscale_dump.zip

# Unzip the backup file
unzip /tmp/pscale_dump.zip -d /tmp/pscale_dump/

# Restore the schema files
mysql $MYSQL_DATABASE < /tmp/pscale_dump/*-schema.sql

# Restore the table data files
for file in /tmp/pscale_dump/*.sql; do
  if [[ "$file" != *-schema.sql ]]; then
    mysql $MYSQL_DATABASE < "$file"
  fi
done
