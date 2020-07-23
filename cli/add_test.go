package cli

import "testing"

func TestAdd(t *testing.T) {
	cli := NewCLI()

	cli.TestCommand("add 0xeF0b04a14e62434a99C4aF28C6dAb52ba9B1C8F3")
}
