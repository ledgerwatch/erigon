// Copyright 2024 The Erigon Authors
// This file is part of Erigon.
//
// Erigon is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Erigon is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with Erigon. If not, see <http://www.gnu.org/licenses/>.

package initial_state

import (
	_ "embed"

	"github.com/erigontech/erigon/cl/phase1/core/state"

	"github.com/erigontech/erigon/cl/clparams"
)

//go:embed mainnet.state.ssz
var mainnetStateSSZ []byte

//go:embed sepolia.state.ssz
var sepoliaStateSSZ []byte

//go:embed gnosis.state.ssz
var gnosisStateSSZ []byte

// Return genesis state
func GetGenesisState(network clparams.NetworkType) (*state.CachingBeaconState, error) {
	_, config := clparams.GetConfigsByNetwork(network)
	returnState := state.New(config)

	switch network {
	case clparams.MainnetNetwork:
		if err := returnState.DecodeSSZ(mainnetStateSSZ, int(clparams.Phase0Version)); err != nil {
			return nil, err
		}
	case clparams.SepoliaNetwork:
		if err := returnState.DecodeSSZ(sepoliaStateSSZ, int(clparams.Phase0Version)); err != nil {
			return nil, err
		}
	case clparams.GnosisNetwork:
		if err := returnState.DecodeSSZ(gnosisStateSSZ, int(clparams.Phase0Version)); err != nil {
			return nil, err
		}
	default:
		return nil, nil
	}
	return returnState, nil
}

func IsGenesisStateSupported(network clparams.NetworkType) bool {
	return network == clparams.MainnetNetwork || network == clparams.SepoliaNetwork || network == clparams.GnosisNetwork
}
