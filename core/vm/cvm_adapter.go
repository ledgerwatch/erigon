package vm

import (
	"fmt"
	"github.com/holiman/uint256"
	"github.com/ledgerwatch/erigon-lib/gointerfaces/starknet"
	"github.com/ledgerwatch/erigon/common"
	"github.com/ledgerwatch/erigon/params"
)

const CairoNotImplemented = "the method is currently not implemented for cvm: %s"

func NewCVMAdapter(ibs IntraBlockState, starknetGrpcClient starknet.CAIROVMClient) *CVMAdapter {
	return &CVMAdapter{
		Cvm: NewCVM(ibs, starknetGrpcClient),
	}
}

type CVMAdapter struct {
	Cvm *CVM
}

func (c *CVMAdapter) Reset(txCtx TxContext, ibs IntraBlockState) {
	c.Cvm.intraBlockState = ibs
}

func (c *CVMAdapter) Create(caller ContractRef, code []byte, gas uint64, value *uint256.Int, salt []byte) (ret []byte, contractAddr common.Address, leftOverGas uint64, err error) {
	leftOverGas = 0

	ret, contractAddr, err = c.Cvm.Create(caller, code, salt)

	return ret, contractAddr, leftOverGas, err
}

func (cvm *CVMAdapter) Call(caller ContractRef, addr common.Address, input []byte, gas uint64, value *uint256.Int, bailout bool) (ret []byte, leftOverGas uint64, err error) {
	return nil, 0, fmt.Errorf(CairoNotImplemented, "Call")
}

func (cvm *CVMAdapter) Config() Config {
	return cvm.Cvm.Config()
}

func (cvm *CVMAdapter) ChainConfig() *params.ChainConfig {
	return params.FermionChainConfig
}

func (cvm *CVMAdapter) ChainRules() params.Rules {
	return params.Rules{}
}

func (cvm *CVMAdapter) Context() BlockContext {
	return BlockContext{}
}

func (cvm *CVMAdapter) IntraBlockState() IntraBlockState {
	return cvm.Cvm.IntraBlockState()
}

func (cvm *CVMAdapter) TxContext() TxContext {
	return TxContext{}
}
