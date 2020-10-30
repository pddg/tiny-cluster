package main

import (
	"fmt"
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"

	"github.com/pddg/tiny-cluster/pkg/boot"
)

func newStartCommand() *cobra.Command {
	var (
		listenPort  int
		bootFileDir string
	)
	startCmd := &cobra.Command{
		Use:   "start",
		Short: "Start server",
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := os.Stat(bootFileDir); err != nil {
				log.Fatalf("%s does not exist.", bootFileDir)
			}
			e := echo.New()

			e.Use(middleware.Logger())

			e.GET("/default.ipxe", boot.IPXEScriptHandler)
			e.Static("/boot", bootFileDir)

			return e.Start(fmt.Sprintf(":%d", listenPort))
		},
	}
	startCmd.Flags().IntVarP(&listenPort, "port", "p", 8080, "Listen port number")
	startCmd.Flags().StringVarP(&bootFileDir, "dist", "d", "/opt/bootserver", "Path to files to distribute")
	return startCmd
}
