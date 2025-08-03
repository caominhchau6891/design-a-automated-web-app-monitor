package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// WebAppMonitor monitors a web application for availability and performance
type WebAppMonitor struct {
	// URL of the web application to monitor
	URL string
	// Interval at which to check the web application
	Interval time.Duration
	// WebSocket connection to send updates to clients
	WsConn *websocket.Conn
}

// NewWebAppMonitor creates a new WebAppMonitor instance
func NewWebAppMonitor(url string, interval time.Duration, wsConn *websocket.Conn) *WebAppMonitor {
	return &WebAppMonitor{
		URL:     url,
		Interval: interval,
		WsConn:  wsConn,
	}
}

// StartMonitoring begins monitoring the web application
func (m *WebAppMonitor) StartMonitoring() {
	ticker := time.NewTicker(m.Interval)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				m.checkWebApp()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

// checkWebApp checks the web application for availability and performance
func (m *WebAppMonitor) checkWebApp() {
	start := time.Now()
	resp, err := http.Get(m.URL)
	if err != nil {
		log.Printf("Error checking web app: %v", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("Web app returned status code %d", resp.StatusCode)
		return
	}
	duration := time.Since(start)
	log.Printf("Web app responded in %v", duration)
	m.sendUpdate(duration)
}

// sendUpdate sends an update to connected clients
func (m *WebAppMonitor) sendUpdate(duration time.Duration) {
	message := fmt.Sprintf("Web app responded in %v", duration)
	m.WsConn.WriteMessage(websocket.TextMessage, []byte(message))
}

func main() {
	wsUpgrader := websocket.Upgrader{}
	wsConn, err := wsUpgrader.Upgrade(nil, nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer wsConn.Close()

	monitor := NewWebAppMonitor("https://example.com", 10*time.Second, wsConn)
	monitor.StartMonitoring()

	select {} // keep the program running
}