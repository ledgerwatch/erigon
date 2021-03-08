package headerdownload

import (
	"container/heap"
	"context"
	"encoding/base64"
	"fmt"
	"math/big"
	"sort"

	"github.com/ledgerwatch/turbo-geth/common"
	"github.com/ledgerwatch/turbo-geth/common/dbutils"
	"github.com/ledgerwatch/turbo-geth/core/rawdb"
	"github.com/ledgerwatch/turbo-geth/core/types"
	"github.com/ledgerwatch/turbo-geth/eth/stagedsync/stages"
	"github.com/ledgerwatch/turbo-geth/ethdb"
	"github.com/ledgerwatch/turbo-geth/log"
	"github.com/ledgerwatch/turbo-geth/rlp"
)

// Implements sort.Interface so we can sort the incoming header in the message by block height
type HeadersByBlockHeight []*types.Header

func (h HeadersByBlockHeight) Len() int {
	return len(h)
}

func (h HeadersByBlockHeight) Less(i, j int) bool {
	// Note - the ordering is the inverse ordering of the block heights
	return h[i].Number.Cmp(h[j].Number) > 0
}

func (h HeadersByBlockHeight) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

// SplitIntoSegments converts message containing headers into a collection of chain segments
func (hd *HeaderDownload) SplitIntoSegments(headersRaw [][]byte, msg []*types.Header) ([]*ChainSegment, Penalty, error) {
	hd.lock.RLock()
	defer hd.lock.RUnlock()
	sort.Sort(HeadersByBlockHeight(msg))
	// Now all headers are order from the highest block height to the lowest
	var segments []*ChainSegment                         // Segments being built
	segmentMap := make(map[common.Hash]int)              // Mapping of the header hash to the index of the chain segment it belongs
	childrenMap := make(map[common.Hash][]*types.Header) // Mapping parent hash to the children
	dedupMap := make(map[common.Hash]struct{})           // Map used for detecting duplicate headers
	for i, header := range msg {
		headerHash := header.Hash()
		if _, bad := hd.badHeaders[headerHash]; bad {
			return nil, BadBlockPenalty, nil
		}
		if _, duplicate := dedupMap[headerHash]; duplicate {
			return nil, DuplicateHeaderPenalty, nil
		}
		dedupMap[headerHash] = struct{}{}
		var segmentIdx int
		children := childrenMap[headerHash]
		for _, child := range children {
			if valid, penalty := hd.childParentValid(child, header); !valid {
				return nil, penalty, nil
			}
		}
		if len(children) == 1 {
			// Single child, extract segmentIdx
			segmentIdx = segmentMap[headerHash]
		} else {
			// No children, or more than one child, create new segment
			segmentIdx = len(segments)
			segments = append(segments, &ChainSegment{})
		}
		segments[segmentIdx].Headers = append(segments[segmentIdx].Headers, header)
		segments[segmentIdx].HeadersRaw = append(segments[segmentIdx].HeadersRaw, headersRaw[i])
		segmentMap[header.ParentHash] = segmentIdx
		siblings := childrenMap[header.ParentHash]
		siblings = append(siblings, header)
		childrenMap[header.ParentHash] = siblings
	}
	return segments, NoPenalty, nil
}

// Checks whether child-parent relationship between two headers is correct
// (excluding Proof Of Work validity)
func (hd *HeaderDownload) childParentValid(child, parent *types.Header) (bool, Penalty) {
	if parent.Number.Uint64()+1 != child.Number.Uint64() {
		return false, WrongChildBlockHeightPenalty
	}
	return true, NoPenalty
}

// SingleHeaderAsSegment converts message containing 1 header into one singleton chain segment
func (hd *HeaderDownload) SingleHeaderAsSegment(headerRaw []byte, header *types.Header) ([]*ChainSegment, Penalty, error) {
	hd.lock.RLock()
	defer hd.lock.RUnlock()
	headerHash := header.Hash()
	if _, bad := hd.badHeaders[headerHash]; bad {
		return nil, BadBlockPenalty, nil
	}
	return []*ChainSegment{{HeadersRaw: [][]byte{headerRaw}, Headers: []*types.Header{header}}}, NoPenalty, nil
}

// FindAnchors attempts to find anchors to which given chain segment can be attached to
func (hd *HeaderDownload) findAnchors(segment *ChainSegment) (found bool, start int) {
	// Walk the segment from children towards parents
	for i, header := range segment.Headers {
		// Check if the header can be attached to an anchor of a working tree
		if _, attaching := hd.anchors[header.Hash()]; attaching {
			return true, i
		}
	}
	return false, 0
}

// FindTip attempts to find tip of a tree that given chain segment can be attached to
// the given chain segment may be found invalid relative to a working tree, in this case penalty for peer is returned
func (hd *HeaderDownload) findTip(segment *ChainSegment, start int) (found bool, end int, penalty Penalty) {
	if _, duplicate := hd.getTip(segment.Headers[start].Hash()); duplicate {
		return false, 0, NoPenalty
	}
	// Walk the segment from children towards parents
	for i, header := range segment.Headers[start:] {
		// Check if the header can be attached to any tips
		if _, attaching := hd.getTip(header.ParentHash); attaching {
			return true, start + i + 1, NoPenalty
		}
	}
	return false, len(segment.Headers), NoPenalty
}

// VerifySeals verifies Proof Of Work for the part of the given chain segment
// It reports first verification error, or returns nil
func (hd *HeaderDownload) verifySeals(segment *ChainSegment, start, end int) error {
	for _, header := range segment.Headers[start:end] {
		if err := hd.verifySealFunc(header); err != nil {
			return err
		}
	}
	return nil
}

func (hd *HeaderDownload) removeUpwards(toRemove []*Tip) {
	for len(toRemove) > 0 {
		removal := toRemove[len(toRemove)-1]
		toRemove = toRemove[:len(toRemove)-1]
		delete(hd.tips, removal.header.Hash())
		hd.tipCount--
		toRemove = append(toRemove, removal.next...)
	}
}

func (hd *HeaderDownload) markPreverified(tip *Tip) {
	// Go through all parent tips that are not preveried and mark them too
	var prevTip *Tip
	for tip != nil && !tip.preverified {
		tip.preverified = true
		if prevTip != nil && len(tip.next) > 1 {
			// Remove all non-canonical tips
			var toRemove []*Tip
			for _, n := range tip.next {
				if n != prevTip {
					toRemove = append(toRemove, n)
				}
			}
			hd.removeUpwards(toRemove)
			tip.next = append(tip.next[:0], prevTip)
		}
		tip = hd.tips[tip.header.ParentHash]
	}
}

// ExtendUp extends a working tree up from the tip, using given chain segment
func (hd *HeaderDownload) extendUp(segment *ChainSegment, start, end int) error {
	// Find attachment tip again
	tipHeader := segment.Headers[end-1]
	if attachmentTip, attaching := hd.getTip(tipHeader.ParentHash); attaching {
		if attachmentTip.preverified && len(attachmentTip.next) > 0 {
			return fmt.Errorf("cannot extendUp from preverified link %d with children", attachmentTip.blockHeight)
		}
		// Iterate over headers backwards (from parents towards children), to be able calculate cumulative difficulty along the way
		prevTip := attachmentTip
		for i := end - 1; i >= start; i-- {
			header := segment.Headers[i]
			if tip, err := hd.addHeaderAsTip(header, false /* persisted */); err == nil {
				prevTip.next = append(prevTip.next, tip)
				prevTip = tip
				if _, hard := hd.hardTips[header.Hash()]; hard {
					hd.markPreverified(tip)
				}
			} else {
				return fmt.Errorf("extendUp addHeaderAsTip for %x: %v", header.Hash(), err)
			}
		}
	} else {
		return fmt.Errorf("extendUp attachment tip not found for %x", tipHeader.ParentHash)
	}
	return nil
}

func (hd *HeaderDownload) Ready() (bool, uint64) {
	hd.lock.Lock()
	defer hd.lock.Unlock()
	return hd.stageReady, hd.stageHeight
}

func (hd *HeaderDownload) StageReadyChannel() chan struct{} {
	return hd.stageReadyCh
}

// ExtendDown extends some working trees down from the anchor, using given chain segment
// it creates a new anchor and collects all the tips from the attached anchors to it
func (hd *HeaderDownload) extendDown(segment *ChainSegment, start, end int) error {
	// Find attachement anchor again
	anchorHeader := segment.Headers[start]
	if anchor, attaching := hd.anchors[anchorHeader.Hash()]; attaching {
		anchorPreverified := false
		for _, tip := range anchor.tips {
			if tip.preverified {
				anchorPreverified = true
				break
			}
		}
		newAnchorHeader := segment.Headers[end-1]
		newAnchor := &Anchor{
			parentHash:  newAnchorHeader.ParentHash,
			timestamp:   0,
			blockHeight: newAnchorHeader.Number.Uint64(),
		}
		hd.anchors[newAnchorHeader.ParentHash] = newAnchor
		heap.Push(hd.anchorQueue, newAnchor)

		delete(hd.anchors, anchor.parentHash)
		// Add all headers in the segments as tips to this anchor
		var prevTip *Tip
		for i := end - 1; i >= start; i-- {
			header := segment.Headers[i]
			if tip, err := hd.addHeaderAsTip(header, false /* pesisted */); err == nil {
				if prevTip == nil {
					newAnchor.tips = append(newAnchor.tips, tip)
				} else {
					prevTip.next = append(prevTip.next, tip)
				}
				prevTip = tip
				if !anchorPreverified {
					if _, hard := hd.hardTips[header.Hash()]; hard {
						hd.markPreverified(tip)
					}
				}
			} else {
				return fmt.Errorf("extendUp addHeaderAsTip for %x: %v", header.Hash(), err)
			}
		}
		prevTip.next = anchor.tips
		anchor.tips = nil
		if anchorPreverified {
			// Mark the entire segment as preverified
			hd.markPreverified(prevTip)
		}
	} else {
		return fmt.Errorf("extendDown attachment anchors not found for %x", anchorHeader.Hash())
	}
	return nil
}

// Connect connects some working trees using anchors of some, and a tip of another
func (hd *HeaderDownload) connect(segment *ChainSegment, start, end int) error {
	// Find attachment tip again
	tipHeader := segment.Headers[end-1]
	// Find attachement anchors again
	anchorHeader := segment.Headers[start]
	attachmentTip, ok1 := hd.getTip(tipHeader.ParentHash)
	if !ok1 {
		return fmt.Errorf("connect attachment tip not found for %x", tipHeader.ParentHash)
	}
	if attachmentTip.preverified && len(attachmentTip.next) > 0 {
		return fmt.Errorf("cannot connect to preverified link %d with children", attachmentTip.blockHeight)
	}
	anchor, ok2 := hd.anchors[anchorHeader.Hash()]
	if !ok2 {
		return fmt.Errorf("connect attachment anchors not found for %x", anchorHeader.Hash())
	}
	anchorPreverified := false
	for _, tip := range anchor.tips {
		if tip.preverified {
			anchorPreverified = true
			break
		}
	}
	delete(hd.anchors, anchor.parentHash)
	// Iterate over headers backwards (from parents towards children), to be able calculate cumulative difficulty along the way
	prevTip := attachmentTip
	for i := end - 1; i >= start; i-- {
		header := segment.Headers[i]
		if tip, err := hd.addHeaderAsTip(header, false /* persisted */); err == nil {
			prevTip.next = append(prevTip.next, tip)
			prevTip = tip
			if !anchorPreverified {
				if _, hard := hd.hardTips[header.Hash()]; hard {
					hd.markPreverified(tip)
				}
			}
		} else {
			return fmt.Errorf("extendUp addHeaderAsTip for %x: %v", header.Hash(), err)
		}
	}
	prevTip.next = anchor.tips
	anchor.tips = nil
	if anchorPreverified {
		// Mark the entire segment as preverified
		hd.markPreverified(prevTip)
	}
	if attachmentTip.persisted {
		tip := hd.tips[tipHeader.Hash()]
		hd.insertList = append(hd.insertList, tip)
	}
	// If we connect to the hard-coded tip, we remove it. Once there is only one hard-coded tip left, it is clear that everything is connected
	delete(hd.hardTips, tipHeader.ParentHash)
	return nil
}

func (hd *HeaderDownload) newAnchor(segment *ChainSegment, start, end int) error {
	anchorHeader := segment.Headers[end-1]

	anchor := &Anchor{
		parentHash:  anchorHeader.ParentHash,
		timestamp:   0,
		blockHeight: anchorHeader.Number.Uint64(),
	}
	hd.anchors[anchorHeader.ParentHash] = anchor
	heap.Push(hd.anchorQueue, anchor)
	// Iterate over headers backwards (from parents towards children), to be able calculate cumulative difficulty along the way
	var prevTip *Tip
	for i := end - 1; i >= start; i-- {
		header := segment.Headers[i]
		if tip, err1 := hd.addHeaderAsTip(header, false /* persisted */); err1 == nil {
			if prevTip == nil {
				anchor.tips = append(anchor.tips, tip)
			} else {
				prevTip.next = append(prevTip.next, tip)
			}
			prevTip = tip
			if _, hard := hd.hardTips[header.Hash()]; hard {
				hd.markPreverified(tip)
			}
		} else {
			return fmt.Errorf("newAnchor addHeaderAsTip for %x: %v", header.Hash(), err1)
		}
	}
	return nil
}

/*
func (hd *HeaderDownload) AnchorState() string {
	hd.lock.RLock()
	defer hd.lock.RUnlock()
	return hd.anchorState()
}
*/

/*
func (hd *HeaderDownload) anchorState() string {
	//nolint:prealloc
	var ss []string
	for anchorParent, anchor := range hd.anchors {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("{%8d", anchor.blockHeight))
		end := anchor.maxTipHeight
		var sbb strings.Builder
		var bs []int
		for _, tipQueueItem := range *anchor.tipQueue {
			bs = append(bs, int(tipQueueItem.height))
		}
		sort.Ints(bs)
		for j, b := range bs {
			if j == 0 {
				sbb.WriteString(fmt.Sprintf("%d", b))
			} else if j == len(bs)-1 {
				if bs[j-1]+1 == b {
					// Close interval
					sbb.WriteString(fmt.Sprintf("-%d", b))
				} else {
					// Standalone
					sbb.WriteString(fmt.Sprintf(" %d", b))
				}
			} else {
				if bs[j-1] == b {
					// Skip
				} else if bs[j-1]+1 == b {
					if b+1 == bs[j+1] {
						// Skip
					} else {
						// Close interval
						sbb.WriteString(fmt.Sprintf("-%d", b))
					}
				} else {
					// Open interval or standalone
					sbb.WriteString(fmt.Sprintf(" %d", b))
				}
			}
		}
		if end == 0 {
			sb.WriteString(fmt.Sprintf(" HardCoded tips=%d tipStretch=%d (%s)}", anchor.tipQueue.Len(), anchor.tipStretch(), sbb.String()))
		} else {
			sb.WriteString(fmt.Sprintf("-%d (%d) tips=%d tipStretch=%d (%s)}", end, end-anchor.blockHeight, anchor.tipQueue.Len(), anchor.tipStretch(), sbb.String()))
		}
		sb.WriteString(fmt.Sprintf(" => %x", anchorParent))
		ss = append(ss, sb.String())
	}
	sort.Strings(ss)
	return strings.Join(ss, "\n")
}
*/

func InitHardCodedTips(network string) map[common.Hash]HeaderRecord {
	var encodings []string
	switch network {
	case "mainnet":
		encodings = mainnetHardCodedHeaders
	default:
		log.Error("Hard coded headers not found for", "network", network)
		return nil
	}

	// Insert hard-coded headers if present
	return DecodeTips(encodings)
}

func DecodeTips(encodings []string) map[common.Hash]HeaderRecord {
	hardTips := make(map[common.Hash]HeaderRecord, len(encodings))

	for _, encoding := range encodings {
		b, err := base64.RawStdEncoding.DecodeString(encoding)
		if err != nil {
			log.Error("Parsing hard coded header", "error", err)
		} else {
			var h types.Header
			if err := rlp.DecodeBytes(b, &h); err != nil {
				log.Error("Parsing hard coded header", "error", err)
			} else {
				hardTips[h.Hash()] = HeaderRecord{Raw: b, Header: &h}
			}
		}
	}

	return hardTips
}

func (hd *HeaderDownload) SetHardCodedTips(hardTips map[common.Hash]HeaderRecord) {
	hd.lock.Lock()
	defer hd.lock.Unlock()
	var maxHeight uint64
	for tipHash, headerRecord := range hardTips {
		height := headerRecord.Header.Number.Uint64()
		if height <= hd.highestInDb {
			// No need for this hard coded header anymore
			continue
		}
		if height <= hd.maxHardTipHeight {
			continue
		}
		if height > maxHeight {
			maxHeight = height
		}
		hd.hardTips[tipHash] = headerRecord
	}
	if maxHeight > hd.maxHardTipHeight {
		hd.maxHardTipHeight = maxHeight
	}
}

func (hd *HeaderDownload) RecoverFromDb(db ethdb.Database) error {
	err := db.(ethdb.HasKV).KV().View(context.Background(), func(tx ethdb.Tx) error {
		c := tx.Cursor(dbutils.HeaderPrefix)
		// Take hd.persistedTipLimit headers (with the highest heights) as tips
		for k, v, err := c.Last(); k != nil && hd.tipQueue.Len() < hd.persistedTipLimit; k, v, err = c.Prev() {
			if err != nil {
				return err
			}
			if len(k) != 40 {
				continue
			}
			var h types.Header
			if err = rlp.DecodeBytes(v, &h); err != nil {
				return err
			}
			if tip, err1 := hd.addHeaderAsTip(&h, true /* persisted */); err1 == nil {
				heap.Push(hd.tipQueue, tip)
			} else {
				return err1
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	hd.highestInDb, err = stages.GetStageProgress(db, stages.Headers)
	if err != nil {
		return err
	}
	return nil
}

func (hd *HeaderDownload) invalidateAnchor(anchor *Anchor) {
	fmt.Printf("invalidateAnchor %d\n", anchor.blockHeight)
	delete(hd.anchors, anchor.parentHash)
	hd.removeUpwards(anchor.tips)
}

func (hd *HeaderDownload) RequestMoreHeaders(currentTime, timeout uint64) *HeaderRequest {
	hd.lock.Lock()
	defer hd.lock.Unlock()

	if hd.anchorQueue.Len() == 0 {
		return nil
	}
	for hd.anchorQueue.Len() > 0 {
		anchor := (*hd.anchorQueue)[0]
		if _, ok := hd.anchors[anchor.parentHash]; ok {
			if anchor.timestamp <= currentTime {
				if anchor.timeouts < 10 {
					anchor.timeouts++
					anchor.timestamp = currentTime + timeout
					heap.Fix(hd.anchorQueue, 0)
					return &HeaderRequest{Hash: anchor.parentHash, Number: anchor.blockHeight - 1, Length: 192, Skip: 0, Reverse: true}
				} else {
					// Ancestors of this anchor seem to be unavailable, invalidate and move on
					hd.invalidateAnchor(anchor)
				}
			} else {
				// Anchor not ready for re-request yet
				return nil
			}
		}
		// Anchor disappered or unavailble, pop from the queue and move on
		heap.Remove(hd.anchorQueue, 0)
	}
	return nil
}

func (hd *HeaderDownload) RequestSkeleton() *HeaderRequest {
	hd.lock.RLock()
	defer hd.lock.RUnlock()
	if len(hd.anchors) > 4 {
		return nil // Need to be below anchor threshold to produce skeleton request
	}
	length := int(hd.topSeenHeight-hd.highestInDb) / 192
	if length > 192 {
		length = 192
	}
	if length == 0 {
		return nil // No need in sketelon request
	}
	return &HeaderRequest{Number: hd.highestInDb, Length: length, Skip: 192, Reverse: false}
}

func (hd *HeaderDownload) InsertHeaders(hf func(header *types.Header, blockHeight uint64) error) error {
	hd.lock.Lock()
	defer hd.lock.Unlock()
	for len(hd.insertList) > 0 {
		tip := hd.insertList[len(hd.insertList)-1]
		if tip.blockHeight <= hd.maxHardTipHeight && !tip.preverified {
			// Header should be preverified, but not yet, try again later
			break
		}
		hd.insertList = hd.insertList[:len(hd.insertList)-1]
		hd.tipCount--
		if !tip.preverified {
			if err := hd.verifySealFunc(tip.header); err != nil {
				log.Error("Verification failed for header", "hash", tip.header.Hash(), "height", tip.blockHeight, "error", err)
				// skip this tip and its children
				continue
			}
		}
		if err := hf(tip.header, tip.blockHeight); err != nil {
			return err
		}
		if tip.blockHeight > hd.highestInDb {
			hd.highestInDb = tip.blockHeight
		}
		tip.persisted = true
		heap.Push(hd.tipQueue, tip)
		if len(tip.next) > 0 {
			hd.insertList = append(hd.insertList, tip.next...)
		}
	}
	for hd.tipQueue.Len() > hd.persistedTipLimit {
		tip := heap.Pop(hd.tipQueue).(*Tip)
		delete(hd.tips, tip.hash)
	}
	return nil
}

func (hd *HeaderDownload) Progress() uint64 {
	hd.lock.RLock()
	defer hd.lock.RUnlock()
	return hd.highestInDb
}

func (hd *HeaderDownload) HasTip(tipHash common.Hash) bool {
	hd.lock.RLock()
	defer hd.lock.RUnlock()
	if _, ok := hd.getTip(tipHash); ok {
		return true
	}
	return false
}

func (hd *HeaderDownload) getTip(tipHash common.Hash) (*Tip, bool) {
	if tip, ok := hd.tips[tipHash]; ok {
		return tip, true
	}
	return nil, false
}

// addHeaderAsTip adds given header as a tip belonging to a given anchorParent
func (hd *HeaderDownload) addHeaderAsTip(header *types.Header, persisted bool) (*Tip, error) {
	height := header.Number.Uint64()
	tipHash := header.Hash()
	tip := &Tip{
		blockHeight: height,
		hash:        tipHash,
		header:      header,
		persisted:   persisted,
	}
	hd.tips[tipHash] = tip
	hd.tipCount++
	return tip, nil
}

func (hi *HeaderInserter) FeedHeader(header *types.Header, blockHeight uint64) error {
	hash := header.Hash()
	if hash == hi.prevHash {
		// Skip duplicates
		return nil
	}
	if blockHeight < hi.prevHeight {
		return fmt.Errorf("[%s] headers are unexpectedly unsorted, got %d after %d", hi.logPrefix, blockHeight, hi.prevHeight)
	}
	if oldH := rawdb.ReadHeader(hi.batch, hash, blockHeight); oldH != nil {
		// Already inserted, skip
		return nil
	}
	// Load parent header
	parent := rawdb.ReadHeader(hi.batch, header.ParentHash, blockHeight-1)
	if parent == nil {
		log.Error(fmt.Sprintf("Could not find parent with hash %x and height %d for header %x %d", header.ParentHash, blockHeight-1, hash, blockHeight))
		// Skip headers without parents
		return nil
	}
	// Parent's total difficulty
	parentTd, err := rawdb.ReadTd(hi.batch, header.ParentHash, blockHeight-1)
	if err != nil {
		return fmt.Errorf("[%s] parent's total difficulty not found with hash %x and height %d for header %x %d: %v", hi.logPrefix, header.ParentHash, blockHeight-1, hash, blockHeight, err)
	}
	// Calculate total difficulty of this header using parent's total difficulty
	td := new(big.Int).Add(parentTd, header.Difficulty)
	// Now we can decide wether this header will create a change in the canonical head
	if td.Cmp(hi.localTd) > 0 {
		hi.newCanonical = true
		// Find the forking point - i.e. the latest header on the canonical chain which is an ancestor of this one
		// Most common case - forking point is the height of the parent header
		var forkingPoint uint64
		ch, err1 := rawdb.ReadCanonicalHash(hi.batch, blockHeight-1)
		if err1 != nil {
			return fmt.Errorf("reading canonical hash for height %d: %w", blockHeight-1, err1)
		}
		if ch == (common.Hash{}) || ch == header.ParentHash {
			forkingPoint = blockHeight - 1
		} else {
			// Going further back
			ancestorHash := parent.ParentHash
			ancestorHeight := blockHeight - 2
			for ch, err = rawdb.ReadCanonicalHash(hi.batch, ancestorHeight); err == nil && ch != ancestorHash; ch, err = rawdb.ReadCanonicalHash(hi.batch, ancestorHeight) {
				ancestor := rawdb.ReadHeader(hi.batch, ancestorHash, ancestorHeight)
				ancestorHash = ancestor.ParentHash
				ancestorHeight--
			}
			if err != nil {
				return fmt.Errorf("[%s] reading canonical hash for %d: %w", hi.logPrefix, ancestorHeight, err)
			}
			// Loop above terminates when either err != nil (handled already) or ch == ancestorHash, therefore ancestorHeight is our forking point
			forkingPoint = ancestorHeight
		}
		if err = rawdb.WriteHeadHeaderHash(hi.batch, hash); err != nil {
			return fmt.Errorf("[%s] marking head header hash as %x: %w", hi.logPrefix, hash, err)
		}
		// See if the forking point affects the unwindPoint (the block number to which other stages will need to unwind before the new canonical chain is applied)
		if forkingPoint < hi.unwindPoint {
			hi.unwindPoint = forkingPoint
		}
	}
	data, err2 := rlp.EncodeToBytes(header)
	if err2 != nil {
		return fmt.Errorf("[%s] failed to RLP encode header: %w", hi.logPrefix, err2)
	}
	if err = rawdb.WriteTd(hi.batch, hash, blockHeight, td); err != nil {
		return fmt.Errorf("[%s] failed to WriteTd: %w", hi.logPrefix, err)
	}
	if err = hi.batch.Put(dbutils.HeaderPrefix, dbutils.HeaderKey(blockHeight, hash), data); err != nil {
		return fmt.Errorf("[%s] failed to store header: %w", hi.logPrefix, err)
	}
	if blockHeight > hi.headerProgress {
		hi.headerProgress = blockHeight
		if err = stages.SaveStageProgress(hi.batch, stages.Headers, blockHeight); err != nil {
			return fmt.Errorf("[%s] saving Headers progress: %w", hi.logPrefix, err)
		}
	}
	hi.prevHash = hash
	if blockHeight > hi.highest {
		hi.highest = blockHeight
	}
	return nil
}

func (hi *HeaderInserter) GetHighest() uint64 {
	return hi.highest
}

func (hi *HeaderInserter) UnwindPoint() uint64 {
	return hi.unwindPoint
}

//nolint:interfacer
func (hd *HeaderDownload) ProcessSegment(segment *ChainSegment) {
	log.Debug("processSegment", "from", segment.Headers[0].Number.Uint64(), "to", segment.Headers[len(segment.Headers)-1].Number.Uint64())
	hd.lock.Lock()
	defer hd.lock.Unlock()
	foundAnchor, start := hd.findAnchors(segment)
	foundTip, end, penalty := hd.findTip(segment, start) // We ignore penalty because we will check it as part of PoW check
	if penalty != NoPenalty {
		log.Error(fmt.Sprintf("FindTip penalty %d", penalty))
		return
	}
	if end == 0 {
		log.Debug("Duplicate segment")
		return
	}
	hd.topSeenHeight = segment.Headers[len(segment.Headers)-1].Number.Uint64()
	if segment.Headers[0].Number.Uint64() > hd.maxHardTipHeight && hd.highestInDb < hd.maxHardTipHeight {
		//TODO: Change it to log.Debug
		log.Info("Segment cannot be inserted yet", "height", segment.Headers[0].Number.Uint64(), "maxHardCoded", hd.maxHardTipHeight, "in Db", hd.highestInDb)
		return
	}
	if hd.tipCount+end-start > hd.tipLimit {
		log.Warn("Too many tips", "count", hd.tipCount, "tried to add", end-start, "limit", hd.tipLimit)
		//return
	}
	// There are 4 cases
	if foundAnchor {
		if foundTip {
			// Connect
			if err := hd.connect(segment, start, end); err != nil {
				log.Error("Connect failed", "error", err)
				return
			}
			log.Debug("Connected", "start", start, "end", end)
		} else {
			// ExtendDown
			if err := hd.extendDown(segment, start, end); err != nil {
				log.Error("ExtendDown failed", "error", err)
				return
			}
			log.Debug("Extended Down", "start", start, "end", end)
		}
	} else if foundTip {
		if end > 0 {
			// ExtendUp
			if err := hd.extendUp(segment, start, end); err != nil {
				log.Error("ExtendUp failed", "error", err)
				return
			}
			log.Debug("Extended Up", "start", start, "end", end)
		}
	} else {
		// NewAnchor
		if err := hd.newAnchor(segment, start, end); err != nil {
			log.Error("NewAnchor failed", "error", err)
			return
		}
		log.Debug("NewAnchor", "start", start, "end", end)
	}
	//log.Info(hd.anchorState())
}

func (hd *HeaderDownload) TopSeenHeight() uint64 {
	hd.lock.RLock()
	defer hd.lock.RUnlock()
	return hd.topSeenHeight
}
