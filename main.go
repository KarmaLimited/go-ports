package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
)

func protocolToString(connType uint32, ip string) string {
	isIPv6 := strings.Contains(ip, ":")
	switch connType {
	case syscall.SOCK_STREAM:
		if isIPv6 {
			return "TCP6"
		}
		return "TCP"
	case syscall.SOCK_DGRAM:
		if isIPv6 {
			return "UDP6"
		}
		return "UDP"
	default:
		return "Unknown"
	}
}

func displayNetworkInfo(s tcell.Screen) {
	s.Clear()

	conns, err := net.Connections("all")
	if err != nil {
		log.Fatal(err)
	}

	titles := []string{"Protocol", "Local Address", "Foreign Address", "State", "PID", "Process Name"}
	columnWidths := []int{10, 30, 30, 20, 10, 25}
printRow(s, titles, columnWidths, 0, tcell.StyleDefault.Bold(true).Foreground(tcell.ColorWhite))


	row := 1
	for _, conn := range conns {
		if conn.Pid == 0 {
			continue
		}

		proc, err := process.NewProcess(conn.Pid)
		if err != nil {
			continue
		}
		procName, err := proc.Name()
		if err != nil {
			continue
		}

		protocol := protocolToString(conn.Type, conn.Laddr.IP)
		localAddr := conn.Laddr.IP + ":" + strconv.Itoa(int(conn.Laddr.Port))
		remoteAddr := conn.Raddr.IP + ":" + strconv.Itoa(int(conn.Raddr.Port))
		values := []string{protocol, localAddr, remoteAddr, conn.Status, strconv.Itoa(int(conn.Pid)), procName}
		bgColor := tcell.ColorBlack
		if row%2 == 1 {
			bgColor = tcell.ColorGray
		}
printRow(s, values, columnWidths, row, tcell.StyleDefault.Background(bgColor).Foreground(tcell.ColorWhite), procName)

		row++
	}

	s.Show()
}

func printRow(s tcell.Screen, cols []string, widths []int, row int, style tcell.Style, procName ...string) {
	x := 0
	for i, col := range cols {
		// Default text color is white
		currentStyle := style.Foreground(tcell.ColorWhite)
		if i == 5 && len(procName) > 0 && cols[i] == procName[0] {
			currentStyle = style.Foreground(tcell.ColorRoyalBlue)
		}
		for _, r := range fmt.Sprintf("%-*s", widths[i], col) {
			s.SetContent(x, row, r, nil, currentStyle)
			x++
		}
		x++ // Space between columns
	}
}

func main() {
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("failed to create screen: %v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("failed to initialize screen: %v", err)
	}
	defer s.Fini()

	s.EnableMouse()
	s.Clear()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		for {
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyCtrlC {
					sigChan <- syscall.SIGINT
				}
			}
		}
	}()

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

loop:
	for {
		select {
		case <-ticker.C:
			displayNetworkInfo(s)
		case <-sigChan:
			break loop
		}
	}
	fmt.Println("\nExiting...")
}
