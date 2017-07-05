# detailer
Detailer is a command-line utility that strips (de-tails) garbage past the end of media files (currently, JPEG and PNG are supported).

# Building
`go build`

# Usage
```
detailer [--json] [--truncate] file
```

* `--json` results in machine-readable run results being outputted to stdout, for example
```
{
        "format": "PNG",
        "size": 52333,
        "data_size": 52187,
        "truncated": true
}
```

* `--truncate` tells Detailer to actually truncate the file, otherwise a dry run will be performed.
