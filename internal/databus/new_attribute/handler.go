package new_attribute

import (
	"context"
	"encoding/json"
	"fmt"
)

func NewAttribute(ctx context.Context, msg []byte) error {
	// Пример десериализации JSON сообщения
	var data map[string]interface{}
	if err := json.Unmarshal(msg, &data); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Обработка сообщения
	fmt.Printf("Received message: %v\n", data)
	return nil
}
