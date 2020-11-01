package archive

import (
    "archive/zip"
    "bytes"
    "compress/flate"
    "fmt"
    "io"
    "io/ioutil"
    "log"
    "os"
)

func ZipReadWrite() {
    // Create a buffer to write our archive to
    buf := new(bytes.Buffer)
    // Create a new zip.archive
    w := zip.NewWriter(buf)

    // Register a custom Deflate compressor.
    w.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
        return flate.NewWriter(out, flate.BestCompression)
    })

    // Add some files to the archive.
    var files = []struct {
        Name, Body string
    }{
        {"readme.txt", "This archive contains some text files."},
        {"gopher.txt", "Gopher names:\nGeorge\nGeoffrey\nGonzo"},
        {"todo.txt", "Get animal handling licence.\nWrite more examples."},
    }

    for _, file := range files {
        f, err := w.Create(file.Name)
        if err != nil {
            log.Fatal(err)
        }
        _, err = f.Write([]byte(file.Body))
        if err != nil {
            log.Fatal()
        }
    }

    err := w.Close()
    if err != nil {
        log.Fatal(err)
    }

    // 写入文件
    err = ioutil.WriteFile("testdata/readme.zip", buf.Bytes(), 0644)
    if err != nil {
        log.Fatal(err)
    }

    r, err := zip.OpenReader("testdata/readme.zip")
    if err != nil {
        log.Fatal(err)
    }

    defer r.Close()

    for _, f := range r.File {
        fmt.Printf("Contents of %s:\n", f.Name)
        rc, err := f.Open()
        if err != nil {
            log.Fatal("open err",err)
        }
        _, err = io.Copy(os.Stdout, rc)
        if err != nil {
            log.Fatal(err)
        }
        rc.Close()
        fmt.Println()
    }
}