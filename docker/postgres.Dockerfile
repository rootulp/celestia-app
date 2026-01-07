# Custom PostgreSQL image with TimescaleDB and Toolkit extensions
FROM timescale/timescaledb:latest-pg15

# Install build dependencies for compiling timescaledb_toolkit (Rust-based)
RUN apk add --no-cache \
    build-base \
    git \
    postgresql-dev \
    clang \
    llvm \
    llvm-dev \
    rust \
    cargo \
    && rm -rf /var/cache/apk/*

# Clone and build timescaledb_toolkit
# The toolkit is written in Rust, so we need cargo to build it
WORKDIR /tmp
RUN git clone --depth 1 https://github.com/timescale/timescaledb-toolkit.git \
    && cd timescaledb-toolkit \
    && PG_CONFIG=/usr/local/bin/pg_config cargo build --release \
    && PG_CONFIG=/usr/local/bin/pg_config cargo install --path . \
    && cd .. \
    && rm -rf timescaledb-toolkit \
    && apk del build-base git clang llvm llvm-dev rust cargo

WORKDIR /
