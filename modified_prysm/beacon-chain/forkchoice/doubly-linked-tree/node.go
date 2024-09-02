package doublylinkedtree

import (
	"context"

	"github.com/prysmaticlabs/prysm/v5/config/params"
	forkchoice2 "github.com/prysmaticlabs/prysm/v5/consensus-types/forkchoice"
	"github.com/prysmaticlabs/prysm/v5/consensus-types/primitives"
	"github.com/prysmaticlabs/prysm/v5/time/slots"
)

// ProcessAttestationsThreshold  is the number of seconds after which we
// process attestations for the current slot
const ProcessAttestationsThreshold = 10

// depth returns the length of the path to the root of Fork Choice
func (n *Node) depth() uint64 {
	ret := uint64(0)
	for node := n.parent; node != nil; node = node.parent {
		ret += 1
	}
	return ret
}

// applyWeightChanges recomputes the weight of the node passed as an argument and all of its descendants,
// using the current balance stored in each node.
//func (n *Node) applyWeightChanges(ctx context.Context) error {
//	// Recursively calling the children to sum their weights.
//	childrenWeight := uint64(0)
//	for _, child := range n.children {
//		if ctx.Err() != nil {
//			return ctx.Err()
//		}
//		if err := child.applyWeightChanges(ctx); err != nil {
//			return err
//		}
//		childrenWeight += child.weight
//	}
//	if n.root == params.BeaconConfig().ZeroHash {
//		return nil
//	}
//	n.weight = n.balance + childrenWeight
//	return nil
//}

// updateBestDescendant updates the best descendant of this node and its
// children.
//func (n *Node) updateBestDescendant(ctx context.Context, store *Store, justifiedEpoch, finalizedEpoch, currentEpoch primitives.Epoch) error {
//	if ctx.Err() != nil {
//		return ctx.Err()
//	}
//	if len(n.children) == 0 {
//		n.bestDescendant = nil
//		return nil
//	}
//	// todo: luxq 用最长链算法替换当前算法，最后更新到 n.bestDescendant
//
//	// 1. 获取所有的叶子节点. tips
//	// 2. filter by viableForHead on tips
//	// 3. get max depth by depth.
//
//	return nil
//}

// viableForHead returns true if the node is viable to head.
// Any node with different finalized or justified epoch than
// the ones in fork choice store should not be viable to head.
func (n *Node) viableForHead(justifiedEpoch, currentEpoch primitives.Epoch) bool {
	justified := justifiedEpoch == n.justifiedEpoch || justifiedEpoch == 0

	return justified
}

//
//func (n *Node) leadsToViableHead(justifiedEpoch, currentEpoch primitives.Epoch) bool {
//	if n.bestDescendant == nil {
//		return n.viableForHead(justifiedEpoch, currentEpoch)
//	}
//	return n.bestDescendant.viableForHead(justifiedEpoch, currentEpoch)
//}

// setNodeAndParentValidated sets the current node and all the ancestors as validated (i.e. non-optimistic).
func (n *Node) setNodeAndParentValidated(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if !n.optimistic {
		return nil
	}
	n.optimistic = false

	if n.parent == nil {
		return nil
	}
	return n.parent.setNodeAndParentValidated(ctx)
}

// arrivedEarly returns whether this node was inserted before the first
// threshold to orphan a block.
// Note that genesisTime has seconds granularity, therefore we use a strict
// inequality < here. For example a block that arrives 3.9999 seconds into the
// slot will have secs = 3 below.
func (n *Node) arrivedEarly(genesisTime uint64) (bool, error) {
	secs, err := slots.SecondsSinceSlotStart(n.slot, genesisTime, n.timestamp)
	votingWindow := params.BeaconConfig().SecondsPerSlot / params.BeaconConfig().IntervalsPerSlot
	return secs < votingWindow, err
}

// arrivedAfterOrphanCheck returns whether this block was inserted after the
// intermediate checkpoint to check for candidate of being orphaned.
// Note that genesisTime has seconds granularity, therefore we use an
// inequality >= here. For example a block that arrives 10.00001 seconds into the
// slot will have secs = 10 below.
func (n *Node) arrivedAfterOrphanCheck(genesisTime uint64) (bool, error) {
	secs, err := slots.SecondsSinceSlotStart(n.slot, genesisTime, n.timestamp)
	return secs >= ProcessAttestationsThreshold, err
}

// nodeTreeDump appends to the given list all the nodes descending from this one
func (n *Node) nodeTreeDump(ctx context.Context, nodes []*forkchoice2.Node) ([]*forkchoice2.Node, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	var parentRoot [32]byte
	if n.parent != nil {
		parentRoot = n.parent.root
	}
	thisNode := &forkchoice2.Node{
		Slot:                     n.slot,
		BlockRoot:                n.root[:],
		ParentRoot:               parentRoot[:],
		JustifiedEpoch:           n.justifiedEpoch,
		FinalizedEpoch:           n.finalizedEpoch,
		UnrealizedJustifiedEpoch: n.justifiedEpoch,
		UnrealizedFinalizedEpoch: n.finalizedEpoch,
		Balance:                  n.balance,
		Weight:                   n.weight,
		ExecutionOptimistic:      n.optimistic,
		ExecutionBlockHash:       n.payloadHash[:],
		Timestamp:                n.timestamp,
	}
	if n.optimistic {
		thisNode.Validity = forkchoice2.Optimistic
	} else {
		thisNode.Validity = forkchoice2.Valid
	}

	nodes = append(nodes, thisNode)
	var err error
	for _, child := range n.children {
		nodes, err = child.nodeTreeDump(ctx, nodes)
		if err != nil {
			return nil, err
		}
	}
	return nodes, nil
}
