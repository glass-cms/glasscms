---
title: Glasscms Completion Powershell
create_time: 1751798081
---
## glasscms completion powershell

Generate the autocompletion script for powershell

### Synopsis

Generate the autocompletion script for powershell.

To load completions in your current shell session:

	glasscms completion powershell | Out-String | Invoke-Expression

To load completions for every new session, add the output of the above command
to your powershell profile.


```
glasscms completion powershell [flags]
```

### Options

```
  -h, --help              help for powershell
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

