package doublylinkedtree

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"

	"github.com/pkg/errors"
	fieldparams "github.com/prysmaticlabs/prysm/v5/config/fieldparams"
	"github.com/prysmaticlabs/prysm/v5/config/params"
	"github.com/prysmaticlabs/prysm/v5/consensus-types/primitives"
	"github.com/prysmaticlabs/prysm/v5/time/slots"
	"go.opencensus.io/trace"
)

// head starts from justified root and then follows the best descendant links
// to find the best block for head.
func (s *Store) head(ctx context.Context) ([32]byte, error) {
	ctx, span := trace.StartSpan(ctx, "doublyLinkedForkchoice.head")
	defer span.End()

	if err := ctx.Err(); err != nil {
		return [32]byte{}, err
	}
	if s.headNode == nil {
		if s.justifiedCheckpoint.Epoch == params.BeaconConfig().GenesisEpoch {
			s.headNode = s.treeRootNode
		} else {
			return [32]byte{}, errors.WithMessage(errUnknownJustifiedRoot, fmt.Sprintf("%#x", s.justifiedCheckpoint.Root))
		}
	}
	return s.headNode.root, nil

	// JustifiedRoot has to be known
	//justifiedNode, ok := s.nodeByRoot[s.justifiedCheckpoint.Root]
	//if !ok || justifiedNode == nil {
	//	// If the justifiedCheckpoint is from genesis, then the root is
	//	// zeroHash. In this case it should be the root of forkchoice
	//	// tree.
	//	if s.justifiedCheckpoint.Epoch == params.BeaconConfig().GenesisEpoch {
	//		justifiedNode = s.treeRootNode
	//	} else {
	//		return [32]byte{}, errors.WithMessage(errUnknownJustifiedRoot, fmt.Sprintf("%#x", s.justifiedCheckpoint.Root))
	//	}
	//}
	//
	//// If the justified node doesn't have a best descendant,
	//// the best node is itself.
	//bestDescendant := justifiedNode.bestDescendant
	//if bestDescendant == nil {
	//	bestDescendant = justifiedNode
	//}
	//currentEpoch := slots.EpochsSinceGenesis(time.Unix(int64(s.genesisTime), 0))
	//if !bestDescendant.viableForHead(s.justifiedCheckpoint.Epoch, currentEpoch) {
	//	s.allTipsAreInvalid = true
	//	return [32]byte{}, fmt.Errorf("head at slot %d with weight %d is not eligible, finalizedEpoch, justified Epoch %d, %d != %d, %d",
	//		bestDescendant.slot, bestDescendant.weight/10e9, bestDescendant.finalizedEpoch, bestDescendant.justifiedEpoch, s.finalizedCheckpoint.Epoch, s.justifiedCheckpoint.Epoch)
	//}
	//s.allTipsAreInvalid = false
	//
	//// Update metrics.
	//if bestDescendant != s.headNode {
	//	headChangesCount.Inc()
	//	headSlotNumber.Set(float64(bestDescendant.slot))
	//	s.headNode = bestDescendant
	//}
	//
	//return bestDescendant.root, nil
}

// insert registers a new block node to the fork choice store's node list.
// It then updates the new node's parent with the best child and descendant node.
func (s *Store) insert(ctx context.Context,
	slot primitives.Slot,
	root, parentRoot, payloadHash [fieldparams.RootLength]byte,
	justifiedEpoch, finalizedEpoch primitives.Epoch) (*Node, error) {
	ctx, span := trace.StartSpan(ctx, "doublyLinkedForkchoice.insert")
	defer span.End()

	// Return if the block has been inserted into Store before.
	if n, ok := s.nodeByRoot[root]; ok {
		return n, nil
	}

	parent := s.nodeByRoot[parentRoot]
	n := &Node{
		slot:           slot,
		root:           root,
		parent:         parent,
		justifiedEpoch: justifiedEpoch,
		finalizedEpoch: finalizedEpoch,
		optimistic:     true,
		payloadHash:    payloadHash,
		timestamp:      uint64(time.Now().Unix()),
	}

	// Set the node's target checkpoint
	if slot%params.BeaconConfig().SlotsPerEpoch == 0 {
		n.target = n
	} else if parent != nil {
		if slots.ToEpoch(slot) == slots.ToEpoch(parent.slot) {
			n.target = parent.target
		} else {
			n.target = parent
		}
	}

	s.nodeByPayload[payloadHash] = n
	s.nodeByRoot[root] = n
	if parent == nil {
		if s.treeRootNode == nil {
			s.treeRootNode = n
			s.headNode = n
			n.stabled = true
		} else {
			return n, errInvalidParentRoot
		}
	} else {
		parent.children = append(parent.children, n)
		// Apply proposer boost
		timeNow := uint64(time.Now().Unix())
		if timeNow < s.genesisTime {
			return n, nil
		}

		currentSlot := slots.CurrentSlot(s.genesisTime)

		// Update best descendants
		jEpoch := s.justifiedCheckpoint.Epoch
		fEpoch := s.finalizedCheckpoint.Epoch

		{
			if count, exist := s.cacheAttCount[parentRoot]; exist {
				if count >= ValidatorPerSlot/2 {
					n.stabled = true
				}
				log.WithFields(logrus.Fields{
					"slot":             slot,
					"root":             fmt.Sprintf("0x%x", root),
					"parent":           fmt.Sprintf("0x%x", parentRoot),
					"parent slot":      parent.slot,
					"parent att count": count,
					"stable":           n.stabled,
				}).Info("Debug ForkChoice insert")
			} else {
				if parent.slot == 0 {
					n.stabled = true
				}
				log.WithFields(logrus.Fields{
					"slot":        slot,
					"root":        fmt.Sprintf("0x%x", root),
					"parent slot": parent.slot,
					"stable":      n.stabled,
				}).Info("Debug ForkChoice insert")
			}
		}

		if err := s.updateBestDescendant(ctx, jEpoch, fEpoch, slots.ToEpoch(currentSlot)); err != nil {
			return n, err
		}
	}

	// Update metrics.
	processedBlockCount.Inc()
	nodeCount.Set(float64(len(s.nodeByRoot)))

	// Only update received block slot if it's within epoch from current time.
	if slot+params.BeaconConfig().SlotsPerEpoch > slots.CurrentSlot(s.genesisTime) {
		s.receivedBlocksLastEpoch[slot%params.BeaconConfig().SlotsPerEpoch] = slot
	}

	return n, nil
}

// pruneFinalizedNodeByRootMap prunes the `nodeByRoot` map
// starting from `node` down to the finalized Node or to a leaf of the Fork
// choice store.
func (s *Store) pruneFinalizedNodeByRootMap(ctx context.Context, node, finalizedNode *Node) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if node == finalizedNode {
		if node.target != node {
			node.target = nil
		}
		return nil
	}
	for _, child := range node.children {
		if err := s.pruneFinalizedNodeByRootMap(ctx, child, finalizedNode); err != nil {
			return err
		}
	}

	node.children = nil
	delete(s.nodeByRoot, node.root)
	delete(s.nodeByPayload, node.payloadHash)
	return nil
}

// prune prunes the fork choice store. It removes all nodes that compete with the finalized root.
// This function does not prune for invalid optimistically synced nodes, it deals only with pruning upon finalization
func (s *Store) prune(ctx context.Context) error {
	ctx, span := trace.StartSpan(ctx, "doublyLinkedForkchoice.Prune")
	defer span.End()

	finalizedRoot := s.finalizedCheckpoint.Root
	finalizedEpoch := s.finalizedCheckpoint.Epoch
	finalizedNode, ok := s.nodeByRoot[finalizedRoot]
	if !ok || finalizedNode == nil {
		return errors.WithMessage(errUnknownFinalizedRoot, fmt.Sprintf("%#x", finalizedRoot))
	}
	// return early if we haven't changed the finalized checkpoint
	if finalizedNode.parent == nil {
		return nil
	}

	// Prune nodeByRoot starting from root
	if err := s.pruneFinalizedNodeByRootMap(ctx, s.treeRootNode, finalizedNode); err != nil {
		return err
	}

	finalizedNode.parent = nil
	s.treeRootNode = finalizedNode

	prunedCount.Inc()
	// Prune all children of the finalized checkpoint block that are incompatible with it
	checkpointMaxSlot, err := slots.EpochStart(finalizedEpoch)
	if err != nil {
		return errors.Wrap(err, "could not compute epoch start")
	}
	if finalizedNode.slot == checkpointMaxSlot {
		return nil
	}

	for _, child := range finalizedNode.children {
		if child != nil && child.slot <= checkpointMaxSlot {
			if err := s.pruneFinalizedNodeByRootMap(ctx, child, finalizedNode); err != nil {
				return errors.Wrap(err, "could not prune incompatible finalized child")
			}
		}
	}
	return nil
}

// tips returns a list of possible heads from fork choice store, it returns the
// roots and the slots of the leaf nodes.
func (s *Store) tips() ([][32]byte, []primitives.Slot) {
	var roots [][32]byte
	var slots []primitives.Slot

	for root, node := range s.nodeByRoot {
		if len(node.children) == 0 {
			roots = append(roots, root)
			slots = append(slots, node.slot)
		}
	}
	return roots, slots
}

// updateBestDescendant updates the best descendant of this node and its
// children.
func (s *Store) updateBestDescendant(ctx context.Context, justifiedEpoch, finalizedEpoch, currentEpoch primitives.Epoch) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}
	if len(s.treeRootNode.children) == 0 || s.treeRootNode == nil {
		s.headNode = nil
		return nil
	}
	// 1. get all leaf nodes. tips
	// 2. filter by viableForHead on tips
	// 3. get max depth by depth.
	leaves, slots := s.tips()
	filters := make([]*Node, 0)
	log.WithFields(logrus.Fields{
		"leaf count": len(leaves),
		"slots":      fmt.Sprintf("0x%v", slots),
	}).Info("Debug ForkChoice tips")
	for i, root := range leaves {
		node := s.nodeByRoot[root]
		log.WithFields(logrus.Fields{
			"slot":           slots[i],
			"justifiedEpoch": justifiedEpoch,
			"currentEpoch":   currentEpoch,
			"viableForHead":  node.viableForHead(justifiedEpoch, currentEpoch),
			"stabled":        node.stabled,
			"root":           fmt.Sprintf("0x%x", node.root),
		}).Info("Debug ForkChoice tips")
		if node != nil && node.viableForHead(justifiedEpoch, currentEpoch) && node.stabled {
			filters = append(filters, node)
		}
	}
	log.WithFields(logrus.Fields{
		"filter": len(filters),
	}).Info("Debug ForkChoice tips after filter")
	if len(filters) == 0 {
		s.headNode = nil
		return nil
	}
	maxNode := filters[0]
	for _, node := range filters {
		if node.depth() > maxNode.depth() {
			maxNode = node
			log.WithFields(logrus.Fields{
				"maxNode": fmt.Sprintf("0x%x-depth:%d", maxNode.root, maxNode.depth()),
				"node":    fmt.Sprintf("0x%x-depth:%d", node.root, node.depth()),
			}).Info("Debug ForkChoice tips after filter")
		} else if node.depth() == maxNode.depth() {
			if node.slot > maxNode.slot {
				maxNode = node
				log.WithFields(logrus.Fields{
					"maxNode": fmt.Sprintf("0x%x-depth:%d", maxNode.root, maxNode.depth()),
					"node":    fmt.Sprintf("0x%x-depth:%d", node.root, node.depth()),
				}).Info("Debug ForkChoice tips after filter when depth equal")
			}
		}
	}
	s.headNode = maxNode

	return nil
}

func (s *Store) UpdateVoted(slot uint64, root [fieldparams.RootLength]byte, count int) {
	if slot > s.cacheSlot || s.cacheAttCount == nil {
		log.WithFields(logrus.Fields{
			"slot": slot,
			"root": fmt.Sprintf("0x%x", root),
		}).Info("Debug ForkChoice UpdateVoted Cache")
		s.cacheAttCount = make(map[[fieldparams.RootLength]byte]int)
	}
	s.cacheAttCount[root] += count
	log.WithFields(logrus.Fields{
		"slot":       slot,
		"root":       fmt.Sprintf("0x%x", root),
		"new count":  count,
		"vote count": s.cacheAttCount[root],
	}).Info("Debug ForkChoice UpdateVoted")
}

// HighestReceivedBlockSlot returns the highest slot received by the forkchoice
func (f *ForkChoice) HighestReceivedBlockSlot() primitives.Slot {
	return primitives.Slot(0)
}

// HighestReceivedBlockSlotDelay returns the number of slots that the highest
// received block was late when receiving it
func (f *ForkChoice) HighestReceivedBlockDelay() primitives.Slot {
	return primitives.Slot(0)
}

// ReceivedBlocksLastEpoch returns the number of blocks received in the last epoch
func (f *ForkChoice) ReceivedBlocksLastEpoch() (uint64, error) {
	count := uint64(0)
	lowerBound := slots.CurrentSlot(f.store.genesisTime)
	var err error
	if lowerBound > fieldparams.SlotsPerEpoch {
		lowerBound, err = lowerBound.SafeSub(fieldparams.SlotsPerEpoch)
		if err != nil {
			return 0, err
		}
	}

	for _, s := range f.store.receivedBlocksLastEpoch {
		if s != 0 && lowerBound <= s {
			count++
		}
	}
	return count, nil
}
