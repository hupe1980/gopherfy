package internal

type InsertionOrderMap struct {
	values map[string]string
	keys   []string
}

func NewInsertionOrderMap() *InsertionOrderMap {
	return &InsertionOrderMap{
		values: make(map[string]string),
		keys:   []string{},
	}
}

func (m *InsertionOrderMap) Keys() []string {
	return m.keys
}

func (m *InsertionOrderMap) Set(k, v string) {
	if !Contains(m.keys, k) {
		m.keys = append(m.keys, k)
	}

	m.values[k] = v
}

func (m *InsertionOrderMap) Get(k string) (string, bool) {
	v, found := m.values[k]
	return v, found
}
