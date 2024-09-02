package doublylinkedtree

// orphanLateBlockProposingEarly determines the maximum threshold that we
// consider the node is proposing early and sure to receive proposer boost
const orphanLateBlockProposingEarly = 2

// ShouldOverrideFCU returns whether the current forkchoice head is weak
// and thus may be reorged when proposing the next block.
// This function should only be called if the following two conditions are
// satisfied:
// 1-   It is immediately after receiving a block that may be subject to a reorg
//
//	or
//
//	It is right after processAttestationsThreshold and we have processed the
//	current slots attestations.
//
// 2- The caller has already called Forkchoice.Head() so that forkchoice has
// been updated.
// 3- The beacon node is serving a validator that will propose during the next
// slot.
//
// This function only applies a heuristic to decide if the beacon will update
// the engine's view of head with the parent block or the incoming block. It
// does not guarantee an attempted reorg. This will only be decided later at
// proposal time by calling GetProposerHead.
func (f *ForkChoice) ShouldOverrideFCU() (override bool) {
	override = false
	return
}

// GetProposerHead returns the block root that has to be used as ParentRoot by a
// proposer. It may not be the actual head of the canonical chain, in certain
// cases it may be its parent, when the last head block has arrived early and is
// considered safe to be orphaned.
//
// This function needs to be called only when proposing a block and all
// attestation processing has already happened.
func (f *ForkChoice) GetProposerHead() [32]byte {
	head := f.store.headNode
	if head == nil {
		return [32]byte{}
	}
	return head.root
}
