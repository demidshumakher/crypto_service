package trigger

import (
	"cryptoserver/domain"
	"sync"
	"time"
)

type Work func() any

type TriggerCfg struct {
	IntervalSeconds int
}

type Trigger struct {
	cfg     TriggerCfg
	enabled bool

	work Work

	updateChan chan struct{}
	cancelChan chan struct{}

	lastUpdate time.Time
	nextUpdate time.Time

	mx *sync.Mutex
}

func NewTrigger(work Work, cfg TriggerCfg) *Trigger {
	return &Trigger{
		cfg:        cfg,
		work:       work,
		updateChan: make(chan struct{}, 1),
		mx:         &sync.Mutex{},
	}
}

func (t *Trigger) GetConfig() domain.ScheduleCfg {
	return domain.ScheduleCfg{
		IntervalSeconds: t.cfg.IntervalSeconds,
		LastUpdate:      t.lastUpdate,
		NextUpdate:      t.nextUpdate,
		Enabled:         t.enabled,
	}
}

func (t *Trigger) Start() {
	cancel := make(chan struct{}, 1)

	t.mx.Lock()

	t.enabled = true
	t.cancelChan = cancel

	t.mx.Unlock()

	go func() {
		t.mx.Lock()
		interval := time.Second * time.Duration(t.cfg.IntervalSeconds)
		t.mx.Unlock()

		ticker := time.NewTicker(interval)

		for {
			select {
			case <-cancel:
				ticker.Stop()
				return
			case <-ticker.C:
				t.lastUpdate = time.Now()
				t.nextUpdate = time.Now().Add(interval)
				t.work()
			case <-t.updateChan:
				ticker.Stop()

				t.mx.Lock()
				interval = time.Second * time.Duration(t.cfg.IntervalSeconds)
				t.mx.Unlock()

				ticker = time.NewTicker(interval)
			}
		}
	}()
}

func (t *Trigger) Stop() {
	if t.cancelChan == nil {
		return
	}
	select {
	case <-t.cancelChan:
	default:
		close(t.cancelChan)
	}
}

func (t *Trigger) Update(cfg TriggerCfg) {
	t.mx.Lock()
	defer t.mx.Unlock()

	t.cfg = cfg
	t.updateChan <- struct{}{}
}

func (t *Trigger) DoWork() any {
	return t.work()
}

func (t *Trigger) IsEnabled() bool {
	return t.enabled
}
