package main

import (
	"github.com/joaolsilva/pghealth/cmd"
	"log"
	"os"
	"syscall"
)

func dropRootPrivileges() (err error) {
	if syscall.Getuid() != 0 {
		return nil
	}
	// Change user / group to nobody (65534)
	err = syscall.Setgid(65534)
	if err != nil {
		return err
	}

	err = syscall.Setuid(65534)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	if err := dropRootPrivileges(); err != nil {
		log.Printf("WARNING: Unable to drop root privileges: %v", err)
	}
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
