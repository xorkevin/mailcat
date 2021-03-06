## mailcat completion fish

generate the autocompletion script for fish

### Synopsis


Generate the autocompletion script for the fish shell.

To load completions in your current shell session:
$ mailcat completion fish | source

To load completions for every new session, execute once:
$ mailcat completion fish > ~/.config/fish/completions/mailcat.fish

You will need to start a new shell for this setup to take effect.


```
mailcat completion fish [flags]
```

### Options

```
  -h, --help              help for fish
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
      --config string   config file (default is $XDG_CONFIG_HOME/.mailcat.yaml)
      --debug           turn on debug output
```

### SEE ALSO

* [mailcat completion](mailcat_completion.md)	 - generate the autocompletion script for the specified shell

