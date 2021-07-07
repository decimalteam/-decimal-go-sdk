# Decimal JS SDK
For detailed explanation on how things work, checkout the:

- Decimal SDK docs.
- Decimal Console site.


## Install
```
go get bitbucket.org/decimalteam/decimal-go-sdk
```

## Usage

### Init decapi
```
import (
    ...
    decapi "bitbucket.org/decimalteam/decimal-go-sdk/api"
    wallet "bitbucket.org/decimalteam/decimal-go-sdk/wallet"
)

const (
    hostURL = "https://mainnet-gate.decimalchain.com/api/" 
    // or for testnet: "https://testnet-gate.decimalchain.com/api/"
)

func main() {
    api := decapi.NewAPI(hostURL)
    account := wallet.NewAccountFromMnemonicWords(mnemonic

    
}
```

This README would normally document whatever steps are necessary to get your application up and running.

### What is this repository for? ###

* Quick summary
* Version
* [Learn Markdown](https://bitbucket.org/tutorials/markdowndemo)

### How do I get set up? ###

* Summary of set up
* Configuration
* Dependencies
* Database configuration
* How to run tests
* Deployment instructions

### Contribution guidelines ###

* Writing tests
* Code review
* Other guidelines

### Who do I talk to? ###

* Repo owner or admin
* Other community or team contact