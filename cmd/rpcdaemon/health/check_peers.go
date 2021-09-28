package health

import (
	"context"
	"fmt"
)

func checkMinPeers(minPeerCount uint, api NetAPI) error {
	if api == nil {
		return fmt.Errorf("no connection to the Erigon server")
	}

	peerCount, err := api.NetPeerCount(context.TODO())
	if err != nil {
		return err
	}

	if peerCount < uint64(minPeerCount) {
		return fmt.Errorf("not enough peers: %d (minimum %d))", peerCount, minPeerCount)
	}

	return nil
}
