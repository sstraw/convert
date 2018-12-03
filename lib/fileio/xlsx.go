package fileio

import (
    "archive/zip"
    "encoding/xml"
    "io"
    "fmt"
    "bytes"
)

type XLSXStream struct {
    zip   *zip.Writer
    sheet io.Writer
    name string
}

type genFunc func() *bytes.Buffer

func NewStream (buf io.Writer, name string) (*XLSXStream, error) {
    var stream XLSXStream
    stream.zip = zip.NewWriter(buf)
    stream.name = name

    // Directory structure. Files needed:
    // out.xlsx
    // |-> [Content_Types].xml
    // |-> _rels/
    // |  |->.rels
    // |-> xl
    // |  |-> workbook.xml
    // |  |-> _rels/
    // |  |  |-> workbook.xml.rels
    // |  |-> worksheets/
    // |  |  |-> %SheetName.xml

    var files = []struct {
        Name string
        Generator genFunc
    }{
        {"[Content_Types].xml",        stream.genContentTypes},
        {"_rels/.rels",                stream.genRels},
        {"xl/_rels/workbook.xml.rels", stream.genWorkbookRels},
        {"xl/workbook.xml",            stream.genWorkbook},
    }

    for _, file := range files {
        f, err := stream.zip.Create(file.Name)
        if err != nil {
            return &stream, err
        }
        _, err = f.Write(file.Generator().Bytes())
        if err != nil {
            return &stream, err
        }
    }

    // Open sheet handle
    if err := stream.openSheet(); err != nil {
        return &stream, err
    }

    return &stream, nil
}

func (w *XLSXStream) WriteRow (cells []string) error {
    if err := w.writeSheet("<row>"); err != nil {
        return err
    }

    for _, cell := range (cells) {
        if err := w.writeSheet("<c t=\"inlineStr\"><is><t>"); err != nil{
            return err
        }
        if err := xml.EscapeText(w.sheet, []byte(cell)); err != nil {
            return err
        }
        if err := w.writeSheet("</t></is></c>"); err != nil {
            return err
        }
    }
    if err := w.writeSheet("</row>\n"); err != nil {
        return err
    }

    return w.zip.Flush()
}

func (w *XLSXStream) writeSheet (s string) error {
    _, err := w.sheet.Write([]byte(s))
    return err
}

func (w *XLSXStream) openSheet () error {
    var err error
    w.sheet, err = w.zip.Create(fmt.Sprintf("xl/worksheets/%v.xml", w.name))
    if err != nil {
        return err
    }

    _, err = w.sheet.Write([]byte("<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"yes\"?><worksheet xmlns=\"http://schemas.openxmlformats.org/spreadsheetml/2006/main\" xmlns:r=\"http://schemas.openxmlformats.org/officeDocument/2006/relationships\"><sheetData>\n"))
    if err != nil{
        return err
    }

    return nil
//    return w.zip.Flush()
}

func (w *XLSXStream) Close () error {
    if err := w.writeSheet("</sheetData>\n</worksheet>\n"); err != nil {
        return err
    }
    return w.zip.Close()
}

// Creates a buffer of the ContentType Stream
func (w *XLSXStream) genContentTypes () *bytes.Buffer {
    var buf bytes.Buffer

    buf.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"yes\"?>\n")
    buf.WriteString("<Types xmlns=\"http://schemas.openxmlformats.org/package/2006/content-types\">\n")
    buf.WriteString("    <Default Extension=\"rels\" ContentType=\"application/vnd.openxmlformats-package.relationships+xml\"/>\n")
    buf.WriteString("    <Default Extension=\"xml\" ContentType=\"application/xml\"/>\n")
    buf.WriteString("    <Override PartName=\"/xl/workbook.xml\" ContentType=\"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet.main+xml\"/>\n")
    buf.WriteString(
        fmt.Sprintf("    <Override PartName=\"/xl/worksheets/%v.xml\" ContentType=\"application/vnd.openxmlformats-officedocument.spreadsheetml.worksheet+xml\"/>\n",
                    w.name))
    buf.WriteString("</Types>\n")

    return &buf
}

// Creates _rels/.rels buffer
func (w *XLSXStream) genRels () *bytes.Buffer {
    var buf bytes.Buffer

    buf.WriteString("<?xml version=\"1.0\" encoding=\"UTF-8\" standalone=\"yes\"?>\n")
    buf.WriteString("<Relationships xmlns=\"http://schemas.openxmlformats.org/package/2006/relationships\">\n")
    buf.WriteString("    <Relationship Id=\"rId1\" Type=\"http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument\" Target=\"xl/workbook.xml\"/>\n")
    buf.WriteString("</Relationships>\n")

    return &buf
}

// Creates xl/workbook.xml buffer
func (w *XLSXStream) genWorkbook () *bytes.Buffer {
    var buf bytes.Buffer

    buf.WriteString("<?xml version=\"1.0\" encoding=\"utf-8\"?>\n")
    buf.WriteString("<x:workbook xmlns:x=\"http://schemas.openxmlformats.org/spreadsheetml/2006/main\">\n")
    buf.WriteString("    <x:sheets>\n")
    buf.WriteString(
        fmt.Sprintf("        <x:sheet name=\"%v\" sheetId=\"1\" r:id=\"rId1\" xmlns:r=\"http://schemas.openxmlformats.org/officeDocument/2006/relationships\" />\n",
                    w.name))
    buf.WriteString("    </x:sheets>\n")
    buf.WriteString("</x:workbook>\n")

    return &buf
}

// Creates xl/_rels/workbook.xml.rels buffer
func (w *XLSXStream) genWorkbookRels () *bytes.Buffer {
    var buf bytes.Buffer

    buf.WriteString("<?xml version=\"1.0\" encoding=\"utf-8\"?>\n")
    buf.WriteString("<Relationships xmlns=\"http://schemas.openxmlformats.org/package/2006/relationships\">\n")
    buf.WriteString(
        fmt.Sprintf("    <Relationship Type=\"http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet\" Target=\"/xl/worksheets/%v.xml\" Id=\"rId1\" />\n",
                    w.name))
    buf.WriteString("</Relationships>\n")

    return &buf
}

// Creates empty worksheet buffer
func (w *XLSXStream) genEmptyWorksheet () *bytes.Buffer {
    var buf bytes.Buffer

    buf.WriteString("<?xml version=\"1.0\" encoding=\"utf-8\"?>\n")
    buf.WriteString("<x:worksheet xmlns:x=\"http://schemas.openxmlformats.org/spreadsheetml/2006/main\">\n")
    buf.WriteString("    <x:sheetData />\n")
    buf.WriteString("</x:worksheet>\n")

    return &buf
}
