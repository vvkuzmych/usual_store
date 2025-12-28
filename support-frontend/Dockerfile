# Build stage
FROM node:18-alpine AS builder

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy source code
COPY . .

# Build the application with PUBLIC_URL=/admin so assets are served from /admin/static/*
RUN PUBLIC_URL=/admin npm run build

# Production stage
FROM nginx:alpine

# Copy built files from builder to root (for /support route)
COPY --from=builder /app/build /usr/share/nginx/html

# Also copy to /admin directory for /admin/users route
COPY --from=builder /app/build /usr/share/nginx/html/admin

# Copy nginx configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expose port
EXPOSE 3005

CMD ["nginx", "-g", "daemon off;"]

