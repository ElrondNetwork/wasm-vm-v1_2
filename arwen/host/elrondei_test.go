package host

import (
	"math/big"
	"testing"

	"github.com/ElrondNetwork/arwen-wasm-vm/arwen"
	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/stretchr/testify/assert"
)

func TestElrondEI_CallValue(t *testing.T) {
	code := GetTestSCCode("elrondei", "../../")

	// 1-byte call value
	host, _ := DefaultTestArwenForCall(t, code)
	input := DefaultTestContractCallInput()
	input.GasProvided = 100000
	input.Function = "test_getCallValue_1byte"
	input.CallValue = big.NewInt(64)

	vmOutput, err := host.RunSmartContractCall(input)
	assert.Nil(t, err)
	assert.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
	assert.Len(t, vmOutput.ReturnData, 3)
	data := vmOutput.ReturnData
	assert.Equal(t, []byte("ok"), data[0])
	assert.Equal(t, []byte{1, 0, 0, 0}, data[1])
	assert.Equal(t, []byte{64}, data[2])

	// 4-byte call value
	host, _ = DefaultTestArwenForCall(t, code)
	input = DefaultTestContractCallInput()
	input.GasProvided = 100000
	input.Function = "test_getCallValue_4bytes"
	input.CallValue = big.NewInt(0).SetBytes([]byte{64, 12, 16, 99})

	vmOutput, err = host.RunSmartContractCall(input)
	assert.Nil(t, err)
	assert.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
	assert.Len(t, vmOutput.ReturnData, 3)
	data = vmOutput.ReturnData
	assert.Equal(t, []byte("ok"), data[0])
	assert.Equal(t, []byte{4, 0, 0, 0}, data[1])
	assert.Equal(t, []byte{64, 12, 16, 99}, data[2])

	// BigInt call value
	host, _ = DefaultTestArwenForCall(t, code)
	input = DefaultTestContractCallInput()
	input.GasProvided = 100000
	input.Function = "test_getCallValue_bigInt_to_Bytes"
	input.CallValue = big.NewInt(19*256 + 233)

	vmOutput, err = host.RunSmartContractCall(input)
	assert.Nil(t, err)
	assert.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
	assert.Len(t, vmOutput.ReturnData, 4)
	data = vmOutput.ReturnData
	assert.Equal(t, []byte("ok"), data[0])
	assert.Equal(t, []byte{2, 0, 0, 0}, data[1])
	assert.Equal(t, []byte{19, 233}, data[2])

	val12345 := big.NewInt(0).SetBytes(data[3])
	assert.Equal(t, big.NewInt(12345), val12345)
}

func TestElrondEI_int64getArgument(t *testing.T) {
	code := GetTestSCCode("elrondei", "../../")
	host, _ := DefaultTestArwenForCall(t, code)
	input := DefaultTestContractCallInput()
	input.GasProvided = 100000
	input.Function = "test_int64getArgument"
	input.Arguments = [][]byte{big.NewInt(12345).Bytes()}

	vmOutput, err := host.RunSmartContractCall(input)
	assert.Nil(t, err)
	assert.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
	assert.Len(t, vmOutput.ReturnData, 3)
	data := vmOutput.ReturnData
	assert.Equal(t, []byte("ok"), data[0])
	assert.Equal(t, []byte{57, 48, 0, 0}, data[1])

	invBytes := arwen.InverseBytes(data[1])
	val12345 := big.NewInt(0).SetBytes(invBytes)
	assert.Equal(t, big.NewInt(12345), val12345)

	i64val12345 := big.NewInt(0).SetBytes(data[2])
	assert.Equal(t, big.NewInt(12345), i64val12345)

	// Take the result of the SC method (the number 12345 as bytes, received from
	// the SC in data[2]) and feed it back into the SC method.
	input.Arguments = [][]byte{data[2]}

	vmOutput, err = host.RunSmartContractCall(input)
	assert.Nil(t, err)
	assert.Equal(t, vmcommon.Ok, vmOutput.ReturnCode)
	assert.Len(t, vmOutput.ReturnData, 3)
	data = vmOutput.ReturnData
	assert.Equal(t, []byte("ok"), data[0])
	assert.Equal(t, []byte{57, 48, 0, 0}, data[1])

	invBytes = arwen.InverseBytes(data[1])
	val12345 = big.NewInt(0).SetBytes(invBytes)
	assert.Equal(t, big.NewInt(12345), val12345)

	i64val12345 = big.NewInt(0).SetBytes(data[2])
	assert.Equal(t, big.NewInt(12345), i64val12345)
}
