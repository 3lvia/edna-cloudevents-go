package ednaevents

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"sync"
)

var metrics = &metricsGuard{
	counterMap: map[string]prometheus.Counter{},
	mux:        &sync.Mutex{},
}

type metricsGuard struct {
	counterMap map[string]prometheus.Counter
	mux        *sync.Mutex
}
func (g *metricsGuard) incCounter(name string, constLabels map[string]string) {
	c := g.getCounter(name, constLabels)
	c.Inc()
}

func (g *metricsGuard) getCounter(name string, constLabels map[string]string) prometheus.Counter {
	key := metricsKey(name, constLabels)
	if c, ok := g.counterMap[key]; ok {
		return c
	}

	g.mux.Lock()
	defer g.mux.Unlock()

	if c, ok := g.counterMap[key]; ok {
		return c
	}

	c := promauto.NewCounter(prometheus.CounterOpts{
		Name:        name,
		ConstLabels: constLabels,
	})

	g.counterMap[key] = c

	return c
}

func metricsKey(name string, constLabels map[string]string) string {
	if constLabels == nil {
		return name
	}
	key := name
	for k, v := range constLabels {
		key = fmt.Sprintf("%s%s%s", key, k, v)
	}
	return key
}