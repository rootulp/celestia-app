# Celestia App Storage Monitoring

This directory contains a complete monitoring stack for investigating storage growth in celestia-app. The setup includes:

- **celestia-app**: A node that state syncs to Mocha testnet
- **Prometheus**: Metrics collection and storage
- **Grafana**: Visualization and dashboards
- **node-exporter**: System-level metrics
- **directory-size-exporter**: Custom exporter for monitoring individual data directory sizes

## Architecture

```
┌─────────────────┐   scrape (26660)   ┌─────────────────┐
│  celestia-app   │ ─────────────────► │   Prometheus    │
│  (port 26660)   │                    │  (port 9090)    │
└─────────────────┘                    └────────┬────────┘
                                                │
┌─────────────────┐   scrape (9100)           │
│ node-exporter   │ ──────────────────────────┤
│  (port 9100)    │                            │
└─────────────────┘                            │
                                                │ data source
┌─────────────────┐   scrape (9101)            │
│directory-exporter│ ──────────────────────────┤
│  (port 9101)    │                            │
└─────────────────┘                            ▼
                                         ┌─────────────┐
                                         │   Grafana   │
                                         │ (port 3000) │
                                         └─────────────┘
```

## Prerequisites

- Docker and Docker Compose installed
- At least 50GB of free disk space (for data storage)
- Network access to Mocha testnet RPC endpoints

## Quick Start

1. **Navigate to this directory:**
   ```bash
   cd investigation/storage-monitoring
   ```

2. **Start all services:**
   ```bash
   docker-compose up -d
   ```

3. **Check service status:**
   ```bash
   docker-compose ps
   ```

4. **View logs:**
   ```bash
   # All services
   docker-compose logs -f

   # Specific service
   docker-compose logs -f celestia-app
   ```

5. **Access Grafana:**
   - Open your browser to `http://localhost:3000`
   - Default credentials: `admin` / `admin`
   - The storage monitoring dashboard should load automatically

## Services

### celestia-app

- **Purpose**: Runs a celestia-app node that state syncs to Mocha testnet
- **Ports**:
  - `26657`: RPC endpoint
  - `26660`: Prometheus metrics
- **Data**: Stored in Docker volume `celestia-app-data`
- **Initialization**: Automatically configured on first run via `scripts/init-celestia-app.sh`

### Prometheus

- **Purpose**: Collects and stores metrics from all exporters
- **Port**: `9090`
- **Access**: `http://localhost:9090`
- **Retention**: 30 days
- **Configuration**: `prometheus/prometheus.yml`

### Grafana

- **Purpose**: Visualizes metrics in dashboards
- **Port**: `3000`
- **Access**: `http://localhost:3000`
- **Default credentials**: `admin` / `admin` (change via `GRAFANA_PASSWORD` env var)
- **Dashboard**: Auto-loaded from `grafana/dashboards/storage-monitoring.json`

### node-exporter

- **Purpose**: Exposes system-level metrics (CPU, memory, disk, etc.)
- **Port**: `9100`
- **Metrics**: Standard Prometheus node exporter metrics

### directory-size-exporter

- **Purpose**: Monitors the size of individual data directories
- **Port**: `9101`
- **Metrics**: `celestia_data_directory_size_bytes{directory="..."}`
- **Monitored directories**:
  - `application.db`
  - `blockstore.db`
  - `cs.wal`
  - `evidence.db`
  - `fibre-store`
  - `snapshots`
  - `state.db`
  - `tx_index.db`

## Dashboard

The Grafana dashboard includes:

1. **Total Data Directory Size**: Sum of all directory sizes
2. **Storage Growth Rate**: Rate of change for total storage
3. **Per-Directory Size Panels**: Individual graphs for each data directory
4. **Directory Growth Rates**: Growth rate per directory over time
5. **System Disk Usage**: Overall filesystem usage from node-exporter
6. **All Directory Sizes (Stacked)**: Stacked view of all directories

## Configuration

### Environment Variables

- `GRAFANA_PASSWORD`: Grafana admin password (default: `admin`)

### Customizing State Sync

Edit `scripts/init-celestia-app.sh` to modify:
- RPC endpoints
- Trust height calculation
- Seeds and peers

### Customizing Monitoring

- **Prometheus scrape interval**: Edit `prometheus/prometheus.yml`
- **Directory exporter interval**: Edit `exporters/directory-size-exporter.py` (default: 30s)
- **Grafana dashboard**: Edit `grafana/dashboards/storage-monitoring.json`

## Troubleshooting

### Services won't start

1. Check Docker is running: `docker ps`
2. Check port conflicts: `netstat -tuln | grep -E '3000|9090|26657|26660|9100|9101'`
3. View logs: `docker-compose logs`

### celestia-app not syncing

1. Check RPC endpoints are accessible
2. View logs: `docker-compose logs celestia-app`
3. Check state sync configuration: `docker-compose exec celestia-app cat /home/celestia/.celestia-app/config/config.toml | grep -A 10 "\[statesync\]"`

### Metrics not appearing in Grafana

1. Verify Prometheus is scraping: `curl http://localhost:9090/api/v1/targets`
2. Check exporter endpoints:
   - celestia-app: `curl http://localhost:26660/metrics`
   - node-exporter: `curl http://localhost:9100/metrics`
   - directory-exporter: `curl http://localhost:9101/metrics`
3. Verify Grafana datasource: Check `http://localhost:3000/datasources`

### Dashboard not loading

1. Check dashboard file exists: `ls -la grafana/dashboards/storage-monitoring.json`
2. Verify JSON is valid: `python3 -m json.tool grafana/dashboards/storage-monitoring.json`
3. Check Grafana logs: `docker-compose logs grafana`

## Data Persistence

All data is stored in Docker volumes:

- `celestia-app-data`: celestia-app data directory
- `celestia-app-config`: celestia-app configuration
- `prometheus-data`: Prometheus time-series database
- `grafana-data`: Grafana dashboards and settings

To remove all data:
```bash
docker-compose down -v
```

## Stopping Services

```bash
# Stop services (keeps data)
docker-compose stop

# Stop and remove containers (keeps data)
docker-compose down

# Stop and remove everything including data
docker-compose down -v
```

## Monitoring Over Time

This setup is designed to run continuously to monitor storage growth. The dashboard will show:

- Real-time directory sizes
- Growth rates over time
- Which directories are growing fastest
- Total storage usage trends

## Contributing

When making changes:

1. Test locally with `docker-compose up`
2. Verify all services start correctly
3. Check metrics are being collected
4. Verify dashboard displays correctly

## Related Issues

This monitoring setup was created to investigate:
- [celestia-core #2606](https://github.com/celestiaorg/celestia-core/issues/2606): Storage growing indefinitely despite pruning
