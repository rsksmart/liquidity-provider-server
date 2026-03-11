#!/bin/bash

# Script to import LPS alert rules to Grafana using Provisioning API
# Usage: ./import-alerts.sh [grafana_url] [username] [password] [folder_uid] [datasource_uid]

# Configuration
GRAFANA_URL=${1:-"http://localhost:3000"}
USERNAME=${2:-"admin"}
PASSWORD=${3:-"test"}
FOLDER_UID=${4:-"LPS"}  # Default to LPS folder
DATASOURCE_UID=${5:-"loki-uid"}  # Default Loki datasource UID
ALERT_EMAIL=${ALERT_RECIPIENT_EMAIL:?"ALERT_RECIPIENT_EMAIL environment variable is required"}
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
    "hot-wallet-low-liquidity-warning.json"
    "hot-wallet-low-liquidity-critical.json"
)

RULE_NAMES=(
    "Node Eclipse Detection Alert"
    "PegIn Out of Liquidity Alert"
    "PegOut Out of Liquidity Alert"
    "LPS Penalization Alert"
    "Hot Wallet Low Liquidity Warning Alert"
    "Hot Wallet Critical Low Liquidity Alert"
)

echo "Importing LPS alert rules to Grafana..."
echo "Target: $GRAFANA_URL"
echo "Folder: $FOLDER_UID"
echo "Loki Datasource UID: $DATASOURCE_UID"
echo "Alert email: $ALERT_EMAIL"
echo "Rules to import: ${#RULE_FILES[@]}"
echo ""

# Check if Grafana is accessible
echo "Checking Grafana connectivity..."
HTTP_STATUS=$(curl -s -o /dev/null -w "%{http_code}" "$GRAFANA_URL/api/health")
if [ "$HTTP_STATUS" != "200" ]; then
    echo "ERROR: Cannot reach Grafana at $GRAFANA_URL (HTTP $HTTP_STATUS)"
    echo "Make sure Grafana is running and accessible."
    exit 1
fi
echo "Grafana is accessible"

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
        echo "Folder created successfully"
    else
        echo "Failed to create folder"
        echo "Response: $CREATE_FOLDER_RESPONSE"
        exit 1
    fi
else
    echo "Folder already exists"
fi

# Set up default email contact point for all alerts
DEFAULT_CONTACT_POINT_NAME="lps-email"
echo ""
echo "Checking if default contact point '$DEFAULT_CONTACT_POINT_NAME' exists..."
CONTACT_POINTS_RESPONSE=$(curl -s -u "$USERNAME:$PASSWORD" "$GRAFANA_URL/api/v1/provisioning/contact-points")
CONTACT_POINT_EXISTS=$(echo "$CONTACT_POINTS_RESPONSE" | grep -o "\"name\":\"$DEFAULT_CONTACT_POINT_NAME\"" || echo "not_found")

if [ "$CONTACT_POINT_EXISTS" == "not_found" ]; then
    echo "Creating default contact point '$DEFAULT_CONTACT_POINT_NAME' (email: $ALERT_EMAIL)..."
    CONTACT_POINT_CREATE_RESPONSE=$(curl -s -X POST \
        -u "$USERNAME:$PASSWORD" \
        -H "Content-Type: application/json" \
        -d "{\"name\":\"$DEFAULT_CONTACT_POINT_NAME\",\"type\":\"email\",\"settings\":{\"addresses\":\"$ALERT_EMAIL\",\"singleEmail\":true}}" \
        "$GRAFANA_URL/api/v1/provisioning/contact-points")

    if echo "$CONTACT_POINT_CREATE_RESPONSE" | grep -q '"uid"'; then
        echo "Default contact point created successfully"
    else
        echo "Failed to create default contact point"
        echo "Response: $CONTACT_POINT_CREATE_RESPONSE"
    fi
else
    echo "Default contact point already exists"
fi

# Import custom contact points from JSON files
echo ""
echo "Importing custom contact points..."

for contact_point_file in "$SCRIPT_DIR"/contact-points/*.json; do
    [ -f "$contact_point_file" ] || continue
    contact_point_basename=$(basename "$contact_point_file")
    echo "  Processing: $contact_point_basename"

    TEMP_CONTACT_POINT_FILE=$(mktemp)
    cp "$contact_point_file" "$TEMP_CONTACT_POINT_FILE"
    "${SED_INPLACE[@]}" "s/__ALERT_EMAIL__/$ALERT_EMAIL/g" "$TEMP_CONTACT_POINT_FILE"

    CONTACT_POINT_NAME=$(python3 -c "import json; print(json.load(open('$TEMP_CONTACT_POINT_FILE'))['name'])")
    CONTACT_POINT_CHECK=$(echo "$CONTACT_POINTS_RESPONSE" | grep -o "\"name\":\"$CONTACT_POINT_NAME\"" || echo "not_found")

    if [ "$CONTACT_POINT_CHECK" == "not_found" ]; then
        CONTACT_POINT_RESULT=$(curl -s -X POST \
            -u "$USERNAME:$PASSWORD" \
            -H "Content-Type: application/json" \
            -d "@$TEMP_CONTACT_POINT_FILE" \
            "$GRAFANA_URL/api/v1/provisioning/contact-points")

        if echo "$CONTACT_POINT_RESULT" | grep -q '"uid"'; then
            echo "  Contact point '$CONTACT_POINT_NAME' created"
            CONTACT_POINTS_RESPONSE=$(curl -s -u "$USERNAME:$PASSWORD" "$GRAFANA_URL/api/v1/provisioning/contact-points")
        else
            echo "  Failed to create contact point '$CONTACT_POINT_NAME'"
            echo "  Response: $CONTACT_POINT_RESULT"
        fi
    else
        echo "  Contact point '$CONTACT_POINT_NAME' already exists"
    fi
    rm "$TEMP_CONTACT_POINT_FILE"
done

# Build child routes from alert files that declare a __contact_point__
echo ""
echo "Building notification routes from alert rules..."
CHILD_ROUTES="[]"

for alert_file in "$SCRIPT_DIR"/alerts/*.json; do
    [ -f "$alert_file" ] || continue
    ROUTE_INFO=$(python3 -c "
import json
data = json.load(open('$alert_file'))
cp = data.get('__contact_point__')
if cp:
    print(data['title'] + '|' + cp)
")
    if [ -n "$ROUTE_INFO" ]; then
        ALERT_TITLE=$(echo "$ROUTE_INFO" | cut -d'|' -f1)
        CONTACT_POINT_NAME=$(echo "$ROUTE_INFO" | cut -d'|' -f2)
        echo "  Route: '$ALERT_TITLE' -> '$CONTACT_POINT_NAME'"
        CHILD_ROUTES=$(echo "$CHILD_ROUTES" | python3 -c "
import json, sys
routes = json.load(sys.stdin)
routes.append({
    'receiver': '$CONTACT_POINT_NAME',
    'object_matchers': [['alertname', '=', '$ALERT_TITLE']],
    'group_by': ['grafana_folder', 'alertname'],
    'group_wait': '10s',
    'group_interval': '1m',
    'repeat_interval': '5m'
})
json.dump(routes, sys.stdout)
")
    fi
done

# Set up notification policy with child routes for custom contact points
echo ""
echo "Configuring notification policy..."
POLICY_JSON=$(python3 -c "
import json, sys
routes = json.loads('$CHILD_ROUTES')
policy = {
    'receiver': '$DEFAULT_CONTACT_POINT_NAME',
    'group_by': ['grafana_folder', 'alertname'],
    'group_wait': '10s',
    'group_interval': '1m',
    'repeat_interval': '5m'
}
if routes:
    policy['routes'] = routes
json.dump(policy, sys.stdout)
")

POLICY_RESPONSE=$(curl -s -X PUT \
    -u "$USERNAME:$PASSWORD" \
    -H "Content-Type: application/json" \
    -d "$POLICY_JSON" \
    "$GRAFANA_URL/api/v1/provisioning/policies")

if echo "$POLICY_RESPONSE" | grep -q '"message"'; then
    echo "Notification policy configured successfully"
else
    echo "Failed to configure notification policy"
    echo "Response: $POLICY_RESPONSE"
fi

# Function to prepare rule file with correct datasource UID
prepare_rule_file() {
    local source_file="$1"
    local temp_file="$2"

    # Copy the source file to temp file first
    cp "$source_file" "$temp_file"

    # Replace datasource UID in the temp file (in-place)
    "${SED_INPLACE[@]}" "s/\"loki-uid\"/\"$DATASOURCE_UID\"/g" "$temp_file"

    # Strip __contact_point__ metadata field (not a valid Grafana API field)
    python3 -c "
import json
data = json.load(open('$temp_file'))
data.pop('__contact_point__', None)
json.dump(data, open('$temp_file', 'w'), indent=2)
"

    # Also update the folderUID if needed
    if [ "$FOLDER_UID" != "LPS" ]; then
        "${SED_INPLACE[@]}" "s/\"folderUID\": \"LPS\"/\"folderUID\": \"$FOLDER_UID\"/g" "$temp_file"
    fi
}

# Function to create individual alert rule
create_alert_rule() {
    local rule_file="$1"
    local rule_name="$2"

    echo "Creating rule: $rule_name"

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
        echo "  Rule created successfully (ID: $RULE_ID)"
        return 0
    elif [ "$HTTP_STATUS" = "409" ]; then
        echo "  Rule already exists (skipping)"
        return 0
    else
        echo "  Failed to create rule (HTTP $HTTP_STATUS)"
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
        echo "ERROR: Rule file not found: $rule_path"
    else
        echo "  Found: $rule_file"
    fi
done

if [ ${#MISSING_FILES[@]} -gt 0 ]; then
    echo "ERROR: ${#MISSING_FILES[@]} rule file(s) missing. Aborting."
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
    echo "All alert rules imported successfully!"
else
    echo "Import completed with $FAILED_RULES failed rules"
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
    echo "Rules verification successful!"
    echo ""
    echo "Current rules in Grafana:"
    echo "$LIST_RESPONSE" | grep -o '"title":"[^"]*"' | sed 's/"title":"/ - /' | sed 's/"$//'
else
    echo "Could not verify rules"
fi

echo ""
echo "Next steps:"
echo "1. Go to $GRAFANA_URL/alerting/list"
echo "2. Check that your rules are visible in the '$FOLDER_UID' folder"
echo "3. Notifications will be sent to $ALERT_EMAIL (default: '$DEFAULT_CONTACT_POINT_NAME', custom contact points for specific alerts)"
echo ""
echo "Import process completed!"
