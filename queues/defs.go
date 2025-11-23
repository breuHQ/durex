package queues

import (
	"encoding/json"
)

type (

	// Signal is a string alias intended for defining groups of workflow signals, "register" , "send_welcome_email" etc.
	// It ensures consistency and code clarity. The Signal type provides methods for conversion and serialization,
	// promoting good developer experience.
	Signal string

	// Query is a string alias intended for defining groups of workflow queries. We could have created an alias for
	// for Signal type, but for some wierd reason, if was causing temporal to panic when marshalling the type to JSON.
	Query string
)

func (s Signal) String() string {
	return string(s)
}

func (s Signal) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(s))
}

func (s *Signal) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	*s = Signal(str)

	return nil
}

func (q Query) String() string {
	return string(q)
}

func (q Query) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(q))
}

func (q *Query) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	*q = Query(str)

	return nil
}
