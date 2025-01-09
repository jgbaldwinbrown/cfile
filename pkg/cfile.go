package cfile

import (
	"io"
	"fmt"
	"unsafe"
)

//#include <stdio.h>
//#include <stdlib.h>
import "C"

type IntErr C.int

func (i IntErr) Error() string {
	return fmt.Sprint(i)
}

type CFile struct {
	Fp *C.FILE
}

func (f CFile) Read(p []byte) (n int, err error) {
	file := f.Fp
	ptr := unsafe.Pointer(&p[0])
	nread := C.fread(ptr, 1, C.size_t(len(p)), file)
	n = int(nread)
	if n == len(p) {
		return n, nil
	}
	if C.feof(file) != 0 {
		return n, io.EOF
	}
	if e := C.ferror(file); e != 0 {
		return n, IntErr(e)
	}
	return n, nil
}

func (f CFile) Write(p []byte) (n int, err error) {
	file := f.Fp
	ptr := unsafe.Pointer(&p[0])
	nwrit := C.fwrite(ptr, 1, C.size_t(len(p)), file)
	n = int(nwrit)
	if n == len(p) {
		return n, nil
	}
	if e := C.ferror(file); e != 0 {
		return n, IntErr(e)
	}
	return n, nil
}

func Open(path string, perm string) CFile {
	cpath := C.CString(path)
	cperm := C.CString(perm)
	fp := C.fopen(cpath, cperm)
	C.free(unsafe.Pointer(cpath))
	C.free(unsafe.Pointer(cperm))
	return CFile{Fp: fp}
}

func Close(fp CFile) {
	C.fclose(fp.Fp)
}
