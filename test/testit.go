package main

import (
	"io"
	"fmt"
	"os"

	"github.com/jgbaldwinbrown/cfile/pkg"
)

func WriteFile(path string) error {
	fp := cfile.Open(path, "w")
	fmt.Fprintln(fp, "Banana")
	cfile.Close(fp)
	return nil
}

func ReadFile(path string) error {
	fp := cfile.Open(path, "r")
	io.Copy(os.Stdout, fp)
	cfile.Close(fp)
	return nil
}

func SeekReadFile(path string, offset int64) error {
	fp := cfile.Open(path, "r")
	fp.Seek(offset, io.SeekStart)
	io.Copy(os.Stdout, fp)
	cfile.Close(fp)
	return nil
}

func main() {
	WriteFile("temp.txt")
	ReadFile("temp.txt")
	SeekReadFile("temp.txt", 2)
}
