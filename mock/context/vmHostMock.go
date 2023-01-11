package mock

import (
	"math/big"

	"github.com/multiversx/mx-chain-core-go/data/vm"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/multiversx/mx-chain-vm-v1_2-go/vmhost"
	"github.com/multiversx/mx-chain-vm-v1_2-go/crypto"
	"github.com/multiversx/mx-chain-vm-v1_2-go/wasmer"
)

var _ arwen.VMHost = (*VMHostMock)(nil)

// VMHostMock is used in tests to check the VMHost interface method calls
type VMHostMock struct {
	BlockChainHook vmcommon.BlockchainHook
	CryptoHook     crypto.VMCrypto

	EthInput []byte

	BlockchainContext arwen.BlockchainContext
	RuntimeContext    arwen.RuntimeContext
	OutputContext     arwen.OutputContext
	MeteringContext   arwen.MeteringContext
	StorageContext    arwen.StorageContext
	BigIntContext     arwen.BigIntContext

	SCAPIMethods  *wasmer.Imports
	IsBuiltinFunc bool
}

// Crypto mocked method
func (host *VMHostMock) Crypto() crypto.VMCrypto {
	return host.CryptoHook
}

// Blockchain mocked method
func (host *VMHostMock) Blockchain() arwen.BlockchainContext {
	return host.BlockchainContext
}

// Runtime mocked method
func (host *VMHostMock) Runtime() arwen.RuntimeContext {
	return host.RuntimeContext
}

// Output mocked method
func (host *VMHostMock) Output() arwen.OutputContext {
	return host.OutputContext
}

// Metering mocked method
func (host *VMHostMock) Metering() arwen.MeteringContext {
	return host.MeteringContext
}

// Storage mocked method
func (host *VMHostMock) Storage() arwen.StorageContext {
	return host.StorageContext
}

// BigInt mocked method
func (host *VMHostMock) BigInt() arwen.BigIntContext {
	return host.BigIntContext
}

// IsArwenV2Enabled mocked method
func (host *VMHostMock) IsArwenV2Enabled() bool {
	return true
}

// IsArwenV3Enabled mocked method
func (host *VMHostMock) IsArwenV3Enabled() bool {
	return true
}

// IsAheadOfTimeCompileEnabled mocked method
func (host *VMHostMock) IsAheadOfTimeCompileEnabled() bool {
	return true
}

// IsDynamicGasLockingEnabled mocked method
func (host *VMHostMock) IsDynamicGasLockingEnabled() bool {
	return true
}

// IsESDTFunctionsEnabled mocked method
func (host *VMHostMock) IsESDTFunctionsEnabled() bool {
	return true
}

// AreInSameShard mocked method
func (host *VMHostMock) AreInSameShard(_ []byte, _ []byte) bool {
	return true
}

// RevertESDTTransfer mocked method
func (host *VMHostMock) RevertESDTTransfer(_ *vmcommon.ContractCallInput) {
}

// ExecuteESDTTransfer mocked method
func (host *VMHostMock) ExecuteESDTTransfer(_ []byte, _ []byte, _ []byte, _ uint64, _ *big.Int, _ vm.CallType, _ bool) (*vmcommon.VMOutput, uint64, error) {
	return nil, 0, nil
}

// CreateNewContract mocked method
func (host *VMHostMock) CreateNewContract(_ *vmcommon.ContractCreateInput) ([]byte, error) {
	return nil, nil
}

// ExecuteOnSameContext mocked method
func (host *VMHostMock) ExecuteOnSameContext(_ *vmcommon.ContractCallInput) (*arwen.AsyncContextInfo, error) {
	return nil, nil
}

// ExecuteOnDestContext mocked method
func (host *VMHostMock) ExecuteOnDestContext(_ *vmcommon.ContractCallInput) (*vmcommon.VMOutput, *arwen.AsyncContextInfo, uint64, error) {
	return nil, nil, 0, nil
}

// InitState mocked method
func (host *VMHostMock) InitState() {
}

// PushState mocked method
func (host *VMHostMock) PushState() {
}

// PopState mocked method
func (host *VMHostMock) PopState() {
}

// ClearStateStack mocked method
func (host *VMHostMock) ClearStateStack() {
}

// GetAPIMethods mocked method
func (host *VMHostMock) GetAPIMethods() *wasmer.Imports {
	return host.SCAPIMethods
}

// GetProtocolBuiltinFunctions mocked method
func (host *VMHostMock) GetProtocolBuiltinFunctions() vmcommon.FunctionNames {
	return make(vmcommon.FunctionNames)
}

// IsBuiltinFunctionName mocked method
func (host *VMHostMock) IsBuiltinFunctionName(_ string) bool {
	return host.IsBuiltinFunc
}
