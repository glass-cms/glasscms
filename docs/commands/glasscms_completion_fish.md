---
title: Glasscms Completion Fish
create_time: 1740900092
---
## glasscms completion fish

Generate the autocompletion script for fish

### Synopsis

Generate the autocompletion script for the fish shell.

To load completions in your current shell session:

	glasscms completion fish | source

To load completions for every new session, execute once:

	glasscms completion fish > ~/.config/fish/completions/glasscms.fish

You will need to start a new shell for this setup to take effect.


```
glasscms completion fish [flags]
```

### Options

```
  -h, --help              help for fish
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
      --logger.format string   Log format (default "TEXT")
      --logger.level string    Log level (default "INFO")
  -v, --verbose                Enable verbose output
```

### SEE ALSO

* [glasscms completion](glasscms_completion.md)	 - Generate the autocompletion script for the specified shell

