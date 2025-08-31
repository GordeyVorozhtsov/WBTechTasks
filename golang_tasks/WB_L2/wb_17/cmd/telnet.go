package cmd

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

var (
	host    string
	port    int
	timeout int
)

var rootCmd = &cobra.Command{
	Use:   "telnet [host]",
	Short: "Утилита сетевых запросов",
	Long:  "Утилита сетевых запросов для tcp подключения к сокету",
	Args:  cobra.ExactArgs(1),
	RunE:  runTelnet,
}

func init() {
	rootCmd.Flags().IntVarP(&port, "port", "p", 0, "port to connect")
	rootCmd.Flags().IntVarP(&timeout, "timeout", "t", 10, "connection timeout in seconds")
	rootCmd.MarkFlagRequired("port")
}

func runTelnet(cmd *cobra.Command, args []string) error {
	host := args[0]
	address := fmt.Sprintf("%s:%d", host, port)

	duration := time.Duration(timeout) * time.Second
	conn, err := net.DialTimeout("tcp", address, duration)
	if err != nil {
		return fmt.Errorf("error connecting to %s: %v", address, err)
	}
	defer conn.Close()

	fmt.Printf("Connected to %s\n", address)
	fmt.Println("Press Ctrl+D to exit")

	done := make(chan struct{})
	var wg sync.WaitGroup

	// обработка сигналов
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	wg.Add(1)
	go func() {
		defer wg.Done()
		select {
		case <-sigCh:
			fmt.Println("\nInterrupt received - closing connection")
			close(done)
		case <-done:
		}
	}()

	// чтение из сокета
	wg.Add(1)
	go func() {
		defer wg.Done()
		reader := bufio.NewReader(conn)
		buf := make([]byte, 1024)

		for {
			select {
			case <-done:
				return
			default:
				conn.SetReadDeadline(time.Now().Add(200 * time.Millisecond))
				n, err := reader.Read(buf)
				if err != nil {
					if err == io.EOF {
						fmt.Println("\nServer closed connection")
						close(done)
						return
					}
					if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
						continue
					}
					fmt.Printf("\nRead error: %v\n", err)
					close(done)
					return
				}
				if n > 0 {
					fmt.Print(string(buf[:n]))
				}
			}
		}
	}()

	// чтение из stdin и запись в сокет
	wg.Add(1)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			select {
			case <-done:
				return
			default:
				text := scanner.Text() + "\n"
				_, err := conn.Write([]byte(text))
				if err != nil {
					fmt.Printf("Write error: %v\n", err)
					close(done)
					return
				}
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Printf("STDIN read error: %v\n", err)
		} else {
			fmt.Println("Ctrl+D pressed - closing connection")
		}
		close(done)
	}()

	wg.Wait()
	fmt.Println("Connection closed")
	return nil
}
