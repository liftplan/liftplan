package fto

import (
	"bytes"
	_ "embed" // used for embeding templates
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/url"

	"github.com/liftplan/liftplan"
	"github.com/liftplan/liftplan/gear"
)

const (
	// TMIncreaseFactor is the percent increase for a
	// training max that is used to calculate a new training max
	// after 3 weeks. In the book the rule of thumb is 5 lbs for upper body
	// and 10 lbs for lower body, but this make the jumps way too high for
	// light lifters, and this facter seems to be a sweet spot for lifting in general.
	TMIncreaseFactor float64 = 0.02

	// MaxTrainingMax is the absolute maximum value that is allowed for any lift.
	MaxTrainingMax float64 = 2000
)

var (
	// ErrInvalidDeloadType represents an invalid DeloadType
	ErrInvalidDeloadType = errors.New("invalid DeloadType")
	// ErrInvalidSetType represents an invalid SetType
	ErrInvalidSetType = errors.New("invalid SetType")
	// ErrInvalidStrategyType represents an invalid StrategyType
	ErrInvalidStrategyType = errors.New("invalid StrategyType")
	//go:embed templates/plan.go.html
	planTemplate string
)

// SetType is used as an ENUM type for Movements
type SetType uint8

const (
	// Working is the primary SetTyp for a 531 workout.
	Working SetType = iota
	// Warmup is an optional SetType to be performed before a working set.
	Warmup
	// Auxiliary is First Set Last, a special type.
	Auxiliary
	// Joker is Joker Sets, which are typically performed after a working set.
	Joker
)

// worker is used as a container for concurrency patterns in calculations
// this allows for massively concurrent processing of calculations.
type worker struct {
	Inc     int
	Set     Set
	Week    Week
	Session Session
	Error   error
}

var stringToSetType = map[string]SetType{
	"Working":   Working,
	"Warmup":    Warmup,
	"Auxiliary": Auxiliary,
	"Joker":     Joker,
}

// String implementation of SetType
func (s SetType) String() string {
	n := []string{
		"Working",
		"Warmup",
		"Auxiliary",
		"Joker",
	}
	if int(s) < len(n) {
		return n[s]
	}
	return ""
}

// MarshalJSON is the json marshaller for SetType
func (s SetType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%v"`, s.String())), nil
}

// UnmarshalJSON is the json unmarshaller for SetType
func (s *SetType) UnmarshalJSON(b []byte) error {
	var st string
	if err := json.Unmarshal(b, &st); err != nil {
		return err
	}
	setType, ok := stringToSetType[st]
	if !ok {
		return ErrInvalidSetType
	}
	*s = setType
	return nil
}

// DeloadType is an ENUM type for the various versions of deload templates
type DeloadType uint

const (
	// Deload1 (5x40%, 5x50%, 5x60%)
	Deload1 DeloadType = iota
	// Deload2 (5x50%, 5x60%, 5x70%)
	Deload2
	// Deload3 (3x65%, 5x75%, 5x85%)
	Deload3
	// Deload4 (10x40%, 8x50%, 6x60%)
	Deload4
	// Deload5 (10x50%, 8x60%, 6x70%)
	Deload5
)

var stringToDeloadType = map[string]DeloadType{
	"deload1": Deload1,
	"deload2": Deload2,
	"deload3": Deload3,
	"deload4": Deload4,
	"deload5": Deload5,
}

// String is the string representation of a deload type
func (d DeloadType) String() string {
	return fmt.Sprintf("deload%v", uint8(d)+1)
}

// MarshalJSON is the json marshaller for DeloadType
func (d DeloadType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, d)), nil
}

// UnmarshalJSON is the json unmarshaller for DeloadType
func (d *DeloadType) UnmarshalJSON(b []byte) error {
	var dt string
	if err := json.Unmarshal(b, &dt); err != nil {
		return err
	}
	deloadType, ok := stringToDeloadType[dt]
	if !ok {
		return ErrInvalidDeloadType
	}
	*d = deloadType
	return nil
}

// StrategyType is a type used to enumerate various strategies for 531.
type StrategyType uint

const (
	// FSLMULTI is a StrategyType for First Set Last (5x8) + Joker Sets
	FSLMULTI StrategyType = iota
	// FSL is a StrategyType for FSL AMRAP in a single set
	FSL
)

// StrategyTypeFromString takes a string and returns a StrategyType and an error
func StrategyTypeFromString(s string) (StrategyType, error) {
	strategyType, ok := stringToStrategyType[s]
	if !ok {
		return 0, ErrInvalidStrategyType
	}
	return strategyType, nil
}

var stringToStrategyType = map[string]StrategyType{
	"FSL Multiple Sets": FSLMULTI,
	"FSL":               FSL,
}

func (s StrategyType) String() string {
	n := []string{"FSL Multiple Sets", "FSL"}
	if int(s) < len(n) {
		return n[s]
	}
	return ""
}

// MarshalJSON is the json marshaller for StrategyType
func (s StrategyType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%v"`, s.String())), nil
}

// UnmarshalJSON is the json unmarshaller for StrategyType
func (s *StrategyType) UnmarshalJSON(b []byte) error {
	var st string
	if err := json.Unmarshal(b, &st); err != nil {
		return err
	}

	strategyType, err := StrategyTypeFromString(st)
	*s = strategyType
	return err
}

// Movement is used to capture the needed info for a 5/3/1 movement, such as Deadlift, Overhead Press, etc.
// It gets a Name, TrainingMax (90% of absolute 1RM) and a unit.
type Movement struct {
	Name        string    `json:"name"`
	TrainingMax float64   `json:"training_max"`
	Unit        gear.Unit `json:"unit"`
	Calculated  bool      `json:"calculated"`
}

func (m Movement) percentOfMax(weight float64, unit gear.Unit) (float64, error) {
	w, err := gear.ConvertFromTo(weight, unit, m.Unit)
	return w / m.TrainingMax * 100, err
}

// Set outlines how a Movement is performed, it includes a percent to calculate from Movement.TrainingMax,
// rep count, it also gets an AMRAP (As Many Reps As Possible) bool, which indicates if a set is a "plus" round. For instance
// a set with 5 reps simply should perform 5 reps, while a set of 5 reps + AMRAP = true, should perform
// a minimum of 5 reps, but should attempt for As Many Reps As Possible(AMRAP). The Type is the SetType for the movement.
type Set struct {
	Movement Movement  `json:"movement"`
	Percent  float64   `json:"percentage"`
	Reps     uint      `json:"reps"`
	AMRAP    bool      `json:"amrap"`
	Type     SetType   `json:"type"`
	Weight   float64   `json:"weight,omitempty"`
	Plates   []float64 `json:"plates,omitempty"`
}

func (s *Set) calculate(recommendPlates bool, g gear.Gear) error {
	floor, err := s.Movement.floor(g)
	if err != nil {
		return err
	}
	if s.Percent < floor {
		s.Percent = floor
	}

	// this error isn't possible because we check both
	// units in .floor()
	max, _ := gear.ConvertFromTo(s.Movement.TrainingMax, s.Movement.Unit, g.Unit)

	c := max * s.Percent / 100
	rounded, err := g.Round(c)
	if err != nil {
		return err
	}
	(*s).Weight = rounded

	if recommendPlates {
		rec, _ := g.Recommend(c)
		(*s).Plates = rec
	}

	return nil
}

// Session represents a slice of Set. This is all sets that should be performed in a workout.
type Session []Set

func (s *Session) setMovement(m Movement) {
	for i := 0; i < len(*s); i++ {
		(*s)[i].Movement = m
	}
}

// floor takes gear as an argument and returns the percentage of the max as the floor based
// on the bar.
func (m Movement) floor(g gear.Gear) (float64, error) {
	min, err := g.Min()
	if err != nil {
		return 0, err
	}
	return m.percentOfMax(min, g.Unit)
}

// addWarmup adds a maximum of 5 warmup rounds to a Session.
// It calculates this by finding the first working set and the
// floor, which is the percent of Training Max that represents
// an empty bar. It then works backwards by subtracting 10 percent
// from the first working set until either the floor is reached
// or the 5 round limit is reached.
func (s *Session) addWarmup(g gear.Gear) error {
	// get the first working set
	f, err := s.first(Working)
	if err != nil {
		return err
	}
	// get the floor as a percentage of Training Max
	// this the the weight of the empty bar.
	floor, err := f.Movement.floor(g)
	if err != nil {
		return err
	}
	var warmupSets []Set
	f.Reps = 5
	f.Type = Warmup

	// subtract 10 from activeSet until either floor is reached
	// or 4 rounds of warmups are reached, and add the barSet
	for {
		f.Percent = f.Percent - 10
		if f.Percent <= floor || len(warmupSets) == 4 {
			barSet := f
			barSet.Percent = floor
			barSet.Reps = 10
			warmupSets = append([]Set{barSet}, warmupSets...)
			break
		}
		warmupSets = append([]Set{f}, warmupSets...)
	}
	*s = append(warmupSets, (*s)...)

	return err

}

func (s *Session) addJokers() error {
	l, err := s.last(Working)
	if err != nil {
		return err
	}
	l.Type = Joker
	l.AMRAP = false
	for l.Percent < 120 {
		l.Percent = l.Percent + 5
		(*s) = append((*s), l)
	}
	return nil
}

// addFSLMulti adds 5 founds of 8 for whatever the first
// working set of a Session was for a movement.
func (s *Session) addFSLMulti() error {
	f, err := s.first(Working)
	if err != nil {
		return err
	}
	f.Type = Auxiliary
	f.Reps = 8
	f.AMRAP = false
	sets := []Set{f, f, f, f, f}
	(*s) = append((*s), sets...)
	return nil
}

func (s *Session) addFSL() error {
	f, err := s.first(Working)
	if err != nil {
		return err
	}
	f.Type = Auxiliary
	f.Reps = 10
	f.AMRAP = true
	(*s) = append((*s), f)
	return nil
}

func (s *Session) calculate(recommendPlates bool, g gear.Gear) error {
	l := len(*s)
	c := make(chan worker, l)
	for i, set := range *s {
		go func(i int, set Set, g gear.Gear) {
			err := set.calculate(recommendPlates, g)
			c <- worker{
				Inc:   i,
				Set:   set,
				Error: err,
			}
		}(i, set, g)
	}
	for i := 0; i < l; i++ {
		sw := <-c
		if sw.Error != nil {
			return sw.Error
		}
		(*s)[sw.Inc] = sw.Set
	}
	return nil
}

// first takes a SetType and returns the first instance
// from a session, or an error if it finds no match for
// the SetType.
func (s Session) first(t SetType) (set Set, err error) {
	for i := 0; i < len(s); i++ {
		if s[i].Type == t {
			return s[i], err
		}
	}
	return set, fmt.Errorf("no set found matching: %v", t)
}

// SetTypeIndex takes a SetType and returns the first
// match as an int for its position in the Session slice
// its returns a -1 if it finds nothing.
func (s Session) SetTypeIndex(t SetType) int {
	for i := 0; i < len(s); i++ {
		if s[i].Type == t {
			return i
		}
	}
	return -1
}

// CountSetType takes a SetType an returns count as and int
// of that type. This is used to see how many sets of a specific
// type exists in a Session
func (s Session) CountSetType(t SetType) int {
	counter := 0
	for i := 0; i < len(s); i++ {
		if s[i].Type == t {
			counter++
		}
	}
	return counter
}

// last takes a SetType and returns the last instance
// from a session or an error if it finds no match for
// the SetType.
func (s Session) last(t SetType) (set Set, err error) {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i].Type == t {
			return s[i], err
		}
	}
	return set, fmt.Errorf("no set found matching: %v", t)
}

func (s Session) copy() Session {
	l := len(s)
	sess := make(Session, l)
	for i, item := range s {
		set := Set{
			Movement: Movement{
				Name:        item.Movement.Name,
				TrainingMax: item.Movement.TrainingMax,
				Unit:        item.Movement.Unit,
				Calculated:  item.Movement.Calculated,
			},
			Percent: item.Percent,
			Reps:    item.Reps,
			AMRAP:   item.AMRAP,
			Type:    item.Type,
		}
		sess[i] = set
	}
	return sess
}

// A Week is a slice of sessions as well as a Deload boolean
type Week struct {
	Sessions        []Session `json:"sessions"`
	Deload          bool      `json:"deload,omitempty"`
	RecommendPlates bool      `json:"recommend_plates,omitempty"`
}

// DisplayNumber shows the week number in human readable form from index
func (w Week) DisplayNumber(n int) int {
	return n + 1
}

func (w *Week) calculate(recommendPlates, warmup, jokersets bool, aux StrategyType, g gear.Gear) error {
	l := len(w.Sessions)
	c := make(chan worker, l)
	(*w).RecommendPlates = recommendPlates
	for i, sess := range w.Sessions {
		go func(i int, sess Session, g gear.Gear) {
			if warmup {
				if err := sess.addWarmup(g); err != nil {
					c <- worker{Error: err}
					return
				}
			}
			if jokersets {
				if err := sess.addJokers(); err != nil {
					c <- worker{Error: err}
					return
				}
			}
			{
				var err error
				switch aux {
				case FSLMULTI:
					err = sess.addFSLMulti()
				case FSL:
					err = sess.addFSL()
				default:
					err = errors.New("strategy type not implemented")
				}
				if err != nil {
					c <- worker{Error: err}
					return
				}
			}

			err := sess.calculate(recommendPlates, g)
			c <- worker{
				Inc:     i,
				Session: sess,
				Error:   err,
			}
		}(i, sess, g)
	}
	for i := 0; i < l; i++ {
		sw := <-c
		if sw.Error != nil {
			return sw.Error
		}
		w.Sessions[sw.Inc] = sw.Session
	}
	return nil
}

// Progression is a slice of Weeks.
type Progression []Week

func (p *Progression) calculate(recommendPlates, warmup, jokersets bool, aux StrategyType, g gear.Gear) error {
	l := len(*p)
	c := make(chan worker, l)

	for i, w := range *p {
		go func(i int, w Week, g gear.Gear) {
			err := w.calculate(recommendPlates, warmup, jokersets, aux, g)
			c <- worker{
				Inc:   i,
				Week:  w,
				Error: err,
			}
		}(i, w, g)
	}
	for i := 0; i < l; i++ {
		ww := <-c
		if ww.Error != nil {
			return ww.Error
		}
		(*p)[ww.Inc] = ww.Week
	}
	return nil
}

// Strategy is a 531 struct that containers movements, gear, and strategy
type Strategy struct {
	Movements       []Movement   `json:"movements"`
	Gear            gear.Gear    `json:"gear"`
	Type            StrategyType `json:"type"`
	Deload          DeloadType   `json:"deload_type"`
	Warmup          bool         `json:"warmup"`
	JokerSets       bool         `json:"joker_sets"`
	RecommendPlates bool         `json:"recommend_plates"`
}

//Plan implements a liftplan.Plan
func (s Strategy) Plan(f liftplan.Format) ([]byte, error) {
	p := newProgression(s.Movements, s.Deload)
	err := p.calculate(s.RecommendPlates, s.Warmup, s.JokerSets, s.Type, s.Gear)
	if err != nil {
		return nil, err
	}

	switch f {
	case liftplan.JSON:
		return json.Marshal(p)
	case liftplan.HTML:
		var b bytes.Buffer
		t, err := template.New("plan").Parse(planTemplate)
		if err != nil {
			return nil, err
		}
		if err := t.Execute(&b, p); err != nil {
			return nil, err
		}
		return b.Bytes(), nil
	default:
		return nil, errors.New("liftplan format not implemented")
	}
}

// Values conforms to the Valuer interface and is part of the LiftPlanner interface
func (s Strategy) Values() (url.Values, error) {
	vals, err := gear.ToValues(s.Gear)
	if err != nil {
		return vals, err
	}
	vals.Set("method", namespace)
	vals.Set(namespace+".warmup", fmt.Sprintf("%v", s.Warmup))
	vals.Set(namespace+".jokersets", fmt.Sprintf("%v", s.JokerSets))
	vals.Set(namespace+".recplates", fmt.Sprintf("%v", s.RecommendPlates))
	vals.Set(namespace+".strategy", s.Type.String())
	// TODO: we need to make sure these movements are exported properly

	for i, m := range s.Movements {
		a, err := gear.ConvertFromTo(m.TrainingMax, m.Unit, s.Gear.Unit)
		if err != nil {
			return vals, err
		}
		vals.Set(namespace+fmt.Sprintf(".%v", i), fmt.Sprintf("%.2f", a))
	}
	return vals, nil
}

// NewProgression Generates a new 7 week progression from a set of movements
func newProgression(movements []Movement, d DeloadType) Progression {
	l := 7
	p := make([]Week, l)
	c := make(chan worker, l)

	for i, w := range p {
		go func(i int, w Week) {
			w.Deload = i == 6
			for _, m := range movements {
				var sess Session
				tmIncrease := m.TrainingMax * TMIncreaseFactor
				if w.Deload {
					m.TrainingMax = m.TrainingMax + tmIncrease
					m.Calculated = true
					sess = deloadTemplate[d].copy()
				} else {
					item := i
					if item > 2 {
						item -= 3
						m.TrainingMax = m.TrainingMax + tmIncrease
						m.Calculated = true
					}
					sess = workingSetTemplate[item].copy()
				}
				// ensure that we never calculate anything larger than absolute max.
				if m.TrainingMax > MaxTrainingMax {
					m.TrainingMax = MaxTrainingMax
				}
				sess.setMovement(m)
				w.Sessions = append(w.Sessions, sess)
			}
			c <- worker{Inc: i, Week: w}
		}(i, w)
	}
	for i := 0; i < l; i++ {
		ww := <-c
		p[ww.Inc] = ww.Week
	}

	return p
}
