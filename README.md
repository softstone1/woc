# Task1

# Weather Forecasts for Major Global Cities

This application provides weather forecasts for major global cities.

## Features

- **Clean Architecture**: Organized in a way that separates concerns (domain, infrastructure, interface, and application layers), enhancing the maintainability and scalability of the application.
- **Environment Configuration**: Uses Viper to manage and load environment variables, making the application configurable and easy to adapt to different environments.
- **Mock Generation**: Utilizes GoMock for generating mocks in unit tests.
- **In-Memory Storage**: Employs in-memory data storage to manage a list of major global cities.
- **Graceful Shutdown**: Implements graceful shutdown processes to handle server terminations smoothly, preserving data integrity and ensuring that all processes are completed before shutdown.
- **Structured Logging**: Uses the `slog` package from the Go standard library for structured logging in JSON format, providing better traceability and readability of logs.
- **Routing with MuxServe**: Uses the `muxserve` library from the Go standard library to manage routing, enhancing the routing capabilities with minimal overhead.
- **Logging and Recovery Middleware**: Integrates Gorilla Handlers to provide logging and recovery middleware.
- **Automated Unit Tests**: Leverages GitHub Copilot to generate unit tests.
- **Go Profiling**: Includes profiling capabilities to optimize performance and troubleshoot bottlenecks in the application.
- **Go Template and HTMX**: Implements a simple web UI using Go's native template system and HTMX for dynamic content without writing JavaScript, enhancing user interaction and page responsiveness.
- **Dependency Injection**: Utilizes dependency injection in the main function to orchestrate application setup. For more complex dependency management, the application integrates Google's Wire package, automating dependency injection and reducing boilerplate.

## Getting Started

### Prerequisites

- Go 1.18 or higher
- Docker (if running via Docker container)
- Properly configured Go environment (if running locally without Docker)

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/softstone1/woc.git
   ```
2. Navigate to the project directory:
   ```bash
   cd woc
   ```

### Configuration

Set the necessary environment variables through a `.env` file or export them directly into your environment. 

### Running the Application

#### Locally

To start the server, run:

```bash
go run cmd/main.go
```

The server will start on the port specified in your environment variables, defaulting to `8080`.

#### With Docker

1. Build the Docker image:

   ```bash
   docker build -t woc .
   ```

2. Run the Docker container:

   ```bash
   docker run -p 8080:8080 woc
   ```

This will start the application inside a Docker container and map the container's port 8080 to the local port 8080.

### Profiling

***Access the profiling data***

Once woc application is running, You can use curl to download the profiles directly from the command line.

```bash
# CPU profile
curl -o cpu.pprof "http://localhost:8080/debug/pprof/profile?seconds=30"
# Heap Profile
curl -o heap.pprof "http://localhost:8080/debug/pprof/heap"

# Block Profile
curl -o block.pprof "http://localhost:8080/debug/pprof/block"

# Thread Creation Profile
curl -o threadcreate.pprof "http://localhost:8080/debug/pprof/threadcreate"

# Goroutine Profile
curl -o goroutine.pprof "http://localhost:8080/debug/pprof/goroutine"
```

***Analyzing the Profiling Data***

To analyze the data, you can use the go tool pprof command. If you downloaded a profile (like the CPU profile), you can analyze it by running:
```bash
go tool pprof path-to-your-profile-file
```

### Using the Application

Access the application through your web browser or API client at `http://localhost:8080`. The homepage will allow you to select a city from a dropdown menu and view the current weather forecast.

## Development

### Testing

To run the unit tests:

```bash
go test ./...
```

### Checking Logs in Docker

To check logs from the running Docker container, first identify the container ID using:

```bash
docker ps
```

Then use the following command to view the logs:

```bash
docker logs <container_id>
```