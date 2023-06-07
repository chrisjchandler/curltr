package main

import (
	"bufio"
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type ProbeResult struct {
	Sequence int
	Address  string
	Success  bool
	Loss     float64
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <target_url>")
		return
	}

	targetURL := os.Args[1]

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter the number of probes to send: ")
	numProbesStr, _ := reader.ReadString('\n')
	numProbes, _ := strconv.Atoi(strings.TrimSpace(numProbesStr))

	reader.Reset(os.Stdin)
	fmt.Print("Enter the timeout in milliseconds: ")
	timeoutStr, _ := reader.ReadString('\n')
	timeout, _ := strconv.ParseInt(strings.TrimSpace(timeoutStr), 10, 64)

	reader.Reset(os.Stdin)
	fmt.Print("Enter the port number: ")
	portStr, _ := reader.ReadString('\n')
	port, _ := strconv.Atoi(strings.TrimSpace(portStr))

	reader.Reset(os.Stdin)
	fmt.Print("Enter the protocol type (TCP/UDP): ")
	protocolType, _ := reader.ReadString('\n')
	protocolType = strings.TrimSpace(protocolType)

	u, err := url.ParseRequestURI(targetURL)
	if err != nil {
		fmt.Printf("Invalid target URL: %s\n", targetURL)
		return
	}

	targetHost := u.Hostname()
	targetPath := u.Path

	if targetPath == "" {
		targetPath = "/"
	}

	tracerouteResults := make([]ProbeResult, numProbes)

	for i := 0; i < numProbes; i++ {
		tracerouteResults[i] = ProbeResult{
			Sequence: i + 1,
			Address:  "",
			Success:  false,
			Loss:     0,
		}

		if protocolType == "TCP" {
			err := performTracerouteTCP(targetHost, port, &tracerouteResults[i])
			if err != nil {
				fmt.Printf("Error performing traceroute: %v\n", err)
			}
		} else if protocolType == "UDP" {
			err := performTracerouteUDP(targetHost, port, &tracerouteResults[i])
			if err != nil {
				fmt.Printf("Error performing traceroute: %v\n", err)
			}
		} else {
			fmt.Println("Invalid protocol type")
			return
		}

		if tracerouteResults[i].Address != "" {
			success, loss := performProbeTCP(tracerouteResults[i].Address, port, targetPath, timeout)
			tracerouteResults[i].Success = success
			tracerouteResults[i].Loss = loss
		}
	}

	fmt.Println("Probe results:")
	for _, result := range tracerouteResults {
		if result.Address == "" {
			fmt.Printf("Probe %d: No response\n", result.Sequence)
		} else if result.Success {
			fmt.Printf("Probe %d: Address: %s, Loss: %f%%\n", result.Sequence, result.Address, result.Loss*100)
		} else {
			fmt.Printf("Probe %d: Address: %s, Connection failed\n", result.Sequence, result.Address)
		}
	}
}

func performTracerouteTCP(targetHost string, port int, result *ProbeResult) error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", targetHost, port))
	if err != nil {
		return err
	}
	defer conn.Close()

	addr := conn.RemoteAddr().(*net.TCPAddr)
	result.Address = addr.IP.String()

	return nil
}

func performTracerouteUDP(targetHost string, port int, result *ProbeResult) error {
	conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", targetHost, port))
	if err != nil {
		return err
	}
	defer conn.Close()

	addr := conn.RemoteAddr().(*net.UDPAddr)
	result.Address = addr.IP.String()

	return nil
}

func performProbeTCP(targetIP string, port int, targetPath string, timeout int64) (bool, float64) {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", targetIP, port), time.Duration(timeout)*time.Millisecond)
	if err != nil {
		return false, 1.0
	}
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Millisecond))
	buf := make([]byte, 1024)
	_, err = conn.Read(buf)
	if err != nil {
		return true, 0.0
	}

	return true, 1.0
}
