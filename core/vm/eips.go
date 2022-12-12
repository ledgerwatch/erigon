// Copyright 2019 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package vm

import (
	"fmt"
	"sort"

	"github.com/holiman/uint256"

	"github.com/ledgerwatch/erigon/params"
)

var activators = map[int]func(*JumpTable){
	3855: enable3855,
	3860: enable3860,
	3529: enable3529,
	3198: enable3198,
	2929: enable2929,
	2200: enable2200,
	1884: enable1884,
	1344: enable1344,
	3540: enable3540,
	3670: enable3670,
	4200: enable4200,
	4750: enable4750,
}

// EnableEIP enables the given EIP on the config.
// This operation writes in-place, and callers need to ensure that the globally
// defined jump tables are not polluted.
func EnableEIP(eipNum int, jt *JumpTable) error {
	enablerFn, ok := activators[eipNum]
	if !ok {
		return fmt.Errorf("undefined eip %d", eipNum)
	}
	enablerFn(jt)
	return nil
}

func ValidEip(eipNum int) bool {
	_, ok := activators[eipNum]
	return ok
}
func ActivateableEips() []string {
	var nums []string //nolint:prealloc
	for k := range activators {
		nums = append(nums, fmt.Sprintf("%d", k))
	}
	sort.Strings(nums)
	return nums
}

// enable1884 applies EIP-1884 to the given jump table:
// - Increase cost of BALANCE to 700
// - Increase cost of EXTCODEHASH to 700
// - Increase cost of SLOAD to 800
// - Define SELFBALANCE, with cost GasFastStep (5)
func enable1884(jt *JumpTable) {
	// Gas cost changes
	jt[SLOAD].constantGas = params.SloadGasEIP1884
	jt[BALANCE].constantGas = params.BalanceGasEIP1884
	jt[EXTCODEHASH].constantGas = params.ExtcodeHashGasEIP1884

	// New opcode
	jt[SELFBALANCE] = &operation{
		execute:     opSelfBalance,
		constantGas: GasFastStep,
		minStack:    minStack(0, 1),
		maxStack:    maxStack(0, 1),
	}
}

func opSelfBalance(pc *uint64, interpreter *EVMInterpreter, callContext *ScopeContext) ([]byte, error) {
	balance := interpreter.evm.IntraBlockState().GetBalance(callContext.Contract.Address())
	callContext.Stack.Push(balance)
	return nil, nil
}

// enable1344 applies EIP-1344 (ChainID Opcode)
// - Adds an opcode that returns the current chain’s EIP-155 unique identifier
func enable1344(jt *JumpTable) {
	// New opcode
	jt[CHAINID] = &operation{
		execute:     opChainID,
		constantGas: GasQuickStep,
		minStack:    minStack(0, 1),
		maxStack:    maxStack(0, 1),
	}
}

// opChainID implements CHAINID opcode
func opChainID(pc *uint64, interpreter *EVMInterpreter, callContext *ScopeContext) ([]byte, error) {
	chainId, _ := uint256.FromBig(interpreter.evm.ChainRules().ChainID)
	callContext.Stack.Push(chainId)
	return nil, nil
}

// enable2200 applies EIP-2200 (Rebalance net-metered SSTORE)
func enable2200(jt *JumpTable) {
	jt[SLOAD].constantGas = params.SloadGasEIP2200
	jt[SSTORE].dynamicGas = gasSStoreEIP2200
}

// enable2929 enables "EIP-2929: Gas cost increases for state access opcodes"
// https://eips.ethereum.org/EIPS/eip-2929
func enable2929(jt *JumpTable) {
	jt[SSTORE].dynamicGas = gasSStoreEIP2929

	jt[SLOAD].constantGas = 0
	jt[SLOAD].dynamicGas = gasSLoadEIP2929

	jt[EXTCODECOPY].constantGas = params.WarmStorageReadCostEIP2929
	jt[EXTCODECOPY].dynamicGas = gasExtCodeCopyEIP2929

	jt[EXTCODESIZE].constantGas = params.WarmStorageReadCostEIP2929
	jt[EXTCODESIZE].dynamicGas = gasEip2929AccountCheck

	jt[EXTCODEHASH].constantGas = params.WarmStorageReadCostEIP2929
	jt[EXTCODEHASH].dynamicGas = gasEip2929AccountCheck

	jt[BALANCE].constantGas = params.WarmStorageReadCostEIP2929
	jt[BALANCE].dynamicGas = gasEip2929AccountCheck

	jt[CALL].constantGas = params.WarmStorageReadCostEIP2929
	jt[CALL].dynamicGas = gasCallEIP2929

	jt[CALLCODE].constantGas = params.WarmStorageReadCostEIP2929
	jt[CALLCODE].dynamicGas = gasCallCodeEIP2929

	jt[STATICCALL].constantGas = params.WarmStorageReadCostEIP2929
	jt[STATICCALL].dynamicGas = gasStaticCallEIP2929

	jt[DELEGATECALL].constantGas = params.WarmStorageReadCostEIP2929
	jt[DELEGATECALL].dynamicGas = gasDelegateCallEIP2929

	// This was previously part of the dynamic cost, but we're using it as a constantGas
	// factor here
	jt[SELFDESTRUCT].constantGas = params.SelfdestructGasEIP150
	jt[SELFDESTRUCT].dynamicGas = gasSelfdestructEIP2929
}

func enable3529(jt *JumpTable) {
	jt[SSTORE].dynamicGas = gasSStoreEIP3529
	jt[SELFDESTRUCT].dynamicGas = gasSelfdestructEIP3529
}

// enable3198 applies EIP-3198 (BASEFEE Opcode)
// - Adds an opcode that returns the current block's base fee.
func enable3198(jt *JumpTable) {
	// New opcode
	jt[BASEFEE] = &operation{
		execute:     opBaseFee,
		constantGas: GasQuickStep,
		minStack:    minStack(0, 1),
		maxStack:    maxStack(0, 1),
	}
}

// opBaseFee implements BASEFEE opcode
func opBaseFee(pc *uint64, interpreter *EVMInterpreter, callContext *ScopeContext) ([]byte, error) {
	baseFee := interpreter.evm.Context().BaseFee
	callContext.Stack.Push(baseFee)
	return nil, nil
}

// enable3855 applies EIP-3855 (PUSH0 opcode)
func enable3855(jt *JumpTable) {
	// New opcode
	jt[PUSH0] = &operation{
		execute:     opPush0,
		constantGas: GasQuickStep,
		minStack:    minStack(0, 1),
		maxStack:    maxStack(0, 1),
	}
}

// opPush0 implements the PUSH0 opcode
func opPush0(pc *uint64, interpreter *EVMInterpreter, scope *ScopeContext) ([]byte, error) {
	scope.Stack.Push(new(uint256.Int))
	return nil, nil
}

// EIP-3860: Limit and meter initcode
// https://eips.ethereum.org/EIPS/eip-3860
func enable3860(jt *JumpTable) {
	jt[CREATE].dynamicGas = gasCreateEip3860
	jt[CREATE2].dynamicGas = gasCreate2Eip3860
}

func enable3540(jt *JumpTable) {
	// Do nothing.
}

func enable3670(jt *JumpTable) {
	// Do nothing.
}

// enable4200 applies EIP-4200 (RJUMP and RJUMPI opcodes)
func enable4200(jt *JumpTable) {
	jt[RJUMP] = &operation{
		execute:     opRjump,
		constantGas: GasFastStep,
		minStack:    minStack(0, 0),
		maxStack:    maxStack(0, 0),
		eof1Only:    true,
	}
	jt[RJUMPI] = &operation{
		execute:     opRjumpi,
		constantGas: GasFastishStep,
		minStack:    minStack(1, 1),
		maxStack:    maxStack(1, 1),
		eof1Only:    true,
	}
}

// opRjump implements the RJUMP opcode
func opRjump(pc *uint64, interpreter *EVMInterpreter, scope *ScopeContext) ([]byte, error) {
	var (
		code = scope.Contract.Container.CodeAt(scope.ActiveSection)
		idx  = *pc + 1
	)

	relativeOffset, err := parseArg(code[idx : idx+2])
	if err != nil {
		return nil, err
	}

	// Move PC past the RJUMP instruction and its immediate argument.
	*pc += 2 + 1

	// Calculate the new PC given the relative offset. Already validated,
	// so no need to verify casts.
	*pc = uint64(int64(*pc)+int64(relativeOffset)) - 1 // pc will also be increased by interpreter loop

	return nil, nil
}

// opRjumpi implements the RJUMPI opcode
func opRjumpi(pc *uint64, interpreter *EVMInterpreter, scope *ScopeContext) ([]byte, error) {
	condition := scope.Stack.Pop()
	if condition.IsZero() {
		// Not branching, just skip over immediate argument.
		*pc += 2
		return nil, nil
	}
	return opRjump(pc, interpreter, scope)
}

// enable4750 applies EIP-4750 (CALLF and RETF opcodes)
func enable4750(jt *JumpTable) {
	jt[JUMP].legacyOnly = true
	jt[JUMPI].legacyOnly = true
	jt[CALLF] = &operation{
		execute:     opCallf,
		constantGas: GasMidStep,
		minStack:    minStack(0, 0),
		maxStack:    maxStack(0, 0),
		eof1Only:    true,
	}
	jt[RETF] = &operation{
		execute:     opRetf,
		constantGas: GasMidStep,
		minStack:    minStack(0, 0),
		maxStack:    maxStack(0, 0),
		eof1Only:    true,
	}
}

// opCallf implements the CALLF opcode.
func opCallf(pc *uint64, interpreter *EVMInterpreter, scope *ScopeContext) ([]byte, error) {
	var (
		code = scope.Contract.Container.CodeAt(scope.ActiveSection)
		idx  = *pc + 1
	)
	section, err := parseArg(code[idx : idx+2])
	if err != nil {
		return nil, err
	}
	caller := scope.RetStack[len(scope.RetStack)-1]
	sig := scope.Contract.Container.types[int(section)]
	if scope.Stack.Len() < int(caller.StackHeight)+int(sig.Input) {
		return nil, fmt.Errorf("too few stack items")
	}
	if len(scope.RetStack) >= 1024 {
		return nil, fmt.Errorf("return stack too deep")
	}
	context := &SubroutineContext{
		Section:     scope.ActiveSection,
		StackHeight: caller.StackHeight + uint64(scope.Stack.Len()),
		Pc:          *pc,
	}
	scope.RetStack = append(scope.RetStack, context)

	*pc = 0
	*pc -= 1 // hacks xD
	scope.ActiveSection = uint64(section)
	scope.Stack.Floor = int(context.StackHeight)

	return nil, nil
}

// opRetf implements the RETF opcode.
func opRetf(pc *uint64, interpreter *EVMInterpreter, scope *ScopeContext) ([]byte, error) {
	sig := scope.Contract.Container.types[scope.ActiveSection]
	if scope.Stack.Len() < int(sig.Output) {
		return nil, fmt.Errorf("too few stack items")
	}

	last := len(scope.RetStack) - 1
	context := scope.RetStack[last]
	scope.RetStack = scope.RetStack[:last]

	*pc = context.Pc - 1
	scope.ActiveSection = context.Section
	scope.Stack.Floor = int(context.StackHeight)

	return nil, nil
}
