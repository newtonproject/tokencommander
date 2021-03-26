
## TokenCommander 

`TokenCommander` project contains the following:
* Deploy TokenCommander contract with custom token name, symbol, decimals or total supply amount
* Support [NRC6](https://github.com/newtonproject/NEPs/blob/master/NEPS/nep-6.md)|[ERC20](https://github.com/ethereum/EIPs/blob/master/EIPS/eip-20.md) and [NRC7](https://github.com/newtonproject/NEPs/blob/master/NEPS/nep-7.md)|[ERC721](https://github.com/ethereum/EIPs/blob/master/EIPS/eip-721.md)
* Get basic information of the contract
* Get balance of the address on the token
* Pay token to address
* Decimal amount support for total supply or payment
* Mint token for NRC7|ERC721

## QuickStart

### Download from releases

Binary archives are published at:
* [NewChain] https://release.cloud.diynova.com/newton/TokenCommander/
* [Ethereum] https://release.cloud.diynova.com/newton/TokenCommander/ethereum/

### Building the source

To get from gitlab via `go get`, this will get source and install dependens(cobra, viper, logrus).

#### Windows

install:

```bash
git clone https://github.com/newtonproject/tokencommander.git && cd tokencommander && make install
```

run TokenCommander:

```bash
%GOPATH%/bin/tokencommander.exe
```

#### Linux or Mac

install:

```bash
git clone https://github.com/newtonproject/tokencommander.git && cd tokencommander && make install
```

run TokenCommander:

```bash
$GOPATH/bin/tokencommander
```

## Usage

### contract

Use commands `solc` and `abigen` to generate SimpleToken.go from [Newton Contracts](https://github.com/newtonproject/contracts).

- BaseToken generate

```
solc-static-linux-v0.8.0 --abi --bin -o contracts/ERC20/build --allow-paths . contracts/contracts/contracts/NRC6/BaseToken.sol --optimize
abigen --abi contracts/ERC20/build/BaseToken.abi --bin contracts/ERC20/build/BaseToken.bin --pkg ERC20 --out contracts/ERC20/BaseToken.go --type BaseToken
```

- NRC7Full generate

```
solc-static-linux-v0.8.0 --abi --bin -o contracts/ERC721/build   --allow-paths . contracts/contracts/contracts/NRC7/NRC7Full.sol --optimize
abigen --abi contracts/ERC721/build/NRC7Full.abi --bin contracts/ERC721/build/NRC7Full.bin --pkg ERC721 --out contracts/ERC721/NRC7Full.go --type NRC7Full
```

- NRC6|ERC20
    - Token name: BaseToken
    - Solidity Version: 0.8.0+commit.c7dfd78e.Linux.g++
    - abigen: abigen version 1.10.1-stable-0f9b9ae5
    - Optimization: enable with 200
    - others: default
- NRC7|ERC721
    - Token name: NRC7Full
    - Solidity Compiler Version: 0.8.0+commit.c7dfd78e.Linux.g++
    - abigen: abigen version 1.10.1-stable-0f9b9ae5
    - Optimization: enable with 200
    - others: default

### commandline client

#### Help

Use command `TokenCommander help` to display the usage.

```bash
Usage:
  TokenCommander [flags]
  TokenCommander [command]

Available Commands:
  account     Manage NewChain accounts
  add         Add custom contract
  balance     Balance of address on Token
  deploy      Deploy NewChain contract
  help        Help about any command
  info        Show contract basic info
  init        Initialize config file
  mint        Command to mint tokenID for address, only for NRC7
  pay         Command about transaction
  version     Get version of TokenCommander CLI

Flags:
  -c, --config path                The path to config file (default "./config.toml")
  -a, --contractAddress address    Contract address
  -f, --from address               the from address who pay gas
  -h, --help                       help for TokenCommander
      --mode string                use NRC6|NRC7 token (default "NRC6")
  -i, --rpcURL url                 NewChain json rpc or ipc url (default "https://rpc1.newchain.newtonproject.org")
  -s, --symbol --contractAddress   the symbol of the contract, this'll overwrite the --contractAddress when load token
  -w, --walletPath directory       Wallet storage directory (default "./wallet/")

Use "TokenCommander [command] --help" for more information about a command.
```

#### Use config.toml

You can use a configuration file to simplify the command line parameters.

One available configuration file `config.toml` is as follows:


```conf
contractaddress = "0x832c0e9Fa5fF7556E357212a42939d9c9D070bAA"
from = "0xeBF02C8C496C76079E2425D64d73030264BEA352"
mode = "NRC7"
rpcurl = "https://rpc1.newchain.newtonproject.org"
walletpath = "./wallet/"
password = "password"
```

#### Initialize config file

```bash
# Initialize config file
$ tokencommander init
```

Just press Enter to use the default configuration, and it's best to create a new user.


```bash
$ tokencommander init
Initialize config file
Enter file in which to save (./config.toml):
Enter the wallet storage directory (./wallet/):
Enter geth json rpc or ipc url (https://rpc1.newchain.newtonproject.org):
Create a default account or not: [Y/n]
Your new account is locked with a password. Please give a password. Do not forget this password.
Enter passphrase (empty for no passphrase):
Enter same passphrase again:
New accout is  0xeBF02C8C496C76079E2425D64d73030264BEA352
Get faucet for 0xeBF02C8C496C76079E2425D64d73030264BEA352
Your configuration has been saved in  ./config.toml
```

#### Create account

```bash
# Create an account with faucet
tokencommander account new --faucet

# Create 10 accounts
tokencommander account new -n 10 --faucet

# Get balance of address on NewChain or Ethereum
tokencommander account balance
```

#### Deploy contract

```bash
# Deploy NRC6 Token 'MyToken'
tokencommander deploy --name MyToken --symbol MT --total 100000000 --decimals 1

# Deploy NRC6 Token 'MyToken' in short
tokencommander deploy -n MyToken -s MT -t 100000000 -d 1

# Deploy NRC6 Token 'MyToken' with decimal total supply
tokencommander deploy -n MyToken -s MT -t 0.1 -d 8

# Deploy NRC7 Token 'MyToken'
tokencommander deploy --name MyToken --symbol MT --mode NRC7
```

#### Get basic information

```bash
# Get basic information
tokencommander info

# Get basic information of contract address
tokencommander info -a 0xdAC17F958D2ee523a2206206994597C13D831ec7

# Get basic information of contract local symbol
tokencommander info -s USDT
```

#### Get balance

```bash
# Get balance of the default address
tokencommander balance

# Get balance of the some address
tokencommander balance 0xc8B5c4cB6DB7254d082b24A96627F143E8A80c31 0xeBF02C8C496C76079E2425D64d73030264BEA352
```

#### transaction

```bash
# Pay 10 NRC6 token to other 
tokencommander pay 10 --to 0xc8B5c4cB6DB7254d082b24A96627F143E8A80c31

# Pay 0.01 NRC6 token to other 
tokencommander pay 0.01 --to 0xc8B5c4cB6DB7254d082b24A96627F143E8A80c31

# Transfer NRC7 tokenID 10 to other
tokencommander pay 10 --to 0xc8B5c4cB6DB7254d082b24A96627F143E8A80c31
```

#### Mint NRC7 Token

```bash
# Mint token for address
tokencommander mint 0xc8B5c4cB6DB7254d082b24A96627F143E8A80c31

# Mint tokenID 10 for address
tokencommander mint 0xc8B5c4cB6DB7254d082b24A96627F143E8A80c31 10

# Mint tokenID 1 for address with tokenURL
tokencommander mint 0xc8B5c4cB6DB7254d082b24A96627F143E8A80c31 10 --url https://www.newtonproject.org/tokencommander/token721/1.json
```


#### Add contract to local

```bash
# Add contract 0xdAC17F958D2ee523a2206206994597C13D831ec7 to local
tokencommander add 0xdAC17F958D2ee523a2206206994597C13D831ec7

# Add contract 0xdAC17F958D2ee523a2206206994597C13D831ec7 to local with custom symbol
tokencommander add 0xdAC17F958D2ee523a2206206994597C13D831ec7 USDT
```

#### Batch pay ERC20 token

```bash
# batch pay base on batch.txt
tokencommander batch batch.txt
tokencommander batchpay batch.txt
```