#!/bin/sh
set -e

# Create runtime-config.js so the browser can read API URL at runtime
RUNTIME_PUBLIC_API_URL=${RUNTIME_PUBLIC_API_URL:-${NEXT_PUBLIC_API_URL:-http://localhost:8000}}

cat > /app/public/runtime-config.js <<EOF
window.__RUNTIME_CONFIG__ = {
  API_URL: "${RUNTIME_PUBLIC_API_URL}"
};
EOF

exec "$@"
