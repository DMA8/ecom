package loms

import (
	"ecom/loms/internal/model"
	"fmt"
)

func (m *Manager) ConvertStatus(status model.OrderStatus) (string, error) {
	value, ok := m.statusMappings[status]
	if !ok {
		return "", fmt.Errorf("unknown status %d", status)
	}

	return value, nil
}
