package memory

import "context"

func (m *Memory) Down(_ context.Context) error {

	if err := m.Snapshot(); err != nil {
		return err
	}

	return nil
}
