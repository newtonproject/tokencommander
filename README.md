
## TokenCommander 

`TokenCommander` project contains the following:
* Deploy TokenCommander contract with custom token name, symbol, decimals or total supply amount
* Support [NRC20|ERC20](https://github.com/ethereum/EIPs/blob/master/EIPS/eip-20.md) and [NRC721|ERC721](https://github.com/ethereum/EIPs/blob/master/EIPS/eip-721.md)
* Get basic information of the contract
* Get balance of the address on the token
* Pay token to address
* Decimal amount support for total supply or payment
* Mint token for NRC721|ERC721

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
go get github.com/newtonproject/tokencommander
```

run TokenCommander:

```bash
%GOPATH%/bin/TokenCommander.exe
```

#### Linux or Mac

install:

```bash
git config --global url."git@gitlab.newtonproject.org:".insteadOf "https://gitlab.newtonproject.org/"
go get github.com/newtonproject/tokencommander
```

run TokenCommander:

```bash
$GOPATH/bin/TokenCommander
```

## Usage

### contract

Use commands `go generate` or `abigen` to generate SimpleToken.go from [SimpleToken.sol](https://github.com/OpenZeppelin/openzeppelin-solidity/blob/master/contracts/token/ERC20/BasicToken.sol).

```bash
abigen --sol contract/SimpleToken.sol --pkg cli --out cli/SimpleToken.go
```

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
  mint        Command to mint tokenID for address, only for NRC721
  pay         Command about transaction
  version     Get version of TokenCommander CLI

Flags:
  -c, --config path                The path to config file (default "./config.toml")
  -a, --contractAddress address    Contract address
  -f, --from address               the from address who pay gas
  -h, --help                       help for TokenCommander
      --mode string                use NRC20|NRC721 token (default "NRC20")
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
mode = "NRC721"
rpcurl = "https://rpc1.newchain.newtonproject.org"
walletpath = "./wallet/"
password = "password"
```

#### Initialize config file

```bash
# Initialize config file
$ TokenCommander init
```

Just press Enter to use the default configuration, and it's best to create a new user.


```bash
$ TokenCommander init
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
TokenCommander account new --faucet

# Create 10 accounts
TokenCommander account new -n 10 --faucet

# Get balance of address on NewChain or Ethereum
TokenCommander account balance
```

#### Deploy contract

```bash
# Deploy NRC20 Token 'MyToken'
TokenCommander deploy --name MyToken --symbol MT --total 100000000 --decimals 1

# Deploy NRC20 Token 'MyToken' in short
TokenCommander deploy -n MyToken -s MT -t 100000000 -d 1

# Deploy NRC20 Token 'MyToken' with decimal total supply
TokenCommander deploy -n MyToken -s MT -t 0.1 -d 8

# Deploy NRC721 Token 'MyToken'
TokenCommander deploy --name MyToken --symbol MT --mode NRC721
```

#### Get basic information

```bash
# Get basic information
TokenCommander info

# Get basic information of contract address
TokenCommander info -a 0xdAC17F958D2ee523a2206206994597C13D831ec7

# Get basic information of contract local symbol
TokenCommander info -s USDT
```

#### Get balance

```bash
# Get balance of the default address
TokenCommander balance

# Get balance of the some address
TokenCommander balance 0xc8B5c4cB6DB7254d082b24A96627F143E8A80c31 0xeBF02C8C496C76079E2425D64d73030264BEA352
```

#### transaction

```bash
# Pay 10 NRC20 token to other 
TokenCommander pay 10 --to 0xc8B5c4cB6DB7254d082b24A96627F143E8A80c31

# Pay 0.01 NRC20 token to other 
TokenCommander pay 0.01 --to 0xc8B5c4cB6DB7254d082b24A96627F143E8A80c31

# Transfer NRC721 tokenID 10 to other
TokenCommander pay 10 --to 0xc8B5c4cB6DB7254d082b24A96627F143E8A80c31
```

#### Mint NRC721 Token

```bash
# Mint token for address
TokenCommander mint 0xc8B5c4cB6DB7254d082b24A96627F143E8A80c31

# Mint tokenID 10 for address
TokenCommander mint 0xc8B5c4cB6DB7254d082b24A96627F143E8A80c31 10
```


#### Add contract to local

```bash
# Add contract 0xdAC17F958D2ee523a2206206994597C13D831ec7 to local
TokenCommander add 0xdAC17F958D2ee523a2206206994597C13D831ec7

# Add contract 0xdAC17F958D2ee523a2206206994597C13D831ec7 to local with custom symbol
TokenCommander add 0xdAC17F958D2ee523a2206206994597C13D831ec7 USDT
```

#### Batch pay ERC20 token

```bash
# batch pay base on batch.txt
TokenCommander batch batch.txt
TokenCommander batchpay batch.txt
```