package slotstrategy

import (
	"fmt"
)

type NumberSlot int64

func (n NumberSlot) StrValue() string {
	return fmt.Sprintf("%d", n)
}

func (n NumberSlot) Compare(slot int64) int {
	if int64(n) > slot {
		return 1
	}
	if int64(n) < slot {
		return -1
	}
	return 0
}
