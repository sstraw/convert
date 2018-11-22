package main

import (
    "os"
    "io"
    "log"
    "encoding/csv"
    "flag"

    "golang.org/x/text/encoding/unicode"
    "golang.org/x/text/transform"
    "github.com/tealeg/xlsx"
)

func main() {

    finPtr := flag.String("i", "", "File to read from.")
    foutPtr := flag.String("o", "", "File to write out to.")

    flag.Parse()

    if *finPtr == "" {
        flag.PrintDefaults()
        log.Fatal("-i required")
    } else if *foutPtr == "" {
        flag.PrintDefaults()
        log.Fatal("-o required")
    }

    file, err := os.Open(*finPtr)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    fout, err := os.Create(*foutPtr)
    if err != nil {
        log.Fatal(err)
    }
    defer fout.Close()

    utf16lebom := unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewDecoder()

    csvReader := csv.NewReader(transform.NewReader(file, utf16lebom))

    xlsxFile := xlsx.NewFile()
    xlsxSheet, err := xlsxFile.AddSheet("Sheet")
    if err != nil {
        log.Fatal(err)
    }
    for {
        fields, err := csvReader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            log.Fatal(err)
        }

        row := xlsxSheet.AddRow()
        _ = row.WriteSlice(&fields, -1) //-1 writes all fields
    }
    err = xlsxFile.Write(fout)
    if err != nil {
        log.Println(err)
    }
}
