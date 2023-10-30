package memory

import (
	"context"
)

func (m *Memory) Up(_ context.Context) error {

	if m.cfg.IsRestore() {
		if err := m.Restore(); err != nil {
			return err
		}
	}

	return nil
}
