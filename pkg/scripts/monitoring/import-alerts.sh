#!/bin/bash

# Script to import LPS alert rules to Grafana using Provisioning API
# Usage: ./import-alerts.sh [grafana_url] [username] [password] [folder_uid]

# Configuration
GRAFANA_URL=${1:-"http://localhost:3000"}
USERNAME=${2:-"admin"}
PASSWORD=${3:-"test"}
FOLDER_UID=${4:-"LPS"}  # Default to LPS folder
SCRIPT_DIR="$(dirname "$0")"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${YELLOW}Importing LPS alert rules to Grafana...${NC}"
echo "Target: $GRAFANA_URL"
echo "Folder: $FOLDER_UID"
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

# Function to create individual alert rule
create_alert_rule() {
    local rule_file="$1"
    local rule_name="$2"

    echo -e "${BLUE}Creating rule: $rule_name${NC}"

    # Make a single API call and capture both response and HTTP status
    TEMP_FILE=$(mktemp)
    HTTP_STATUS=$(curl -s -w "%{http_code}" -X POST \
        -u "$USERNAME:$PASSWORD" \
        -H "Content-Type: application/json" \
        -d "@$rule_file" \
        -o "$TEMP_FILE" \
        "$GRAFANA_URL/api/v1/provisioning/alert-rules")

    CREATE_RESPONSE=$(cat "$TEMP_FILE")
    rm "$TEMP_FILE"

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

# Create alert rules
echo ""
echo "Creating alert rules..."

# Check if rule files exist
if [ ! -f "$SCRIPT_DIR/alerts/node-eclipse-detection.json" ]; then
    echo -e "${RED}ERROR: Node Eclipse Detection rule file not found: $SCRIPT_DIR/alerts/node-eclipse-detection.json${NC}"
    exit 1
fi

if [ ! -f "$SCRIPT_DIR/alerts/pegin-out-of-liquidity.json" ]; then
    echo -e "${RED}ERROR: PegIn Liquidity rule file not found: $SCRIPT_DIR/alerts/pegin-out-of-liquidity.json${NC}"
    exit 1
fi

if [ ! -f "$SCRIPT_DIR/alerts/pegout-out-of-liquidity.json" ]; then
    echo -e "${RED}ERROR: PegOut Liquidity rule file not found: $SCRIPT_DIR/alerts/pegout-out-of-liquidity.json${NC}"
    exit 1
fi

if [ ! -f "$SCRIPT_DIR/alerts/lps-penalization.json" ]; then
    echo -e "${RED}ERROR: LPS Penalization rule file not found: $SCRIPT_DIR/alerts/lps-penalization.json${NC}"
    exit 1
fi

# Import rules
FAILED_RULES=0

if ! create_alert_rule "$SCRIPT_DIR/alerts/node-eclipse-detection.json" "Node Eclipse Detection Alert"; then
    FAILED_RULES=$((FAILED_RULES + 1))
fi

if ! create_alert_rule "$SCRIPT_DIR/alerts/pegin-out-of-liquidity.json" "PegIn Out of Liquidity Alert"; then
    FAILED_RULES=$((FAILED_RULES + 1))
fi

if ! create_alert_rule "$SCRIPT_DIR/alerts/pegout-out-of-liquidity.json" "PegOut Out of Liquidity Alert"; then
    FAILED_RULES=$((FAILED_RULES + 1))
fi

if ! create_alert_rule "$SCRIPT_DIR/alerts/lps-penalization.json" "LPS Penalization Alert"; then
    FAILED_RULES=$((FAILED_RULES + 1))
fi

# Summary
echo ""
if [ $FAILED_RULES -eq 0 ]; then
    echo -e "${GREEN}All alert rules imported successfully!${NC}"
else
    echo -e "${YELLOW}Import completed with $FAILED_RULES failed rules${NC}"
fi

# Verify imported rules
echo ""
echo "Verifying imported rules..."
LIST_RESPONSE=$(curl -s -u "$USERNAME:$PASSWORD" "$GRAFANA_URL/api/v1/provisioning/alert-rules")

if echo "$LIST_RESPONSE" | grep -q "Node Eclipse Detection Alert\|PegIn Out of Liquidity Alert\|PegOut Out of Liquidity Alert\|LPS Penalization Alert"; then
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
