package posseg

import (
	"cmp"
	"fmt"
	"slices"
)

type probState struct {
	prob  float64
	state uint16
}

func (ps probState) String() string {
	return fmt.Sprintf("(%v: %f)", ps.state, ps.prob)
}

func viterbi(obs []rune) []tag {
	obsLength := len(obs)
	V := make([]map[uint16]float64, obsLength)
	V[0] = make(map[uint16]float64)
	memPath := make([]map[uint16]uint16, obsLength)
	memPath[0] = make(map[uint16]uint16)
	ys := charStateTab.get(obs[0]) // default is all_states
	for _, y := range ys {
		V[0][y] = probEmit[y].get(obs[0]) + probStart[y]
		memPath[0][y] = 0
	}

	for t := 1; t < obsLength; t++ {
		// Pre-allocate with estimated capacity
		prevStates := make([]uint16, 0, len(memPath[t-1]))
		for x := range memPath[t-1] {
			if len(probTrans[x]) > 0 {
				prevStates = append(prevStates, x)
			}
		}

		// Use map[string]struct{} as set (zero memory overhead)
		prevStatesExpectNext := make(map[uint16]struct{}, len(prevStates)*2)
		for _, x := range prevStates {
			for y := range probTrans[x] {
				prevStatesExpectNext[y] = struct{}{}
			}
		}

		tmpObsStates := charStateTab.get(obs[t])

		obsStates := make([]uint16, 0, len(tmpObsStates))
		for _, state := range tmpObsStates {
			if _, ok := prevStatesExpectNext[state]; ok {
				obsStates = append(obsStates, state)
			}
		}

		if len(obsStates) == 0 {
			for key := range prevStatesExpectNext {
				obsStates = append(obsStates, key)
			}
		}
		if len(obsStates) == 0 {
			obsStates = probTransKeys
		}

		memPath[t] = make(map[uint16]uint16)
		V[t] = make(map[uint16]float64)

		for _, y := range obsStates {
			var max probState
			for i, y0 := range prevStates {
				ps := probState{
					prob:  V[t-1][y0] + probTrans[y0].Get(y) + probEmit[y].get(obs[t]),
					state: y0,
				}
				if i == 0 || ps.prob > max.prob || (ps.prob == max.prob && ps.state > max.state) {
					max = ps
				}
			}
			V[t][y] = max.prob
			memPath[t][y] = max.state
		}
	}

	// Build last slice with pre-allocated capacity
	last := make([]probState, 0, len(memPath[obsLength-1]))
	for y := range memPath[obsLength-1] {
		last = append(last, probState{prob: V[obsLength-1][y], state: y})
	}

	// Use slices.SortFunc instead of sort.Interface
	slices.SortFunc(last, func(a, b probState) int {
		if a.prob == b.prob {
			return cmp.Compare(b.state, a.state) // Reverse state order
		}
		return cmp.Compare(b.prob, a.prob) // Descending prob
	})

	state := last[0].state
	route := make([]tag, obsLength)

	for i := obsLength - 1; i >= 0; i-- {
		route[i] = tag(state)
		state = memPath[i][state]
	}
	return route
}
