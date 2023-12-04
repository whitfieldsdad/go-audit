package monitor

type Subscription struct {
	ProcessFilter *ProcessFilter `json:"processes,omitempty"`
}

func NewSubscription() *Subscription {
	return &Subscription{
		ProcessFilter: NewProcessFilter(),
	}
}

func (s *Subscription) Merge(other Subscription) {
	if s.ProcessFilter == nil {
		s.ProcessFilter = other.ProcessFilter
	} else {
		s.ProcessFilter.Merge(*other.ProcessFilter)
	}
}
