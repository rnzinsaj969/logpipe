// Package score provides a priority scoring mechanism for log entries.
//
// A Scorer holds a set of Rules, each mapping a field/value pair to a
// numeric weight. When an entry is evaluated, all matching rule weights
// are summed to produce a final score.
//
// Typical usage:
//
//	s, err := score.New([]score.Rule{
//		{Field: "level",   Value: "error",  Weight: 10},
//		{Field: "service", Value: "payment", Weight: 5},
//	})
//	if err != nil { /* handle */ }
//
//	// Score without mutating the entry:
//	v := s.Score(entry)
//
//	// Or attach the score to entry.Extra["_score"]:
//	enriched := s.Apply(entry)
package score
