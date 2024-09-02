package doublylinkedtree

import (
	"sync"

	"github.com/prysmaticlabs/prysm/v5/beacon-chain/forkchoice"
	forkchoicetypes "github.com/prysmaticlabs/prysm/v5/beacon-chain/forkchoice/types"
	fieldparams "github.com/prysmaticlabs/prysm/v5/config/fieldparams"
	"github.com/prysmaticlabs/prysm/v5/consensus-types/primitives"
)

// ForkChoice defines the overall fork choice store which includes all block nodes, validator's latest votes and balances.
type ForkChoice struct {
	sync.RWMutex
	store               *Store
	votes               []Vote                      // tracks individual validator's last vote.
	balances            []uint64                    // tracks individual validator's balances last accounted in votes.
	justifiedBalances   []uint64                    // tracks individual validator's last justified balances.
	numActiveValidators uint64                      // tracks the total number of active validators.
	balancesByRoot      forkchoice.BalancesByRooter // handler to obtain balances for the state with a given root
}

// Store defines the fork choice store which includes block nodes and the last view of checkpoint information.
type Store struct {
	justifiedCheckpoint     *forkchoicetypes.Checkpoint            // latest justified epoch in store.
	prevJustifiedCheckpoint *forkchoicetypes.Checkpoint            // previous justified checkpoint in store.
	finalizedCheckpoint     *forkchoicetypes.Checkpoint            // latest finalized epoch in store.
	treeRootNode            *Node                                  // the root node of the store tree.
	headNode                *Node                                  // last head Node
	nodeByRoot              map[[fieldparams.RootLength]byte]*Node // nodes indexed by roots.
	nodeByPayload           map[[fieldparams.RootLength]byte]*Node // nodes indexed by payload Hash
	originRoot              [fieldparams.RootLength]byte           // The genesis block root
	genesisTime             uint64
	receivedBlocksLastEpoch [fieldparams.SlotsPerEpoch]primitives.Slot        // Using `highestReceivedSlot`. The slot of blocks received in the last epoch.
	allTipsAreInvalid       bool                                              // tracks if all tips are not viable for head
	votedSlotBlock          map[uint64]map[[fieldparams.RootLength]byte]*Node // attest slot => (block root => ChainNode)
	mux                     sync.Mutex
}

// Node defines the individual block which includes its block parent, ancestor and how much weight accounted for it.
// This is used as an array based stateful DAG for efficient fork choice look up.
type Node struct {
	slot           primitives.Slot              // slot of the block converted to the node.
	root           [fieldparams.RootLength]byte // root of the block converted to the node.
	payloadHash    [fieldparams.RootLength]byte // payloadHash of the block converted to the node.
	parent         *Node                        // parent index of this node.
	target         *Node                        // target checkpoint for
	children       []*Node                      // the list of direct children of this Node
	justifiedEpoch primitives.Epoch             // justifiedEpoch of this node.
	finalizedEpoch primitives.Epoch             // finalizedEpoch of this node.
	balance        uint64                       // the balance that voted for this node directly
	weight         uint64                       // weight of this node: the total balance including children
	//bestDescendant   *Node                        // bestDescendant node of this node.
	optimistic       bool   // whether the block has been fully validated or not
	timestamp        uint64 // The timestamp when the node was inserted.
	stabled          bool
	validatorIndices []uint64 // all validator index vote for current block.
}

// Vote defines an individual validator's vote.
type Vote struct {
	currentRoot [fieldparams.RootLength]byte // current voting root.
	nextRoot    [fieldparams.RootLength]byte // next voting root.
	nextEpoch   primitives.Epoch             // epoch of next voting period.
}
