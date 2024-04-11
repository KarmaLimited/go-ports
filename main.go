package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
)

func protocolToString(connType uint32, ip string) string {
    isIPv6 := strings.Contains(ip, ":")
    switch connType {
    case syscall.SOCK_STREAM: // TCP
        if isIPv6 {
            return "TCP6"
        }
        return "TCP"
    case syscall.SOCK_DGRAM: // UDP
        if isIPv6 {
            return "UDP6"
        }
        return "UDP"
    default:
        return "Unknown"
    }
}

func displayNetworkInfo() {
    // Get network connections
    conns, err := net.Connections("all")
    if err != nil {
        log.Fatal(err)
    }

    // Clear the terminal screen
    fmt.Print("\033[H\033[2J")

    // Print the headers for the table
    fmt.Printf("%-9s %-25s %-30s %-20s %-10s %-25s\n", "Protocol", "Local Address", "Foreign Address", "State", "PID", "Process Name")

    for _, conn := range conns {
        // Filter out connections with no associated process
        if conn.Pid == 0 {
            continue
        }

        // Get process name
        proc, err := process.NewProcess(conn.Pid)
        if err != nil {
            log.Println(err)
            continue
        }
        procName, err := proc.Name()
        if err != nil {
            log.Println(err)
            continue
        }

        // Convert the protocol to a string
        protocol := protocolToString(conn.Type, conn.Laddr.IP)

        // Construct local and remote addresses
        localAddr := conn.Laddr.IP + ":" + strconv.Itoa(int(conn.Laddr.Port))
        remoteAddr := conn.Raddr.IP + ":" + strconv.Itoa(int(conn.Raddr.Port))

        // Check OS and apply color if not Windows
        var coloredProcName string
        if runtime.GOOS == "linux" || runtime.GOOS == "darwin" { // darwin is for macOS
            blue := "\033[34m"
            reset := "\033[0m"
            coloredProcName = fmt.Sprintf("%s%s%s", blue, procName, reset)
        } else {
            // For Windows or other OSes, print without color
            coloredProcName = procName
        }

        // Print the connection details
        fmt.Printf("%-9s %-25s %-30s %-20s %-10d %-25s\n", protocol, localAddr, remoteAddr, conn.Status, conn.Pid, coloredProcName)
    }
}

func clearScreen() {
    //fmt.Print("\033[H\033[2J")
		fmt.Print("\033[H\033[2J\033[3J")
}

func main() {
    // Setup a channel to listen for SIGINT signals (Control+C)
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT)

    // Set up a ticker for periodic refresh, adjust the duration for refresh rate
    ticker := time.NewTicker(5 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            clearScreen()
            displayNetworkInfo()

        case <-sigChan:
						clearScreen()
            // If Control+C is pressed, exit the program
            fmt.Println("\nExiting...")
            return
        }
    }
}
