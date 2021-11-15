## mailcat fmt

Formats plaintext mail output

### Synopsis

Formats plaintext mail output

```
mailcat fmt [flags]
```

### Options

```
  -a, --add stringArray      specify header values to be added (HEADER:VALUE); may be specified multiple times
  -b, --body                 input is body instead of a full RFC5322 message with headers
  -m, --crlf                 output with CRLF line endings
  -e, --edit                 output in editor convenient format
  -z, --empty                do not read from stdin and instead use empty reader
  -s, --header stringArray   set default header value (HEADER:VALUE); may be specified multiple times
  -h, --help                 help for fmt
  -y, --msgid string         set default generated message id domain (default "mail.example.com")
```

### Options inherited from parent commands

```
      --config string   config file (default is $XDG_CONFIG_HOME/.mailcat.yaml)
      --debug           turn on debug output
```

### SEE ALSO

* [mailcat](mailcat.md)	 - A mail and smtp test tool

