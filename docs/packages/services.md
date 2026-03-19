# ⚙️ Services

Additional services to enhance your server capabilities.

## Overview

| Service | Purpose | Use Case |
|---------|---------|----------|
| **Docker** | Containerization | Microservices, isolated apps |
| **Portainer** | Container Management | Docker GUI management |

---

## Docker

Platform for developing, shipping, and running applications in containers.

### Features
- Container isolation
- Image-based deployment
- Docker Hub integration
- Docker Compose support
- Resource management

### Install
```bash
hostkit install docker
```

### What's Installed
- Docker Engine
- Docker CLI
- Docker Compose (plugin)
- Containerd

### Basic Commands

#### Images
```bash
# List images
docker images

# Pull image
docker pull nginx:latest

# Remove image
docker rmi nginx:latest

# Build image
docker build -t myapp:latest .
```

#### Containers
```bash
# List running containers
docker ps

# List all containers
docker ps -a

# Run container
docker run -d --name mynginx -p 80:80 nginx

# Stop container
docker stop mynginx

# Start container
docker start mynginx

# Remove container
docker rm mynginx

# View logs
docker logs mynginx

# Execute command in container
docker exec -it mynginx bash
```

#### Docker Compose
```bash
# Start services
docker compose up -d

# Stop services
docker compose down

# View logs
docker compose logs -f

# Build and start
docker compose up -d --build
```

### Docker Compose Example
Create `docker-compose.yml`:

```yaml
version: '3.8'

services:
  web:
    image: nginx:latest
    ports:
      - "80:80"
    volumes:
      - ./html:/usr/share/nginx/html
    restart: unless-stopped

  db:
    image: mariadb:latest
    environment:
      MYSQL_ROOT_PASSWORD: secret
      MYSQL_DATABASE: myapp
    volumes:
      - db_data:/var/lib/mysql
    restart: unless-stopped

volumes:
  db_data:
```

Run:
```bash
docker compose up -d
```

### User Permissions
Add your user to docker group to run without sudo:

```bash
sudo usermod -aG docker $USER
newgrp docker
```

### Common Use Cases

#### Run WordPress
```yaml
version: '3.8'

services:
  wordpress:
    image: wordpress:latest
    ports:
      - "8080:80"
    environment:
      WORDPRESS_DB_HOST: db:3306
      WORDPRESS_DB_USER: wordpress
      WORDPRESS_DB_PASSWORD: wordpress
      WORDPRESS_DB_NAME: wordpress
    volumes:
      - wordpress_data:/var/www/html

  db:
    image: mariadb:latest
    environment:
      MYSQL_ROOT_PASSWORD: root_secret
      MYSQL_DATABASE: wordpress
      MYSQL_USER: wordpress
      MYSQL_PASSWORD: wordpress
    volumes:
      - db_data:/var/lib/mysql

volumes:
  wordpress_data:
  db_data:
```

#### Run Redis
```bash
docker run -d --name redis -p 6379:6379 redis:latest
```

#### Run PostgreSQL
```bash
docker run -d \
  --name postgres \
  -e POSTGRES_PASSWORD=secret \
  -p 5432:5432 \
  -v pgdata:/var/lib/postgresql/data \
  postgres:latest
```

### Resource Limits
```bash
docker run -d \
  --name limited-app \
  --memory="512m" \
  --cpus="1.0" \
  nginx:latest
```

---

## Portainer

Web-based UI for managing Docker containers and environments.

### Features
- Visual container management
- Image management
- Volume management
- Network management
- User management
- Multi-environment support

### Install
```bash
hostkit install portainer
```

### Access
```
URL: http://YOUR_IP:9000
```

### First Setup
1. Navigate to `http://YOUR_IP:9000`
2. Create admin account
3. Select Docker environment
4. Connect to local Docker socket

### Basic Usage

#### Dashboard
- View all containers
- Resource usage statistics
- Quick actions (start/stop/restart)

#### Container Management
- Create containers from UI
- View logs in real-time
- Access console/terminal
- View resource usage

#### Stack Deployment
Deploy full applications using Docker Compose:

1. Go to "Stacks" in sidebar
2. Click "Add stack"
3. Choose:
   - Web editor (paste compose file)
   - Upload (upload compose file)
   - Git repository
4. Click "Deploy the stack"

#### Templates
Portainer includes pre-built templates for common applications:
- WordPress
- MySQL
- Nginx
- Redis
- And more...

### Security
```bash
# Change default port
docker run -d -p 8000:8000 -p 9443:9443 --name portainer \
  --restart=always \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v portainer_data:/data \
  portainer/portainer-ce:latest
```

### Multi-Environment Management
Portainer can manage multiple Docker hosts:

1. Go to "Environments"
2. Click "Add environment"
3. Choose type:
   - Docker Standalone
   - Docker Swarm
   - Kubernetes
4. Follow setup instructions

---

## Docker Best Practices

### 1. Use Official Images
```bash
# Good
docker pull nginx:latest

# Avoid
docker pull randomuser/nginx
```

### 2. Pin Image Versions
```bash
# Better for production
docker pull nginx:1.25.3

# Instead of
docker pull nginx:latest
```

### 3. Use .dockerignore
Create `.dockerignore`:
```
.git
.env
node_modules
vendor
*.log
```

### 4. Minimize Layers
```dockerfile
# Bad
RUN apt-get update
RUN apt-get install -y curl
RUN apt-get install -y nginx

# Good
RUN apt-get update && apt-get install -y \
    curl \
    nginx \
    && rm -rf /var/lib/apt/lists/*
```

### 5. Use Multi-Stage Builds
```dockerfile
# Build stage
FROM golang:1.21 AS builder
WORKDIR /app
COPY . .
RUN go build -o myapp

# Runtime stage
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/myapp .
CMD ["./myapp"]
```

### 6. Don't Run as Root
```dockerfile
# Create non-root user
RUN useradd -m myuser
USER myuser
```

### 7. Resource Limits
Always set memory and CPU limits in production.

### 8. Health Checks
```dockerfile
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:80/ || exit 1
```

### 9. Logging
```bash
# View logs
docker logs container_name

# Follow logs
docker logs -f container_name

# Limit log size in compose
driver: "json-file"
options:
  max-size: "10m"
  max-file: "3"
```

### 10. Cleanup
```bash
# Remove unused containers
docker container prune

# Remove unused images
docker image prune

# Remove unused volumes
docker volume prune

# Remove everything unused
docker system prune -a
```
