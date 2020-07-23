package cli

import "testing"

func TestBalance(t *testing.T) {
	cli := NewCLI()

	cli.TestCommand("balance 0xeF0b04a14e62434a99C4aF28C6dAb52ba9B1C8F3 0xDC8F76075Db000Fa70fdA3AA2c95d63F22A10a67 0x6a038842f9E9010624eAeB5f30ec5004C05EE21D")
	cli.TestCommand("balance")
}
