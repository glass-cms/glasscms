---
title: Glasscms Auth Init
create_time: 1751798081
---
## glasscms auth init

Initialize a new authentication token

### Synopsis

Initialize a new authentication token for API access.

This command creates a new authentication token that can be used to authenticate 
requests to the GlassCMS API. By default, the token is valid for 24 hours.

The token is displayed only once upon creation and should be stored securely.
It cannot be retrieved later, so make sure to save it in a secure location.


```
glasscms auth init [flags]
```

### Examples

```
# Create a new token with default settings
glasscms auth init

# Create a new token with a specific database driver and DSN
glasscms auth init --driver postgres --dsn "postgres://user:password@localhost:5432/glasscms"

```

### Options

```
      --database.driver string              The database driver to use (e.g., postgres, mysql, sqlite)
      --database.dsn string                 The data source name (DSN) for the database connection
      --database.max_connections int        The maximum number of open connections to the database (default 5)
      --database.max_idle_connections int   The maximum number of idle connections in the connection pool (default 1)
  -h, --help                                help for init
```

### Options inherited from parent commands

```
      --logger.format string   Log format (default "TEXT")
      --logger.level string    Log level (default "INFO")
  -v, --verbose                Enable verbose output
      --version                Show version information
```

### SEE ALSO

* [glasscms auth](glasscms_auth.md)	 - 

