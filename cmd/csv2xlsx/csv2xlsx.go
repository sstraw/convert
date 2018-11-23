package main

import (
    "os"
    "io"
    "log"
    "encoding/csv"
    "flag"
    "bufio"

    "golang.org/x/text/encoding/unicode"
    "golang.org/x/text/transform"
    "github.com/tealeg/xlsx"
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
    var fin *(os.File)
    var err error
    if *finPtr == "-" {
        fin = os.Stdin
        log.Println("Reading from stdin")
    } else if *finPtr != ""{
        fin, err = os.Open(*finPtr)
        if err != nil {
            log.Fatal(err)
        }
        defer fin.Close()
    } else {
        stat, _ := os.Stdin.Stat()
        if (stat.Mode() & os.ModeCharDevice) == 0 {
            //Checks stickychar to see if data is being piped
            fin = os.Stdin
        } else {
            log.Println("No pipe detected and no input file specified.")
            log.Fatal("Use -i - to force reading from stdin")
        }
    }

    //File output opening
    var fout *(os.File)
    if *foutPtr == "-" {
        fout, err = os.Stdout, nil
    } else if *force {
        fout, err = os.OpenFile(*foutPtr, os.O_RDWR | os.O_CREATE, 0666)
    } else {
        fout, err = os.OpenFile(*foutPtr, os.O_RDWR | os.O_CREATE | os.O_EXCL, 0666)
    }

    if err != nil {
        log.Fatal(err)
    }
    defer fout.Close()


    //Input reading
    enc, err:= EncHandle(fin)
    if err == io.EOF {
        log.Fatal("Reached EOF - Not enough data")
    } else if err != nil {
        log.Fatal(err)
    }

    csvReader := csv.NewReader(enc)

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

func EncHandle (fhandle io.Reader) (io.Reader, error) {
    /*
    Attempts to detect encoding (specifically, decide if
    utf-16le, utf-16be or utf-8 with or without bom)
    and returns a reader object for it, translating 
    if necessary
    */
    buf := bufio.NewReader(fhandle)

    head, err := buf.Peek(4)

    if err!= nil {
        return nil, err
    }

    switch {
        //UTF-16-le-like with no bom
        case head[0] != 0x00 && head[1] == 0x00:
            log.Println("Assuming utf-16-LE")
            fallthrough
        //UTF-16-le with bom
        case head[0] == 0xFF && head[1] == 0xFE:
            dec := unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewDecoder()
            return transform.NewReader(buf, dec), nil
        //UTF-16-be-like with no bom
        case head[0] != 0x00 && head[1] == 0x00:
            log.Println("Assuming utf-16-LE")
            fallthrough
        //UTF-16-be with bom
        case head[0] == 0xFE && head[1] == 0xFF:
            dec := unicode.UTF16(unicode.BigEndian,    unicode.UseBOM).NewDecoder()
            return transform.NewReader(buf, dec), nil
   }
    return buf, nil
}
