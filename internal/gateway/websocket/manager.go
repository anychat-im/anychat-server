package websocket

import (
	"encoding/json"
	"sync"

	"github.com/anychat/server/pkg/logger"
	"go.uber.org/zap"
)

// Manager WebSocket连接管理器
type Manager struct {
	clients map[string]map[string]*Client // userID -> deviceID -> client
	mu      sync.RWMutex
}

// NewManager 创建WebSocket连接管理器
func NewManager() *Manager {
	return &Manager{
		clients: make(map[string]map[string]*Client),
	}
}

// Register 注册新客户端；同一用户同一设备重复连接时替换旧连接
func (m *Manager) Register(client *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.clients[client.UserID]; !exists {
		m.clients[client.UserID] = make(map[string]*Client)
	}

	if old, exists := m.clients[client.UserID][client.DeviceID]; exists {
		close(old.Done)
		logger.Info("Replaced existing WebSocket connection",
			zap.String("userID", client.UserID),
			zap.String("deviceID", client.DeviceID))
	}

	m.clients[client.UserID][client.DeviceID] = client
	logger.Info("WebSocket client registered",
		zap.String("userID", client.UserID),
		zap.String("deviceID", client.DeviceID),
		zap.Int("deviceCount", len(m.clients[client.UserID])))
}

// Unregister 注销客户端（仅当传入的client是当前活跃client时才注销）
func (m *Manager) Unregister(client *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if userClients, exists := m.clients[client.UserID]; exists {
		if current, ok := userClients[client.DeviceID]; ok && current == client {
			delete(userClients, client.DeviceID)
			if len(userClients) == 0 {
				delete(m.clients, client.UserID)
			}
			logger.Info("WebSocket client unregistered",
				zap.String("userID", client.UserID),
				zap.String("deviceID", client.DeviceID))
		}
	}
}

// SendToUser 向指定用户发送原始消息，返回是否成功
func (m *Manager) SendToUser(userID string, data []byte) bool {
	m.mu.RLock()
	userClients, exists := m.clients[userID]
	if !exists || len(userClients) == 0 {
		m.mu.RUnlock()
		return false
	}
	type targetClient struct {
		deviceID string
		client   *Client
	}
	targets := make([]targetClient, 0, len(userClients))
	for deviceID, client := range userClients {
		targets = append(targets, targetClient{
			deviceID: deviceID,
			client:   client,
		})
	}
	m.mu.RUnlock()

	sent := false
	for _, target := range targets {
		select {
		case target.client.Send <- data:
			sent = true
		default:
			logger.Warn("WebSocket send buffer full, dropping message",
				zap.String("userID", userID),
				zap.String("deviceID", target.deviceID))
		}
	}
	return sent
}

// SendMessageToUser 向指定用户发送结构化消息
func (m *Manager) SendMessageToUser(userID string, msg *Message) bool {
	data, err := json.Marshal(msg)
	if err != nil {
		logger.Error("Failed to marshal WebSocket message", zap.Error(err))
		return false
	}
	return m.SendToUser(userID, data)
}

// IsOnline 检查用户是否在线
func (m *Manager) IsOnline(userID string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	userClients, exists := m.clients[userID]
	return exists && len(userClients) > 0
}

// OnlineCount 获取当前在线用户数
func (m *Manager) OnlineCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.clients)
}
