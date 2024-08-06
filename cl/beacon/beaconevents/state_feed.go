package beaconevents

import (
	ethevent "github.com/erigontech/erigon/event"
)

type stateFeed struct {
	feed *ethevent.Feed
}

func newStateFeed() *stateFeed {
	return &stateFeed{
		feed: &ethevent.Feed{},
	}
}

func (f *stateFeed) Subscribe(channel chan *EventStream) ethevent.Subscription {
	return f.feed.Subscribe(channel)
}

func (f *stateFeed) SendHead(value *HeadData) int {
	return f.feed.Send(&EventStream{
		Event: StateHead,
		Data:  value,
	})
}

// The node has received a block (from P2P or API) that is successfully imported on the fork-choice on_block handler
func (f *stateFeed) SendBlock(value *BlockData) int {
	return f.feed.Send(&EventStream{
		Event: StateBlock,
		Data:  value,
	})
}

// The node has received a block (from P2P or API) that passes validation rules of the beacon_block topic
func (f *stateFeed) SendBlockGossip(value *BlockGossipData) int {
	return f.feed.Send(&EventStream{
		Event: StateBlockGossip,
		Data:  value,
	})
}

func (f *stateFeed) SendFinalizedCheckpoint(value *FinalizedCheckpointData) int {
	return f.feed.Send(&EventStream{
		Event: StateFinalizedCheckpoint,
		Data:  value,
	})
}
