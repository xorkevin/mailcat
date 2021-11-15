## mailcat completion powershell

generate the autocompletion script for powershell

### Synopsis


Generate the autocompletion script for powershell.

To load completions in your current shell session:
PS C:\> mailcat completion powershell | Out-String | Invoke-Expression

To load completions for every new session, add the output of the above command
to your powershell profile.


```
mailcat completion powershell [flags]
```

### Options

```
  -h, --help              help for powershell
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
      --config string   config file (default is $XDG_CONFIG_HOME/.mailcat.yaml)
      --debug           turn on debug output
```

### SEE ALSO

* [mailcat completion](mailcat_completion.md)	 - generate the autocompletion script for the specified shell

