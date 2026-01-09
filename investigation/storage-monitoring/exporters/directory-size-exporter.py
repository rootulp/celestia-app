#!/usr/bin/env python3
"""
Directory Size Exporter for Prometheus

This exporter monitors the size of specific directories in the celestia-app data directory
and exposes them as Prometheus metrics.
"""

import os
import subprocess
import time
from http.server import HTTPServer
from prometheus_client import Gauge, start_http_server

# Define the directories to monitor
DIRECTORIES = [
    "application.db",
    "blockstore.db",
    "cs.wal",
    "evidence.db",
    "fibre-store",
    "snapshots",
    "state.db",
    "tx_index.db",
]

# Create a Prometheus Gauge metric
directory_size_bytes = Gauge(
    'celestia_data_directory_size_bytes',
    'Size of celestia-app data directories in bytes',
    ['directory']
)

def get_directory_size(directory_path):
    """
    Get the size of a directory in bytes using du -sb.
    Returns 0 if the directory doesn't exist or can't be accessed.
    """
    if not os.path.exists(directory_path):
        return 0

    try:
        # Use du -sb to get size in bytes
        result = subprocess.run(
            ['du', '-sb', directory_path],
            capture_output=True,
            text=True,
            timeout=30,
            check=False
        )
        if result.returncode == 0:
            # du -sb output format: "size\tpath"
            size_str = result.stdout.split('\t')[0]
            return int(size_str)
        else:
            print(f"Warning: Failed to get size for {directory_path}: {result.stderr}")
            return 0
    except (subprocess.TimeoutExpired, ValueError, OSError) as e:
        print(f"Error getting size for {directory_path}: {e}")
        return 0

def collect_metrics(data_dir):
    """
    Collect metrics for all monitored directories.
    """
    for directory in DIRECTORIES:
        directory_path = os.path.join(data_dir, directory)
        size_bytes = get_directory_size(directory_path)
        directory_size_bytes.labels(directory=directory).set(size_bytes)
        print(f"Updated {directory}: {size_bytes} bytes ({size_bytes / (1024**3):.2f} GB)")

def main():
    """
    Main function to start the exporter.
    """
    # Get data directory from environment variable or use default
    data_dir = os.environ.get('DATA_DIR', '/data')

    if not os.path.exists(data_dir):
        print(f"Warning: Data directory {data_dir} does not exist yet. Will continue monitoring...")

    print(f"Starting directory size exporter for: {data_dir}")
    print(f"Monitoring directories: {', '.join(DIRECTORIES)}")

    # Start Prometheus HTTP server on port 9101
    start_http_server(9101)
    print("Metrics server started on port 9101")
    print("Metrics available at http://localhost:9101/metrics")

    # Collect metrics every 30 seconds
    while True:
        try:
            collect_metrics(data_dir)
        except Exception as e:
            print(f"Error collecting metrics: {e}")

        time.sleep(30)

if __name__ == '__main__':
    main()
