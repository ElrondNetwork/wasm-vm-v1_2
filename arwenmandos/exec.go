package arwenmandos

import (
	"fmt"
	"path/filepath"

	"github.com/ElrondNetwork/wasm-vm-v1_2/arwen"
	arwenHost "github.com/ElrondNetwork/wasm-vm-v1_2/arwen/host"
	"github.com/ElrondNetwork/wasm-vm-v1_2/config"
	mc "github.com/ElrondNetwork/wasm-vm-v1_2/mandos-go/controller"
	er "github.com/ElrondNetwork/wasm-vm-v1_2/mandos-go/expression/reconstructor"
	fr "github.com/ElrondNetwork/wasm-vm-v1_2/mandos-go/fileresolver"
	mj "github.com/ElrondNetwork/wasm-vm-v1_2/mandos-go/json/model"
	worldhook "github.com/ElrondNetwork/wasm-vm-v1_2/mock/world"
	logger "github.com/ElrondNetwork/elrond-go-logger"
	vmi "github.com/ElrondNetwork/elrond-vm-common"
)

var log = logger.GetOrCreate("arwen/mandos")

// TestVMType is the VM type argument we use in tests.
var TestVMType = []byte{0, 0}

// ArwenTestExecutor parses, interprets and executes both .test.json tests and .scen.json scenarios with Arwen.
type ArwenTestExecutor struct {
	World                   *worldhook.MockWorld
	vm                      vmi.VMExecutionHandler
	checkGas                bool
	arwenmandosPath         string
	mandosGasScheduleLoaded bool
	fileResolver            fr.FileResolver
	exprReconstructor       er.ExprReconstructor
}

var _ mc.TestExecutor = (*ArwenTestExecutor)(nil)
var _ mc.ScenarioExecutor = (*ArwenTestExecutor)(nil)

// NewArwenTestExecutor prepares a new ArwenTestExecutor instance.
func NewArwenTestExecutor(arwenmandosPath string) (*ArwenTestExecutor, error) {
	world := worldhook.NewMockWorld()

	gasScheduleMap := config.MakeGasMapForTests()
	err := world.InitBuiltinFunctions(gasScheduleMap)
	if err != nil {
		return nil, err
	}

	blockGasLimit := uint64(10000000)
	vm, err := arwenHost.NewArwenVM(world, &arwen.VMHostParameters{
		VMType:                   TestVMType,
		BlockGasLimit:            blockGasLimit,
		GasSchedule:              gasScheduleMap,
		ProtocolBuiltinFunctions: world.GetBuiltinFunctionNames(),
		ElrondProtectedKeyPrefix: []byte(ElrondProtectedKeyPrefix),
	})
	if err != nil {
		return nil, err
	}

	return &ArwenTestExecutor{
		World:                   world,
		vm:                      vm,
		checkGas:                true,
		arwenmandosPath:         arwenmandosPath,
		mandosGasScheduleLoaded: false,
		fileResolver:            nil,
		exprReconstructor:       er.ExprReconstructor{},
	}, nil
}

// GetVM yields a reference to the VMExecutionHandler used.
func (ae *ArwenTestExecutor) GetVM() vmi.VMExecutionHandler {
	return ae.vm
}

func (ae *ArwenTestExecutor) gasScheduleMapFromMandos(mandosGasSchedule mj.GasSchedule) (config.GasScheduleMap, error) {
	switch mandosGasSchedule {
	case mj.GasScheduleDefault:
		return arwenHost.LoadGasScheduleConfig(filepath.Join(ae.arwenmandosPath, "gasSchedules/gasScheduleV3.toml"))
	case mj.GasScheduleDummy:
		return config.MakeGasMapForTests(), nil
	case mj.GasScheduleV1:
		return arwenHost.LoadGasScheduleConfig(filepath.Join(ae.arwenmandosPath, "gasSchedules/gasScheduleV1.toml"))
	case mj.GasScheduleV2:
		return arwenHost.LoadGasScheduleConfig(filepath.Join(ae.arwenmandosPath, "gasSchedules/gasScheduleV2.toml"))
	case mj.GasScheduleV3:
		return arwenHost.LoadGasScheduleConfig(filepath.Join(ae.arwenmandosPath, "gasSchedules/gasScheduleV3.toml"))
	default:
		return nil, fmt.Errorf("unknown mandos GasSchedule: %d", mandosGasSchedule)
	}
}

// SetMandosGasSchedule updates the gas costs based on the mandos scenario config
// only changes the gas schedule once,
// this prevents subsequent gasSchedule declarations in externalSteps to overwrite
func (ae *ArwenTestExecutor) SetMandosGasSchedule(newGasSchedule mj.GasSchedule) error {
	if ae.mandosGasScheduleLoaded {
		return nil
	}
	ae.mandosGasScheduleLoaded = true
	gasSchedule, err := ae.gasScheduleMapFromMandos(newGasSchedule)
	if err != nil {
		return err
	}
	ae.vm.GasScheduleChange(gasSchedule)
	return nil
}
