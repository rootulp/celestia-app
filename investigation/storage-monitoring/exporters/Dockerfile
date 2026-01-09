FROM python:3.11-alpine

# Install du command (part of coreutils)
RUN apk add --no-cache coreutils

# Install Python dependencies
RUN pip install --no-cache-dir prometheus_client

# Copy the exporter script
COPY directory-size-exporter.py /app/directory-size-exporter.py

# Make the script executable
RUN chmod +x /app/directory-size-exporter.py

# Expose the metrics port
EXPOSE 9101

# Run the exporter
CMD ["python3", "/app/directory-size-exporter.py"]
