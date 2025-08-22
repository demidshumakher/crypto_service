package rest

import (
	"cryptoserver/domain"
	"encoding/json"
	"net/http"
	"time"
)

type ScheduleService interface {
	GetCfg() domain.ScheduleCfg
	UpdateCfg(enabled bool, intervalSeconds int) error
	Update() (int, error)
}

type ScheduleHandler struct {
	ss ScheduleService
}

func NewScheduleHandler(ss ScheduleService, mx *Router) {
	sh := &ScheduleHandler{
		ss: ss,
	}

	mx.Handle("GET /schedule", http.HandlerFunc(sh.GetCfg))
	mx.Handle("PUT /schedule", http.HandlerFunc(sh.UpdateCfg))
	mx.Handle("POST /schedule/trigger", http.HandlerFunc(sh.Update))
}

func (sh *ScheduleHandler) GetCfg(w http.ResponseWriter, r *http.Request) {
	res := sh.ss.GetCfg()
	json.NewEncoder(w).Encode(res)
}

type scheduleUpdateConfigRequest struct {
	Enabled         bool `json:"enabled"`
	IntervalSeconds int  `json:"interval_seconds"`
}

func (sh *ScheduleHandler) UpdateCfg(w http.ResponseWriter, r *http.Request) {
	cfg := &scheduleUpdateConfigRequest{}

	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(cfg)
	if err != nil {
		WriteError(w, err)
		return
	}

	err = sh.ss.UpdateCfg(cfg.Enabled, cfg.IntervalSeconds)
	if err != nil {
		WriteError(w, err)
		return
	}
	json.NewEncoder(w).Encode(cfg)
}

type scheduleUpdateResponse struct {
	UpdatedCount int       `json:"updated_count"`
	Timestamp    time.Time `json:"timestamp"`
}

func (sh *ScheduleHandler) Update(w http.ResponseWriter, r *http.Request) {
	cnt, err := sh.ss.Update()
	if err != nil {
		WriteError(w, err)
		return
	}

	res := scheduleUpdateResponse{
		UpdatedCount: cnt,
		Timestamp:    time.Now(),
	}

	json.NewEncoder(w).Encode(res)
}
