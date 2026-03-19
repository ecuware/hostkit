# 🗄️ Databases

HostKit supports 5 popular database servers with easy installation and management.

## Overview

| Database | Type | Best For | Size |
|----------|------|----------|------|
| **MariaDB** | Relational | General purpose, MySQL compatible | 200MB |
| **MySQL** | Relational | Enterprise, Oracle compatibility | 400MB |
| **PostgreSQL** | Relational | Complex queries, extensions | 300MB |
| **MongoDB** | Document | NoSQL, JSON data | 800MB |
| **Redis** | Key-Value | Caching, real-time data | 50MB |

---

## MariaDB

MySQL-compatible relational database, developed by the original MySQL creators.

### Features
- Drop-in replacement for MySQL
- Better performance on some workloads
- More open development
- Galera clustering

### Install
```bash
hostkit install mariadb
```

### Version
Latest stable version from official MariaDB repository.

### Access
```bash
# Connect as root
sudo mysql -u root

# Or with password
mysql -u root -p
```

### Default Configuration
- Root password: Set during installation
- Port: 3306
- Bind address: 127.0.0.1 (local only)

### Post-Install Security
```bash
sudo mysql_secure_installation
```

---

## MySQL

World's most popular open-source relational database.

### Features
- Enterprise-grade reliability
- Extensive tooling ecosystem
- Oracle compatibility
- Group replication

### Install
```bash
hostkit install mysql
```

### Version
Latest stable version from official MySQL repository.

### Access
```bash
# Connect as root
sudo mysql -u root

# Or with password
mysql -u root -p
```

### Default Configuration
- Root password: Set during installation
- Port: 3306
- Bind address: 127.0.0.1 (local only)

### Post-Install Security
```bash
sudo mysql_secure_installation
```

---

## PostgreSQL

Advanced open-source relational database with extensibility.

### Features
- Complex SQL support
- JSON/JSONB data types
- Custom extensions
- Advanced indexing
- Full ACID compliance

### Install
```bash
hostkit install postgresql
```

### Version
Latest stable version from official PostgreSQL repository.

### Access
```bash
# Switch to postgres user
sudo -u postgres psql

# Or create database
sudo -u postgres createdb mydb
```

### Default Configuration
- Superuser: postgres
- Port: 5432
- Authentication: peer (local), md5 (remote)

### Useful Commands
```bash
# Create user
sudo -u postgres createuser -P myuser

# Create database
sudo -u postgres createdb -O myuser mydb

# List databases
\l

# List tables
\dt
```

---

## MongoDB

Popular document-oriented NoSQL database.

### Features
- Flexible JSON-like documents
- Horizontal scaling
- Rich query language
- Indexing on any field
- Aggregation framework

### Install
```bash
hostkit install mongodb
```

### Version
Latest stable version from official MongoDB repository.

### Access
```bash
# Connect to MongoDB shell
mongosh

# Or specify database
mongosh mydb
```

### Default Configuration
- Port: 27017
- Bind address: 127.0.0.1 (local only)
- Authentication: Disabled by default

### Enable Authentication
```bash
# Connect to MongoDB
mongosh

# Create admin user
use admin
db.createUser({
  user: "admin",
  pwd: "your_password",
  roles: [ { role: "userAdminAnyDatabase", db: "admin" } ]
})
```

---

## Redis

In-memory data structure store, used as database, cache, and message broker.

### Features
- Blazing fast performance
- Rich data types (strings, hashes, lists, sets)
- Pub/Sub messaging
- Persistence options
- Clustering support

### Install
```bash
hostkit install redis
```

### Version
Latest stable version from distribution repositories.

### Access
```bash
# Connect to Redis CLI
redis-cli

# Test connection
ping
# Response: PONG
```

### Default Configuration
- Port: 6379
- Bind address: 127.0.0.1 (local only)
- No authentication by default

### Enable Password
```bash
# Edit configuration
sudo nano /etc/redis/redis.conf

# Add or modify:
requirepass your_strong_password

# Restart Redis
sudo systemctl restart redis
```

### Useful Commands
```bash
# Set key
SET mykey "value"

# Get key
GET mykey

# Check all keys
KEYS *

# Get info
INFO

# Monitor in real-time
MONITOR
```

---

## Choosing the Right Database

### For Web Applications
1. **MariaDB** - General purpose, widely supported
2. **MySQL** - Enterprise needs, Oracle compatibility

### For Complex Applications
1. **PostgreSQL** - Advanced features, extensions
2. **MySQL** - Large scale deployments

### For Modern/NoSQL
1. **MongoDB** - Flexible schema, JSON data
2. **Redis** - Caching, real-time features

### For Caching
1. **Redis** - In-memory speed
2. **MariaDB** - Query cache

### For Small Projects
1. **MariaDB** - Easy setup, lightweight
2. **SQLite** (manual install) - Serverless

---

## Multiple Databases

You can install multiple databases on the same server:

```bash
# Install MariaDB for main application
hostkit install mariadb

# Install Redis for caching
hostkit install redis
```

!!! tip "Resource Consideration"
    Each database consumes RAM and CPU. Monitor your server's resources when running multiple databases.

---

## Database Management Tools

### Command Line
- mysql/mariadb - MySQL/MariaDB CLI
- psql - PostgreSQL CLI
- mongosh - MongoDB CLI
- redis-cli - Redis CLI

### GUI Tools (not included)
- phpMyAdmin - Web-based MySQL/MariaDB
- pgAdmin - PostgreSQL management
- MongoDB Compass - MongoDB GUI
- RedisInsight - Redis GUI

---

## Backup & Restore

### MariaDB / MySQL
```bash
# Backup
mysqldump -u root -p database_name > backup.sql

# Restore
mysql -u root -p database_name < backup.sql
```

### PostgreSQL
```bash
# Backup
pg_dump -U postgres database_name > backup.sql

# Restore
psql -U postgres database_name < backup.sql
```

### MongoDB
```bash
# Backup
mongodump --db database_name --out /backup/

# Restore
mongorestore --db database_name /backup/database_name/
```

### Redis
```bash
# Backup (automatic with RDB)
cp /var/lib/redis/dump.rdb /backup/redis-backup.rdb

# Restore
systemctl stop redis
cp /backup/redis-backup.rdb /var/lib/redis/dump.rdb
systemctl start redis
```
