package service

import (
	"cryptoserver/domain"
	"cryptoserver/pkg/coingecko"
	"cryptoserver/pkg/trigger"
)

type ScheduleService struct {
	t     trigger.Trigger
	cr    CryptoRepository
	gecko *coingecko.CoinGeckoClient
}

func NewScheduleService(cr CryptoRepository, gecko *coingecko.CoinGeckoClient, triggerCfg trigger.TriggerCfg) *ScheduleService {
	ss := &ScheduleService{
		cr:    cr,
		gecko: gecko,
	}
	ss.t = *trigger.NewTrigger(ss.updateAllCrypto, triggerCfg)

	ss.t.Start()

	return ss
}

func (ss *ScheduleService) GetCfg() domain.ScheduleCfg {
	return ss.t.GetConfig()
}

func (ss *ScheduleService) UpdateCfg(enabled bool, intervalSeconds int) error {
	if intervalSeconds < 10 || intervalSeconds > 3600 {
		return domain.ErrBadRequest
	}

	ss.t.Update(trigger.TriggerCfg{
		IntervalSeconds: intervalSeconds,
	})

	ss.t.Stop()

	if enabled {
		ss.t.Start()
	}

	return nil
}

func (ss *ScheduleService) Update() (int, error) {
	res := ss.t.DoWork().(int)
	return res, nil
}

func (ss *ScheduleService) updateAllCrypto() any {
	arr, err := ss.cr.GetAll()
	if err != nil {
		return 0
	}

	symbols := make([]string, len(arr))

	for index, el := range arr {
		symbols[index] = el.Symbol
	}

	data, err := ss.gecko.GetDataSymbols(symbols...)
	if err != nil {
		return 0
	}

	for _, el := range data {
		ss.cr.Update(el.Symbol, el.Name, el.Current_price, el.Last_updated)
	}

	return len(data)
}
