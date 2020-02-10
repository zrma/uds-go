package browser

import "os/exec"

// Open help to oauth with a browser on specific os platform
func Open(url string) error {
	return exec.Command("xdg-open", url).Start()
}
