package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// SystemStats holds the CPU and memory usage information.
type SystemStats struct {
	CPUUsage    string
	MemoryUsage string
}

func getStats() (SystemStats, error) {
	cmd := exec.Command("top", "-l", "1", "-n", "0")
	output, err := cmd.Output()
	if err != nil {
		return SystemStats{}, err
	}

	lines := strings.Split(string(output), "\n")
	stats := SystemStats{}

	for _, line := range lines {
		if strings.HasPrefix(line, "CPU usage:") {
			stats.CPUUsage = line
		} else if strings.HasPrefix(line, "PhysMem:") {
			stats.MemoryUsage = line
		}
	}

	if stats.CPUUsage == "" || stats.MemoryUsage == "" {
		return SystemStats{}, fmt.Errorf("failed to get system stats")
	}

	return stats, nil
}

func clearConsole() {
	cmd := exec.Command("clear") // For Unix/Linux/Mac
	if strings.Contains(os.Getenv("OS"), "Windows") {
		cmd = exec.Command("cmd", "/c", "cls") // For Windows
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func printStats(stats SystemStats) {
	clearConsole()
	fmt.Println("CPU Usage:", stats.CPUUsage)
	fmt.Println("Memory Usage:", stats.MemoryUsage)
}

func main() {
	// Parsing command-line flags
	refreshInterval := flag.Int("s", 1, "Refresh interval in seconds")
	flag.Parse()

	statsChan := make(chan SystemStats)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT)

	go func() {
		for {
			stats, err := getStats()
			if err != nil {
				fmt.Println("Error fetching system stats:", err)
				return
			}
			statsChan <- stats
			time.Sleep(time.Duration(*refreshInterval) * time.Second)
		}
	}()

	for {
		select {
		case <-signalChan:
			fmt.Println("\nExiting...")
			return
		case stats := <-statsChan:
			printStats(stats)
		}
	}
}
