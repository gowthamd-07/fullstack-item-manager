# Item Manager - Full-Stack Kubernetes Application

A modern full-stack CRUD application for managing inventory items, deployed on Kubernetes using Kind (Kubernetes in Docker). The application consists of a Go backend API, React frontend, and PostgreSQL database.

## 🏗️ Architecture

```
┌─────────────┐
│   Frontend  │  React + Vite (Port 8080)
│  (Nginx)    │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   Backend   │  Go + Chi Router (Port 8000)
│     API     │
└──────┬──────┘
       │
       ▼
┌─────────────┐
│  PostgreSQL │  Database (Port 5432)
└─────────────┘
```

## 📦 Tech Stack

### Backend

- **Language**: Go 1.23
- **Framework**: Chi Router
- **Database**: PostgreSQL 15
- **ORM**: pgx/v5 (PostgreSQL driver)
- **Features**: RESTful API, CORS support, graceful shutdown

### Frontend

- **Framework**: React 19
- **Build Tool**: Vite 7
- **UI Library**: Lucide React (icons)
- **Web Server**: Nginx (Alpine)

### Infrastructure

- **Container Runtime**: Docker
- **Orchestration**: Kubernetes (Kind)
- **Database**: PostgreSQL (Helm Chart)
- **Namespace**: `application`

## 📁 Project Structure

```
project #1/
├── backend/                 # Go backend API
│   ├── cmd/
│   │   └── api/
│   │       └── main.go     # Application entry point
│   ├── internal/
│   │   ├── config/         # Configuration management
│   │   ├── db/             # Database connection
│   │   ├── handler/        # HTTP handlers
│   │   ├── models/         # Data models
│   │   └── repository/     # Data access layer
│   ├── migrations/         # Database migrations
│   ├── kubernetes/         # K8s manifests
│   │   ├── deployment.yml
│   │   ├── service.yml
│   │   └── secrets.yml
│   ├── test/              # API tests
│   ├── Dockerfile
│   └── go.mod
├── frontend/               # React frontend
│   ├── src/
│   │   ├── App.jsx         # Main component
│   │   └── components/
│   ├── kubernetes/         # K8s manifests
│   │   ├── deployment.yml
│   │   └── service.yml
│   ├── Dockerfile
│   ├── nginx.conf
│   └── package.json
├── docker/                 # Docker Compose for local dev
│   └── docker-compose.yml
└── postgres-values.yml     # Helm values for PostgreSQL
```

## 🚀 Quick Start

### Prerequisites

- Docker & Docker Compose
- Kubernetes (Kind) cluster
- kubectl configured
- Helm 3.x
- Go 1.23+ (for local backend development)
- Node.js 20+ (for local frontend development)

### 1. Create Kind Cluster

```bash
# Navigate to kind directory
cd ../../kind

# Create cluster with local storage
kind create cluster --config kind.yml --name kind-node-local-storage

# Verify cluster
kubectl cluster-info --context kind-kind-node-local-storage
```

### 2. Setup PostgreSQL

```bash
# Add Bitnami Helm repository
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update

# Create namespace
kubectl create namespace postgres

# Install PostgreSQL
helm install go-postgresql bitnami/postgresql \
  --namespace postgres \
  --values postgres-values.yml

# Wait for PostgreSQL to be ready
kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=postgresql \
  -n postgres --timeout=300s

# Get database credentials (if needed)
kubectl get secret go-postgresql -n postgres -o jsonpath="{.data.postgres-password}" | base64 -d
```

### 3. Create Application Namespace

```bash
kubectl create namespace application
```

### 4. Deploy Backend

```bash
# Install backend chart
helm install backend backend/backend -n application

# Or upgrade if already installed
helm upgrade backend backend/backend -n application

# Check deployment status
kubectl get pods -n application -l app=backend
kubectl get svc -n application
```

### 5. Deploy Frontend

```bash
# Install frontend chart
helm install frontend frontend/frontend -n application

# Or upgrade if already installed
helm upgrade frontend frontend/frontend -n application

# Check deployment status
kubectl get pods -n application -l app=frontend
```

### 6. Access the Application

```bash
# Port forward to access services
kubectl port-forward -n application svc/frontend 8080:8080
kubectl port-forward -n application svc/backend 8000:8000

# Access frontend
open http://localhost:8080

# Test backend API
curl http://localhost:8000/health
curl http://localhost:8000/api/items
```

## 🛠️ Local Development

### Using Docker Compose

```bash
cd docker
docker-compose up -d

# Access services
# Frontend: http://localhost:8080
# Backend: http://localhost:8000
# PostgreSQL: localhost:5432
```

### Backend Development

```bash
cd backend

# Install dependencies
go mod download

# Run migrations (requires PostgreSQL running)
# Set DB_URL environment variable
export DB_URL="postgres://user:password@localhost:5432/cruddb?sslmode=disable"

# Run application
go run cmd/api/main.go
```

### Frontend Development

```bash
cd frontend

# Install dependencies
npm install

# Start development server
npm run dev

# Build for production
npm run build
```

## 📡 API Endpoints

### Health Checks

- `GET /health` - Liveness probe
- `GET /health/ready` - Readiness probe

### Items API

- `GET /api/items` - List all items
- `GET /api/items/{id}` - Get item by ID
- `POST /api/items` - Create new item
  ```json
  {
    "name": "Item Name",
    "price": 99.99
  }
  ```
- `PUT /api/items/{id}` - Update item
  ```json
  {
    "name": "Updated Name",
    "price": 149.99
  }
  ```
- `DELETE /api/items/{id}` - Delete item

## 🗄️ Database

### Schema

```sql
CREATE TABLE items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

### Migrations

Migrations are located in `backend/migrations/`:

- `000001_create_items_table.up.sql` - Create items table
- `000001_create_items_table.down.sql` - Drop items table

## 🐳 Docker Images

The application uses pre-built Docker images:

- Backend: `gowthamd07101999/go-crud-app-app:docker-backend`
- Frontend: `gowthamd07101999/go-crud-app-app:docker-frontend`

### Building Images

```bash
# Build backend
cd backend
docker build -t gowthamd07101999/go-crud-app-app:docker-backend .

# Build frontend
cd frontend
docker build -t gowthamd07101999/go-crud-app-app:docker-frontend .
```

## ☸️ Kubernetes Configuration

### Backend Deployment

- **Replicas**: 2
- **Resources**:
  - Requests: 100m CPU, 250Mi memory
  - Limits: 250m CPU, 500Mi memory
- **Health Probes**: Liveness, Readiness, Startup
- **Port**: 8000

### Frontend Deployment

- **Replicas**: 2
- **Resources**:
  - Requests: 100m CPU, 250Mi memory
  - Limits: 250m CPU, 500Mi memory
- **Health Probes**: Liveness, Readiness, Startup
- **Port**: 8080

### Services

- Backend service exposes port 8000
- Frontend service exposes port 8080

## 🔧 Configuration

### Environment Variables

**Backend:**

- `DB_URL` - PostgreSQL connection string
- `PORT` - Server port (default: 8000)
- `LOG_LEVEL` - Logging level (default: info)

**Frontend:**

- `VITE_API_URL` - Backend API URL (default: `/api`)

### Database Connection

In Kubernetes, the backend connects to PostgreSQL using:

```
postgres://admin:admin123@go-postgresql-hl.postgres.svc.cluster.local:5432/crud?sslmode=disable
```

## 🧪 Testing

### Backend API Tests

```bash
cd backend/test

# Run test script
bash test_api.sh

# Or use Postman collection
# Import postman_collection.json
```

### Manual Testing

```bash
# Create an item
curl -X POST http://localhost:8000/api/items \
  -H "Content-Type: application/json" \
  -d '{"name": "Test Item", "price": 29.99}'

# List items
curl http://localhost:8000/api/items

# Get item by ID
curl http://localhost:8000/api/items/{id}

# Update item
curl -X PUT http://localhost:8000/api/items/{id} \
  -H "Content-Type: application/json" \
  -d '{"name": "Updated Item", "price": 39.99}'

# Delete item
curl -X DELETE http://localhost:8000/api/items/{id}
```

## 📊 Monitoring & Debugging

### Check Pod Status

```bash
# All pods in application namespace
kubectl get pods -n application

# Pod logs
kubectl logs -n application -l app=backend --tail=100
kubectl logs -n application -l app=frontend --tail=100

# Describe pod
kubectl describe pod <pod-name> -n application
```

### Check Services

```bash
kubectl get svc -n application
kubectl describe svc backend -n application
kubectl describe svc frontend -n application
```

### Database Access

```bash
# Port forward to PostgreSQL
kubectl port-forward -n postgres svc/go-postgresql-hl 5432:5432

# Connect using psql
psql -h localhost -U admin -d crud
```

## 🧹 Cleanup

### Remove Application

```bash
# Uninstall Helm releases
helm uninstall frontend -n application
helm uninstall backend -n application
kubectl delete namespace application
```

### Remove PostgreSQL

```bash
# Uninstall Helm release
helm uninstall go-postgresql -n postgres
kubectl delete namespace postgres
```

### Delete Kind Cluster

```bash
kind delete cluster --name kind-node-local-storage
```

## 🔐 Security Notes

⚠️ **Important**: This is a development/demo setup. For production:

1. Use Kubernetes Secrets for sensitive data (database passwords, API keys)
2. Enable TLS/HTTPS for all services
3. Implement proper authentication and authorization
4. Use Network Policies to restrict pod-to-pod communication
5. Scan Docker images for vulnerabilities
6. Use non-root users in containers
7. Implement resource quotas and limits
8. Enable audit logging

## 📝 Notes

- The application uses rolling updates for zero-downtime deployments
- Health probes ensure only healthy pods receive traffic
- Database migrations should be run before deploying the backend
- Frontend is served via Nginx with optimized configuration
- All services run in the `application` namespace

## 🐛 Troubleshooting

### Pods not starting

```bash
# Check pod events
kubectl describe pod <pod-name> -n application

# Check logs
kubectl logs <pod-name> -n application
```

### Database connection issues

```bash
# Verify PostgreSQL is running
kubectl get pods -n postgres

# Check backend environment variables
kubectl exec -n application <backend-pod> -- env | grep DB_URL

# Test database connectivity
kubectl exec -n application <backend-pod> -- wget -O- http://go-postgresql-hl.postgres.svc.cluster.local:5432
```

### Frontend can't reach backend

```bash
# Verify services
kubectl get svc -n application

# Check service endpoints
kubectl get endpoints -n application

# Test from frontend pod
kubectl exec -n application <frontend-pod> -- wget -O- http://backend:8000/health
```

## 📚 Additional Resources

- [Kubernetes Documentation](https://kubernetes.io/docs/)
- [Kind Documentation](https://kind.sigs.k8s.io/)
- [Go Chi Router](https://github.com/go-chi/chi)
- [React Documentation](https://react.dev/)
- [PostgreSQL Helm Chart](https://github.com/bitnami/charts/tree/main/bitnami/postgresql)

---

**Project**: Item Manager  
**Stack**: Go + React + PostgreSQL + Kubernetes  
**Deployment**: Kind (Kubernetes in Docker)
