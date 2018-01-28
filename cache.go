package main

type Cache interface {
	Set() error
	Get() interface{}
	Del() error
}

type CacheType int

const (
	MarketCacheType CacheType = 1 << iota
	MatchCacheType
)

