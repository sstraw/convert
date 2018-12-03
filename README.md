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

I used the following as reference when developing the xlsx creator.

* https://docs.microsoft.com/en-us/office/open-xml/structure-of-a-spreadsheetml-document
* https://www.loc.gov/preservation/digital/formats/fdd/fdd000398.shtml
* http://officeopenxml.com/anatomyofOOXML-xlsx.php
* http://officeopenxml.com/SScontentOverview.php
