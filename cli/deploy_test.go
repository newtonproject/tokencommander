package cli

import "testing"

func TestDeploy(t *testing.T) {
	cli := NewCLI()

	cli.TestCommand("deploy --name MyToken --symbol MT --total 100000000 --decimals 1")
}
