#!/usr/bin/env bash
# ============================================================================
# restore.sh — Disaster Recovery for SWU_OSR
#
# Restores PostgreSQL from a base backup + WAL replay (point-in-time recovery).
# Can download backups from Google Drive before restoring.
#
# Usage:
#   ./restore.sh list                       # List available local backups
#   ./restore.sh restore <backup_name>      # Restore specific backup
#   ./restore.sh restore latest             # Restore latest backup
#   ./restore.sh download                   # Download backups from Google Drive
#
# ⚠️  WARNING: This will REPLACE the current database!
#     Always verify your backup first.
# ============================================================================

set -euo pipefail

# ── Resolve paths ─────────────────────────────────────────────────────────────
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ENV_FILE="${SCRIPT_DIR}/.env"

# ── Load configuration ───────────────────────────────────────────────────────
if [ ! -f "$ENV_FILE" ]; then
    echo "[ERROR] Config not found: ${ENV_FILE}"
    exit 1
fi

# shellcheck disable=SC1090
source "$ENV_FILE"

# ── Defaults ──────────────────────────────────────────────────────────────────
BACKUP_LOCAL_DIR="${BACKUP_LOCAL_DIR:-/backups}"
PROJECT_DIR="${PROJECT_DIR:-$(dirname "$SCRIPT_DIR")}"
PG_CONTAINER="${PG_CONTAINER:-swu_osr-postgres-1}"
PG_USER="${PG_USER:-postgres}"
PG_DB="${PG_DB:-swu_osr}"

# ── Colors ────────────────────────────────────────────────────────────────────
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m'

log_info()  { echo -e "[$(date -u +'%Y-%m-%dT%H:%M:%SZ')] ${GREEN}[INFO]${NC}  $*"; }
log_warn()  { echo -e "[$(date -u +'%Y-%m-%dT%H:%M:%SZ')] ${YELLOW}[WARN]${NC}  $*"; }
log_error() { echo -e "[$(date -u +'%Y-%m-%dT%H:%M:%SZ')] ${RED}[ERROR]${NC} $*"; }
log_step()  { echo -e "[$(date -u +'%Y-%m-%dT%H:%M:%SZ')] ${CYAN}[STEP]${NC}  $*"; }

# ══════════════════════════════════════════════════════════════════════════════
# COMMAND: list — List available base backups
# ══════════════════════════════════════════════════════════════════════════════
cmd_list() {
    echo ""
    echo -e "${CYAN}Available Base Backups:${NC}"
    echo ""

    if [ ! -d "${BACKUP_LOCAL_DIR}/base" ]; then
        echo "  No backups found in ${BACKUP_LOCAL_DIR}/base"
        echo "  Run: ./restore.sh download   (to fetch from Google Drive)"
        return
    fi

    local count=0
    # shellcheck disable=SC2012
    for dir in $(ls -dt "${BACKUP_LOCAL_DIR}/base"/*/ 2>/dev/null); do
        local name
        name=$(basename "$dir")
        local size
        size=$(du -sh "$dir" 2>/dev/null | cut -f1)
        local file_count
        file_count=$(find "$dir" -type f | wc -l)

        if [ $count -eq 0 ]; then
            echo -e "  ${GREEN}→ ${name}${NC}  (${size}, ${file_count} files) ${YELLOW}[latest]${NC}"
        else
            echo -e "    ${name}  (${size}, ${file_count} files)"
        fi
        ((count++))
    done

    if [ $count -eq 0 ]; then
        echo "  No base backups found."
        echo "  Run: ./restore.sh download   (to fetch from Google Drive)"
    fi

    # WAL info
    echo ""
    echo -e "${CYAN}WAL Segments:${NC}"
    local wal_count
    wal_count=$(find "${BACKUP_LOCAL_DIR}/wal" -type f 2>/dev/null | wc -l)
    local wal_size
    wal_size=$(du -sh "${BACKUP_LOCAL_DIR}/wal" 2>/dev/null | cut -f1)
    echo "  ${wal_count} segments (${wal_size})"
    echo ""
}

# ══════════════════════════════════════════════════════════════════════════════
# COMMAND: restore — Restore from base backup + WAL
# ══════════════════════════════════════════════════════════════════════════════
cmd_restore() {
    local backup_name="${1:-}"

    if [ -z "$backup_name" ]; then
        log_error "Usage: $0 restore <backup_name|latest>"
        cmd_list
        return 1
    fi

    # Resolve "latest"
    if [ "$backup_name" = "latest" ]; then
        # shellcheck disable=SC2012
        backup_name=$(ls -dt "${BACKUP_LOCAL_DIR}/base"/*/ 2>/dev/null | head -1 | xargs basename 2>/dev/null || echo "")
        if [ -z "$backup_name" ]; then
            log_error "No backups found! Run: ./restore.sh download"
            return 1
        fi
        log_info "Resolved 'latest' to: ${backup_name}"
    fi

    local backup_path="${BACKUP_LOCAL_DIR}/base/${backup_name}"
    if [ ! -d "$backup_path" ]; then
        log_error "Backup not found: ${backup_path}"
        cmd_list
        return 1
    fi

    # ── Safety confirmation ───────────────────────────────────────────────
    echo ""
    echo -e "${RED}╔══════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${RED}║                    ⚠️  DANGER ZONE ⚠️                       ║${NC}"
    echo -e "${RED}╠══════════════════════════════════════════════════════════════╣${NC}"
    echo -e "${RED}║  This will REPLACE the current database with:              ║${NC}"
    printf "${RED}║  Backup: %-51s║${NC}\n" "${backup_name}"
    echo -e "${RED}║                                                            ║${NC}"
    echo -e "${RED}║  ALL CURRENT DATA WILL BE LOST!                            ║${NC}"
    echo -e "${RED}╚══════════════════════════════════════════════════════════════╝${NC}"
    echo ""
    read -rp "  Type 'RESTORE' to confirm: " confirm
    if [ "$confirm" != "RESTORE" ]; then
        log_info "Restore cancelled."
        return 0
    fi

    # ── Step 1: Stop application services ─────────────────────────────────
    log_step "[1/5] Stopping application services..."
    cd "$PROJECT_DIR"
    docker compose stop backend frontend nginx 2>/dev/null || true

    # ── Step 2: Stop PostgreSQL ───────────────────────────────────────────
    log_step "[2/5] Stopping PostgreSQL..."
    docker compose stop postgres

    # ── Step 3: Replace data directory ────────────────────────────────────
    log_step "[3/5] Replacing PostgreSQL data with backup..."

    # Get the Docker volume mountpoint
    local pgdata_volume
    pgdata_volume=$(docker volume inspect swu_osr_pgdata --format '{{ .Mountpoint }}' 2>/dev/null)

    if [ -z "$pgdata_volume" ]; then
        log_error "Cannot find pgdata volume. Is Docker running?"
        return 1
    fi

    # Clear existing data
    sudo rm -rf "${pgdata_volume:?}"/*

    # Extract base backup
    log_info "Extracting base backup to ${pgdata_volume}..."
    sudo tar xzf "${backup_path}/base.tar.gz" -C "$pgdata_volume"

    # Copy WAL segments for replay
    if [ -d "${BACKUP_LOCAL_DIR}/wal" ] && [ "$(find "${BACKUP_LOCAL_DIR}/wal" -type f | wc -l)" -gt 0 ]; then
        log_info "Copying WAL segments for replay..."
        local pg_wal_dir="${pgdata_volume}/pg_wal"
        sudo mkdir -p "$pg_wal_dir"
        sudo cp "${BACKUP_LOCAL_DIR}/wal/"* "$pg_wal_dir/" 2>/dev/null || true
    fi

    # Create recovery signal file (PostgreSQL 12+ uses this instead of recovery.conf)
    sudo touch "${pgdata_volume}/recovery.signal"

    # Create postgresql.auto.conf with restore_command if not present
    if ! sudo grep -q "restore_command" "${pgdata_volume}/postgresql.auto.conf" 2>/dev/null; then
        echo "restore_command = 'cp /backups/wal/%f %p'" | sudo tee -a "${pgdata_volume}/postgresql.auto.conf" >/dev/null
    fi

    # ── Step 4: Start PostgreSQL for recovery ─────────────────────────────
    log_step "[4/5] Starting PostgreSQL for WAL recovery..."
    docker compose start postgres

    # Wait for recovery to complete
    log_info "Waiting for recovery to complete..."
    local retry=0
    local max_retry=60
    until docker exec "${PG_CONTAINER}" pg_isready -U "${PG_USER}" &>/dev/null || [ $retry -ge $max_retry ]; do
        retry=$((retry + 1))
        echo -ne "\r  Recovery in progress... (${retry}/${max_retry})"
        sleep 2
    done
    echo ""

    if [ $retry -ge $max_retry ]; then
        log_error "PostgreSQL did not become ready after recovery!"
        log_error "Check logs: docker compose logs postgres"
        return 1
    fi

    log_info "PostgreSQL recovery completed."

    # ── Step 5: Restart all services ──────────────────────────────────────
    log_step "[5/5] Restarting all services..."
    docker compose start backend frontend nginx

    echo ""
    log_info "Restore completed successfully from backup: ${backup_name}"
    log_info "Verify data integrity by checking the application."
    echo ""
}

# ══════════════════════════════════════════════════════════════════════════════
# COMMAND: download — Download backups from Google Drive
# ══════════════════════════════════════════════════════════════════════════════
cmd_download() {
    if [ "${GDRIVE_ENABLED}" != "true" ]; then
        log_error "Google Drive is not configured (GDRIVE_ENABLED=false)"
        return 1
    fi

    if ! command -v rclone &>/dev/null; then
        log_error "rclone not found! Run: backup/setup-gdrive.sh"
        return 1
    fi

    log_step "Downloading backups from Google Drive: ${GDRIVE_REMOTE}:${GDRIVE_PATH}"

    local rclone_opts=""
    if [ -n "${RCLONE_CONFIG:-}" ] && [ -f "${RCLONE_CONFIG}" ]; then
        rclone_opts="--config ${RCLONE_CONFIG}"
    fi

    mkdir -p "${BACKUP_LOCAL_DIR}"

    rclone sync \
        ${rclone_opts} \
        --transfers=4 \
        --progress \
        "${GDRIVE_REMOTE}:${GDRIVE_PATH}/" \
        "${BACKUP_LOCAL_DIR}/"

    log_info "Download from Google Drive completed."
    cmd_list
}

# ══════════════════════════════════════════════════════════════════════════════
# MAIN
# ══════════════════════════════════════════════════════════════════════════════
main() {
    local command="${1:-}"
    shift || true

    case "$command" in
        list)
            cmd_list
            ;;
        restore)
            cmd_restore "$@"
            ;;
        download)
            cmd_download
            ;;
        *)
            echo "Usage: $0 {list|restore|download}"
            echo ""
            echo "Commands:"
            echo "  list              List available local backups"
            echo "  restore <name>    Restore from a base backup (+ WAL replay)"
            echo "  restore latest    Restore from the most recent backup"
            echo "  download          Download backups from Google Drive"
            echo ""
            echo "Recovery workflow:"
            echo "  1. ./restore.sh download         # fetch from Google Drive"
            echo "  2. ./restore.sh list              # see available backups"
            echo "  3. ./restore.sh restore latest    # restore"
            exit 1
            ;;
    esac
}

main "$@"
