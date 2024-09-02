package doublylinkedtree

import fieldparams "github.com/prysmaticlabs/prysm/v5/config/fieldparams"

const (
	TotalValidatorCount = 256
	ValidatorPerSlot    = TotalValidatorCount / 32
)

func changeUnstableToStable(unstable int64) int64 {
	/*
			(16, 3), (17, 3), (18, 4), (19, 5), (20, 6),
		    (21, 7), (22, 8), (23, 9), (24, 10), (25, 10),
		    (26, 11), (27, 12), (28, 13), (29, 14), (30, 16),
		    (31, 17), (32, 18), (33, 19), (34, 20), (35, 21),
		    (36, 22), (37, 23), (38, 24), (39, 26), (40, 27),
		    (41, 28), (42, 29), (43, 30), (44, 32), (45, 33),
		    (46, 34), (47, 35), (48, 37), (49, 38), (50, 39),
		    (51, 40), (52, 42), (53, 43), (54, 44), (55, 45),
		    (56, 47), (57, 48), (58, 49), (59, 51), (60, 52),
		    (61, 53), (62, 55), (63, 56), (64, 57)
	*/
	if unstable < 16 {
		return 0
	}
	switch unstable {
	case 16, 17:
		return 3
	case 18:
		return 4
	case 19:
		return 5
	case 20:
		return 6
	case 21:
		return 7
	case 22:
		return 8
	case 23:
		return 9
	case 24:
		return 10
	case 25:
		return 10
	case 26:
		return 11
	case 27:
		return 12
	case 28:
		return 13
	case 29:
		return 14
	case 30:
		return 16
	case 31:
		return 17
	case 32:
		return 18
	case 33:
		return 19
	case 34:
		return 20
	case 35:
		return 21
	case 36:
		return 22
	case 37:
		return 23
	case 38:
		return 24
	case 39:
		return 26
	case 40:
		return 27
	case 41:
		return 28
	case 42:
		return 29
	case 43:
		return 30
	case 44:
		return 32
	case 45:
		return 33
	case 46:
		return 34
	case 47:
		return 35
	case 48:
		return 37
	case 49:
		return 38
	case 50:
		return 39
	case 51:
		return 40
	case 52:
		return 42
	case 53:
		return 43
	case 54:
		return 44
	case 55:
		return 45
	case 56:
		return 47
	case 57:
		return 48
	case 58:
		return 49
	case 59:
		return 51
	case 60:
		return 52
	case 61:
		return 53
	case 62:
		return 55
	case 63:
		return 56
	case 64:
		return 57
	default:
		return 0
	}
}

func (n *Node) UpdateVoted(store *Store, slot uint64) {
	if len(n.validatorIndices) >= ValidatorPerSlot/3 {
		n.stabled = true
		store.mux.Lock()
		// store block status to votedSlotBlock
		if _, ok := store.votedSlotBlock[slot]; !ok {
			store.votedSlotBlock[slot] = make(map[[fieldparams.RootLength]byte]*Node)
		}
		store.votedSlotBlock[slot][n.root] = n
		store.mux.Unlock()
	}
}
