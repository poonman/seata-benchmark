package handler

import (
	"github.com/poonman/seata-benchmark/config"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

type Handler struct {
	conf *config.Config
	wg   sync.WaitGroup
}

func NewHandler(conf *config.Config) (b *Handler) {
	return &Handler{
		conf: conf,
	}
}

type Stats struct {
	Success             int
	Failure             int
	MaxCostDuration     time.Duration
	MinCostDuration     time.Duration
	TotalCostDuration   time.Duration
	SuccessCostDuration time.Duration
	FailureCostDuration time.Duration
}

type Report struct {
	Concurrency         int
	Success             int
	Failure             int
	TPS                 float32
	MaxCostDuration     time.Duration
	MinCostDuration     time.Duration
	AvgCostDuration     time.Duration
	TotalCostDuration   time.Duration
	SuccessCostDuration time.Duration
	FailureCostDuration time.Duration
}

func (b *Handler) Run() {

	var statsSet []*Stats

	b.wg.Add(b.conf.Benchmark.Concurrency)

	for i := 1; i <= b.conf.Benchmark.Concurrency; i++ {
		stats := &Stats{
			Success:             0,
			Failure:             0,
			MaxCostDuration:     0,
			MinCostDuration:     1000 * time.Second,
			TotalCostDuration:   0,
			SuccessCostDuration: 0,
			FailureCostDuration: 0,
		}

		statsSet = append(statsSet, stats)

		go b.Request(int64(i), stats)
	}

	b.wg.Wait()

	// 统计

	b.Stats(statsSet)
}

func (b *Handler) Stats(statsSet []*Stats) {
	rep := &Report{
		Concurrency:         b.conf.Benchmark.Concurrency,
		Success:             0,
		Failure:             0,
		TPS:                 0,
		MaxCostDuration:     0,
		MinCostDuration:     1000 * time.Second,
		AvgCostDuration:     0,
		TotalCostDuration:   0,
		SuccessCostDuration: 0,
		FailureCostDuration: 0,
	}

	for _, s := range statsSet {
		if s.MinCostDuration < rep.MinCostDuration {
			rep.MinCostDuration = s.MinCostDuration
		}

		if s.MaxCostDuration > rep.MaxCostDuration {
			rep.MaxCostDuration = s.MaxCostDuration
		}

		rep.Success += s.Success
		rep.Failure += s.Failure

		rep.TotalCostDuration += s.TotalCostDuration
		rep.SuccessCostDuration += s.SuccessCostDuration
		rep.FailureCostDuration += s.FailureCostDuration
	}

	//log.Infof("rep:%+v", *rep)

	rep.TPS = float32(b.conf.Benchmark.Concurrency) * float32(rep.Success) / (float32(rep.TotalCostDuration) / float32(time.Second))
	rep.AvgCostDuration = rep.TotalCostDuration / time.Duration(rep.Success+rep.Failure)

	log.Debugf("report: ")
	log.Infof("%20s: %20d", "Concurrency", rep.Concurrency)
	log.Infof("%20s: %20d", "Success", rep.Success)
	log.Infof("%20s: %20d", "Failure", rep.Failure)
	log.Infof("%20s: %20.3f", "TPS", rep.TPS)
	log.Infof("%20s: %20v", "MaxCostDuration", rep.MaxCostDuration)
	log.Infof("%20s: %20v", "MinCostDuration", rep.MinCostDuration)
	log.Infof("%20s: %20v", "AvgCostDuration", rep.AvgCostDuration)
	log.Infof("%20s: %20v", "TotalCostDuration", rep.TotalCostDuration)
}

func (b *Handler) Request(uid int64, stats *Stats) {
	//ctx := context.TODO()

	for i := 0; i < b.conf.Benchmark.RequestNumPerCon; i++ {
		before := time.Now()
		err := Request(b.conf.Url)
		cost := time.Now().Sub(before)
		if err != nil {
			stats.Failure++
		} else {
			stats.Success++
		}

		if cost > stats.MaxCostDuration {
			stats.MaxCostDuration = cost
		}

		if cost < stats.MinCostDuration {
			stats.MinCostDuration = cost
		}

		stats.TotalCostDuration += cost
	}

	b.wg.Done()
}
