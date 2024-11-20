# Centralized Configuration Management System

This project focuses on implementing a centralized configuration management system

## Components

### Core Components
1. **Web Service**: Manages user requests for configurations and groups.
2. **Database**: Stores configuration data and metadata.

### Auxiliary Components
1. **Log and Trace Module**: Manages logging and trace analysis.
2. **Metrics Module**: Handles the collection and review of system metrics.

---

## Technologies

- **Programming Language**: Go (Golang)
- **Containerization**: Docker
- **API Documentation**: Swagger
- **Database**: Consul (NoSQL)
- **Tracing**: Jaeger
- **CI/CD**: GitHub Actions

---

## Features

### Web Service Functionalities
- **Create Configuration**: Adds new configurations using JSON input.
- **Create Configuration Group**: Groups multiple configurations under a common identifier.
- **Retrieve Configuration**: Fetches configuration details using a unique identifier.
- **Retrieve Configuration Group**: Gets details of a configuration group by ID.
- **Delete Configuration**: Removes configurations by their unique identifier.
- **Delete Configuration Group**: Deletes a group and its associated configurations.
- **Update Configuration Group**: Adds or modifies configurations within a group.
- **Label-Based Filtering**: Supports advanced operations based on a flexible label system.

### Label System
- Configurations can have labels as `key:value` pairs (e.g., `env:prod;type:api`).
- Labels enable precise filtering and searching.
- Label-based deletion and querying ensure all conditions match.

---

## Additional Capabilities

1. **Versioning**:
   - Configurations and groups support multiple versions.
   - Clients must specify the desired version when querying.

2. **Immutability**:
   - Configurations can only be replaced entirely, ensuring no partial updates.

3. **Idempotent Requests**:
   - Repeated identical requests yield the same result.
   - Metadata for idempotence is stored in the database.

4. **UUIDs**:
   - Unique identifiers for configurations and groups use UUIDs.

---

## Database

- All configurations and groups are stored in a **Consul** NoSQL database for scalability and reliability.

---

## Deployment

- **Containerized Setup**:
  - The application and database are fully containerized using Docker.
  - Multi-stage builds reduce image size and improve performance.
  - Docker Compose is used for seamless orchestration of components.

- **Tracing and Metrics**:
  - Integrated with Jaeger for distributed tracing.
  - Metrics are exposed for system monitoring.

- **CI/CD**:
  - GitHub Actions workflows automatically run for changes in the `main` branch.
  - Follows GitFlow branching strategy for development.

