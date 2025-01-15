package main

import (
	"io"
	"fmt"
	"os"
	"unsafe"

	"github.com/jgbaldwinbrown/cfile/pkg"
)

//#include <stdio.h>
//#include <stdlib.h>
import "C"

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

func ReadFile2(path string) error {
	cstr := C.CString(path)
	defer C.free(unsafe.Pointer(cstr))
	rstr := C.CString("r")
	defer C.free(unsafe.Pointer(rstr))
	cfp := C.fopen(cstr, rstr)
	io.Copy(os.Stdout, cfile.Wrap(unsafe.Pointer(cfp)))
	return nil
}

func main() {
	WriteFile("temp.txt")
	ReadFile("temp.txt")
	SeekReadFile("temp.txt", 2)
	ReadFile2("temp.txt")
}
