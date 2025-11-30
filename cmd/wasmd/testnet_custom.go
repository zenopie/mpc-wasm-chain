package main

// Custom testnet start that enables vote extensions
// Uses SDK network but with extended timeouts for vote extensions

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	cmtservice "github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
	"github.com/cosmos/cosmos-sdk/testutil/network"

	"mpc-wasm-chain/app"
)

// startTestnetCustom starts a testnet with vote extensions enabled
// It uses the SDK network but with extended timeouts
func startTestnetCustom(cmd *cobra.Command, args startArgs) error {
	networkConfig := network.DefaultConfig(app.NewTestNetworkFixture)

	if args.chainID != "" {
		networkConfig.ChainID = args.chainID
	}
	networkConfig.SigningAlgo = args.algo
	networkConfig.MinGasPrices = args.minGasPrices
	networkConfig.NumValidators = args.numValidators
	networkConfig.EnableLogging = args.enableLogging
	networkConfig.RPCAddress = args.rpcAddress
	networkConfig.APIAddress = args.apiAddress
	networkConfig.GRPCAddress = args.grpcAddress
	networkConfig.PrintMnemonic = args.printMnemonic
	networkConfig.TimeoutCommit = args.timeoutCommit

	networkLogger := network.NewCLILogger(cmd)

	baseDir := fmt.Sprintf("%s/%s", args.outputDir, networkConfig.ChainID)
	if _, err := os.Stat(baseDir); !os.IsNotExist(err) {
		return fmt.Errorf(
			"testnet directory already exists for chain-id '%s': %s, please remove or select a new --chain-id",
			networkConfig.ChainID, baseDir)
	}

	cmd.Println("Starting testnet with vote extensions support...")
	cmd.Printf("Validators: %d, Commit Timeout: %s\n", args.numValidators, args.timeoutCommit)
	cmd.Println("(Genesis will be patched to enable vote extensions at height 2)")

	// Start the network - SDK will call our AppConstructor which patches genesis
	testnet, err := network.New(networkLogger, baseDir, networkConfig)
	if err != nil {
		// The SDK's network.New has a hardcoded 5s timeout which may not be enough
		// with vote extensions. Let's check if the network is actually running.
		cmd.PrintErrln("SDK network.New returned error:", err)
		cmd.Println("\nChecking if network is actually running...")

		// The network object might still be valid even if WaitForHeight timed out
		// Try to connect and check
		if testnet != nil && len(testnet.Validators) > 0 {
			cmd.Println("Network object exists with validators, attempting to wait longer...")

			if waitErr := waitForNetworkWithTimeout(testnet, 60*time.Second); waitErr != nil {
				cmd.PrintErrln("Extended wait also failed:", waitErr)
				testnet.Cleanup()
				return fmt.Errorf("network failed to start: %w", err)
			}

			cmd.Println("\n=== Network is running! ===")
		} else {
			return fmt.Errorf("testnet creation failed: %w", err)
		}
	}

	cmd.Println("\n=== Testnet started ===")
	cmd.Printf("Chain ID: %s\n", networkConfig.ChainID)
	cmd.Printf("Validators: %d\n", networkConfig.NumValidators)
	cmd.Println("\nVote extensions enabled at height 2")
	cmd.Println("\nEndpoints:")
	cmd.Printf("  RPC: %s\n", networkConfig.RPCAddress)
	cmd.Printf("  API: %s\n", networkConfig.APIAddress)
	cmd.Printf("  gRPC: %s\n", networkConfig.GRPCAddress)
	cmd.Println("\nPress Enter to stop...")

	// Wait for user input before cleanup
	if _, err := fmt.Scanln(); err != nil {
		// Ignore scan errors (e.g., EOF)
	}

	cmd.Println("Shutting down testnet...")
	testnet.Cleanup()

	return nil
}

// waitForNetworkWithTimeout waits for the network to produce a block with a custom timeout
func waitForNetworkWithTimeout(n *network.Network, timeout time.Duration) error {
	if len(n.Validators) == 0 {
		return fmt.Errorf("no validators")
	}

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	val := n.Validators[0]
	queryClient := cmtservice.NewServiceClient(val.ClientCtx)

	for {
		select {
		case <-timer.C:
			return fmt.Errorf("timeout waiting for first block after %v", timeout)
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			res, err := queryClient.GetLatestBlock(ctx, &cmtservice.GetLatestBlockRequest{})
			cancel()

			if err == nil && res != nil && res.SdkBlock != nil {
				height := res.SdkBlock.Header.Height
				if height > 0 {
					fmt.Printf("Network reached height %d\n", height)
					return nil
				}
			}
			fmt.Print(".")
		}
	}
}
