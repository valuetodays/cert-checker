package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type BarkNotifier struct {
	Enabled   bool   `yaml:"enabled"`
	ServerUrl string `yaml:"server_url"`
	DeviceKey string `yaml:"device_key"`
	Group     string `yaml:"group"`
}

func (n *BarkNotifier) Send(msg AlertMessage) error {
	if !n.Enabled {
		return nil
	}
	payload := map[string]interface{}{
		"device_key": n.DeviceKey,
		"title":      "SSL证书到期提醒: " + msg.Domain,
		"body":       fmt.Sprintf("SSL证书到期剩余天数: %d", msg.DaysLeft),
		"group":      n.Group,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal bark payload error: %v", err)
	}
	resp, err := http.Post(n.ServerUrl, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("send bark message error: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bark returned non-200 status: %d", resp.StatusCode)
	}
	return nil
}
func (n *BarkNotifier) Name() string {
	return "Bark"
}

func (n *BarkNotifier) IsEnabled() bool {
	return n.Enabled
}
