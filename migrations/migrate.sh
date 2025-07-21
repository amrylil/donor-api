#!/bin/bash

set -e

# Pastikan file .env ada di root folder jika menggunakan ini
if [ -f .env ]; then
  export $(cat .env | grep -v '#' | awk '/=/ {print $1}')
fi

# Cek apakah variabel database sudah di-set
if [ -z "${DB_USER}" ] || [ -z "${DB_PASSWORD}" ] || [ -z "${DB_HOST}" ] || [ -z "${DB_PORT}" ] || [ -z "${DB_NAME}" ]; then
  echo "ERROR: Pastikan environment variables berikut sudah di-set: DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME"
  exit 1
fi

# 2. Buat DSN (Database Source Name) untuk PostgreSQL
DATABASE_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"
MIGRATIONS_PATH="migrations"

# 3. Baca Perintah dari Argumen
CMD=$1
MIGRATION_NAME=$2

# Fungsi untuk menampilkan cara penggunaan
usage() {
  echo "Usage: $0 {up|down|force|create}"
  echo "  up                - Menjalankan semua migrasi."
  echo "  down [n]          - Membatalkan (rollback) 'n' migrasi terakhir (default: 1)."
  echo "  force <version>   - Memaksa migrasi ke versi tertentu (untuk memperbaiki state 'dirty')."
  echo "  create <name>     - Membuat file migrasi baru (contoh: $0 create add_role_to_users)."
  exit 1
}

# Menjalankan perintah yang sesuai
case "$CMD" in
  "up")
    echo "🚀 Menjalankan migrasi..."
    migrate -database "$DATABASE_URL" -path "$MIGRATIONS_PATH" up
    echo "✅ Migrasi selesai."
    ;;
  "down")
    # Rollback n steps, default to 1
    STEPS=${MIGRATION_NAME:-1}
    echo "⏪ Membatalkan $STEPS migrasi terakhir..."
    migrate -database "$DATABASE_URL" -path "$MIGRATIONS_PATH" down "$STEPS"
    echo "✅ Rollback selesai."
    ;;
  "force")
    if [ -z "$MIGRATION_NAME" ]; then
      echo "ERROR: Mohon masukkan versi migrasi."
      usage
    fi
    echo "⚠️  Memaksa migrasi ke versi $MIGRATION_NAME..."
    migrate -database "$DATABASE_URL" -path "$MIGRATIONS_PATH" force "$MIGRATION_NAME"
    echo "✅ Force migrate selesai."
    ;;
  "create")
    if [ -z "$MIGRATION_NAME" ]; then
      echo "ERROR: Mohon masukkan nama untuk migrasi baru."
      usage
    fi
    echo "✨ Membuat file migrasi baru: $MIGRATION_NAME"
    migrate create -ext sql -dir "$MIGRATIONS_PATH" -seq "$MIGRATION_NAME"
    echo "✅ File migrasi berhasil dibuat di folder '$MIGRATIONS_PATH'."
    ;;
  *)
    usage
    ;;
esac