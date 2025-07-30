# LPS Alert Rules

This directory contains the alert rules for the Liquidity Provider Server (LPS) monitoring system.

## Files

### Alert Definitions (in alerts/ subdirectory)
- `alerts/node-eclipse-detection.json` - Detects potential eclipse attacks on Bitcoin nodes
- `alerts/pegin-out-of-liquidity.json` - Alerts when PegIn operations run out of Bitcoin liquidity
- `alerts/pegout-out-of-liquidity.json` - Alerts when PegOut operations run out of BTC liquidity
- `alerts/lps-penalization.json` - Alerts when the LPS has been penalized

### Import Script
- `import-alerts.sh` - Script to import alert rules into any Grafana instance

## Usage

### Import Alerts to Grafana

```bash
# From project root - basic usage (localhost:3000 with admin/test)
./docker-compose/monitoring/scripts/import-alerts.sh

# Custom Grafana instance
./docker-compose/monitoring/scripts/import-alerts.sh http://grafana.example.com:3000 admin password

# Different folder
./docker-compose/monitoring/scripts/import-alerts.sh http://localhost:3000 admin test ALERTS

# Custom datasource UID
./docker-compose/monitoring/scripts/import-alerts.sh http://localhost:3000 admin test LPS my-loki-uid

# Or run from the monitoring scripts directory
cd docker-compose/monitoring/scripts
./import-alerts.sh
```

### Parameters
1. `grafana_url` - Grafana instance URL (default: http://localhost:3000)
2. `username` - Grafana username (default: admin)
3. `password` - Grafana password (default: test)
4. `folder_uid` - Folder UID for alerts (default: LPS)
5. `datasource_uid` - Loki datasource UID (default: loki-uid)

## Directory Structure

```
docker-compose/monitoring/
├── scripts/
│   ├── import-alerts.sh          # Main import script
│   └── README.md                 # This documentation
└── alerts/
    ├── node-eclipse-detection.json     # Eclipse attack alert
    ├── pegin-out-of-liquidity.json     # PegIn liquidity alert
    ├── pegout-out-of-liquidity.json    # PegOut liquidity alert
    └── lps-penalization.json           # LPS penalization alert
```

## Alert Details

### Node Eclipse Detection Alert
- **Trigger**: When log contains "Alert! - Subject: Node Eclipse Detected"
- **Purpose**: Detects potential eclipse attacks on Bitcoin nodes

### PegIn Out of Liquidity Alert
- **Trigger**: When log contains "Alert! - Subject: PegIn: Out of liquidity"
- **Purpose**: Alerts when insufficient liquidity for PegIn operations

### PegOut Out of Liquidity Alert
- **Trigger**: When log contains "Alert! - Subject: PegOut: Out of liquidity"
- **Purpose**: Alerts when insufficient liquidity for PegOut operations

### LPS Penalization Alert
- **Trigger**: When log contains "Alert! - Subject: LPS has been penalized"
- **Purpose**: Alerts when the Liquidity Provider has been penalized for failing to fulfill quote commitments

## Configuration Details

### Alert Rule Settings
All alert rules are configured with:
- `noDataState: OK` - Prevents "DatasourceNoData" alerts when no logs match
- `execErrState: OK` - Prevents error alerts on query execution issues
- `for: "0s"` - Fires immediately when condition is met (no pending period)

### Datasource UID Configuration
The alert JSON files use `"datasourceUid": "loki-uid"` by default.

**Key Differences:**
- **Alert Rules (API Import)**: Use concrete UIDs like `"loki-uid"`
- **Dashboard Templates**: Use template variables like `"${DS_LOKI}"`

The import script automatically replaces the datasource UID if you specify a different one via the `datasource_uid` parameter, making it portable across different Grafana instances.

## Requirements

- Grafana with unified alerting enabled
- Loki datasource configured (default UID: `loki-uid`, customizable via script parameter)
- curl command available
- jq for JSON processing (optional, for verification)

## Notes

- Rules are created in the specified folder (default: LPS)
- Script automatically creates the folder if it doesn't exist
- Duplicate rules are skipped (no error)
- Uses Grafana Provisioning API for reliable imports
- Script looks for JSON files in the `alerts/` subdirectory relative to its location
- Rules created using this script cannot be edited in the Grafana UI
- To avoid a second email being sent once the rule is resolved, the contact point must be changed to "Disable resolved message" in Grafana UI
