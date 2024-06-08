package rpchelper

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"reflect"
	"sync"
	"sync/atomic"
	"time"

	libcommon "github.com/ledgerwatch/erigon-lib/common"
	"github.com/ledgerwatch/erigon-lib/common/concurrent"
	"github.com/ledgerwatch/erigon-lib/gointerfaces"
	"github.com/ledgerwatch/erigon-lib/gointerfaces/grpcutil"
	remote "github.com/ledgerwatch/erigon-lib/gointerfaces/remoteproto"
	txpool "github.com/ledgerwatch/erigon-lib/gointerfaces/txpoolproto"
	txpool2 "github.com/ledgerwatch/erigon-lib/txpool"
	"github.com/ledgerwatch/log/v3"
	"google.golang.org/grpc"

	"github.com/ledgerwatch/erigon/core/types"
	"github.com/ledgerwatch/erigon/eth/filters"
	"github.com/ledgerwatch/erigon/rlp"
)

// Filters holds the state for managing subscriptions to various Ethereum events.
// It allows for the subscription and management of events such as new blocks, pending transactions,
// logs, and other Ethereum-related activities.
type Filters struct {
	mu sync.RWMutex

	pendingBlock *types.Block

	headsSubs        *concurrent.SyncMap[HeadsSubID, Sub[*types.Header]]
	pendingLogsSubs  *concurrent.SyncMap[PendingLogsSubID, Sub[types.Logs]]
	pendingBlockSubs *concurrent.SyncMap[PendingBlockSubID, Sub[*types.Block]]
	pendingTxsSubs   *concurrent.SyncMap[PendingTxsSubID, Sub[[]types.Transaction]]
	logsSubs         *LogsFilterAggregator
	logsRequestor    atomic.Value
	onNewSnapshot    func()

	storeMu            sync.Mutex
	logsStores         *concurrent.SyncMap[LogsSubID, []*types.Log]
	pendingHeadsStores *concurrent.SyncMap[HeadsSubID, []*types.Header]
	pendingTxsStores   *concurrent.SyncMap[PendingTxsSubID, [][]types.Transaction]
	logger             log.Logger
}

// New creates a new Filters instance, initializes it, and starts subscription goroutines for Ethereum events.
// It requires a context, Ethereum backend, transaction pool client, mining client, snapshot callback function,
// and a logger for logging events.
func New(ctx context.Context, ethBackend ApiBackend, txPool txpool.TxpoolClient, mining txpool.MiningClient, onNewSnapshot func(), logger log.Logger) *Filters {
	logger.Info("rpc filters: subscribing to Erigon events")

	ff := &Filters{
		headsSubs:          concurrent.NewSyncMap[HeadsSubID, Sub[*types.Header]](),
		pendingTxsSubs:     concurrent.NewSyncMap[PendingTxsSubID, Sub[[]types.Transaction]](),
		pendingLogsSubs:    concurrent.NewSyncMap[PendingLogsSubID, Sub[types.Logs]](),
		pendingBlockSubs:   concurrent.NewSyncMap[PendingBlockSubID, Sub[*types.Block]](),
		logsSubs:           NewLogsFilterAggregator(),
		onNewSnapshot:      onNewSnapshot,
		logsStores:         concurrent.NewSyncMap[LogsSubID, []*types.Log](),
		pendingHeadsStores: concurrent.NewSyncMap[HeadsSubID, []*types.Header](),
		pendingTxsStores:   concurrent.NewSyncMap[PendingTxsSubID, [][]types.Transaction](),
		logger:             logger,
	}

	go func() {
		if ethBackend == nil {
			return
		}
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			if err := ethBackend.Subscribe(ctx, ff.OnNewEvent); err != nil {
				select {
				case <-ctx.Done():
					return
				default:
				}
				if grpcutil.IsEndOfStream(err) || grpcutil.IsRetryLater(err) {
					time.Sleep(3 * time.Second)
					continue
				}
				logger.Warn("rpc filters: error subscribing to events", "err", err)
			}
		}
	}()

	go func() {
		if ethBackend == nil {
			return
		}
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			if err := ethBackend.SubscribeLogs(ctx, ff.OnNewLogs, &ff.logsRequestor); err != nil {
				select {
				case <-ctx.Done():
					return
				default:
				}
				if grpcutil.IsEndOfStream(err) || grpcutil.IsRetryLater(err) {
					time.Sleep(3 * time.Second)
					continue
				}
				logger.Warn("rpc filters: error subscribing to logs", "err", err)
			}
		}
	}()

	if txPool != nil {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				default:
				}
				if err := ff.subscribeToPendingTransactions(ctx, txPool); err != nil {
					select {
					case <-ctx.Done():
						return
					default:
					}
					if grpcutil.IsEndOfStream(err) || grpcutil.IsRetryLater(err) || grpcutil.ErrIs(err, txpool2.ErrPoolDisabled) {
						time.Sleep(3 * time.Second)
						continue
					}
					logger.Warn("rpc filters: error subscribing to pending transactions", "err", err)
				}
			}
		}()

		if !reflect.ValueOf(mining).IsNil() { //https://groups.google.com/g/golang-nuts/c/wnH302gBa4I
			go func() {
				for {
					select {
					case <-ctx.Done():
						return
					default:
					}
					if err := ff.subscribeToPendingBlocks(ctx, mining); err != nil {
						select {
						case <-ctx.Done():
							return
						default:
						}
						if grpcutil.IsEndOfStream(err) || grpcutil.IsRetryLater(err) {
							time.Sleep(3 * time.Second)
							continue
						}
						logger.Warn("rpc filters: error subscribing to pending blocks", "err", err)
					}
				}
			}()
			go func() {
				for {
					select {
					case <-ctx.Done():
						return
					default:
					}
					if err := ff.subscribeToPendingLogs(ctx, mining); err != nil {
						select {
						case <-ctx.Done():
							return
						default:
						}
						if grpcutil.IsEndOfStream(err) || grpcutil.IsRetryLater(err) {
							time.Sleep(3 * time.Second)
							continue
						}
						logger.Warn("rpc filters: error subscribing to pending logs", "err", err)
					}
				}
			}()
		}
	}

	return ff
}

// LastPendingBlock returns the last pending block that was received.
func (ff *Filters) LastPendingBlock() *types.Block {
	ff.mu.RLock()
	defer ff.mu.RUnlock()
	return ff.pendingBlock
}

// subscribeToPendingTransactions subscribes to pending transactions using the given transaction pool client.
// It listens for new transactions and processes them as they arrive.
func (ff *Filters) subscribeToPendingTransactions(ctx context.Context, txPool txpool.TxpoolClient) error {
	subscription, err := txPool.OnAdd(ctx, &txpool.OnAddRequest{}, grpc.WaitForReady(true))
	if err != nil {
		return err
	}
	for {
		event, err := subscription.Recv()
		if errors.Is(err, io.EOF) {
			ff.logger.Debug("rpcdaemon: the subscription to pending transactions channel was closed")
			break
		}
		if err != nil {
			return err
		}

		ff.OnNewTx(event)
	}
	return nil
}

// subscribeToPendingBlocks subscribes to pending blocks using the given mining client.
// It listens for new pending blocks and processes them as they arrive.
func (ff *Filters) subscribeToPendingBlocks(ctx context.Context, mining txpool.MiningClient) error {
	subscription, err := mining.OnPendingBlock(ctx, &txpool.OnPendingBlockRequest{}, grpc.WaitForReady(true))
	if err != nil {
		return err
	}
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		event, err := subscription.Recv()
		if errors.Is(err, io.EOF) {
			ff.logger.Debug("rpcdaemon: the subscription to pending blocks channel was closed")
			break
		}
		if err != nil {
			return err
		}

		ff.HandlePendingBlock(event)
	}
	return nil
}

// HandlePendingBlock handles a new pending block received from the mining client.
// It updates the internal state and notifies subscribers about the new block.
func (ff *Filters) HandlePendingBlock(reply *txpool.OnPendingBlockReply) {
	b := &types.Block{}
	if reply == nil || len(reply.RplBlock) == 0 {
		return
	}
	if err := rlp.Decode(bytes.NewReader(reply.RplBlock), b); err != nil {
		ff.logger.Warn("OnNewPendingBlock rpc filters, unprocessable payload", "err", err)
	}

	ff.mu.Lock()
	defer ff.mu.Unlock()
	ff.pendingBlock = b

	ff.pendingBlockSubs.Range(func(k PendingBlockSubID, v Sub[*types.Block]) error {
		v.Send(b)
		return nil
	})
}

// subscribeToPendingLogs subscribes to pending logs using the given mining client.
// It listens for new pending logs and processes them as they arrive.
func (ff *Filters) subscribeToPendingLogs(ctx context.Context, mining txpool.MiningClient) error {
	subscription, err := mining.OnPendingLogs(ctx, &txpool.OnPendingLogsRequest{}, grpc.WaitForReady(true))
	if err != nil {
		return err
	}
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}
		event, err := subscription.Recv()
		if errors.Is(err, io.EOF) {
			ff.logger.Debug("rpcdaemon: the subscription to pending logs channel was closed")
			break
		}
		if err != nil {
			return err
		}

		ff.HandlePendingLogs(event)
	}
	return nil
}

// HandlePendingLogs handles new pending logs received from the mining client.
// It updates the internal state and notifies subscribers about the new logs.
func (ff *Filters) HandlePendingLogs(reply *txpool.OnPendingLogsReply) {
	if len(reply.RplLogs) == 0 {
		return
	}
	l := []*types.Log{}
	if err := rlp.Decode(bytes.NewReader(reply.RplLogs), &l); err != nil {
		ff.logger.Warn("OnNewPendingLogs rpc filters, unprocessable payload", "err", err)
	}
	ff.pendingLogsSubs.Range(func(k PendingLogsSubID, v Sub[types.Logs]) error {
		v.Send(l)
		return nil
	})
}

// SubscribeNewHeads subscribes to new block headers and returns a channel to receive the headers
// and a subscription ID to manage the subscription.
func (ff *Filters) SubscribeNewHeads(size int) (<-chan *types.Header, HeadsSubID) {
	id := HeadsSubID(generateSubscriptionID())
	sub := newChanSub[*types.Header](size)
	ff.headsSubs.Put(id, sub)
	return sub.ch, id
}

// UnsubscribeHeads unsubscribes from new block headers using the given subscription ID.
// It returns true if the unsubscription was successful, otherwise false.
func (ff *Filters) UnsubscribeHeads(id HeadsSubID) bool {
	ch, ok := ff.headsSubs.Get(id)
	if !ok {
		return false
	}
	ch.Close()
	if _, ok = ff.headsSubs.Delete(id); !ok {
		return false
	}
	ff.pendingHeadsStores.Delete(id)
	return true
}

// SubscribePendingLogs subscribes to pending logs and returns a channel to receive the logs
// and a subscription ID to manage the subscription. It uses the specified filter criteria.
func (ff *Filters) SubscribePendingLogs(size int) (<-chan types.Logs, PendingLogsSubID) {
	id := PendingLogsSubID(generateSubscriptionID())
	sub := newChanSub[types.Logs](size)
	ff.pendingLogsSubs.Put(id, sub)
	return sub.ch, id
}

// UnsubscribePendingLogs unsubscribes from pending logs using the given subscription ID.
func (ff *Filters) UnsubscribePendingLogs(id PendingLogsSubID) {
	ch, ok := ff.pendingLogsSubs.Get(id)
	if !ok {
		return
	}
	ch.Close()
	ff.pendingLogsSubs.Delete(id)
}

// SubscribePendingBlock subscribes to pending blocks and returns a channel to receive the blocks
// and a subscription ID to manage the subscription.
func (ff *Filters) SubscribePendingBlock(size int) (<-chan *types.Block, PendingBlockSubID) {
	id := PendingBlockSubID(generateSubscriptionID())
	sub := newChanSub[*types.Block](size)
	ff.pendingBlockSubs.Put(id, sub)
	return sub.ch, id
}

// UnsubscribePendingBlock unsubscribes from pending blocks using the given subscription ID.
func (ff *Filters) UnsubscribePendingBlock(id PendingBlockSubID) {
	ch, ok := ff.pendingBlockSubs.Get(id)
	if !ok {
		return
	}
	ch.Close()
	ff.pendingBlockSubs.Delete(id)
}

// SubscribePendingTxs subscribes to pending transactions and returns a channel to receive the transactions
// and a subscription ID to manage the subscription.
func (ff *Filters) SubscribePendingTxs(size int) (<-chan []types.Transaction, PendingTxsSubID) {
	id := PendingTxsSubID(generateSubscriptionID())
	sub := newChanSub[[]types.Transaction](size)
	ff.pendingTxsSubs.Put(id, sub)
	return sub.ch, id
}

// UnsubscribePendingTxs unsubscribes from pending transactions using the given subscription ID.
// It returns true if the unsubscription was successful, otherwise false.
func (ff *Filters) UnsubscribePendingTxs(id PendingTxsSubID) bool {
	ch, ok := ff.pendingTxsSubs.Get(id)
	if !ok {
		return false
	}
	ch.Close()
	if _, ok = ff.pendingTxsSubs.Delete(id); !ok {
		return false
	}
	ff.pendingTxsStores.Delete(id)
	return true
}

// SubscribeLogs subscribes to logs using the specified filter criteria and returns a channel to receive the logs
// and a subscription ID to manage the subscription.
func (ff *Filters) SubscribeLogs(size int, crit filters.FilterCriteria) (<-chan *types.Log, LogsSubID) {
	sub := newChanSub[*types.Log](size)
	id, f := ff.logsSubs.insertLogsFilter(sub)
	f.addrs = concurrent.NewSyncMap[libcommon.Address, int]()
	if len(crit.Addresses) == 0 {
		f.allAddrs = 1
	} else {
		for _, addr := range crit.Addresses {
			f.addrs.Put(addr, 1)
		}
	}
	f.topics = concurrent.NewSyncMap[libcommon.Hash, int]()
	if len(crit.Topics) == 0 {
		f.allTopics = 1
	} else {
		for _, topics := range crit.Topics {
			for _, topic := range topics {
				f.topics.Put(topic, 1)
			}
		}
	}
	f.topicsOriginal = crit.Topics
	ff.logsSubs.addLogsFilters(f)

	// if any filter in the aggregate needs all addresses or all topics then the global log subscription needs to
	// allow all addresses or topics through
	lfr := ff.logsSubs.createFilterRequest()
	addresses, topics := ff.logsSubs.getAggMaps()
	for addr := range addresses {
		lfr.Addresses = append(lfr.Addresses, gointerfaces.ConvertAddressToH160(addr))
	}
	for topic := range topics {
		lfr.Topics = append(lfr.Topics, gointerfaces.ConvertHashToH256(topic))
	}

	loaded := ff.loadLogsRequester()
	if loaded != nil {
		if err := loaded.(func(*remote.LogsFilterRequest) error)(lfr); err != nil {
			ff.logger.Warn("Could not update remote logs filter", "err", err)
			ff.logsSubs.removeLogsFilter(id)
		}
	}

	return sub.ch, id
}

// loadLogsRequester loads the current logs requester and returns it.
func (ff *Filters) loadLogsRequester() any {
	ff.mu.Lock()
	defer ff.mu.Unlock()
	return ff.logsRequestor.Load()
}

// UnsubscribeLogs unsubscribes from logs using the given subscription ID.
// It returns true if the unsubscription was successful, otherwise false.
func (ff *Filters) UnsubscribeLogs(id LogsSubID) bool {
	isDeleted := ff.logsSubs.removeLogsFilter(id)
	// if any filters in the aggregate need all addresses or all topics then the request to the central
	// log subscription needs to honour this
	lfr := ff.logsSubs.createFilterRequest()

	addresses, topics := ff.logsSubs.getAggMaps()

	for addr := range addresses {
		lfr.Addresses = append(lfr.Addresses, gointerfaces.ConvertAddressToH160(addr))
	}
	for topic := range topics {
		lfr.Topics = append(lfr.Topics, gointerfaces.ConvertHashToH256(topic))
	}
	loaded := ff.loadLogsRequester()
	if loaded != nil {
		if err := loaded.(func(*remote.LogsFilterRequest) error)(lfr); err != nil {
			ff.logger.Warn("Could not update remote logs filter", "err", err)
			return isDeleted || ff.logsSubs.removeLogsFilter(id)
		}
	}

	ff.deleteLogStore(id)

	return isDeleted
}

// deleteLogStore deletes the log store associated with the given subscription ID.
func (ff *Filters) deleteLogStore(id LogsSubID) {
	ff.logsStores.Delete(id)
}

// OnNewEvent is called when there is a new event from the remote and processes it.
func (ff *Filters) OnNewEvent(event *remote.SubscribeReply) {
	err := ff.onNewEvent(event)
	if err != nil {
		ff.logger.Warn("OnNewEvent Filters", "event", event.Type, "err", err)
	}
}

// onNewEvent processes the given event from the remote and updates the internal state.
func (ff *Filters) onNewEvent(event *remote.SubscribeReply) error {
	switch event.Type {
	case remote.Event_HEADER:
		return ff.onNewHeader(event)
	case remote.Event_NEW_SNAPSHOT:
		ff.onNewSnapshot()
		return nil
	case remote.Event_PENDING_LOGS:
		return ff.onPendingLog(event)
	case remote.Event_PENDING_BLOCK:
		return ff.onPendingBlock(event)
	default:
		return fmt.Errorf("unsupported event type")
	}
}

// TODO: implement?
// onPendingLog handles a new pending log event from the remote.
func (ff *Filters) onPendingLog(event *remote.SubscribeReply) error {
	//	payload := event.Data
	//	var logs types.Logs
	//	err := rlp.Decode(bytes.NewReader(payload), &logs)
	//	if err != nil {
	//		// ignoring what we can't unmarshal
	//		log.Warn("OnNewEvent rpc filters (pending logs), unprocessable payload", "err", err)
	//	} else {
	//		for _, v := range ff.pendingLogsSubs {
	//			v <- logs
	//		}
	//	}
	return nil
}

// TODO: implement?
// onPendingBlock handles a new pending block event from the remote.
func (ff *Filters) onPendingBlock(event *remote.SubscribeReply) error {
	//	payload := event.Data
	//	var block types.Block
	//	err := rlp.Decode(bytes.NewReader(payload), &block)
	//	if err != nil {
	//		// ignoring what we can't unmarshal
	//		log.Warn("OnNewEvent rpc filters (pending txs), unprocessable payload", "err", err)
	//	} else {
	//		for _, v := range ff.pendingBlockSubs {
	//			v <- &block
	//		}
	//	}
	return nil
}

// onNewHeader handles a new block header event from the remote and updates the internal state.
func (ff *Filters) onNewHeader(event *remote.SubscribeReply) error {
	payload := event.Data
	var header types.Header
	if len(payload) == 0 {
		return nil
	}
	err := rlp.Decode(bytes.NewReader(payload), &header)
	if err != nil {
		return fmt.Errorf("unprocessable payload: %w", err)
	}
	return ff.headsSubs.Range(func(k HeadsSubID, v Sub[*types.Header]) error {
		v.Send(&header)
		return nil
	})
}

// OnNewTx handles a new transaction event from the transaction pool and processes it.
func (ff *Filters) OnNewTx(reply *txpool.OnAddReply) {
	txs := make([]types.Transaction, len(reply.RplTxs))
	for i, rlpTx := range reply.RplTxs {
		var decodeErr error
		if len(rlpTx) == 0 {
			continue
		}
		txs[i], decodeErr = types.DecodeTransaction(rlpTx)
		if decodeErr != nil {
			// ignoring what we can't unmarshal
			ff.logger.Warn("OnNewTx rpc filters, unprocessable payload", "err", decodeErr, "data", hex.EncodeToString(rlpTx))
			break
		}
	}
	ff.pendingTxsSubs.Range(func(k PendingTxsSubID, v Sub[[]types.Transaction]) error {
		v.Send(txs)
		return nil
	})
}

// OnNewLogs handles a new log event from the remote and processes it.
func (ff *Filters) OnNewLogs(reply *remote.SubscribeLogsReply) {
	ff.logsSubs.distributeLog(reply)
}

// AddLogs adds logs to the store associated with the given subscription ID.
func (ff *Filters) AddLogs(id LogsSubID, logs *types.Log) {
	ff.logsStores.DoAndStore(id, func(st []*types.Log, ok bool) []*types.Log {
		if !ok {
			st = make([]*types.Log, 0)
		}
		st = append(st, logs)
		return st
	})
}

// ReadLogs reads logs from the store associated with the given subscription ID.
// It returns the logs and a boolean indicating whether the logs were found.
func (ff *Filters) ReadLogs(id LogsSubID) ([]*types.Log, bool) {
	res, ok := ff.logsStores.Delete(id)
	if !ok {
		return res, false
	}
	return res, true
}

// AddPendingBlock adds a pending block header to the store associated with the given subscription ID.
func (ff *Filters) AddPendingBlock(id HeadsSubID, block *types.Header) {
	ff.pendingHeadsStores.DoAndStore(id, func(st []*types.Header, ok bool) []*types.Header {
		if !ok {
			st = make([]*types.Header, 0)
		}
		st = append(st, block)
		return st
	})
}

// ReadPendingBlocks reads pending block headers from the store associated with the given subscription ID.
// It returns the block headers and a boolean indicating whether the headers were found.
func (ff *Filters) ReadPendingBlocks(id HeadsSubID) ([]*types.Header, bool) {
	res, ok := ff.pendingHeadsStores.Delete(id)
	if !ok {
		return res, false
	}
	return res, true
}

// AddPendingTxs adds pending transactions to the store associated with the given subscription ID.
func (ff *Filters) AddPendingTxs(id PendingTxsSubID, txs []types.Transaction) {
	ff.pendingTxsStores.DoAndStore(id, func(st [][]types.Transaction, ok bool) [][]types.Transaction {
		if !ok {
			st = make([][]types.Transaction, 0)
		}
		st = append(st, txs)
		return st
	})
}

// ReadPendingTxs reads pending transactions from the store associated with the given subscription ID.
// It returns the transactions and a boolean indicating whether the transactions were found.
func (ff *Filters) ReadPendingTxs(id PendingTxsSubID) ([][]types.Transaction, bool) {
	res, ok := ff.pendingTxsStores.Delete(id)
	if !ok {
		return res, false
	}
	return res, true
}
