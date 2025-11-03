# --- STAGE 1: "The Builder" ---
# Start from the official postgres image to get build environment
FROM postgres:18 AS builder

# 1. Install build dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    git \
    postgresql-server-dev-18 \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# 2. Clone and build pg_cron  
RUN git clone --branch v1.6.7 --depth 1 https://github.com/citusdata/pg_cron.git /pg_cron \
    && cd /pg_cron \
    && make \
    && make install

# --- STAGE 2: "The Final Image" ---
# Start fresh from the clean official postgres image
FROM postgres:18

# 1. Copy *only* the compiled extension files from the "builder" stage
# These are the files that 'make install' created
COPY --from=builder /usr/lib/postgresql/18/lib/pg_cron.so /usr/lib/postgresql/18/lib/
COPY --from=builder /usr/share/postgresql/18/extension/pg_cron.control /usr/share/postgresql/18/extension/
COPY --from=builder /usr/share/postgresql/18/extension/pg_cron--*.sql /usr/share/postgresql/18/extension/
