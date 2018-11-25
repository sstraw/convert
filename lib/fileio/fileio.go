package fileio

import (
    "os"
    "io"
    "bufio"
    "errors"
    "golang.org/x/text/encoding/unicode"
    "golang.org/x/text/transform"
)

var NoPipe = errors.New("No pipe detected")

func OpenInput (fname *string) (*os.File, error) {
    if *fname  == "-" {
        return os.Stdin, nil
    }

    if *fname == "" {
        stat, _ := os.Stdin.Stat()
        if (stat.Mode() & os.ModeCharDevice) == 0 {
            return os.Stdin, nil
        }

        return nil, NoPipe
    }

    return os.Open(*fname)
}

func OpenOutput (fname *string, force *bool) (*os.File, error) {
    if *fname == "-" {
        return os.Stdout, nil
    }

    if *force {
        return os.OpenFile(*fname, os.O_RDWR | os.O_CREATE, 0666)
    }

    return os.OpenFile(*fname, os.O_RDWR | os.O_CREATE | os.O_EXCL, 0666)
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
            fallthrough
        //UTF-16-le with bom
        case head[0] == 0xFF && head[1] == 0xFE:
            dec := unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewDecoder()
            return transform.NewReader(buf, dec), nil
        //UTF-16-be-like with no bom
        case head[0] == 0x00 && head[1] != 0x00:
            fallthrough
        //UTF-16-be with bom
        case head[0] == 0xFE && head[1] == 0xFF:
            dec := unicode.UTF16(unicode.BigEndian,    unicode.UseBOM).NewDecoder()
            return transform.NewReader(buf, dec), nil
   }
    return buf, nil
}
