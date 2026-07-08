#!/usr/bin/env bash
# ============================================================================
# setup-gdrive.sh — Setup rclone for Google Drive backup
#
# This script will:
#   1. Install rclone if not present
#   2. Guide you through Google Drive OAuth setup
#   3. Verify the connection with a test upload
#
# Run this ONCE on the server before enabling GDRIVE_ENABLED=true
# ============================================================================

set -euo pipefail

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

echo -e "${CYAN}"
echo "╔══════════════════════════════════════════════════════════════╗"
echo "║        SWU OSR — Google Drive Backup Setup                  ║"
echo "╚══════════════════════════════════════════════════════════════╝"
echo -e "${NC}"

# ── Step 1: Install rclone ────────────────────────────────────────────────────
echo -e "${YELLOW}[1/4] Checking rclone installation...${NC}"

if command -v rclone &>/dev/null; then
    echo -e "${GREEN}  ✓ rclone $(rclone version | head -1 | awk '{print $2}') already installed${NC}"
else
    echo -e "${CYAN}  Installing rclone...${NC}"
    curl -fsSL https://rclone.org/install.sh | sudo bash
    if command -v rclone &>/dev/null; then
        echo -e "${GREEN}  ✓ rclone installed successfully${NC}"
    else
        echo -e "${RED}  ✗ rclone installation failed!${NC}"
        echo "  Try manual install: https://rclone.org/install/"
        exit 1
    fi
fi

# ── Step 2: Configure Google Drive remote ─────────────────────────────────────
echo ""
echo -e "${YELLOW}[2/4] Configuring Google Drive remote...${NC}"
echo ""
echo -e "${CYAN}  This will open an interactive configuration wizard.${NC}"
echo -e "${CYAN}  When prompted:${NC}"
echo -e "${CYAN}    1. Choose 'n' for new remote${NC}"
echo -e "${CYAN}    2. Name it: ${GREEN}gdrive${NC}"
echo -e "${CYAN}    3. Type: ${GREEN}drive${NC} (Google Drive)${NC}"
echo -e "${CYAN}    4. Leave client_id and client_secret blank (use rclone's)${NC}"
echo -e "${CYAN}    5. Scope: ${GREEN}1${NC} (Full access)${NC}"
echo -e "${CYAN}    6. Follow the OAuth flow in your browser${NC}"
echo ""
echo -e "${YELLOW}  ⚠ If this is a HEADLESS server (no browser):${NC}"
echo -e "${YELLOW}    - Choose 'n' when asked about auto config${NC}"
echo -e "${YELLOW}    - It will give you a URL to open on your local machine${NC}"
echo -e "${YELLOW}    - Paste the verification code back here${NC}"
echo ""

read -rp "  Press Enter to start rclone config... "
rclone config

# ── Step 3: Verify connection ─────────────────────────────────────────────────
echo ""
echo -e "${YELLOW}[3/4] Verifying Google Drive connection...${NC}"

# Load remote name from backup config if exists
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
GDRIVE_REMOTE="gdrive"
GDRIVE_PATH="SWU_OSR_Backups"

if [ -f "${SCRIPT_DIR}/.env" ]; then
    # shellcheck disable=SC1090
    source "${SCRIPT_DIR}/.env"
    GDRIVE_REMOTE="${GDRIVE_REMOTE:-gdrive}"
    GDRIVE_PATH="${GDRIVE_PATH:-SWU_OSR_Backups}"
fi

# Test by listing the remote
if rclone lsd "${GDRIVE_REMOTE}:" &>/dev/null; then
    echo -e "${GREEN}  ✓ Google Drive connection successful!${NC}"
else
    echo -e "${RED}  ✗ Cannot connect to Google Drive remote '${GDRIVE_REMOTE}'${NC}"
    echo -e "${RED}    Make sure you named the remote '${GDRIVE_REMOTE}' during setup${NC}"
    echo -e "${RED}    Run 'rclone config' to reconfigure${NC}"
    exit 1
fi

# Create backup folder on Google Drive
echo -e "${CYAN}  Creating backup folder: ${GDRIVE_REMOTE}:${GDRIVE_PATH}/${NC}"
rclone mkdir "${GDRIVE_REMOTE}:${GDRIVE_PATH}" 2>/dev/null || true

# Test write
echo "SWU_OSR backup test — $(date -u)" > /tmp/swu_osr_backup_test.txt
if rclone copyto /tmp/swu_osr_backup_test.txt "${GDRIVE_REMOTE}:${GDRIVE_PATH}/.backup_test" &>/dev/null; then
    echo -e "${GREEN}  ✓ Write test passed!${NC}"
    rclone deletefile "${GDRIVE_REMOTE}:${GDRIVE_PATH}/.backup_test" &>/dev/null || true
else
    echo -e "${RED}  ✗ Write test failed! Check permissions.${NC}"
    exit 1
fi
rm -f /tmp/swu_osr_backup_test.txt

# ── Step 4: Show config location ─────────────────────────────────────────────
echo ""
echo -e "${YELLOW}[4/4] Setup complete!${NC}"
echo ""

RCLONE_CONF_PATH=$(rclone config file | tail -1)
echo -e "${GREEN}╔══════════════════════════════════════════════════════════════╗${NC}"
echo -e "${GREEN}║  Google Drive backup setup complete!                         ║${NC}"
echo -e "${GREEN}╠══════════════════════════════════════════════════════════════╣${NC}"
echo -e "${GREEN}║                                                              ║${NC}"
printf  "${GREEN}║  Config:  %-50s║${NC}\n" "${RCLONE_CONF_PATH}"
printf  "${GREEN}║  Remote:  %-50s║${NC}\n" "${GDRIVE_REMOTE}:${GDRIVE_PATH}/"
echo -e "${GREEN}║                                                              ║${NC}"
echo -e "${GREEN}║  Next steps:                                                 ║${NC}"
echo -e "${GREEN}║  1. Set GDRIVE_ENABLED=true in backup/.env                   ║${NC}"
printf  "${GREEN}║  2. Set RCLONE_CONFIG=%-38s║${NC}\n" "${RCLONE_CONF_PATH}"
echo -e "${GREEN}║  3. Run: ./backup/backup.sh sync                             ║${NC}"
echo -e "${GREEN}║                                                              ║${NC}"
echo -e "${GREEN}╚══════════════════════════════════════════════════════════════╝${NC}"
echo ""
