package gcexporter

import (
	"fmt"
	"github.com/mailgun/groupcache/v2"
	"github.com/prometheus/client_golang/prometheus"
)

type Exporter struct {
	groups              []*groupcache.Group
	groupGets           *prometheus.Desc
	groupCacheHits      *prometheus.Desc
	groupPeerLoads      *prometheus.Desc
	groupPeerErrors     *prometheus.Desc
	groupLoads          *prometheus.Desc
	groupLoadsDeduped   *prometheus.Desc
	groupLocalLoads     *prometheus.Desc
	groupLocalLoadErrs  *prometheus.Desc
	groupServerRequests *prometheus.Desc
	cacheBytes          *prometheus.Desc
	cacheItems          *prometheus.Desc
	cacheGets           *prometheus.Desc
	cacheHits           *prometheus.Desc
	cacheEvictions      *prometheus.Desc
}

func NewExporter(namespace, subsystem, prefix string, groups ...*groupcache.Group) *Exporter {
	return &Exporter{
		groups: groups,
		groupGets: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, fmt.Sprintf("%sgets_total", prefix)),
			"todo",
			[]string{"group"},
			nil,
		),
		groupCacheHits: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, fmt.Sprintf("%shits_total", prefix)),
			"todo",
			[]string{"group"},
			nil,
		),
		groupPeerLoads: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, fmt.Sprintf("%speer_loads_total", prefix)),
			"todo",
			[]string{"group"},
			nil,
		),
		groupPeerErrors: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, fmt.Sprintf("%speer_errors_total", prefix)),
			"todo",
			[]string{"group"},
			nil,
		),
		groupLoads: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, fmt.Sprintf("%sloads_total", prefix)),
			"todo",
			[]string{"group"},
			nil,
		),
		groupLoadsDeduped: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, fmt.Sprintf("%sloads_deduped_total", prefix)),
			"todo",
			[]string{"group"},
			nil,
		),
		groupLocalLoads: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, fmt.Sprintf("%slocal_load_total", prefix)),
			"todo",
			[]string{"group"},
			nil,
		),
		groupLocalLoadErrs: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, fmt.Sprintf("%slocal_load_errs_total", prefix)),
			"todo",
			[]string{"group"},
			nil,
		),
		groupServerRequests: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, fmt.Sprintf("%sserver_requests_total", prefix)),
			"todo",
			[]string{"group"},
			nil,
		),
		cacheBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, fmt.Sprintf("%scache_bytes", prefix)),
			"todo",
			[]string{"group", "type"},
			nil,
		),
		cacheItems: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, fmt.Sprintf("%scache_items", prefix)),
			"todo",
			[]string{"group", "type"},
			nil,
		),
		cacheGets: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, fmt.Sprintf("%scache_gets_total", prefix)),
			"todo",
			[]string{"group", "type"},
			nil,
		),
		cacheHits: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, fmt.Sprintf("%scache_hits_total", prefix)),
			"todo",
			[]string{"group", "type"},
			nil,
		),
		cacheEvictions: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, fmt.Sprintf("%scache_evictions_total", prefix)),
			"todo",
			[]string{"group", "type"},
			nil,
		),
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.groupGets
	ch <- e.groupCacheHits
	ch <- e.groupPeerLoads
	ch <- e.groupPeerErrors
	ch <- e.groupLoads
	ch <- e.groupLoadsDeduped
	ch <- e.groupLocalLoads
	ch <- e.groupLocalLoadErrs
	ch <- e.groupServerRequests
	ch <- e.cacheBytes
	ch <- e.cacheItems
	ch <- e.cacheGets
	ch <- e.cacheHits
	ch <- e.cacheEvictions
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	for _, group := range e.groups {
		e.collectFromGroup(ch, group)
	}
}

func (e *Exporter) collectFromGroup(ch chan<- prometheus.Metric, g *groupcache.Group) {
	e.collectStats(ch, g)
	e.collectCacheStats(ch, groupcache.HotCache, g)
	e.collectCacheStats(ch, groupcache.MainCache, g)
}

func (e *Exporter) collectStats(ch chan<- prometheus.Metric, g *groupcache.Group) {
	ch <- prometheus.MustNewConstMetric(e.groupGets, prometheus.CounterValue, float64(g.Stats.Gets.Get()), g.Name())
	ch <- prometheus.MustNewConstMetric(e.groupCacheHits, prometheus.CounterValue, float64(g.Stats.CacheHits.Get()), g.Name())
	ch <- prometheus.MustNewConstMetric(e.groupPeerLoads, prometheus.CounterValue, float64(g.Stats.PeerLoads.Get()), g.Name())
	ch <- prometheus.MustNewConstMetric(e.groupPeerErrors, prometheus.CounterValue, float64(g.Stats.PeerErrors.Get()), g.Name())
	ch <- prometheus.MustNewConstMetric(e.groupLoads, prometheus.CounterValue, float64(g.Stats.Loads.Get()), g.Name())
	ch <- prometheus.MustNewConstMetric(e.groupLoadsDeduped, prometheus.CounterValue, float64(g.Stats.LoadsDeduped.Get()), g.Name())
	ch <- prometheus.MustNewConstMetric(e.groupLocalLoads, prometheus.CounterValue, float64(g.Stats.LocalLoads.Get()), g.Name())
	ch <- prometheus.MustNewConstMetric(e.groupLocalLoadErrs, prometheus.CounterValue, float64(g.Stats.LocalLoadErrs.Get()), g.Name())
	ch <- prometheus.MustNewConstMetric(e.groupServerRequests, prometheus.CounterValue, float64(g.Stats.ServerRequests.Get()), g.Name())
}

func (e *Exporter) collectCacheStats(ch chan<- prometheus.Metric, t groupcache.CacheType, g *groupcache.Group) {
	s := g.CacheStats(t)
	n := g.Name()
	tn := cacheTypeToLabel(t)

	ch <- prometheus.MustNewConstMetric(e.cacheItems, prometheus.GaugeValue, float64(s.Items), n, tn)
	ch <- prometheus.MustNewConstMetric(e.cacheBytes, prometheus.GaugeValue, float64(s.Bytes), n, tn)
	ch <- prometheus.MustNewConstMetric(e.cacheGets, prometheus.CounterValue, float64(s.Gets), n, tn)
	ch <- prometheus.MustNewConstMetric(e.cacheHits, prometheus.CounterValue, float64(s.Hits), n, tn)
	ch <- prometheus.MustNewConstMetric(e.cacheEvictions, prometheus.CounterValue, float64(s.Evictions), n, tn)
}

func cacheTypeToLabel(cacheType groupcache.CacheType) string {
	if cacheType == groupcache.MainCache {
		return "main"
	}
	return "hot"
}
