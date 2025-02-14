FROM node:18-alpine

# Set working directory to monitoring project
WORKDIR /app/monitoring

# Copy only the monitoring project files
COPY docker-compose/monitoring ./

# Install dependencies and build
RUN npm install && npm run build

# Expose port 8080 for metrics
EXPOSE 8080

# Run the monitoring service
CMD ["npx", "ts-node", "src/index.ts"] 