---
title: Glasscms Sync
create_time: 1751798081
---
## glasscms sync

Synchronize content items from a source to the GlassCMS server

### Synopsis

Synchronize content items from a source to the GlassCMS API server.

The sync command allows you to import and update content items from external 
sources into your GlassCMS instance. It compares the items in the source with 
those on the server and performs the necessary create, update, or delete operations 
to keep them in sync.

Sources are external content repositories that contain structured content items.
Each source has a specific format and organization, which GlassCMS can interpret
and import into its content management system.

Supported source types:
- filesystem: Read items from a directory on the local filesystem. Items should be
  organized in a directory structure with JSON or YAML files representing content items.
  Each file should contain metadata and content according to the GlassCMS schema.

When run in preview mode (default), the command will show what changes would be made
without actually applying them. Use the --live flag to apply the changes.


```
glasscms sync [source-type] [source-path] [flags]
```

### Examples

```
# Preview synchronization from a filesystem source
glasscms sync filesystem /path/to/items

# Perform live synchronization with server authentication
glasscms sync filesystem /path/to/items --live --token "your-auth-token"

# Synchronize to a specific server
glasscms sync filesystem /path/to/items --server "https://cms.example.com" --token "your-auth-token"

# Specify a front matter property to determine if an item is hidden
glasscms sync filesystem /path/to/items --hidden-property "draft" --hidden-value true

```

### Options

```
  -h, --help                     help for sync
      --hidden-property string   Front matter property name to determine if an item is hidden (e.g., 'draft', 'hidden', 'private')
      --hidden-value             Value of the hidden property that indicates an item is hidden 
                                 		(true = truthy values are hidden, false = falsy values are hidden) (default true)
      --live                     When live mode is enabled, items are synchronized to the server, otherwise changes are only previewed
      --parse-wikilinks          Parse wikilinks in the content (default true)
      --server string            The URL of the server to synchronize items to (default "http://localhost:8080")
      --token string             Bearer token for server authentication
```

### Options inherited from parent commands

```
      --logger.format string   Log format (default "TEXT")
      --logger.level string    Log level (default "INFO")
  -v, --verbose                Enable verbose output
      --version                Show version information
```

### SEE ALSO

* [glasscms](glasscms.md)	 - glasscms is a headless CMS powered by markdown

