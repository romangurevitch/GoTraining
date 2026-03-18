package receivers

// Immutable struct
type immutable struct {
	Value string `json:"value"`
}

func (i immutable) SetString(s string) *immutable {
	i.Value = s
	return &i
}

func (i immutable) String() string {
	return i.Value
}

func (i immutable) MarshalJSON() ([]byte, error) {
	return []byte(`{"value":"yes it will be changed!"}`), nil
}

// Mutable struct
type mutable struct {
	Value string `json:"value"`
}

func (m *mutable) SetString(s string) *mutable {
	m.Value = s
	return m
}

func (m *mutable) String() string {
	return m.Value
}

func (m *mutable) MarshalJSON() ([]byte, error) {
	return []byte(`{"value":"no it will not change!"}`), nil
}
