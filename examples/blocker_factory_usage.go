package main

import (
	"fmt"
	"log"

	"github.com/d3m0k1d/BanForge/internal/blocker"
	"github.com/d3m0k1d/BanForge/internal/logger"
)

func main() {
	// Initialize logger
	appLogger := logger.New(true)

	// Create factory
	factory := blocker.NewBlockerFactory(appLogger)

	// Example 1: List all available blockers
	fmt.Println("Available blockers:")
	for _, name := range blocker.ListAvailable(appLogger) {
		fmt.Printf("  - %s\n", name)
	}

	// Example 2: Use nftables
	fmt.Println("\n=== NFTables Example ===")
	nftBlocker, err := factory.Create(blocker.BlockerTypeNftables, "/etc/nftables.conf")
	if err != nil {
		log.Fatalf("failed to create nftables blocker: %v", err)
	}

	// Check if available
	if !nftBlocker.IsAvailable() {
		fmt.Println("NFTables is not available on this system")
	} else {
		fmt.Printf("Blocker: %s\n", nftBlocker.Name())

		// Setup
		if err := nftBlocker.Setup(); err != nil {
			fmt.Printf("Failed to setup: %v\n", err)
		}

		// Ban an IP
		if err := nftBlocker.Ban("192.168.1.100"); err != nil {
			fmt.Printf("Failed to ban IP: %v\n", err)
		}

		// List banned IPs
		bannedIPs, err := nftBlocker.List()
		if err != nil {
			fmt.Printf("Failed to list IPs: %v\n", err)
		} else {
			fmt.Println("Banned IPs:")
			for _, ip := range bannedIPs {
				fmt.Printf("  - %s\n", ip)
			}
		}

		// Cleanup
		if err := nftBlocker.Close(); err != nil {
			fmt.Printf("Failed to close: %v\n", err)
		}
	}

	// Example 3: Use UFW
	fmt.Println("\n=== UFW Example ===")
	ufwBlocker, err := factory.Create(blocker.BlockerTypeUfw, "")
	if err != nil {
		log.Fatalf("failed to create ufw blocker: %v", err)
	}

	if !ufwBlocker.IsAvailable() {
		fmt.Println("UFW is not available on this system")
	} else {
		fmt.Printf("Blocker: %s\n", ufwBlocker.Name())
		// UFW operations...
	}

	// Example 4: Create from string type
	fmt.Println("\n=== String Type Example ===")

blocker, err := factory.CreateFromString("ufw", "")
	if err != nil {
		log.Fatalf("failed to create blocker: %v", err)
	}
	fmt.Printf("Created blocker: %s\n", blocker.Name())
}
