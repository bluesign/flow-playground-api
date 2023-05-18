package blockchain

import (
	"context"
	"fmt"
	"github.com/onflow/cadence"
	jsoncdc "github.com/onflow/cadence/encoding/json"
	"github.com/onflow/cadence/runtime/common"
	"github.com/onflow/cadence/runtime/parser"
	kit "github.com/onflow/flow-cli/flowkit"
	"github.com/onflow/flow-cli/flowkit/accounts"
	"github.com/onflow/flow-cli/flowkit/config"
	"github.com/onflow/flow-cli/flowkit/gateway"
	"github.com/onflow/flow-cli/flowkit/output"
	"github.com/onflow/flow-cli/flowkit/tests"
	"github.com/onflow/flow-cli/flowkit/transactions"
	emu "github.com/onflow/flow-emulator"
	"github.com/onflow/flow-emulator/storage/memstore"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/onflow/flow-go-sdk/templates"
	"github.com/pkg/errors"
)

// blockchain interface defines an abstract API for communication with the blockchain. It hides complexity from the
// consumer and communicates using flow native types.
type blockchain interface {
	// executeTransaction builds and executes a transaction and uses provided authorizers for signing.
	executeTransaction(
		script string,
		arguments []string,
		authorizers []flow.Address,
	) (*flow.Transaction, *flow.TransactionResult, error)

	// executeScript executes a provided script with the arguments.
	executeScript(script string, arguments []string) (cadence.Value, error)

	// createAccount creates a new account and returns it along with transaction and result.
	createAccount() (*flow.Account, error)

	// getAccount gets an account by the address and also returns its storage.
	getAccount(address flow.Address) (*flow.Account, *emu.AccountStorage, error)

	// deployContract deploys a contract on the provided address and returns transaction and result.
	deployContract(address flow.Address, script string) (*flow.Transaction, *flow.TransactionResult, error)

	// removeContract removes specified contract from provided address and returns transaction and result.
	removeContract(address flow.Address, contractName string) (*flow.Transaction, *flow.TransactionResult, error)

	// getLatestBlock height from the network.
	getLatestBlockHeight() (int, error)
}

var _ blockchain = &flowKit{}

type flowKit struct {
	blockchain *kit.Flowkit
}

func newFlowkit() (*flowKit, error) {
	readerWriter, _ := tests.ReaderWriter()
	state, err := kit.Init(readerWriter, crypto.ECDSA_P256, crypto.SHA3_256)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create flow-kit state")
	}

	gw := gateway.NewEmulatorGatewayWithOpts(
		&gateway.EmulatorKey{
			PublicKey: emu.DefaultServiceKey().PublicKey,
			SigAlgo:   emu.DefaultServiceKeySigAlgo,
			HashAlgo:  emu.DefaultServiceKeyHashAlgo,
		},
		gateway.WithEmulatorOptions(
			emu.WithStore(memstore.New()),
			emu.WithTransactionValidationEnabled(false),
			emu.WithSimpleAddresses(),
			emu.WithStorageLimitEnabled(false),
			emu.WithTransactionFeesEnabled(false),
			emu.WithContractRemovalEnabled(true),
		),
	)

	return &flowKit{
		blockchain: kit.NewFlowkit(
			state,
			config.EmulatorNetwork,
			gw,
			output.NewStdoutLogger(output.NoneLog)),
	}, nil
}

func (fk *flowKit) executeTransaction(
	script string,
	arguments []string,
	authorizers []flow.Address,
) (*flow.Transaction, *flow.TransactionResult, error) {
	tx := &flow.Transaction{}
	tx.Script = []byte(script)

	args, err := parseArguments(arguments)
	if err != nil {
		return nil, nil, err
	}
	tx.Arguments = args

	return fk.sendTransaction(tx, authorizers)
}

func (fk *flowKit) executeScript(script string, arguments []string) (cadence.Value, error) {
	cadenceArgs := make([]cadence.Value, len(arguments))
	for i, arg := range arguments {
		val, err := cadence.NewValue(arg)
		if err != nil {
			return nil, err
		}
		cadenceArgs[i] = val
	}

	return fk.blockchain.ExecuteScript(
		context.Background(),
		kit.Script{
			Code:     []byte(script),
			Args:     cadenceArgs,
			Location: "",
		},
		kit.LatestScriptQuery)
}

func (fk *flowKit) createAccount() (*flow.Account, error) {
	state, err := fk.blockchain.State()
	if err != nil {
		return nil, err
	}

	service, err := state.EmulatorServiceAccount()
	if err != nil {
		return nil, err
	}
	serviceKey, err := service.Key.PrivateKey()
	if err != nil {
		return nil, err
	}

	account, _, err := fk.blockchain.CreateAccount(
		context.Background(),
		service,
		[]accounts.PublicKey{{
			Public:   (*serviceKey).PublicKey(),
			Weight:   flow.AccountKeyWeightThreshold,
			SigAlgo:  crypto.ECDSA_P256,
			HashAlgo: crypto.SHA3_256,
		}},
	)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (fk *flowKit) getAccount(address flow.Address) (*flow.Account, *emu.AccountStorage, error) {
	account, err := fk.blockchain.GetAccount(context.Background(), address)
	if err != nil {
		return nil, nil, err
	}
	// TODO: How run Cadence script to get account storage
	return account, nil, nil
}

func (fk *flowKit) deployContract(
	address flow.Address,
	script string,
) (*flow.Transaction, *flow.TransactionResult, error) {
	contractName, err := parseContractName(script)
	if err != nil {
		return nil, nil, err
	}

	tx := templates.AddAccountContract(address, templates.Contract{
		Name:   contractName,
		Source: script,
	})

	//fk.blockchain.AddContract(context.Background()) // TODO: fk.blockchain.AddContract ???
	return fk.sendTransaction(tx, nil)
}

func (fk *flowKit) removeContract(
	address flow.Address,
	contractName string,
) (*flow.Transaction, *flow.TransactionResult, error) {
	tx := templates.RemoveAccountContract(address, contractName)
	return fk.sendTransaction(tx, nil)
}

func (fk *flowKit) sendTransaction(
	tx *flow.Transaction,
	authorizers []flow.Address,
) (*flow.Transaction, *flow.TransactionResult, error) {
	state, err := fk.blockchain.State()
	if err != nil {
		return nil, nil, err
	}

	service, err := state.EmulatorServiceAccount()
	if err != nil {
		return nil, nil, err
	}

	var accountRoles transactions.AccountRoles
	accountRoles.Payer = *service
	accountRoles.Proposer = *service

	for _, auth := range authorizers {
		acc, _ := state.Accounts().ByAddress(auth)
		accountRoles.Authorizers = append(accountRoles.Authorizers, *acc)
	}

	args := make([]cadence.Value, len(tx.Arguments))
	for i := range tx.Arguments {
		arg, err := tx.Argument(i)
		if err != nil {
			return nil, nil, err
		}
		args[i] = arg
	}

	return fk.blockchain.SendTransaction(
		context.Background(),
		accountRoles,
		kit.Script{
			Code:     tx.Script,
			Args:     args,
			Location: "", // TODO: Do we need this?
		},
		tx.GasLimit,
	)
}

func (fk *flowKit) getLatestBlockHeight() (int, error) {
	block, err := fk.blockchain.Gateway().GetLatestBlock()
	if err != nil {
		return 0, err
	}
	return int(block.BlockHeader.Height), nil
}

// parseEventAddress gets an address out of the account creation events payloads
func parseEventAddress(events []flow.Event) flow.Address {
	for _, event := range events {
		if event.Type == flow.EventAccountCreated {
			addressValue := event.Value.Fields[0].(cadence.Address)
			return flow.HexToAddress(addressValue.Hex())
		}
	}
	return flow.EmptyAddress
}

// parseArguments converts string arguments list in cadence-JSON format into a byte serialised list
func parseArguments(args []string) ([][]byte, error) {
	encodedArgs := make([][]byte, len(args))
	for i, arg := range args {
		// decode and then encode again to ensure the value is valid
		val, err := jsoncdc.Decode(nil, []byte(arg))
		if err != nil {
			return nil, errors.Wrap(err, "failed to decode argument")
		}
		enc, _ := jsoncdc.Encode(val)
		encodedArgs[i] = enc
	}

	return encodedArgs, nil
}

// parseContractName extracts contract name from its source
func parseContractName(code string) (string, error) {
	program, err := parser.ParseProgram(nil, []byte(code), parser.Config{})
	if err != nil {
		return "", err
	}
	if len(program.CompositeDeclarations())+len(program.InterfaceDeclarations()) != 1 {
		return "", errors.New("the code must declare exactly one contract or contract interface")
	}

	for _, compositeDeclaration := range program.CompositeDeclarations() {
		if compositeDeclaration.CompositeKind == common.CompositeKindContract {
			return compositeDeclaration.Identifier.Identifier, nil
		}
	}

	for _, interfaceDeclaration := range program.InterfaceDeclarations() {
		if interfaceDeclaration.CompositeKind == common.CompositeKindContract {
			return interfaceDeclaration.Identifier.Identifier, nil
		}
	}

	return "", fmt.Errorf("unable to determine contract name")
}
