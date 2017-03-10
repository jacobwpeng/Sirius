package engine

import "sort"

type ArrayRankUnit RankUnit
type ArrayRankUnitSlice []ArrayRankUnit

func (es ArrayRankUnitSlice) Len() int {
	return len(es)
}

func (es ArrayRankUnitSlice) Swap(i, j int) {
	es[i], es[j] = es[j], es[i]
}

func (es ArrayRankUnitSlice) Less(i, j int) bool {
	return es[i].Key > es[j].Key
}

type ArrayRankEngine struct {
	maxSize uint32
	data    ArrayRankUnitSlice
}

func NewArrayRankEngine(config RankEngineConfig) *ArrayRankEngine {
	return &ArrayRankEngine{maxSize: config.MaxSize}
}

func (e *ArrayRankEngine) Size() uint32 {
	return uint32(e.data.Len())
}

func (e *ArrayRankEngine) Get(id uint64) (bool, uint32, RankUnit) {
	for index := 0; index < e.data.Len(); index++ {
		if e.data[index].ID == id {
			return true, uint32(index), RankUnit(e.data[index])
		}
	}

	return false, 0, RankUnit{}
}

func (e *ArrayRankEngine) GetByRank(pos uint32) (bool, RankUnit) {
	if pos >= e.Size() {
		return false, RankUnit{}
	}
	return true, RankUnit(e.data[pos])
}

func (e *ArrayRankEngine) GetRange(pos, num uint32) []RankUnit {
	if pos >= e.Size() {
		return nil
	}
	n := num
	if pos+num >= e.Size() {
		n = e.Size() - pos
	}
	result := make([]RankUnit, n)
	for i := uint32(0); i < n; i++ {
		result[i] = RankUnit(e.data[pos+i])
	}
	return result
}

func (e *ArrayRankEngine) Update(u RankUnit) (bool, RankUnit) {
	aru := ArrayRankUnit(u)
	exist, index, old := e.Get(u.ID)
	if exist {
		e.data[index] = aru
	} else if e.maxSize != 0 && e.Size() >= e.maxSize {
		// 已经达到最大上限, 淘汰最后一个
		last := e.Size() - 1
		e.data[last] = aru
	} else {
		e.data = append(e.data, aru)
	}
	sort.Stable(e.data)
	return exist, old
}

func (e *ArrayRankEngine) Delete(id uint64) (bool, uint32, RankUnit) {
	exist, pos, u := e.Get(id)
	if exist {
		e.data = append(e.data[:pos], e.data[pos+1:]...)
	}
	return exist, pos, u
}
