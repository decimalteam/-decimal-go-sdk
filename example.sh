FLAGS="-check-coins -check-validators -check-proposals -check-nft -check-multisig -check-stakes -check-transaction"
#FLAGS=""
./decimal-go-sdk -id  testnet-gate -log log-testnet-gate.log $FLAGS
./decimal-go-sdk -id  devnet-gate -log log-devnet-gate.log $FLAGS
./decimal-go-sdk -id  devnet-local -log log-devnet-local.log $FLAGS
./decimal-go-sdk -id  testnet-local -log log-testnet-local.log $FLAGS
