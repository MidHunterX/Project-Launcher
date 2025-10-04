package util

import (
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

func OpenBrowser(url, browser string) {
	if url == "" {
		return
	}

	if browser == "" {
		browser = "xdg-open"
	} else {
		cmdParts := strings.Fields(browser)
		if !CommandExists(cmdParts[0]) {
			Log("❌ Error: Browser command not found: " + cmdParts[0])
			os.Exit(1)
		}
	}

	Log("🔗 Waiting for server to start...")
	waitForServer(url)

	Log("🚀 Launching browser at " + url)
	var cmd *exec.Cmd
	if browser == "xdg-open" {
		cmd = exec.Command(browser, url)
	} else {
		parts := strings.Fields(browser)
		parts = append(parts, url)
		cmd = exec.Command(parts[0], parts[1:]...)
	}
	cmd.Start()
}

func waitForServer(url string) {
	timeout := 20 * time.Second
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		resp, err := http.Get(url)
		if err == nil {
			resp.Body.Close()
			return // Server is up
		}
		time.Sleep(1 * time.Second)
	}

	Log("⚠️ Warning: Timed out waiting for server. Browser will not open.")
}
