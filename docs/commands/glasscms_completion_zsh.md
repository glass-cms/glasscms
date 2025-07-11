---
title: Glasscms Completion Zsh
create_time: 1751798081
---
## glasscms completion zsh

Generate the autocompletion script for zsh

### Synopsis

Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need
to enable it.  You can execute the following once:

	echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions in your current shell session:

	source <(glasscms completion zsh)

To load completions for every new session, execute once:

#### Linux:

	glasscms completion zsh > "${fpath[1]}/_glasscms"

#### macOS:

	glasscms completion zsh > $(brew --prefix)/share/zsh/site-functions/_glasscms

You will need to start a new shell for this setup to take effect.


```
glasscms completion zsh [flags]
```

### Options

```
  -h, --help              help for zsh
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

