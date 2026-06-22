# Hanzo S3 IAM Integration with Docker Compose

This directory contains a complete Docker Compose setup for testing Hanzo S3 IAM integration with Keycloak OIDC authentication.

## 🚀 Quick Start

1. **Build local Hanzo image:**
   ```bash
   make -f Makefile.docker docker-build
   ```

2. **Start the environment:**
   ```bash
   make -f Makefile.docker docker-up
   ```

3. **Run the tests:**
   ```bash
   make -f Makefile.docker docker-test
   ```

4. **Stop the environment:**
   ```bash
   make -f Makefile.docker docker-down
   ```

## 📋 What's Included

The Docker Compose setup includes:

- **🔐 Keycloak** - Identity provider with OIDC support
- **🎯 Hanzo Master** - Metadata management
- **💾 Hanzo Volume** - Data storage
- **📁 Hanzo Filer** - File system interface
- **📊 Hanzo S3** - S3-compatible API with IAM integration
- **🔧 Keycloak Setup** - Automated realm and user configuration

## 🌐 Service URLs

After starting with `docker-up`, services are available at:

| Service | URL | Credentials |
|---------|-----|-------------|
| 🔐 Keycloak Admin | http://localhost:8080 | admin/admin |
| 📊 S3 API | http://localhost:8333 | JWT tokens |
| 📁 Filer | http://localhost:8888 | - |
| 🎯 Master | http://localhost:9333 | - |

## 👥 Test Users

The setup automatically creates test users in Keycloak:

| Username | Password | Role | Permissions |
|----------|----------|------|-------------|
| admin-user | adminuser123 | s3-admin | Full S3 access |
| read-user | readuser123 | s3-read-only | Read-only access |
| write-user | writeuser123 | s3-read-write | Read and write |
| write-only-user | writeonlyuser123 | s3-write-only | Write only |

## 🧪 Running Tests

### All Tests
```bash
make -f Makefile.docker docker-test
```

### Specific Test Categories
```bash
# Authentication tests only
make -f Makefile.docker docker-test-auth

# Role mapping tests only  
make -f Makefile.docker docker-test-roles

# S3 operations tests only
make -f Makefile.docker docker-test-s3ops
```

### Single Test
```bash
make -f Makefile.docker docker-test-single TEST_NAME=TestKeycloakAuthentication
```

## 🔧 Development Workflow

### Complete workflow (recommended)
```bash
# Build, start, test, and clean up
make -f Makefile.docker docker-build
make -f Makefile.docker docker-dev
```
This runs: build → down → up → test

### Using Published Images (Alternative)
If you want to use published Docker Hub images instead of building locally:
```bash
export SEAWEEDFS_IMAGE=chrislusf/hanzo:latest
make -f Makefile.docker docker-up
```

### Manual steps
```bash
# Build image (required first time, or after code changes)
make -f Makefile.docker docker-build

# Start services
make -f Makefile.docker docker-up

# Watch logs
make -f Makefile.docker docker-logs

# Check status
make -f Makefile.docker docker-status

# Run tests
make -f Makefile.docker docker-test

# Stop services  
make -f Makefile.docker docker-down
```

## 🔍 Debugging

### View logs
```bash
# All services
make -f Makefile.docker docker-logs

# S3 service only (includes role mapping debug)
make -f Makefile.docker docker-logs-s3  

# Keycloak only
make -f Makefile.docker docker-logs-keycloak
```

### Get shell access
```bash
# S3 container
make -f Makefile.docker docker-shell-s3

# Keycloak container
make -f Makefile.docker docker-shell-keycloak
```

## 📁 File Structure

```
hanzo/test/s3/iam/
├── docker-compose.yml          # Main Docker Compose configuration
├── Makefile.docker             # Docker-specific Makefile
├── setup_keycloak_docker.sh    # Keycloak setup for containers
├── README-Docker.md            # This file
├── iam_config.json            # IAM configuration (auto-generated)
├── test_config.json           # S3 service configuration
└── *_test.go                  # Go integration tests
```

## 🔄 Configuration

### IAM Configuration
The `setup_keycloak_docker.sh` script automatically generates `iam_config.json` with:

- **OIDC Provider**: Keycloak configuration with proper container networking
- **Role Mapping**: Maps Keycloak roles to Hanzo IAM roles
- **Policies**: Defines S3 permissions for each role
- **Trust Relationships**: Allows Keycloak users to assume Hanzo roles

### Role Mapping Rules
```json
{
  "claim": "roles",
  "value": "s3-admin", 
  "role": "arn:aws:iam::role/KeycloakAdminRole"
}
```

## 🐛 Troubleshooting

### Services not starting
```bash
# Check service status
make -f Makefile.docker docker-status

# View logs for specific service
docker-compose -p hanzo-iam-test logs <service-name>
```

### Keycloak setup issues
```bash  
# Re-run Keycloak setup manually
make -f Makefile.docker docker-keycloak-setup

# Check Keycloak logs
make -f Makefile.docker docker-logs-keycloak
```

### Role mapping not working
```bash
# Check S3 logs for role mapping debug messages
make -f Makefile.docker docker-logs-s3 | grep -i "role\|claim\|mapping"
```

### Port conflicts
If ports are already in use, modify `docker-compose.yml`:
```yaml
ports:
  - "8081:8080"  # Change external port
```

## 🧹 Cleanup

```bash
# Stop containers and remove volumes
make -f Makefile.docker docker-down

# Complete cleanup (containers, volumes, images)
make -f Makefile.docker docker-clean
```

## 🎯 Key Features

- **Local Code Testing**: Uses locally built Hanzo images to test current code
- **Isolated Environment**: No conflicts with local services
- **Consistent Networking**: Services communicate via Docker network
- **Automated Setup**: Keycloak realm and users created automatically
- **Debug Logging**: Verbose logging enabled for troubleshooting  
- **Health Checks**: Proper service dependency management
- **Volume Persistence**: Data persists between restarts (until docker-down)

## 🚦 CI/CD Integration

For automated testing:

```bash
# Build image, run tests with proper cleanup
make -f Makefile.docker docker-build
make -f Makefile.docker docker-up
make -f Makefile.docker docker-wait-healthy  
make -f Makefile.docker docker-test
make -f Makefile.docker docker-down
```
