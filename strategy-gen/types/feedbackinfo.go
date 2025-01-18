package types

type FeedBackInfo struct {
	ReorgCount           int
	ImpactValidatorCount int
}

type FeedBack struct {
	Uid  string
	Info FeedBackInfo
}

type FeedBacker interface {
	WaitFeedback(uid string, ch chan FeedBack)
}
