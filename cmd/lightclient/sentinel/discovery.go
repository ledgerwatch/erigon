/*
   Copyright 2022 Erigon-Lightclient contributors
   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at
       http://www.apache.org/licenses/LICENSE-2.0
   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package sentinel

import (
	"context"
	"time"

	"github.com/ledgerwatch/erigon/cmd/lightclient/clparams"
	"github.com/ledgerwatch/erigon/p2p/enode"
	"github.com/ledgerwatch/erigon/p2p/enr"
	"github.com/ledgerwatch/log/v3"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"github.com/pkg/errors"
	"github.com/prysmaticlabs/go-bitfield"
)

func (s *Sentinel) connectWithPeer(ctx context.Context, info peer.AddrInfo) error {
	if info.ID == s.host.ID() {
		return nil
	}
	if s.peers.IsBadPeer(info.ID) {
		return errors.New("refused to connect to bad peer")
	}
	ctxWithTimeout, cancel := context.WithTimeout(ctx, clparams.MaxDialTimeout)
	defer cancel()
	if err := s.host.Connect(ctxWithTimeout, info); err != nil {
		s.peers.Penalize(info.ID)
		return err
	}
	return nil
}

func (s *Sentinel) connectWithAllPeers(multiAddrs []multiaddr.Multiaddr) error {
	addrInfos, err := peer.AddrInfosFromP2pAddrs(multiAddrs...)
	if err != nil {
		return err
	}
	for _, peerInfo := range addrInfos {
		go func(peerInfo peer.AddrInfo) {
			if err := s.connectWithPeer(s.ctx, peerInfo); err != nil {
				log.Debug("Could not connect with peer", "err", err)
			}
		}(peerInfo)
	}
	return nil
}

func (s *Sentinel) listenForPeers() {
	iterator := s.listener.RandomNodes()
	defer iterator.Close()
	for {

		if s.ctx.Err() != nil {
			break
		}
		if s.HasTooManyPeers() {
			log.Trace("Not looking for peers, at peer limit")
			time.Sleep(100 * time.Millisecond)
			continue
		}
		exists := iterator.Next()
		if !exists {
			break
		}
		node := iterator.Node()
		peerInfo, _, err := convertToAddrInfo(node)
		if err != nil {
			log.Error("Could not convert to peer info", "err", err)
			continue
		}

		go func(peerInfo *peer.AddrInfo) {
			if err := s.connectWithPeer(s.ctx, *peerInfo); err != nil {
				log.Debug("Could not connect with peer", "err", err)
			}
		}(peerInfo)
	}
}

func (s *Sentinel) connectToBootnodes() error {
	for i := range s.cfg.DiscoverConfig.Bootnodes {
		if err := s.cfg.DiscoverConfig.Bootnodes[i].Record().Load(enr.WithEntry("tcp", new(enr.TCP))); err != nil {
			if !enr.IsNotFound(err) {
				log.Error("Could not retrieve tcp port")
			}
			continue
		}
	}
	multiAddresses := convertToMultiAddr(s.cfg.DiscoverConfig.Bootnodes)
	s.connectWithAllPeers(multiAddresses)
	return nil
}

func (s *Sentinel) setupENR(
	node *enode.LocalNode,
) (*enode.LocalNode, error) {
	// TODO(Giulio2002): Implement fork id.

	// Setup subnets key
	node.Set(enr.WithEntry(s.cfg.NetworkConfig.AttSubnetKey, bitfield.NewBitvector64().Bytes()))
	node.Set(enr.WithEntry(s.cfg.NetworkConfig.SyncCommsSubnetKey, bitfield.Bitvector4{byte(0x00)}.Bytes()))
	return node, nil
}
