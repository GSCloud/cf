package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

const (
	VERSION   = "0.0.4"
	NAME      = "Cloudflare Wrangler Proxy"
	URL       = "https://github.com/GSCloud/cf"
	BOLD      = "\033[1m"
	UNDERLINE = "\033[4m"
	RESET     = "\033[0m"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	arg := os.Args[1]
	switch arg {
	case "-V", "--version":
		fmt.Printf(NAME+" v%s\n", VERSION)
		return
	case "-U", "--update":
		runUpdate1()
		if runtime.GOOS != "linux" || runtime.GOARCH != "amd64" {
			fmt.Printf("❌ Error: Self-updater is only supported on Linux amd64. Current: %s %s\n", runtime.GOOS, runtime.GOARCH)
			return
		}
		runUpdate2()
		return
	case "-h", "--help", "help":
		printHelp()
		return
	case "docs":
		fmt.Println("🌐 Opening Cloudflare Docs on your host system ...")
		openBrowser("https://developers.cloudflare.com/workers/wrangler/commands/")
		return
	case "purgecache":
		// future logic
		return
	case "purgeallcache":
		// future logic
		return
	}

	// Docker command
	cwd, _ := os.Getwd()
	dockerImage := "gscloudcz/wrangler-proxy:latest"

	// Docker parameters
	args := []string{
		"run", "--rm", "-it",
		"-v", cwd + ":/app",
		"-w", "/app",
		"-e", "CLOUDFLARE_API_TOKEN",
		"-e", "CLOUDFLARE_ACCOUNT_ID",
		"--network", "host",
		dockerImage,
	}

	// add all args
	if len(os.Args) > 1 {
		args = append(args, os.Args[1:]...)
	}
	cmd := exec.Command("docker", args...)

	// connector
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		os.Exit(0)
	}()

	if err := cmd.Run(); err != nil {
		// ignore exit status 130 = Ctrl+C
		if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() == 130 {
			os.Exit(0)
		}
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// open browser
func openBrowser(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	default:
		fmt.Printf("Please open this URL in your browser: %s\n", url)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open browser: %v\n", err)
	}
}

// print help
func printHelp() {
	fmt.Printf("%s"+NAME+" v%s%s\n", BOLD, VERSION, RESET)
	fmt.Printf("%s%s%s\n\n", UNDERLINE, URL, RESET)
	fmt.Println("Usage: cf [command] [options]")
	fmt.Println("\nGlobal options:")
	fmt.Println("  -U, --update     🌞 Update the Docker image and Go binary")
	fmt.Println("  -V, --version    💁 Show version information")
	fmt.Println("  -h, --help       📄 Show this help information")
	fmt.Println("\nCustom commands:")
	fmt.Println("  docs \t\t\t Open Cloudflare documentation in a browser")
	fmt.Println("  purgecache <domain> \t Purge specific cache (TBD)")
	fmt.Println("  purgecacheall \t Purge all caches (TBD)")
	fmt.Println("\nAll other commands are passed directly to Cloudflare Wrangler.")
}

// self-updater, part 1
func runUpdate1() {
	fmt.Printf("🚀 Updating %s ...\n", NAME)
	cmd := exec.Command("docker", "pull", "gscloudcz/wrangler-proxy:latest")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	fmt.Println("✅ Done.")
}

// self-updater, part 2
func runUpdate2() {
	fmt.Println("📡 Updating Go binary ...")
	updateURL := fmt.Sprintf("https://github.com/GSCloud/cf/raw/refs/heads/master/cf?t=%d", time.Now().Unix())
	if err := doSelfUpdate(updateURL); err != nil {
	} else {
		fmt.Println("♥️ Binary updated to the latest version.")
	}
}

// self-updater, main
func doSelfUpdate(url string) error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	fmt.Printf("Binary updater path: %s.\n", exePath)

	tempPath := exePath + ".tmp"

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("❌ Download error: %v\n", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned %d", resp.StatusCode)
	}

	f, err := os.OpenFile(tempPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Printf("❌ Error: Cannot create temp file %s. Try running with sudo.\n", tempPath)
		return err
	}

	if _, err := io.Copy(f, resp.Body); err != nil {
		f.Close()
		return err
	}
	f.Close()

	os.Remove(exePath)
	if err := os.Rename(tempPath, exePath); err != nil {
		fmt.Printf("❌ Error: Final rename failed: %v\n", err)
		return err
	}
	return nil
}
