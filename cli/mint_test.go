package cli

import "testing"

func TestMint(t *testing.T) {
	cli := NewCLI()

	cli.TestCommand("mint 1 0x8bBc8efCE7Ac8CC7F3D954d2966C8e92E66eE4A8")

}
