FROM nginx:alpine

# Install supervisor for process management
RUN apk add --no-cache supervisor

# Copy pre-built frontend static files
COPY frontend/build/ /usr/share/nginx/html/

# Copy pre-built Go binary
COPY backend/golang/quotes-api /usr/local/bin/quotes-api
RUN chmod +x /usr/local/bin/quotes-api

# Copy Nginx configuration
COPY nginx.conf /etc/nginx/nginx.conf

# Copy supervisor configuration
COPY supervisord.conf /etc/supervisor/conf.d/supervisord.conf

# Copy startup script
COPY start.sh /start.sh
RUN chmod +x /start.sh

EXPOSE 80

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost/health || exit 1

CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"]