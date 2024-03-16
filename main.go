package main

import (
    "fmt"
    "net/http"
    "syscall"
)

// ping uses an HTTP GET request to check if a website is accessible.
func ping(url string, c chan string) {
    _, err := http.Get(url)
    if err != nil {
        c <- fmt.Sprintf("Error: %s is not accessible", url)
        return
    }
    c <- fmt.Sprintf("Success: %s is accessible", url)
}

func main() {
    // Create a channel to communicate between goroutines
    c := make(chan string)

    websites := []string{
        "https://www.google.com",
        "https://www.github.com",
        "https://www.stackoverflow.com",
    }

    // Fetch the PID using syscall
    pid := syscall.Getpid()
    fmt.Printf("Running syscalls example with PID: %d\n", pid)

    // Start a goroutine for each website to ping
    for _, site := range websites {
        go ping(site, c)
    }

    // Collect the results
    for i := 0; i < len(websites); i++ {
        fmt.Println(<-c)
    }

    fmt.Println("Done!")
}

