package main

import (
    "os"
    "io"
    "log"
    "encoding/csv"
    "flag"

//    "github.com/tealeg/xlsx"
    "github.com/sstraw/convert/lib/fileio"
)

func main() {

    finPtr := flag.String("i", "", "File to read from. If a pipe is detected,\n" +
                          "reads from stdin. Use - to explicitly read\n" +
                          "from stdin.")

    foutPtr := flag.String("o", "", "File to write out to. Use - to write to stdout")
    force := flag.Bool("f", false, "Force - overwrites existing files")

    flag.Parse()

    if *foutPtr == "" {
        log.Fatal("-o is required")
    }

    //File input opening
    fin, err := fileio.OpenInput(finPtr)
    if err == fileio.NoPipe {
        log.Fatal("No pipe detected on stdin. Force stdin with \"-i -\"")
    } else if err != nil {
        log.Fatal(err)
    }
    defer fin.Close()

    //File output opening
    var fout *(os.File)
    fout, err = fileio.OpenOutput(foutPtr, force)
    if os.IsExist(err) {
        log.Fatal("Output file already exists. Force overwrite with -f")
    } else if err != nil {
        log.Fatal(err)
    }
    defer fout.Close()

    //Input reading
    enc, err := fileio.EncHandle(fin)
    if err == io.EOF {
        log.Fatal("Reached EOF - Not enough data")
    } else if err != nil {
        log.Fatal(err)
    }

    csvReader := csv.NewReader(enc)

    //Allows variable number of fields
    csvReader.FieldsPerRecord = -1

    xlsxStream, err := fileio.NewStream(fout, "sheet")
    if  err != nil {
        log.Fatal(err)
    }
    defer xlsxStream.Close()

    for {
        fields, err := csvReader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            log.Fatal(err)
        }

        if err := xlsxStream.WriteRow(fields); err != nil {
            log.Fatal(err)
        }
    }
    if err != nil {
        log.Println(err)
    }
}
