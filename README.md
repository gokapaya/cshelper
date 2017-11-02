# cshelper

### WARNING! Under construction.

`cshelper` helps us organize the secret santa over on [/r/ClosetSanta](https://reddit.com/r/closetsanta).

## Features

```
Usage:
  cshelper [flags]
  cshelper [command]
Available Commands:
  help        Help about any command
  match       Generate a list of pairings
  pm          Send PMs to user(s)
Flags:
      --debug   print debug logs
  -h, --help    help for cshelper
```

## Using

To use `cshelper` you need to take a few steps of preparation:

Create a file called `cshelper.toml` in a directory `.cshelper`:

```
debug = false

[bot]
useragent = ""
username = ""
password = ""
clientId = ""
clientSecret = ""
```

The bot configuration contains the same elements as described [here](https://github.com/turnage/graw/wiki/agent-files).

## Building

`cshelper` relies on a modified version of [clphub/munkres](https://github.com/clyphub/munkres). The change is can be
found in the [patch here](./munkres.patch).

To patch the dependency in the vendor/ directory run the following commands:

```
$ dep ensure
$ go generate && go build
```
