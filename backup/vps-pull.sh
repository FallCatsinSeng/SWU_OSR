#!/usr/bin/env bash
# ============================================================================
# vps-pull.sh — Script untuk VPS: Pull backup dari Google Drive
#
# Jalankan script ini di VPS pribadi kamu.
# VPS akan download backup terbaru dari Google Drive menggunakan rclone.
#
# Setup (satu kali di VPS):
#   1. Install rclone:  curl -fsSL https://rclone.org/install.sh | sudo bash
#   2. Setup Google Drive:  rclone config
#      - Name: gdrive
#      - Type: drive (Google Drive)
#      - Ikuti instruksi OAuth
#   3. Copy script ini ke VPS:
#      scp backup/vps-pull.sh user@your-vps:/opt/swu_osr_backup/
#   4. Edit konfigurasi di bawah
#   5. Tambahkan ke crontab:
#      crontab -e
#      # Pull backup dari Google Drive setiap hari jam 04:00 WIB
#      0 21 * * *  /opt/swu_osr_backup/vps-pull.sh >> /var/log/swu_osr_pull.log 2>&1
#
# Usage:
#   ./vps-pull.sh              # Pull backup terbaru
#   ./vps-pull.sh status       # Lihat status backup lokal
# ============================================================================

set -euo pipefail

# ══════════════════════════════════════════════════════════════════════════════
# KONFIGURASI — Edit sesuai kebutuhan
# ══════════════════════════════════════════════════════════════════════════════

# Folder di VPS untuk menyimpan backup
LOCAL_BACKUP_DIR="/var/backups/swu_osr"

# rclone remote name (sama seperti yang di-setup saat rclone config)
GDRIVE_REMOTE="gdrive"

# Folder di Google Drive (harus sama dengan GDRIVE_PATH di production)
GDRIVE_PATH="SWU_OSR_Backups"

# Path ke rclone config (opsional, kosongkan untuk default)
RCLONE_CONFIG=""

# ══════════════════════════════════════════════════════════════════════════════

# ── Colors ────────────────────────────────────────────────────────────────────
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

log_info()  { echo -e "[$(date -u +'%Y-%m-%dT%H:%M:%SZ')] ${GREEN}[INFO]${NC}  $*"; }
log_error() { echo -e "[$(date -u +'%Y-%m-%dT%H:%M:%SZ')] ${RED}[ERROR]${NC} $*"; }
log_step()  { echo -e "[$(date -u +'%Y-%m-%dT%H:%M:%SZ')] ${CYAN}[STEP]${NC}  $*"; }

# ── Pull command ──────────────────────────────────────────────────────────────
cmd_pull() {
    log_step "Pulling backups from Google Drive..."

    # Check rclone
    if ! command -v rclone &>/dev/null; then
        log_error "rclone not found!"
        echo "  Install: curl -fsSL https://rclone.org/install.sh | sudo bash"
        echo "  Setup:   rclone config (pilih Google Drive)"
        exit 1
    fi

    # Ensure local directory exists
    mkdir -p "${LOCAL_BACKUP_DIR}"

    # Build rclone options
    local rclone_opts=""
    if [ -n "${RCLONE_CONFIG}" ] && [ -f "${RCLONE_CONFIG}" ]; then
        rclone_opts="--config ${RCLONE_CONFIG}"
    fi

    # rclone sync — only downloads new/changed files (bandwidth efficient!)
    log_info "Source: ${GDRIVE_REMOTE}:${GDRIVE_PATH}/"
    log_info "Target: ${LOCAL_BACKUP_DIR}/"

    if rclone sync \
        ${rclone_opts} \
        --transfers=4 \
        --checkers=8 \
        --low-level-retries=3 \
        --stats-one-line \
        --stats=10s \
        "${GDRIVE_REMOTE}:${GDRIVE_PATH}/" \
        "${LOCAL_BACKUP_DIR}/"; then
        log_info "Pull completed successfully!"
    else
        log_error "Pull FAILED!"
        exit 1
    fi

    # Show summary
    cmd_status
}

# ── Status command ────────────────────────────────────────────────────────────
cmd_status() {
    echo ""
    echo -e "${CYAN}╔══════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${CYAN}║          SWU OSR — VPS Backup Mirror Status                 ║${NC}"
    echo -e "${CYAN}╚══════════════════════════════════════════════════════════════╝${NC}"
    echo ""

    if [ ! -d "${LOCAL_BACKUP_DIR}" ]; then
        echo "  No backups found. Run: ./vps-pull.sh"
        return
    fi

    # Base backups
    echo -e "${YELLOW}Base Backups:${NC}"
    local base_count
    base_count=$(find "${LOCAL_BACKUP_DIR}/base" -mindepth 1 -maxdepth 1 -type d 2>/dev/null | wc -l)
    echo "  Count: ${base_count}"
    if [ "$base_count" -gt 0 ]; then
        # shellcheck disable=SC2012
        ls -dt "${LOCAL_BACKUP_DIR}/base"/*/ 2>/dev/null | head -3 | while read -r dir; do
            local size
            size=$(du -sh "$dir" 2>/dev/null | cut -f1)
            echo "    $(basename "$dir") — ${size}"
        done
    fi

    # WAL segments
    echo ""
    echo -e "${YELLOW}WAL Segments:${NC}"
    local wal_count
    wal_count=$(find "${LOCAL_BACKUP_DIR}/wal" -type f 2>/dev/null | wc -l)
    local wal_size
    wal_size=$(du -sh "${LOCAL_BACKUP_DIR}/wal" 2>/dev/null | cut -f1)
    echo "  Count: ${wal_count} segments (${wal_size})"

    # Upload snapshots
    echo ""
    echo -e "${YELLOW}Upload Snapshots:${NC}"
    local uploads_count
    uploads_count=$(find "${LOCAL_BACKUP_DIR}/uploads" -name "*.tar.gz" 2>/dev/null | wc -l)
    local uploads_size
    uploads_size=$(du -sh "${LOCAL_BACKUP_DIR}/uploads" 2>/dev/null | cut -f1)
    echo "  Count: ${uploads_count} (${uploads_size})"

    # Total
    echo ""
    echo -e "${YELLOW}Total Size:${NC}"
    local total
    total=$(du -sh "${LOCAL_BACKUP_DIR}" 2>/dev/null | cut -f1)
    echo "  ${total}"
    echo ""
}

# ── Main ──────────────────────────────────────────────────────────────────────
case "${1:-pull}" in
    pull)
        cmd_pull
        ;;
    status)
        cmd_status
        ;;
    *)
        echo "Usage: $0 {pull|status}"
        echo ""
        echo "Commands:"
        echo "  pull      Download latest backups from Google Drive"
        echo "  status    Show local backup status"
        echo ""
        echo "Setup crontab (daily pull jam 04:00 WIB):"
        echo "  0 21 * * *  $(readlink -f "$0") >> /var/log/swu_osr_pull.log 2>&1"
        exit 1
        ;;
esac
