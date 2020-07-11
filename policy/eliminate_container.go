package policy

import (
	"sort"
	"sync"
)

type EliminateContainer struct {
	shortTermContainer sync.Map
	longTermContainer  sync.Map
}


type Entry struct {
	key   string
	value int64
}
type EntrySlice []*Entry

func (e EntrySlice) Len() int           { return len(e) }
func (e EntrySlice) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }
func (e EntrySlice) Less(i, j int) bool { return e[i].value < e[j].value }

type AddPolicyFunc func(*sync.Map, string, int64)

func (e *EliminateContainer) addKey(key string, ttl int64, shortTermPolicy AddPolicyFunc, longTermPolicy AddPolicyFunc) {
	if ttl != -1 {
		shortTermPolicy(&e.shortTermContainer, key, ttl)
	} else {
		longTermPolicy(&e.shortTermContainer, key, ttl)
	}
}

func (e *EliminateContainer) removeKey(key string) {
	if _, ok := e.shortTermContainer.Load(key); ok {
		e.shortTermContainer.Delete(key)
	} else {
		e.longTermContainer.Delete(key)
	}
}

func (e *EliminateContainer) accessKey(key string, ttl int64, shortTermPolicy AddPolicyFunc, longTermPolicy AddPolicyFunc) {
	e.addKey(key, ttl, shortTermPolicy, longTermPolicy)
}


func (e *EliminateContainer) eliminate(eliminateRate float32) []string {
	removeEntries := make(EntrySlice, 0)
	shortTermEntries := make(EntrySlice, 0)
	longTermEntries := make(EntrySlice, 0)

	e.shortTermContainer.Range(func(key, value interface{}) bool {
		shortTermEntries = append(shortTermEntries, &Entry{
			key:   key.(string),
			value: value.(int64),
		})
		return true
	})

	e.longTermContainer.Range(func(key, value interface{}) bool {
		longTermEntries = append(longTermEntries, &Entry{
			key:   key.(string),
			value: value.(int64),
		})
		return true
	})


	shortTermEntriesLen := len(shortTermEntries)
	longTermEntriesLen := len(longTermEntries)
	totalEntriesLen := shortTermEntriesLen + longTermEntriesLen

	sort.Sort(shortTermEntries)
	sort.Sort(longTermEntries)

	rate1 := float32(shortTermEntriesLen / totalEntriesLen)
	if rate1 < eliminateRate {
		removeEntries = append(removeEntries, shortTermEntries...)
		removeLongTermKeyLen := int((eliminateRate - rate1) * float32(totalEntriesLen))
		if removeLongTermKeyLen >= longTermEntriesLen {
			removeEntries = append(removeEntries, longTermEntries...)
		} else {
			removeEntries = append(removeEntries, longTermEntries[0:removeLongTermKeyLen]...)
		}
	} else {
		removeKeyLen := int((rate1) * float32(totalEntriesLen))
		if removeKeyLen >= shortTermEntriesLen {
			removeEntries = append(removeEntries, shortTermEntries[0:shortTermEntriesLen]...)
		} else {
			removeEntries = append(removeEntries, shortTermEntries[0:removeKeyLen]...)
		}

	}

	keys := make([]string, len(removeEntries))
	for _, entry := range removeEntries {
		keys = append(keys, entry.key)
	}
	return keys
}

func (e *EliminateContainer) needEliminate(maxNumber int) bool {
	shortTermContainerLen := 0
	longTermContainerLen := 0
	e.shortTermContainer.Range(func(key, value interface{}) bool {
		shortTermContainerLen++
		return true
	})

	e.longTermContainer.Range(func(key, value interface{}) bool {
		longTermContainerLen++
		return true
	})


	return shortTermContainerLen + longTermContainerLen >= maxNumber
}
