package symtab

import (
	"github.com/grafana/agent/component/pyroscope/ebpf/ebpfspy/metrics"
	"github.com/grafana/agent/component/pyroscope/ebpf/ebpfspy/symtab/elf"
	lru "github.com/hashicorp/golang-lru/v2"
)

type ElfCache struct {
	buildID2Symbols *lru.Cache[elf.BuildID, SymbolNameResolver]
	roundCache      map[elf.BuildID]SymbolNameResolver
	metrics         *metrics.Metrics
}

func NewElfCache(sz int, metrics *metrics.Metrics) (*ElfCache, error) {
	buildID2Symbols, err := lru.New[elf.BuildID, SymbolNameResolver](sz)
	if err != nil {
		return nil, err
	}
	return &ElfCache{
		buildID2Symbols: buildID2Symbols,
		roundCache:      make(map[elf.BuildID]SymbolNameResolver),
		metrics:         metrics,
	}, nil
}

func (e *ElfCache) GetSymbolsByBuildID(buildID elf.BuildID) SymbolNameResolver {
	if buildID.Empty() {
		return nil
	}
	entry, ok := e.buildID2Symbols.Get(buildID)
	if ok && entry != nil {
		return entry
	}
	entry, ok = e.roundCache[buildID]
	if ok && entry != nil {
		return entry
	}
	return nil
}

func (e *ElfCache) CacheByBuildID(buildID elf.BuildID, v SymbolNameResolver) {
	if buildID.Empty() || v == nil {
		return
	}
	e.buildID2Symbols.Add(buildID, v)
	e.roundCache[buildID] = v
}

func (e *ElfCache) Resize(size int) {
	e.buildID2Symbols.Resize(size)
}

func (e *ElfCache) Cleanup() {
	keys := e.buildID2Symbols.Keys()
	for _, pid := range keys {
		tab, ok := e.buildID2Symbols.Peek(pid)
		if !ok || tab == nil {
			continue
		}
		tab.Cleanup()
	}
	for _, resolver := range e.roundCache {
		resolver.Cleanup()
	}
	e.roundCache = make(map[elf.BuildID]SymbolNameResolver)
}