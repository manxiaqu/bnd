package main

import (
	"bytes"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/params"
)

var (
	genesisPath = "genesis.json"
)

// MakeGenesis creates a new genesis struct based on config.
// Referrer https://github.com/ethereum/go-ethereum/blob/master/cmd/puppeth/wizard_genesis.go#L39
func MakeGenesis(config *GenesisConfig, path string) *core.Genesis {
	// Construct a default genesis block
	genesis := &core.Genesis{
		Timestamp:  uint64(time.Now().Unix()),
		GasLimit:   4700000,
		Difficulty: big.NewInt(524288),
		Alloc:      make(core.GenesisAlloc),
		Config: &params.ChainConfig{
			HomesteadBlock:      big.NewInt(0),
			EIP150Block:         big.NewInt(0),
			EIP155Block:         big.NewInt(0),
			EIP158Block:         big.NewInt(0),
			ByzantiumBlock:      big.NewInt(0),
			ConstantinopleBlock: big.NewInt(0),
			PetersburgBlock:     big.NewInt(0),
			IstanbulBlock:       big.NewInt(0),
		},
	}

	// apply consensus spefic.
	switch config.Consensus {
	case "ethash":
		applyEthash(genesis)
	case "clique":
		applyClique(genesis, &config.CliqueConfig)
	default:
		panic("unsupport consensus")
	}

	// apply common.
	applyCommon(genesis, &config.CommonConfig)
	// always pre-fund for pre-compile contract address.
	prefundPrecompile(genesis)

	// store the genesis and flush to disk
	FlushGenesis(genesis, path)

	return genesis
}

func applyEthash(genesis *core.Genesis) {
	// In case of ethash, we're pretty much done
	genesis.Config.Ethash = new(params.EthashConfig)
	genesis.ExtraData = make([]byte, 32)
}

func applyClique(genesis *core.Genesis, config *CliqueConfig) {
	// In the case of clique, configure the consensus parameters
	genesis.Difficulty = big.NewInt(1)
	genesis.Config.Clique = &params.CliqueConfig{
		Period: config.Period,
		Epoch:  config.Epoch,
	}

	// Sort the signers and embed into the extra-data section
	signers := config.Signers
	for i := 0; i < len(signers); i++ {
		for j := i + 1; j < len(signers); j++ {
			if bytes.Compare(signers[i][:], signers[j][:]) > 0 {
				signers[i], signers[j] = signers[j], signers[i]
			}
		}
	}
	genesis.ExtraData = make([]byte, 32+len(signers)*common.AddressLength+65)
	for i, signer := range signers {
		copy(genesis.ExtraData[32+i*common.AddressLength:], signer[:])
	}
}

func applyCommon(genesis *core.Genesis, config *CommonConfig) {
	// Read the address of the account to fund
	for _, address := range config.Alloc {
		genesis.Alloc[address] = core.GenesisAccount{
			Balance: new(big.Int).Lsh(big.NewInt(1), 256-7), // 2^256 / 128 (allow many pre-funds without balance overflows)
		}
	}

	// Set chain config.
	genesis.Config.ChainID = new(big.Int).SetUint64(config.ChainID.Uint64())
}

func prefundPrecompile(genesis *core.Genesis) {
	// Add a batch of precompile balances to avoid them getting deleted
	for i := int64(0); i < 256; i++ {
		genesis.Alloc[common.BigToAddress(big.NewInt(i))] = core.GenesisAccount{Balance: big.NewInt(1)}
	}
}
