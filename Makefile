VERSION=0.1.0

check-fyne:
	fyne --version || go install fyne.io/fyne/v2/fyne@latest

package-win: export GOOS=windows
package-win: export GOARCH=amd64
package-win: check-fyne
	fyne package
	go run build/build.go $(VERSION)
