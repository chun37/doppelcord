#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
ENV_FILE="$PROJECT_ROOT/.env"

if [ ! -f "$ENV_FILE" ]; then
    echo "Error: .env file not found at $ENV_FILE"
    exit 1
fi

set -a
source "$ENV_FILE"
set +a

if [ -z "$DB_HOST" ] || [ -z "$DB_PORT" ] || [ -z "$DB_USER" ] || [ -z "$DB_PASSWORD" ] || [ -z "$DB_NAME" ]; then
    echo "Error: Missing required database environment variables"
    echo "Required: DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME"
    exit 1
fi

DB_SSLMODE="${DB_SSLMODE:-disable}"

DATABASE_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}"

MIGRATIONS_DIR="$PROJECT_ROOT/migrations"

usage() {
    echo "Usage: $0 <command> [args]"
    echo ""
    echo "Commands:"
    echo "  up              Apply all pending migrations"
    echo "  down            Rollback the last migration"
    echo "  version         Show current migration version"
    echo "  force VERSION   Force set migration version (use with caution)"
    echo ""
    exit 1
}

if [ $# -lt 1 ]; then
    usage
fi

COMMAND="$1"
shift

case "$COMMAND" in
    up)
        echo "Applying migrations..."
        migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" up
        ;;
    down)
        echo "Rolling back migration..."
        migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" down 1
        ;;
    version)
        migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" version
        ;;
    force)
        if [ $# -lt 1 ]; then
            echo "Error: force requires a version number"
            exit 1
        fi
        VERSION="$1"
        echo "Forcing version to $VERSION..."
        migrate -path "$MIGRATIONS_DIR" -database "$DATABASE_URL" force "$VERSION"
        ;;
    *)
        echo "Error: Unknown command '$COMMAND'"
        usage
        ;;
esac
