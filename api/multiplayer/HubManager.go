package multiplayer

import (
	"davidmultiplayersnake/utils"
	"errors"
	"sync"
)

//HubManager manages the game hubs
type HubManager struct {
	hubs   map[string]*Hub
	Start  chan *Hub
	End    chan *Hub
	logger *utils.Logger
	lock   sync.Mutex
}

// NewHubManager allocates and returns a *HubManager
func NewHubManager(logger *utils.Logger) *HubManager {
	return &HubManager{
		hubs:   make(map[string]*Hub),
		Start:  make(chan *Hub),
		End:    make(chan *Hub),
		logger: logger,
	}
}

// HubList returns a slice of *Hub
func (hm *HubManager) HubList() (hubs []*Hub) {
	for _, hub := range hm.hubs {
		hubs = append(hubs, hub)
	}

	return
}

// RegisterHub allocates and register a new Hub under a hub name
func (hm *HubManager) RegisterHub(hubName string) error {
	hm.lock.Lock()
	defer hm.lock.Unlock()

	hm.logger.LogChan <- hubName

	_, ok := hm.hubs[hubName]

	if ok {
		return errors.New("Hub name is already registered")
	}

	h := NewHub(hm, hm.logger)
	h.Name = hubName
	hm.hubs[hubName] = h
	hm.Start <- h

	return nil
}

// GetHub gets a specific hub
func (hm *HubManager) GetHub(hubName string) *Hub {
	return hm.hubs[hubName]
}

// Run concurrently starts hubs and
func (hm *HubManager) Run() {
	for {
		select {
		case h := <-hm.Start:
			hm.logger.LogChan <- "Start " + h.Name
			go h.Run()
		case h := <-hm.End:
			hm.lock.Lock()
			delete(hm.hubs, h.Name)
			hm.logger.LogChan <- "End " + h.Name
			hm.lock.Unlock()
		}
	}
}
