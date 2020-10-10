package main

import (
	"log"
)

func main() {
	rootCmd := newRootComand()
	rootCmd.AddCommand(newStartCommand())
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
