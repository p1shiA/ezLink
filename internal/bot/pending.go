package bot

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

type PendingFile struct {
	ID           string
	TelegramID   int64
	ChannelID    int64
	FileID       string
	FileName     string
	FileType     string
	FileSizeByte int64
	CreatedAt    time.Time
}

type PendingUploadCache struct {
	mu    sync.RWMutex
	items map[string]*PendingFile
	ttl   time.Duration
}

var (
	GlobalPendingCache *PendingUploadCache
	oncePendingCache   sync.Once
)

func GetPendingUploadCache() *PendingUploadCache {
	oncePendingCache.Do(func() {
		GlobalPendingCache = NewPendingUploadCache(15 * time.Minute)
	})
	return GlobalPendingCache
}

func NewPendingUploadCache(ttl time.Duration) *PendingUploadCache {
	c := &PendingUploadCache{
		items: make(map[string]*PendingFile),
		ttl:   ttl,
	}
	go c.cleanupLoop()
	return c
}

func (c *PendingUploadCache) Set(file *PendingFile) string {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Generate 8-byte random hex ID (16 chars) to stay safely within Telegram's 64-byte callback limit
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	id := hex.EncodeToString(b)

	file.ID = id
	file.CreatedAt = time.Now()
	c.items[id] = file

	return id
}

func (c *PendingUploadCache) Get(id string) (*PendingFile, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, exists := c.items[id]
	if !exists {
		return nil, false
	}

	if time.Since(item.CreatedAt) > c.ttl {
		return nil, false
	}

	return item, true
}

func (c *PendingUploadCache) Delete(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, id)
}

func (c *PendingUploadCache) cleanupLoop() {
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		c.mu.Lock()
		now := time.Now()
		for id, item := range c.items {
			if now.Sub(item.CreatedAt) > c.ttl {
				delete(c.items, id)
			}
		}
		c.mu.Unlock()
	}
}
