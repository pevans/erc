package cmd

import (
	"bufio"
	"fmt"
	"net"
	"os"

	"github.com/pevans/erc/mcp"
	"github.com/spf13/cobra"
)

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Start MCP bridge",
	Long:  "Bridge MCP protocol from stdio to a running erc instance via TCP",
	Run: func(cmd *cobra.Command, args []string) {
		runMCPBridge()
	},
}

func init() {
	rootCmd.AddCommand(mcpCmd)
}

func runMCPBridge() {
	conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", mcp.Port))
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not connect to erc: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close() //nolint:errcheck

	done := make(chan struct{})

	// TCP -> stdout (line by line with explicit flush)
	go func() {
		defer close(done)
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	// stdin -> TCP (line by line)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()

		_, err = fmt.Fprintln(conn, line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not write to connection: %v\n", err)
			os.Exit(1)
		}
	}

	// Close write side so server knows we're done
	if tcpConn, ok := conn.(*net.TCPConn); ok {
		// we're ok with ignoring an error here
		_ = tcpConn.CloseWrite()
	}

	<-done
}
