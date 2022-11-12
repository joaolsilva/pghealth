package main

import (
	"github.com/joaolsilva/pghealth/cmd"
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
		panic(err)
	}
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
