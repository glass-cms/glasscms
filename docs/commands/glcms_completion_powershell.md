---
title: Glcms Completion Powershell
create_timestamp: 1715116017
---
## glcms completion powershell

Generate the autocompletion script for powershell

### Synopsis

Generate the autocompletion script for powershell.

To load completions in your current shell session:

	glcms completion powershell | Out-String | Invoke-Expression

To load completions for every new session, add the output of the above command
to your powershell profile.


```
glcms completion powershell [flags]
```

### Options

```
  -h, --help              help for powershell
      --no-descriptions   disable completion descriptions
```

### SEE ALSO

* [glcms completion]()	 - Generate the autocompletion script for the specified shell

###### Auto generated by spf13/cobra on 7-May-2024