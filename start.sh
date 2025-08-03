#!/bin/sh

# Set default environment variables
export PORT=${PORT:-8080}

# Create log directories
mkdir -p /var/log/supervisor

# Start supervisord
exec /usr/bin/supervisord -c /etc/supervisor/conf.d/supervisord.conf