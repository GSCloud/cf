package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
)

const (
	VERSION = "0.0.1"
	NAME    = "GSC Cloudflare Wrangler Proxy"
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
		runUpdate()
		return
	case "-h", "--help", "help":
		printHelp()
		return
	case "docs":
		fmt.Println("🌐 Opening Cloudflare Docs on your host system...")
		openBrowser("https://developers.cloudflare.com/workers/wrangler/commands/")
		return
	case "purgecache":
		// future logic
		return
	case "purgeallcache":
		// future logic
		return
	}

	// 2. Příprava Docker příkazu
	cwd, _ := os.Getwd()
	dockerImage := "gscloudcz/wrangler-proxy:latest"

	// Docker params
	args := []string{
		"run", "--rm", "-it",
		"-v", cwd + ":/app",
		"-w", "/app",
		"-e", "CLOUDFLARE_API_TOKEN", // Docker automaticky vezme hodnotu z hostitele
		"-e", "CLOUDFLARE_ACCOUNT_ID",
		"--network", "host",
		dockerImage,
	}

	// add all args
	if len(os.Args) > 1 {
		args = append(args, os.Args[1:]...)
	}
	cmd := exec.Command("docker", args...)

	// coonect input/output
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

func runUpdate() {
	fmt.Printf("🚀 Updating %s ...\n", NAME)
	cmd := exec.Command("docker", "pull", "gscloudcz/wrangler-proxy:latest")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	fmt.Println("✅ Done.")
}

func printHelp() {
	fmt.Printf(NAME+" v%s\n", VERSION)
	fmt.Println("Usage: cf [command] [options]")
	fmt.Println("\nGlobal options:")
	fmt.Println("  -U, --update     Update the Go binary and the Docker image")
	fmt.Println("  -V, --version    Show version information")
	fmt.Println("  -h, --help       Show this help message")
	fmt.Println("\nCustom commands:")
	fmt.Println("  docs             Open Cloudflare documentation in a browser")
	fmt.Println("  purgecache       Purge specific cache (planned)")
	fmt.Println("  purgeallcache    Purge all caches (planned)")
	fmt.Println("\nAll other commands are passed directly to Cloudflare Wrangler.")
}
