---
title: Glasscms Server Start
create_time: 1751798081
---
## glasscms server start

Start the GlassCMS API server

### Synopsis

Start the GlassCMS API server with the specified configuration.

This command initializes and starts the CMS server with database connectivity
and all required services. It sets up the HTTP server with appropriate middleware
for authentication, content negotiation, and request tracking.

The server will continue running until it receives a termination signal.


```
glasscms server start [flags]
```

### Options

```
      --database.driver string              The name of the database driver
      --database.dsn string                 The data source name (DSN) for the database
      --database.max_connections int        The maximum number of connections that can be opened to the database (default 5)
      --database.max_idle_connections int   The maximum number of idle connections that can be maintained (default 1)
  -h, --help                                help for start
```

### Options inherited from parent commands

```
      --logger.format string   Log format (default "TEXT")
      --logger.level string    Log level (default "INFO")
  -v, --verbose                Enable verbose output
      --version                Show version information
```

### SEE ALSO

* [glasscms server](glasscms_server.md)	 - Server management commands

