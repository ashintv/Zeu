package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/ashintv/Zeu/internal/harness"
)

type CLI struct {
	agent *harness.Agent
}

func NewCLI(agent *harness.Agent) *CLI {
	return &CLI{
		agent: agent,
	}
}

func (c *CLI) Run() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Zeu CLI - Interactive Agent Terminal")
	fmt.Println("Type 'exit' or 'quit' to exit. Press Ctrl+C during execution to cancel.")
	fmt.Println("--------------------------------------------------")

	for {
		fmt.Print("\nZeu > ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("\nError reading input:", err)
			return
		}

		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		if input == "exit" || input == "quit" {
			fmt.Println("Exiting. Goodbye!")
			return
		}

		c.executePrompt(input)
	}
}


func (c *CLI) executePrompt(prompt string) {
	// Set up cancellation context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Capture interrupt signal (Ctrl+C)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigChan)

	// Goroutine to monitor Ctrl+C
	go func() {
		select {
		case <-sigChan:
			fmt.Println("\n[Interrupt received, cancelling execution...]")
			cancel()
		case <-ctx.Done():
			// Exit goroutine cleanly if normal completion
		}
	}()

	resChan := c.agent.Invoke(ctx, prompt)
	for res := range resChan {
		fmt.Print(res)
	}
	fmt.Println()
}
