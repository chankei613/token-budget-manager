// Package events はbudget:alertイベントを、購読中のSSEクライアントへ
// リアルタイム配信するための小さなpub/subを提供する。
package events

import (
	"sync"
	"time"
)

const EventBudgetAlert = "budget:alert"

// Event はSSEで配信される1件分のイベント。
type Event struct {
	Type     string    `json:"type"`
	BudgetID string    `json:"budget_id"`
	Level    string    `json:"level"`
	At       time.Time `json:"at"`
}

// Broker は購読者へのファンアウトのみを担当する。
type Broker struct {
	mu   sync.RWMutex
	subs map[chan Event]struct{}
}

func NewBroker() *Broker {
	return &Broker{subs: make(map[chan Event]struct{})}
}

// Subscribe は新しい購読チャネルと、その解除関数を返す。
func (b *Broker) Subscribe() (<-chan Event, func()) {
	ch := make(chan Event, 32)

	b.mu.Lock()
	b.subs[ch] = struct{}{}
	b.mu.Unlock()

	var once sync.Once
	unsubscribe := func() {
		once.Do(func() {
			b.mu.Lock()
			delete(b.subs, ch)
			b.mu.Unlock()
			close(ch)
		})
	}
	return ch, unsubscribe
}

// Publish は全購読者に配信する。遅い購読者のバッファが埋まっていたら捨てる（ブロックしない）。
func (b *Broker) Publish(e Event) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for ch := range b.subs {
		select {
		case ch <- e:
		default:
		}
	}
}
