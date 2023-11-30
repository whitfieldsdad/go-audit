package monitor

type Subscription struct {
	Filters []Filter `json:"filters"`
}

func NewSubscription() *Subscription {
	return &Subscription{}
}

func (s *Subscription) AddFilters(f ...Filter) {
	s.Filters = append(s.Filters, f...)
}
