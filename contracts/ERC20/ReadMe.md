
## BaseToken generate

```
solc-static-linux-v0.8.0 --abi --bin -o contracts/ERC20/build --allow-paths . contracts/contracts/contracts/NRC6/BaseToken.sol --optimize
abigen --abi contracts/ERC20/build/BaseToken.abi --bin contracts/ERC20/build/BaseToken.bin --pkg ERC20 --out contracts/ERC20/BaseToken.go --type BaseToken
```

### Version

```
$ solc-static-linux-v0.8.0 --version
solc, the solidity compiler commandline interface
Version: 0.8.0+commit.c7dfd78e.Linux.g++
```

```
$ abigen --version
abigen version 1.10.1-stable-0f9b9ae5
```
