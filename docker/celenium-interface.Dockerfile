# Multi-stage build for celenium-interface
# This Dockerfile clones and builds celenium-interface from source

FROM node:20-alpine AS builder

# Install git for cloning
RUN apk add --no-cache git

# Clone the celenium-interface repository
WORKDIR /build
RUN git clone --depth 1 https://github.com/celenium-io/celenium-interface.git .

# Install dependencies and build
# Use corepack to enable pnpm (comes with Node 20+)
RUN corepack enable && corepack prepare pnpm@latest --activate
RUN pnpm install --frozen-lockfile
# Increase Node.js heap size to avoid out of memory errors during build
ENV NODE_OPTIONS="--max-old-space-size=4096"
RUN pnpm build

# Production stage
FROM node:20-alpine

WORKDIR /app

# Copy built application from builder
# Nuxt 3 outputs everything needed in .output directory
COPY --from=builder /build/.output /app/.output

# Expose port
EXPOSE 3000

# Set NODE_ENV
ENV NODE_ENV=production

# Start the application
CMD ["node", ".output/server/index.mjs"]
