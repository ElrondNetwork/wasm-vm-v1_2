package arwen

// // Declare the function signatures (see [cgo](https://golang.org/cmd/cgo/)).
//
// #include <stdlib.h>
// typedef unsigned char uint8_t;
// typedef int int32_t;
// typedef int uint32_t;
// typedef unsigned long long uint64_t;
//
// extern int32_t loadFunctionName(void *context, int32_t functionOffset);
// extern int32_t getNumArguments(void *context);
// extern void loadArgumentAsBigInt(void *context, int32_t id, int32_t destination);
// extern int32_t loadArgumentAsBytes(void *context, int32_t id, int32_t argOffset);
// extern long long getArgumentAsInt64(void *context, int32_t id);
//
// extern void loadOwner(void *context, int32_t resultOffset);
// extern void loadCaller(void *context, int32_t resultOffset);
// extern void loadCallValue(void *context, int32_t destination);
// extern void loadBalance(void *context, int32_t addressOffset, int32_t result);
// extern int32_t loadBlockHash(void *context, long long nonce, int32_t resultOffset);
// extern long long getBlockTimestamp(void *context);
//
// extern int32_t sendTransaction(void *context, long long gasLimit, int32_t dstOffset, int32_t valueRef, int32_t dataOffset, int32_t dataLength);
//
// extern int32_t storageStoreAsBytes(void *context, int32_t keyOffset, int32_t dataOffset, int32_t dataLength);
// extern int32_t storageLoadAsBytes(void *context, int32_t keyOffset, int32_t dataOffset);
// extern int32_t storageStoreAsBigInt(void *context, int32_t keyOffset, int32_t source);
// extern int32_t storageLoadAsBigInt(void *context, int32_t keyOffset, int32_t destination);
// extern int32_t storageStoreAsInt64(void *context, int32_t keyOffset, long long value);
// extern long long storageLoadAsInt64(void *context, int32_t keyOffset);
//
// extern void returnBigInt(void* context, int32_t reference);
// extern void returnInt32(void* context, int32_t value);
// extern void signalError(void* context);
// extern void writeLog(void *context, int32_t pointer, int32_t length, int32_t topicPtr, int32_t numTopics);
//
// extern int32_t bigIntNew(void* context, int32_t smallValue);
// extern int32_t bigIntByteLength(void* context, int32_t reference);
// extern int32_t bigIntGetBytes(void* context, int32_t reference, int32_t byteOffset);
// extern void bigIntSetBytes(void* context, int32_t destination, int32_t byteOffset, int32_t byteLength);
// extern int32_t bigIntIsInt64(void* context, int32_t reference);
// extern long long bigIntGetInt64(void* context, int32_t reference);
// extern void bigIntSetInt64(void* context, int32_t destination, long long value);
// extern void bigIntAdd(void* context, int32_t destination, int32_t op1, int32_t op2);
// extern void bigIntSub(void* context, int32_t destination, int32_t op1, int32_t op2);
// extern void bigIntMul(void* context, int32_t destination, int32_t op1, int32_t op2);
// extern int32_t bigIntCmp(void* context, int32_t op1, int32_t op2);
//
// extern void debugPrintBigInt(void* context, int32_t value);
// extern void debugPrintInt32(void* context, int32_t value);
// extern void debugPrintBytes(void* context, int32_t byteOffset, int32_t byteLength);
// extern void debugPrintString(void* context, int32_t byteOffset, int32_t byteLength);
import "C"

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"unsafe"

	vmcommon "github.com/ElrondNetwork/elrond-vm-common"
	"github.com/ElrondNetwork/go-ext-wasm/wasmer"
)

// BigIntHandle is the type we use to represent a reference to a big int in the host.
type BigIntHandle = int32

// HostContext abstracts away the blockchain functionality from wasmer.
type HostContext interface {
	Arguments() []*big.Int
	Function() string
	AccountExists(addr []byte) bool
	GetStorage(addr []byte, key []byte) []byte
	SetStorage(addr []byte, key []byte, value []byte) int32
	LoadBalance(addr []byte, destination BigIntHandle)
	GetCodeSize(addr []byte) int
	BlockHash(nonce int64) []byte
	GetCodeHash(addr []byte) []byte
	GetCode(addr []byte) []byte
	SelfDestruct(addr []byte, beneficiary []byte)
	GetVMInput() vmcommon.VMInput
	GetSCAddress() []byte
	WriteLog(addr []byte, topics [][]byte, data []byte)
	SendTransaction(destination []byte, value *big.Int, input []byte, gas int64) (gasLeft int64, err error)
	SignalUserError()

	BigInsertInt64(smallValue int64) BigIntHandle
	BigUpdate(destination BigIntHandle, newValue *big.Int)
	BigGet(reference BigIntHandle) *big.Int
	BigByteLength(reference BigIntHandle) int32
	BigGetBytes(reference BigIntHandle) []byte
	BigSetBytes(destination BigIntHandle, bytes []byte)
	BigIsInt64(destination BigIntHandle) bool
	BigGetInt64(destination BigIntHandle) int64
	BigSetInt64(destination BigIntHandle, value int64)
	BigAdd(destination, op1, op2 BigIntHandle)
	BigSub(destination, op1, op2 BigIntHandle)
	BigMul(destination, op1, op2 BigIntHandle)
	BigCmp(op1, op2 BigIntHandle) int

	ReturnBigInt(reference BigIntHandle)
	ReturnInt32(value int32)
	DebugPrintBig(value BigIntHandle)
}

func ElrondEImports() (*wasmer.Imports, error) {
	imports := wasmer.NewImports()

	imports, err := imports.Append("loadOwner", loadOwner, C.loadOwner)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("loadBalance", loadBalance, C.loadBalance)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("loadBlockHash", loadBlockHash, C.loadBlockHash)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("sendTransaction", sendTransaction, C.sendTransaction)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("loadArgumentAsBytes", loadArgumentAsBytes, C.loadArgumentAsBytes)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("loadArgumentAsBigInt", loadArgumentAsBigInt, C.loadArgumentAsBigInt)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getArgumentAsInt64", getArgumentAsInt64, C.getArgumentAsInt64)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("loadFunctionName", loadFunctionName, C.loadFunctionName)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getNumArguments", getNumArguments, C.getNumArguments)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("storageStoreAsBytes", storageStoreAsBytes, C.storageStoreAsBytes)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("storageLoadAsBytes", storageLoadAsBytes, C.storageLoadAsBytes)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("storageStoreAsBigInt", storageStoreAsBigInt, C.storageStoreAsBigInt)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("storageLoadAsBigInt", storageLoadAsBigInt, C.storageLoadAsBigInt)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("storageStoreAsInt64", storageStoreAsInt64, C.storageStoreAsInt64)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("storageLoadAsInt64", storageLoadAsInt64, C.storageLoadAsInt64)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("loadCaller", loadCaller, C.loadCaller)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("loadCallValue", loadCallValue, C.loadCallValue)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("writeLog", writeLog, C.writeLog)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("getBlockTimestamp", getBlockTimestamp, C.getBlockTimestamp)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("signalError", signalError, C.signalError)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntNew", bigIntNew, C.bigIntNew)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntByteLength", bigIntByteLength, C.bigIntByteLength)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntGetBytes", bigIntGetBytes, C.bigIntGetBytes)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntSetBytes", bigIntSetBytes, C.bigIntSetBytes)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntIsInt64", bigIntIsInt64, C.bigIntIsInt64)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntGetInt64", bigIntGetInt64, C.bigIntGetInt64)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntSetInt64", bigIntSetInt64, C.bigIntSetInt64)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntAdd", bigIntAdd, C.bigIntAdd)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntSub", bigIntSub, C.bigIntSub)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntMul", bigIntMul, C.bigIntMul)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("bigIntCmp", bigIntCmp, C.bigIntCmp)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("returnBigInt", returnBigInt, C.returnBigInt)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("returnInt32", returnInt32, C.returnInt32)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("debugPrintBigInt", debugPrintBigInt, C.debugPrintBigInt)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("debugPrintInt32", debugPrintInt32, C.debugPrintInt32)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("debugPrintBytes", debugPrintBytes, C.debugPrintBytes)
	if err != nil {
		return nil, err
	}

	imports, err = imports.Append("debugPrintString", debugPrintString, C.debugPrintString)
	if err != nil {
		return nil, err
	}

	return imports, nil
}

// Write the implementation of the functions, and export it (for cgo).

//export loadOwner
func loadOwner(context unsafe.Pointer, resultOffset int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	owner := hostContext.GetSCAddress()
	err := storeBytes(instCtx.Memory(), resultOffset, owner)
	if err != nil {
		fmt.Println("loadOwner error: " + err.Error())
	}
	fmt.Println("loadOwner " + hex.EncodeToString(owner))
}

//export signalError
func signalError(context unsafe.Pointer) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	hostContext.SignalUserError()
	fmt.Println("signalError called")
}

//export loadBalance
func loadBalance(context unsafe.Pointer, addressOffset int32, result int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	address := loadBytes(instCtx.Memory(), addressOffset, addressLen)
	hostContext.LoadBalance(address, result)
}

//export loadBlockHash
func loadBlockHash(context unsafe.Pointer, nonce int64, resultOffset int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	hash := hostContext.BlockHash(nonce)
	err := storeBytes(instCtx.Memory(), resultOffset, hash)
	if err != nil {
		fmt.Println("loadBlockHash error: " + err.Error())
		return 1
	}
	fmt.Println("loadBlockHash " + hex.EncodeToString(hash))
	return 0
}

//export sendTransaction
func sendTransaction(context unsafe.Pointer, gasLimit int64, destOffset int32, valueRef int32, dataOffset int32, dataLength int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	dest := loadBytes(instCtx.Memory(), destOffset, addressLen)
	value := hostContext.BigGet(valueRef)
	data := loadBytes(instCtx.Memory(), dataOffset, dataLength)

	fmt.Printf("sendTransaction to: %s value: %d data: %s\n",
		hex.EncodeToString(dest),
		value,
		data,
	)

	_, err := hostContext.SendTransaction(dest, value, data, gasLimit)
	if err != nil {
		fmt.Println("sendTransaction error: " + err.Error())
		return 1
	}

	fmt.Println("sendTransaction succeed")
	return 0
}

//export loadArgumentAsBytes
func loadArgumentAsBytes(context unsafe.Pointer, id int32, argOffset int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	args := hostContext.Arguments()
	if int32(len(args)) <= id {
		fmt.Println("getArgument id invalid")
		return -1
	}

	err := storeBytes(instCtx.Memory(), argOffset, args[id].Bytes())
	if err != nil {
		fmt.Println("getArgument error " + err.Error())
		return -1
	}

	fmt.Printf("argument #%d (bytes): %s\n", id, hex.EncodeToString(args[id].Bytes()))
	return int32(len(args[id].Bytes()))
}

//export loadArgumentAsBigInt
func loadArgumentAsBigInt(context unsafe.Pointer, id int32, destination int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	args := hostContext.Arguments()
	if int32(len(args)) <= id {
		fmt.Println("getArgument id invalid")
		return
	}

	hostContext.BigUpdate(destination, args[id])

	fmt.Printf("argument #%d (big int): %d\n", id, args[id])
}

//export getArgumentAsInt64
func getArgumentAsInt64(context unsafe.Pointer, id int32) int64 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	args := hostContext.Arguments()
	if int32(len(args)) <= id {
		fmt.Println("getArgument id invalid")
		return -1
	}

	fmt.Printf("argument #%d (int64): %d\n", id, args[id].Int64())
	return args[id].Int64()
}

//export loadFunctionName
func loadFunctionName(context unsafe.Pointer, functionOffset int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	function := hostContext.Function()
	err := storeBytes(instCtx.Memory(), functionOffset, []byte(function))
	if err != nil {
		fmt.Println("loadFunctionName error: ", err.Error())
		return -1
	}

	fmt.Println("loadFunctionName name: " + function)
	return int32(len(function))
}

//export getNumArguments
func getNumArguments(context unsafe.Pointer) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	fmt.Println("getNumArguments ", len(hostContext.Arguments()))
	return int32(len(hostContext.Arguments()))
}

//export storageStoreAsBytes
func storageStoreAsBytes(context unsafe.Pointer, keyOffset int32, dataOffset int32, dataLength int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	key := loadBytes(instCtx.Memory(), keyOffset, hashLen)
	data := loadBytes(instCtx.Memory(), dataOffset, dataLength)

	fmt.Printf("storageStoreAsBytes key: %s  value (bytes): %s\n", hex.EncodeToString(key), hex.EncodeToString(data))
	return hostContext.SetStorage(hostContext.GetSCAddress(), key, data)
}

//export storageLoadAsBytes
func storageLoadAsBytes(context unsafe.Pointer, keyOffset int32, dataOffset int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	key := loadBytes(instCtx.Memory(), keyOffset, hashLen)
	data := hostContext.GetStorage(hostContext.GetSCAddress(), key)

	err := storeBytes(instCtx.Memory(), dataOffset, data)
	if err != nil {
		fmt.Println("storageLoadAsBytes error: " + err.Error())
		return -1
	}

	fmt.Printf("storageLoadAsBytes key: %s  value (bytes): %s\n", hex.EncodeToString(key), hex.EncodeToString(data))
	return int32(len(data))
}

//export storageStoreAsBigInt
func storageStoreAsBigInt(context unsafe.Pointer, keyOffset int32, source int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	key := loadBytes(instCtx.Memory(), keyOffset, hashLen)
	bytes := hostContext.BigGetBytes(source)

	fmt.Printf("storageStoreAsBytes key: %s  value (big int): %s\n", hex.EncodeToString(key), hex.EncodeToString(bytes))
	return hostContext.SetStorage(hostContext.GetSCAddress(), key, bytes)
}

//export storageLoadAsBigInt
func storageLoadAsBigInt(context unsafe.Pointer, keyOffset int32, destination int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	key := loadBytes(instCtx.Memory(), keyOffset, hashLen)
	bytes := hostContext.GetStorage(hostContext.GetSCAddress(), key)

	hostContext.BigSetBytes(destination, bytes)

	fmt.Printf("storageLoadAsBytes key: %s  value (big int): %s\n", hex.EncodeToString(key), hex.EncodeToString(bytes))
	return int32(len(bytes))
}

//export storageStoreAsInt64
func storageStoreAsInt64(context unsafe.Pointer, keyOffset int32, value int64) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	key := loadBytes(instCtx.Memory(), keyOffset, hashLen)
	data := big.NewInt(value)

	fmt.Printf("storageStoreAsInt64 key: %s  value (int64): %x\n", hex.EncodeToString(key), data.Int64())
	return hostContext.SetStorage(hostContext.GetSCAddress(), key, data.Bytes())
}

//export storageLoadAsInt64
func storageLoadAsInt64(context unsafe.Pointer, keyOffset int32) int64 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	key := loadBytes(instCtx.Memory(), keyOffset, hashLen)
	data := hostContext.GetStorage(hostContext.GetSCAddress(), key)

	bigInt := big.NewInt(0).SetBytes(data)
	fmt.Printf("storageLoadAsInt64 key: %s  value (int64): %x\n", hex.EncodeToString(key), bigInt.Int64())

	return bigInt.Int64()
}

//export loadCaller
func loadCaller(context unsafe.Pointer, resultOffset int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	caller := hostContext.GetVMInput().CallerAddr

	err := storeBytes(instCtx.Memory(), resultOffset, caller)
	if err != nil {
		fmt.Println("loadCaller error: " + err.Error())
	}
	fmt.Println("loadCaller " + string(caller))
}

//export loadCallValue
func loadCallValue(context unsafe.Pointer, destination int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())
	hostContext.BigUpdate(destination, hostContext.GetVMInput().CallValue)
	fmt.Printf("loadCallValue %d\n", hostContext.GetVMInput().CallValue)
}

//export writeLog
func writeLog(context unsafe.Pointer, pointer int32, length int32, topicPtr int32, numTopics int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	log := loadBytes(instCtx.Memory(), pointer, length)

	topics := make([][]byte, numTopics)
	fmt.Println("writeLog: ")
	for i := int32(0); i < numTopics; i++ {
		topics[i] = loadBytes(instCtx.Memory(), topicPtr+i*hashLen, hashLen)
		fmt.Println("topics: " + string(topics[i]))
	}

	fmt.Print("log: " + string(log))
	hostContext.WriteLog(hostContext.GetSCAddress(), topics, log)
}

//export getBlockTimestamp
func getBlockTimestamp(context unsafe.Pointer) int64 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	fmt.Println("getBlockTimestamp ", hostContext.GetVMInput().Header.Timestamp.Int64())
	return hostContext.GetVMInput().Header.Timestamp.Int64()
}

//export bigIntNew
func bigIntNew(context unsafe.Pointer, smallValue int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())
	return hostContext.BigInsertInt64(int64(smallValue))
}

//export bigIntByteLength
func bigIntByteLength(context unsafe.Pointer, reference int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())
	return hostContext.BigByteLength(reference)
}

//export bigIntGetBytes
func bigIntGetBytes(context unsafe.Pointer, reference int32, byteOffset int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	bytes := hostContext.BigGetBytes(reference)

	err := storeBytes(instCtx.Memory(), byteOffset, bytes)
	if err != nil {
		fmt.Println("bigIntGetBytes error: " + err.Error())
	}

	return int32(len(bytes))
}

//export bigIntSetBytes
func bigIntSetBytes(context unsafe.Pointer, destination int32, byteOffset int32, byteLength int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())

	bytes := loadBytes(instCtx.Memory(), byteOffset, byteLength)
	hostContext.BigSetBytes(destination, bytes)
}

//export bigIntIsInt64
func bigIntIsInt64(context unsafe.Pointer, destination int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())
	if hostContext.BigIsInt64(destination) {
		return 1
	}
	return 0
}

//export bigIntGetInt64
func bigIntGetInt64(context unsafe.Pointer, destination int32) int64 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())
	return hostContext.BigGetInt64(destination)
}

//export bigIntSetInt64
func bigIntSetInt64(context unsafe.Pointer, destination int32, value int64) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())
	hostContext.BigSetInt64(destination, value)
}

//export bigIntAdd
func bigIntAdd(context unsafe.Pointer, destination, op1, op2 int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())
	hostContext.BigAdd(destination, op1, op2)
}

//export bigIntSub
func bigIntSub(context unsafe.Pointer, destination, op1, op2 int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())
	hostContext.BigSub(destination, op1, op2)
}

//export bigIntMul
func bigIntMul(context unsafe.Pointer, destination, op1, op2 int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())
	hostContext.BigMul(destination, op1, op2)
}

//export bigIntCmp
func bigIntCmp(context unsafe.Pointer, op1, op2 int32) int32 {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())
	return int32(hostContext.BigCmp(op1, op2))
}

//export returnBigInt
func returnBigInt(context unsafe.Pointer, reference int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())
	hostContext.ReturnBigInt(reference)
}

//export returnInt32
func returnInt32(context unsafe.Pointer, value int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())
	hostContext.ReturnInt32(value)
}

//export debugPrintBigInt
func debugPrintBigInt(context unsafe.Pointer, handle int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	hostContext := getHostContext(instCtx.Data())
	hostContext.DebugPrintBig(handle)
}

//export debugPrintInt32
func debugPrintInt32(context unsafe.Pointer, value int32) {
	fmt.Printf(">>> Int32: %d\n", value)
}

//export debugPrintBytes
func debugPrintBytes(context unsafe.Pointer, byteOffset int32, byteLength int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	bytes := loadBytes(instCtx.Memory(), byteOffset, byteLength)
	fmt.Printf(">>> Bytes: %s\n", hex.EncodeToString(bytes))
}

//export debugPrintString
func debugPrintString(context unsafe.Pointer, byteOffset int32, byteLength int32) {
	instCtx := wasmer.IntoInstanceContext(context)
	bytes := loadBytes(instCtx.Memory(), byteOffset, byteLength)
	fmt.Printf(">>> String: \"%s\"\n", string(bytes))
}
