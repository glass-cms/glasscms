---
title: Glasscms Completion Zsh
createTime: 1715500010
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
  -v, --verbose   Enable verbose output
```

### SEE ALSO

* [glasscms completion]()	 - Generate the autocompletion script for the specified shell

###### Auto generated by spf13/cobra on 12-May-2024