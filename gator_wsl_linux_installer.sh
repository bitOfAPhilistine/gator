#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd)

echo "Moving executable into PATH..."
sudo mv ${SCRIPT_DIR}/gator ${PATH%%:*}

echo "Checking for PostgreSQL..."

# Check whether PostgreSQL is already installed
if command -v psql >/dev/null 2>&1; then
    echo "PostgreSQL is already installed."
    psql --version
else
    echo "PostgreSQL not found. Installing..."

    # Make sure apt is available
    if ! command -v apt-get >/dev/null 2>&1; then
        echo "Error: This script requires a Debian/Ubuntu-based Linux distribution with apt."
        exit 1
    fi

    sudo apt-get update
    sudo apt-get install -y postgresql postgresql-contrib

    echo "PostgreSQL installed successfully."
fi

# Start PostgreSQL if possible
if command -v service >/dev/null 2>&1; then
    echo "Starting PostgreSQL..."
    sudo service postgresql start || true
fi

echo "PostgreSQL setup complete."
echo "Version:"
psql --version

echo "Setting up gator database..."

# Create the database if it doesn't already exist
if sudo -u postgres psql -tAc "SELECT 1 FROM pg_database WHERE datname='gator'" | grep -q 1; then
    echo "Database 'gator' already exists."
else
    echo "Creating database 'gator'..."
    sudo -u postgres createdb gator
fi

# Set the postgres password
read -p "Enter database password: " PASSWORD
sudo -u postgres psql -c "ALTER USER postgres PASSWORD '"$PASSWORD"';" >/dev/null
CONNECTION_URL="postgres://postgres:"$PASSWORD"@localhost:5432/gator"

# Initialize the database
echo "Initializing database schema..."
if ! command -v goose >/dev/null 2>&1; then
    echo "Goose not found, installing..."
    sudo curl -fsSL https://raw.githubusercontent.com/pressly/goose/master/install.sh | sudo sh
    echo "Goose installed."
    echo "Version:"
    goose -version
fi
goose postgres $CONNECTION_URL -dir ${SCRIPT_DIR}/sql/schema up

echo "Gator database setup complete."

echo "Initializing config file..."
gator init $CONNECTION_URL
echo "Installastion complete!"