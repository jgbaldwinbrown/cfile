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
	fp *C.FILE
}

func (f CFile) Read(p []byte) (n int, err error) {
	file := f.fp
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
	file := f.fp
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

func toWhence(whence int) C.int {
	switch whence {
	case io.SeekStart:
		return C.SEEK_SET
	case io.SeekEnd:
		return C.SEEK_END
	case io.SeekCurrent:
		return C.SEEK_CUR
	default:
		panic(fmt.Errorf("toWhence: invalid whence %v", whence))
		return -1
	}
}

func (f CFile) Seek(offset int64, whence int) (int64, error) {
	file := f.fp
	cWhence := toWhence(whence)
	e := C.fseek(file, C.long(offset), cWhence);
	if e != 0 {
		return 0, IntErr(e)
	}
	pos := C.ftell(file)
	if pos == -1 {
		return 0, IntErr(pos)
	}
	return int64(pos), nil
}

func Open(path string, perm string) CFile {
	cpath := C.CString(path)
	cperm := C.CString(perm)
	fp := C.fopen(cpath, cperm)
	C.free(unsafe.Pointer(cpath))
	C.free(unsafe.Pointer(cperm))
	return CFile{fp: fp}
}

func Close(fp CFile) {
	C.fclose(fp.fp)
}

func Wrap(fp unsafe.Pointer) CFile {
	return CFile{fp: (*C.FILE)(fp)}
}
