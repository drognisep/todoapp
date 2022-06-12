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
		return errors.New("must pass at least one filename to zipFiles")
	}
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	out, err := os.Create(zipFile)
	if err != nil {
		panic(err)
	}
	defer out.Close()
	archive := zip.NewWriter(out)
	defer archive.Close()

	for _, filename := range filenames {
		func() {
			in, err := os.Open(filename)
			if err != nil {
				panic(err)
			}
			defer in.Close()
			record, err := archive.Create(filename)
			if err != nil {
				panic(err)
			}

			_, err = io.Copy(record, in)
			if err != nil {
				panic(err)
			}
		}()
	}
	return nil
}
