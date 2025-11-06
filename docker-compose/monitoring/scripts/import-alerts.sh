#!/bin/bash

# Script to import LPS alert rules to Grafana using Provisioning API
# Usage: ./import-alerts.sh [grafana_url] [username] [password] [folder_uid] [datasource_uid]

# Configuration
GRAFANA_URL=${1:-"http://localhost:3000"}
USERNAME=${2:-"admin"}
PASSWORD=${3:-"test"}
FOLDER_UID=${4:-"LPS"}  # Default to LPS folder
DATASOURCE_UID=${5:-"loki-uid"}  # Default Loki datasource UID
SCRIPT_DIR="$(dirname "$0")"

# Detect OS for cross-platform sed compatibility
OS_TYPE="$(uname)"
if [[ "$OS_TYPE" == "Darwin" ]]; then
    # macOS - requires backup extension even if empty
    SED_INPLACE=("sed" "-i" "")
elif [[ "$OS_TYPE" == "Linux" ]]; then
    # Linux - no backup extension needed
    SED_INPLACE=("sed" "-i")
else
    echo "Warning: Unsupported OS: $OS_TYPE"
     exit 1
fi

# Alert rules configuration
RULE_FILES=(
    "node-eclipse-detection.json"
    "pegin-out-of-liquidity.json"
    "pegout-out-of-liquidity.json"
    "lps-penalization.json"
)

RULE_NAMES=(
    "Node Eclipse Detection Alert"
    "PegIn Out of Liquidity Alert"
    "PegOut Out of Liquidity Alert"
    "LPS Penalization Alert"
)

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Importing LPS alert rules to Grafana...${NC}"
echo "Target: $GRAFANA_URL"
echo "Folder: $FOLDER_UID"
echo "Loki Datasource UID: $DATASOURCE_UID"
echo "Rules to import: ${#RULE_FILES[@]}"
echo ""

# Check if Grafana is accessible
echo "Checking Grafana connectivity..."
HTTP_STATUS=$(curl -s -o /dev/null -w "%{http_code}" "$GRAFANA_URL/api/health")
if [ "$HTTP_STATUS" != "200" ]; then
    echo -e "${RED}ERROR: Cannot reach Grafana at $GRAFANA_URL (HTTP $HTTP_STATUS)${NC}"
    echo "Make sure Grafana is running and accessible."
    exit 1
fi
echo -e "${GREEN}Grafana is accessible${NC}"

# Check if folder exists, create if not
echo ""
echo "Checking if folder '$FOLDER_UID' exists..."
FOLDER_RESPONSE=$(curl -s -u "$USERNAME:$PASSWORD" "$GRAFANA_URL/api/folders/$FOLDER_UID")
FOLDER_STATUS=$(echo "$FOLDER_RESPONSE" | grep -o '"id"' || echo "not_found")

if [ "$FOLDER_STATUS" == "not_found" ]; then
    echo "Creating folder '$FOLDER_UID'..."
    CREATE_FOLDER_RESPONSE=$(curl -s -X POST \
        -u "$USERNAME:$PASSWORD" \
        -H "Content-Type: application/json" \
        -d "{\"uid\":\"$FOLDER_UID\",\"title\":\"LPS\"}" \
        "$GRAFANA_URL/api/folders")

    if echo "$CREATE_FOLDER_RESPONSE" | grep -q '"uid"'; then
        echo -e "${GREEN}Folder created successfully${NC}"
    else
        echo -e "${RED}Failed to create folder${NC}"
        echo "Response: $CREATE_FOLDER_RESPONSE"
        exit 1
    fi
else
    echo -e "${GREEN}Folder already exists${NC}"
fi

# Function to prepare rule file with correct datasource UID
prepare_rule_file() {
    local source_file="$1"
    local temp_file="$2"

    # Copy the source file to temp file first
    cp "$source_file" "$temp_file"

    # Replace datasource UID in the temp file (in-place)
    "${SED_INPLACE[@]}" "s/\"loki-uid\"/\"$DATASOURCE_UID\"/g" "$temp_file"

    # Also update the folderUID if needed
    if [ "$FOLDER_UID" != "LPS" ]; then
        "${SED_INPLACE[@]}" "s/\"folderUID\": \"LPS\"/\"folderUID\": \"$FOLDER_UID\"/g" "$temp_file"
    fi
}

# Function to create individual alert rule
create_alert_rule() {
    local rule_file="$1"
    local rule_name="$2"

    echo -e "${BLUE}Creating rule: $rule_name${NC}"

    # Prepare temporary file with correct datasource UID
    TEMP_RULE_FILE=$(mktemp)
    prepare_rule_file "$rule_file" "$TEMP_RULE_FILE"

    # Make a single API call and capture both response and HTTP status
    TEMP_RESPONSE_FILE=$(mktemp)
    HTTP_STATUS=$(curl -s -w "%{http_code}" -X POST \
        -u "$USERNAME:$PASSWORD" \
        -H "Content-Type: application/json" \
        -d "@$TEMP_RULE_FILE" \
        -o "$TEMP_RESPONSE_FILE" \
        "$GRAFANA_URL/api/v1/provisioning/alert-rules")

    CREATE_RESPONSE=$(cat "$TEMP_RESPONSE_FILE")

    # Cleanup temporary files
    rm "$TEMP_RULE_FILE" "$TEMP_RESPONSE_FILE"

    if [ "$HTTP_STATUS" = "201" ]; then
        RULE_ID=$(echo "$CREATE_RESPONSE" | grep -o '"id":[0-9]*' | cut -d':' -f2)
        echo -e "${GREEN}  Rule created successfully (ID: $RULE_ID)${NC}"
        return 0
    elif [ "$HTTP_STATUS" = "409" ]; then
        echo -e "${YELLOW}  Rule already exists (skipping)${NC}"
        return 0
    else
        echo -e "${RED}  Failed to create rule (HTTP $HTTP_STATUS)${NC}"
        echo "  Response: $CREATE_RESPONSE"
        return 1
    fi
}

# Check if all rule files exist before proceeding
echo ""
echo "Validating rule files..."
MISSING_FILES=()

for rule_file in "${RULE_FILES[@]}"; do
    rule_path="$SCRIPT_DIR/alerts/$rule_file"
    if [ ! -f "$rule_path" ]; then
        MISSING_FILES=("${MISSING_FILES[@]}" "$rule_path")
        echo -e "${RED}ERROR: Rule file not found: $rule_path${NC}"
    else
        echo -e "${GREEN}  Found: $rule_file${NC}"
    fi
done

if [ ${#MISSING_FILES[@]} -gt 0 ]; then
    echo -e "${RED}ERROR: ${#MISSING_FILES[@]} rule file(s) missing. Aborting.${NC}"
    exit 1
fi

# Import rules using iteration
echo ""
echo "Creating alert rules..."
FAILED_RULES=0
SUCCESSFUL_RULES=0

# Use index-based loop instead of associative arrays
for i in $(seq 0 $((${#RULE_FILES[@]} - 1))); do
    rule_file="${RULE_FILES[$i]}"
    rule_name="${RULE_NAMES[$i]}"
    rule_path="$SCRIPT_DIR/alerts/$rule_file"

    if create_alert_rule "$rule_path" "$rule_name"; then
        SUCCESSFUL_RULES=$((SUCCESSFUL_RULES + 1))
    else
        FAILED_RULES=$((FAILED_RULES + 1))
    fi
done

# Summary
echo ""
echo "Import Summary:"
echo "  Total rules: ${#RULE_FILES[@]}"
echo "  Successful: $SUCCESSFUL_RULES"
echo "  Failed: $FAILED_RULES"

if [ $FAILED_RULES -eq 0 ]; then
    echo -e "${GREEN}All alert rules imported successfully!${NC}"
else
    echo -e "${YELLOW}Import completed with $FAILED_RULES failed rules${NC}"
fi

# Verify imported rules
echo ""
echo "Verifying imported rules..."
LIST_RESPONSE=$(curl -s -u "$USERNAME:$PASSWORD" "$GRAFANA_URL/api/v1/provisioning/alert-rules")

# Build verification pattern from rule names
VERIFICATION_PATTERN=""
for rule_name in "${RULE_NAMES[@]}"; do
    if [ -z "$VERIFICATION_PATTERN" ]; then
        VERIFICATION_PATTERN="$rule_name"
    else
        VERIFICATION_PATTERN="$VERIFICATION_PATTERN\\|$rule_name"
    fi
done

if echo "$LIST_RESPONSE" | grep -q "$VERIFICATION_PATTERN"; then
    echo -e "${GREEN}Rules verification successful!${NC}"
    echo ""
    echo "Current rules in Grafana:"
    echo "$LIST_RESPONSE" | grep -o '"title":"[^"]*"' | sed 's/"title":"/ - /' | sed 's/"$//'
else
    echo -e "${YELLOW}Could not verify rules${NC}"
fi

echo ""
echo -e "${BLUE}Next steps:${NC}"
echo "1. Go to $GRAFANA_URL/alerting/list"
echo "2. Check that your rules are visible in the '$FOLDER_UID' folder"
echo "3. Configure contact points and notification policies if needed"
echo ""
echo -e "${GREEN}Import process completed!${NC}"
