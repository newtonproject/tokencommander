package cli

import "testing"

func TestTx(t *testing.T) {
	cli := NewCLI()

	cli.TestCommand("tx pay 5 --to 0x6a038842f9E9010624eAeB5f30ec5004C05EE21D")

}
