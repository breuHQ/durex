package queues

type (
	WorkflowSignal string
)

func (s WorkflowSignal) String() string {
	return string(s)
}

func (s WorkflowSignal) MarshalJSON() ([]byte, error) {
	return []byte(s.String()), nil
}

func (s *WorkflowSignal) UnmarshalJSON(data []byte) error {
	*s = WorkflowSignal(data)
	return nil
}
