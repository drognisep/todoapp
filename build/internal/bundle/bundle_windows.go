package bundle

import (
	"fmt"
	"runtime"
)

const (
	appName = "TODO"
)

func CreateBundle(version string) (err error) {
	bundleName := fmt.Sprintf("%s_%s_%s-%s.zip", appName, version, runtime.GOOS, runtime.GOARCH)
	return zipFiles(bundleName, "TODO.exe")
}
