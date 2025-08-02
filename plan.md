Project Brief: Themed Quote Collections
This document outlines the project goals and implementation plan for creating "Themed Quote Collections," a web application designed to allow users to discover quotes based on various themes. The primary technical objective is to build and learn a powerful full-text search system using pg_search on a serverless PostgreSQL database.

1. Project Goals & Vision
Primary Goal: To build a functional full-stack application that implements efficient, real-world full-text search.

User Experience: Create a clean, fast, and intuitive interface where users can instantly find quotes related to themes like "courage," "innovation," or "perseverance."

Learning Objectives:

Master the fundamentals of PostgreSQL's full-text search capabilities (tsvector, to_tsquery, GIN indexes).

Gain practical experience building a robust backend API in Go.

Develop a modern, reactive frontend using SvelteKit and the Skeleton UI component library.

Understand the workflow of connecting a frontend application to a backend API, including handling CORS.

2. Core Features
A theme-based search bar for user input.

A dynamic results view that displays quotes in a clean, card-based layout.

Visual feedback for loading states and "no results found" scenarios.

A "Copy to Clipboard" button for each quote to allow for easy sharing.

A fully responsive design that works on both desktop and mobile devices.

3. Tech Stack
Backend: Go

Frontend: SvelteKit

UI Library: Skeleton UI (with Tailwind CSS)

Database: Neon (Serverless PostgreSQL with pg_search)

4. Implementation Plan

Phase 1: Backend Foundation & Database Setup ✅
**Status: COMPLETED**
- ✅ Neon PostgreSQL database setup with quotes table
- ✅ Go HTTP server with pgx driver and environment configuration
- ✅ Health check endpoint and database connection validation

Phase 2: Search Logic Implementation ✅  
**Status: COMPLETED**
- ✅ ParadeDB BM25 full-text search indexes created
- ✅ Search API endpoint with relevance scoring and text highlighting
- ✅ CORS middleware for frontend integration
- ✅ PostgreSQL array handling and nullable field support

Phase 3: Frontend Application ✅
**Status: COMPLETED** 
- ✅ SvelteKit project with TypeScript and Tailwind CSS
- ✅ Responsive search interface with modern design
- ✅ Real-time API integration with loading states and error handling
- ✅ Quote cards with copy-to-clipboard and toast notifications

Phase 4: Frontend-Backend Integration ✅
**Status: COMPLETED**
- ✅ Fetch API integration with Go backend
- ✅ State management for search results, loading, and errors  
- ✅ Dynamic result rendering with tags, categories, and relevance scores

Phase 5: Kubernetes Deployment (2-3 hours)
**Goal: Deploy the full-stack application to Civo Kubernetes cluster with public IP access**

## Chosen Architecture: Simplified Single-Pod Deployment

**Best Practices for Single-Pod Architecture:**
- **Nginx Reverse Proxy**: Frontend served as static files, API requests proxied to Go backend
- **Resource Efficiency**: Lower memory/CPU footprint, cost-effective for small-medium traffic
- **Simplified Networking**: Single service endpoint, easier debugging and monitoring
- **Container Optimization**: Multi-stage builds to minimize image size
- **Health Checks**: Proper liveness/readiness probes for both frontend and backend
- **Graceful Shutdown**: Proper signal handling for zero-downtime deployments

## Implementation Steps:

**Container Registry Setup (15 minutes)**:
- **Free Private Options**: GitHub Container Registry (ghcr.io) - 500MB free private storage
- **Fallback**: Docker Hub public repository if private storage exceeded
- Configure registry authentication for Kubernetes

**Containerization (45 minutes)**:
- Multi-stage Dockerfile combining SvelteKit build + Go backend + Nginx
- Nginx configuration for static file serving and API proxying
- Docker Compose for local full-stack testing
- Image optimization and security scanning

**Kubernetes Manifests (1 hour)**:
- Deployment with single pod containing frontend+backend
- LoadBalancer Service for public IP exposure
- ConfigMap for Nginx configuration and environment variables
- Secret for Neon database credentials
- Health check endpoints and probes

**Manual Deployment (30 minutes)**:
- kubectl apply commands and verification
- Service IP and port testing
- Log monitoring and troubleshooting
- Basic smoke testing of search functionality

## Build vs Deploy Strategy:
**Separate build-time from runtime for optimal Docker images:**
- **Local Build Process**: `make build-frontend build-backend build-docker`
- **Runtime Container**: Only pre-built static files + Go binary (no build tools)
- **Benefits**: Faster builds, smaller images, CI/CD friendly

## Container Registry Recommendation:
**GitHub Container Registry (ghcr.io)** - Free 500MB private storage, integrates with your existing GitHub repo

## Load Balancing Options:

### Option A: Civo LoadBalancer (Current Default)
```yaml
type: LoadBalancer  # Cloud provider managed
```
**Pros**: ✅ Fully managed, ✅ Simple setup, ✅ Production ready
**Cons**: ❌ $10-20/month per service
**Access**: Direct public IP

### Option B: NodePort (Free Alternative)
```yaml
type: NodePort
nodePort: 30080  # Port 30000-32767 range
```
**Pros**: ✅ Free, ✅ No extra components, ✅ Simple
**Cons**: ❌ Non-standard ports, ❌ Must manage node IPs
**Access**: `http://NODE_IP:30080`

### Option C: Ingress Controller (Production Best Practice)
```yaml
type: ClusterIP + Ingress
# nginx-ingress, traefik, or kong
```
**Pros**: ✅ One LoadBalancer for all apps, ✅ SSL termination, ✅ Advanced routing
**Cons**: ❌ More complexity, ❌ Requires setup
**Cost**: $10-20/month total (shared across all applications)
**Access**: Standard HTTP/HTTPS with custom domains

---

Phase 6: Production Enhancements (3-4 hours)
**Goal: Add production-ready features for reliability and automation**

**Domain & SSL (1 hour)**:
- Custom domain configuration with DNS
- cert-manager for automated SSL certificate management
- HTTPS redirect and security headers

**Auto-scaling & Performance (1 hour)**:
- Horizontal Pod Autoscaler (HPA) based on CPU/memory metrics
- Resource requests and limits optimization
- Ingress controller with load balancing

**CI/CD Automation (1.5 hours)**:
- GitHub Actions workflow for automated deployment
- Build, test, and deploy pipeline
- Rolling updates and rollback strategies
- Environment-specific deployments (staging/production)

**Database Resilience (30 minutes)**:
- Neon connection pool optimization
- Database health checks and retry logic
- Graceful error handling with user-friendly messages
- Connection timeout and circuit breaker patterns
