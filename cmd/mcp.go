package cmd

import (
	"bufio"
	"fmt"
	"io"
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
	defer conn.Close()

	done := make(chan struct{})

	// TCP -> stdout
	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		close(done)
	}()

	// stdin -> TCP
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "read error: %v\n", err)
			os.Exit(1)
		}
		conn.Write(line)
	}

	// Close write side so server knows we're done, then wait for response
	conn.(*net.TCPConn).CloseWrite()
	<-done
}
