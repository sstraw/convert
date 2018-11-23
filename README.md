# Convert
A collection of tools for converting between various formats written in Go

## csv2xlsx
Reads a csv and writes it to an xlsx. Tries to detect utf-16 as well.
```
$ ./csv2xlsx -h
Usage of ./csv2xlsx:
  -f    Force - overwrites existing files
  -i string
        File to read from. If a pipe is detected,
        reads from stdin. Use - to explicitly read
        from stdin.
  -o string
        File to write out to. Use - to write to stdout
```
