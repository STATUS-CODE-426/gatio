/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Generated prod code for api-gateway",
	Long: `Generated prod code for api gateway from your proto files, graphql schemes,
	openapi schemes or asyncapi schemes. 
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		clientProtocol, err := cmd.Flags().GetString("client")
		if err != nil {
			return err
		}
		serverProtocol, err := cmd.Flags().GetString("server")
		if err != nil {
			return err
		}
		fmt.Println("client's using protocol", clientProtocol)
		fmt.Println("service's using protocol", serverProtocol)
		return nil
	},
}

func init() {
	createCmd.Flags().StringP("client", "c", "http1.1", "client's using protocol")
	createCmd.Flags().StringP("server", "s", "grpc", "service's using protocol")
	rootCmd.AddCommand(createCmd)
}
