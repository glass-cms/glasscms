---
title: Glasscms Completion Bash
create_time: 1751798081
---
## glasscms completion bash

Generate the autocompletion script for bash

### Synopsis

Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

	source <(glasscms completion bash)

To load completions for every new session, execute once:

#### Linux:

	glasscms completion bash > /etc/bash_completion.d/glasscms

#### macOS:

	glasscms completion bash > $(brew --prefix)/etc/bash_completion.d/glasscms

You will need to start a new shell for this setup to take effect.


```
glasscms completion bash
```

### Options

```
  -h, --help              help for bash
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
      --logger.format string   Log format (default "TEXT")
      --logger.level string    Log level (default "INFO")
  -v, --verbose                Enable verbose output
      --version                Show version information
```

### SEE ALSO

* [glasscms completion](glasscms_completion.md)	 - Generate the autocompletion script for the specified shell

