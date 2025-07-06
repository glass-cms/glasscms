# GlassCMS

A headless CMS powered by markdown that seamlessly turns your markdown files into a REST API.

## Features

- **Markdown-based Content**: Transform your existing markdown files into a headless CMS
- **REST API**: Access your content via a clean REST API
- **Authentication**: Built-in token-based authentication system
- **File System Integration**: Sync content directly from your file system
- **Database Support**: PostgreSQL and SQLite support
- **OpenAPI Specification**: Well-documented API with OpenAPI 3.0
- **CLI Interface**: Comprehensive command-line interface for management

## Installation

```bash
go install github.com/glass-cms/glasscms@latest
```

## Quick Start

1. **Initialize authentication**:
   ```bash
   glasscms auth init
   ```

2. **Start the server**:
   ```bash
   glasscms server start
   ```

3. **Sync your markdown files**:
   ```bash
   glasscms sync --path /path/to/your/markdown/files
   ```

## Commands

- `glasscms auth init` - Initialize authentication system
- `glasscms server start` - Start the API server
- `glasscms sync` - Sync markdown files to the database
- `glasscms convert` - Convert between different formats
- `glasscms migrate` - Run database migrations
- `glasscms docs` - Generate documentation

## Configuration

GlassCMS can be configured via:
- Configuration file (`config.yaml`)
- Environment variables (prefixed with `GLASS_`)
- Command-line flags

### Environment Variables

- `GLASS_LOG_LEVEL` - Log level (default: INFO)
- `GLASS_LOG_FORMAT` - Log format (default: TEXT)
- `GLASS_VERBOSE` - Enable verbose output

## API

The API follows REST conventions and provides endpoints for:
- **Items**: Manage content items (`/items`)
- **Authentication**: Token-based authentication

See the OpenAPI specification in `openapi.yaml` for complete API documentation.

## Development

### Prerequisites

- Go 1.23+
- Task runner (optional)

### Building

```bash
# Using Go
go build -o glasscms

# Using Task
task build
```

### Testing

```bash
# Run tests
task test

# Run tests with coverage
task coverage
```

### Linting

```bash
# Run linter
task lint

# Fix linting issues
task lint-fix
```

## Database

GlassCMS supports:
- PostgreSQL
- SQLite

Migrations are managed automatically via the `migrate` command.

## License

See [LICENSE.md](LICENSE.md) for details.