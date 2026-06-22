package main

import (
	"flag"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/ashintv/Zeu/internal/ai"
	"github.com/ashintv/Zeu/internal/harness"
	"github.com/ashintv/Zeu/internal/logger"
	"github.com/ashintv/Zeu/internal/tools"
	"github.com/ashintv/Zeu/internal/tui"
)

func main() {
	// Define command-line flags
	verbose := flag.Bool("verbose", false, "Enable verbose logging (INFO level)")
	debug := flag.Bool("debug", false, "Enable debug logging (DEBUG level)")
	flag.Parse()

	// Configure logger based on flags
	if *debug {
		logger.SetEnabled(true)
		logger.SetDefaultLevel(logger.DebugLevel)
	} else if *verbose {
		logger.SetEnabled(true)
		logger.SetDefaultLevel(logger.InfoLevel)
	}
	// If neither flag is set, logger remains disabled by default

	reg := tools.NewToolRegistry()
	ai := ai.NewAI()

	agent := harness.CreateAgent(ai, reg)

	cli := tui.GetNewCli(agent)

	p := tea.NewProgram(cli)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
