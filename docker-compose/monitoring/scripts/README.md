# LPS Alert Rules

This directory contains the alert rules for the Liquidity Provider Server (LPS) monitoring system.

## Files

### Alert Definitions (in alerts/ subdirectory)
- `alerts/node-eclipse-detection.json` - Detects potential eclipse attacks on Bitcoin nodes
- `alerts/pegin-out-of-liquidity.json` - Alerts when PegIn operations run out of Bitcoin liquidity
- `alerts/pegout-out-of-liquidity.json` - Alerts when PegOut operations run out of BTC liquidity
- `alerts/lps-penalization.json` - Alerts when the LPS has been penalized
- `alerts/hot-wallet-low-liquidity-warning.json` - Alerts when hot wallet liquidity drops below warning threshold
- `alerts/hot-wallet-low-liquidity-critical.json` - Alerts when hot wallet liquidity is critically low

### Custom Contact Points (in contact-points/ subdirectory)
- `contact-points/low-liquidity.json` - Custom email format for low liquidity alerts (used by both warning and critical alerts)

### Import Script
- `import-alerts.sh` - Script to import alert rules, contact points, and notification policies into any Grafana instance

## Usage

### Prerequisites
- `ALERT_RECIPIENT_EMAIL` environment variable must be set (e.g., via `source .env.regtest`)
- Grafana with unified alerting enabled
- Loki datasource configured (default UID: `loki-uid`, customizable via script parameter)
- `curl` and `python3` available

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

### Environment Variables
- `ALERT_RECIPIENT_EMAIL` (**required**) - Recipient email for alert notifications. Read from the environment (set in `.env` files).

## Directory Structure

```
docker-compose/monitoring/scripts/
├── import-alerts.sh                            # Main import script
├── README.md                                   # This documentation
├── alerts/
│   ├── node-eclipse-detection.json             # Eclipse attack alert
│   ├── pegin-out-of-liquidity.json             # PegIn liquidity alert
│   ├── pegout-out-of-liquidity.json            # PegOut liquidity alert
│   ├── lps-penalization.json                   # LPS penalization alert
│   ├── hot-wallet-low-liquidity-warning.json   # Hot wallet low liquidity warning (regex extraction)
│   └── hot-wallet-low-liquidity-critical.json  # Hot wallet critical low liquidity (regex extraction)
└── contact-points/
    └── low-liquidity.json                      # Custom email format for low liquidity alerts
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

### Hot Wallet Low Liquidity Warning Alert
- **Trigger**: When log contains "Alert! - Subject: Hot wallet: Low liquidity, refill recommended"
- **Purpose**: Alerts when the hot wallet liquidity is below the warning threshold
- **Dynamic extraction**: Uses LogQL `regexp` to extract `network`, `current`, and `threshold` from the log body and includes them in the notification via `{{ $labels.xxx }}`
- **Custom contact point**: Routed to `lps-email-low-liquidity` via `__contact_point__` for a tailored email format

### Hot Wallet Critical Low Liquidity Alert
- **Trigger**: When log contains "Alert! - Subject: Hot wallet: Critical low liquidity, refill required"
- **Purpose**: Alerts when the hot wallet liquidity is critically low and an immediate refill is required
- **Dynamic extraction**: Same as the warning alert -- extracts `network`, `current`, and `threshold` via `regexp`
- **Custom contact point**: Shares `lps-email-low-liquidity` with the warning alert via `__contact_point__`

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

### Contact Points and Notification Policy

In Grafana, a **contact point** bundles together the delivery channel (email, Slack, etc.), the recipient address, and the message template (subject, body format). In our setup the recipient is always the same (`ALERT_RECIPIENT_EMAIL`), so custom contact points are used solely to provide different **email formats** per alert type -- not different recipients.

The import script automatically:
- Creates a default `lps-email` contact point for all alerts, using the `ALERT_RECIPIENT_EMAIL` environment variable as the recipient
- Imports custom contact points from the `contact-points/` directory (e.g., `lps-email-low-liquidity` with a tailored subject/message format)
- Builds child routes from alert rules that declare a `__contact_point__` field, routing them to the named contact point
- Configures the notification policy with:
  - A root route that sends all alerts to `lps-email` (default Grafana email format)
  - Child routes that match specific alerts to their custom contact points (matched by `alertname`)
- Sets notification timing: `group_wait: 10s`, `group_interval: 1m`, `repeat_interval: 5m`

### Routing Alerts to a Custom Contact Point

To give an alert a custom email format:
1. Create a contact point JSON in `contact-points/` with `__ALERT_EMAIL__` as a placeholder for the recipient address (or reuse an existing one)
2. In the alert rule JSON (in `alerts/`), add a `"__contact_point__": "<contact-point-name>"` field referencing the contact point's `name`
3. The script will strip `__contact_point__` before sending to Grafana and automatically create a child route matching the alert to that contact point
4. Multiple alerts can share the same contact point by referencing the same name

## Notes

- Rules are created in the specified folder (default: LPS)
- Script automatically creates the folder if it doesn't exist
- Duplicate rules and contact points are skipped (no error)
- Uses Grafana Provisioning API for reliable imports
- Script looks for JSON files in the `alerts/` and `contact-points/` subdirectories relative to its location
- Rules created using this script cannot be edited in the Grafana UI
- Custom contact points have `disableResolveMessage` set to prevent a second email when the alert resolves
