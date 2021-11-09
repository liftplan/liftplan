package fto

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/liftplan/liftplan/gear"
)

const (
	namespace = "fto"
)

// FromValues takes a `url.Values` and builds and returns a strategy an error.
func FromValues(v url.Values) (s Strategy, err error) {
	g, err := gear.FromValues(v)
	if err != nil {
		return s, err
	}

	strategy, ok := v[namespace+".strategy"]
	if !ok {
		return s, errors.New("missing strategy in query")
	}

	t, err := StrategyTypeFromString(strategy[0])
	if err != nil {
		return s, err
	}

	movements := []string{"deadlift", "bench press", "overhead press", "squat"}
	m := make([]Movement, len(movements))

	for i := 0; i < len(movements); i++ {
		k := fmt.Sprintf(namespace+".%v", i)
		x, ok := v[k]
		if !ok {
			return s, fmt.Errorf("movement %v not found", k)
		}
		tm, err := strconv.ParseFloat(x[0], 64)
		if err != nil {
			return s, fmt.Errorf("unable to convert %v to float", x[0])
		}

		m[i] = Movement{
			Name:        movements[i],
			TrainingMax: tm,
			Unit:        g.Unit,
		}
	}

	s = Strategy{
		Movements:       m,
		Gear:            g,
		Type:            t,
		Warmup:          isChecked(namespace+".warmup", v),
		JokerSets:       isChecked(namespace+".jokersets", v),
		RecommendPlates: isChecked(namespace+".recplates", v),
	}
	return s, nil
}

func isChecked(key string, vals url.Values) bool {
	v, ok := vals[key]
	if !ok {
		return false
	}
	if len(v) == 0 {
		return false
	}
	return v[0] == "true"
}
