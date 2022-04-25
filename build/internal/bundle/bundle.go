package bundle

import (
	"archive/zip"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"os"
)

func zipFiles(zipFile string, filenames ...string) (err error) {
	if len(filenames) == 0 {
		return errors.New("must pass more than one filename to zipFiles")
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	out, err := os.Create(zipFile)
	if err != nil {
		panic(err.Error())
	}
	defer out.Close()

	for _, filename := range filenames {
		in, err := os.Open(filename)
		if err != nil {
			panic(err.Error())
		}
		defer in.Close()
		archive := zip.NewWriter(out)
		record, err := archive.Create(filename)
		if err != nil {
			panic(err)
		}

		_, err = io.Copy(record, in)
		if err != nil {
			panic(err.Error())
		}
	}
	return nil
}
