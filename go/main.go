package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

const (
	VERSION = "0.0.1"
)

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "-U" || os.Args[1] == "--update") {
		runUpdate()
		return
	}

	if len(os.Args) > 1 && (os.Args[1] == "-V" || os.Args[1] == "--version") {
		println("Version: " + VERSION)
		os.Exit(0)
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

	// Přidáme argumenty z CLI (vše za 'cf')
	if len(os.Args) > 1 {
		args = append(args, os.Args[1:]...)
	}

	cmd := exec.Command("docker", args...)

	// Propojíme standardní vstupy/výstupy
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// 3. Ošetření signálů (Graceful Shutdown)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		// Pokud uživatel zmáčkne Ctrl+C, Docker kontejner díky --rm zmizí
		// ale my můžeme přidat extra cleanup logiku zde
		os.Exit(0)
	}()

	// 4. Start!
	if err := cmd.Run(); err != nil {
		// Tady ignorujeme chybu exit status 130, což je běžný Ctrl+C exit
		if exitError, ok := err.(*exec.ExitError); ok && exitError.ExitCode() == 130 {
			os.Exit(0)
		}
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runUpdate() {
	fmt.Println("🚀 Updating GS CLOUD Wrangler environment...")
	// Zde pak přidáme tvůj GitHub self-update
	cmd := exec.Command("docker", "pull", "gscloudcz/wrangler-proxy:latest")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
	fmt.Println("✅ Done.")
}
