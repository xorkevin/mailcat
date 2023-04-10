## mailcat completion zsh

Generate the autocompletion script for zsh

### Synopsis

Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need
to enable it.  You can execute the following once:

	echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions in your current shell session:

	source <(mailcat completion zsh)

To load completions for every new session, execute once:

#### Linux:

	mailcat completion zsh > "${fpath[1]}/_mailcat"

#### macOS:

	mailcat completion zsh > $(brew --prefix)/share/zsh/site-functions/_mailcat

You will need to start a new shell for this setup to take effect.


```
mailcat completion zsh [flags]
```

### Options

```
  -h, --help              help for zsh
      --no-descriptions   disable completion descriptions
```

### SEE ALSO

* [mailcat completion](mailcat_completion.md)	 - Generate the autocompletion script for the specified shell

