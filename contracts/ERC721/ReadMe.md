
## NRC7Full generate

```
solc-static-linux-v0.8.0 --abi --bin -o contracts/ERC721/build   --allow-paths . contracts/contracts/contracts/NRC7/NRC7Full.sol --optimize
abigen --abi contracts/ERC721/build/NRC7Full.abi --bin contracts/ERC721/build/NRC7Full.bin --pkg ERC721 --out contracts/ERC721/NRC7Full.go --type NRC7Full
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
