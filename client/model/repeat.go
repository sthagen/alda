package model

import (
	"alda.io/client/json"
	"github.com/mohae/deepcopy"
)

// A Repeat expression repeats an event a number of times.
type Repeat struct {
	SourceContext AldaSourceContext
	Event         ScoreUpdate
	Times         int32
}

// GetSourceContext implements HasSourceContext.GetSourceContext.
func (repeat Repeat) GetSourceContext() AldaSourceContext {
	return repeat.SourceContext
}

// JSON implements RepresentableAsJSON.JSON.
func (repeat Repeat) JSON() *json.Container {
	return json.Object(
		"type", "repeat",
		"value", json.Object(
			"event", repeat.Event.JSON(),
			"times", repeat.Times,
		),
	)
}

// UpdateScore implements ScoreUpdate.UpdateScore by repeatedly updating the
// score with an event a specified number of times.
func (repeat Repeat) UpdateScore(score *Score) error {
	for repetition := int32(1); repetition <= repeat.Times; repetition++ {
		for _, part := range score.CurrentParts {
			part.currentRepetition = repetition
		}

		if err := score.Update(repeat.Event); err != nil {
			return err
		}
	}

	return nil
}

// DurationMs implements ScoreUpdate.DurationMs by returning the total duration
// of the event being repeated the specified number of times.
func (repeat Repeat) DurationMs(part *Part) float64 {
	durationMs := 0.0

	for repetition := int32(1); repetition <= repeat.Times; repetition++ {
		part.currentRepetition = repetition
		durationMs += repeat.Event.DurationMs(part)
	}

	return durationMs
}

// VariableValue implements ScoreUpdate.VariableValue by returning a version of
// the repeat where the value of the event to be repeated is captured.
func (repeat Repeat) VariableValue(score *Score) (ScoreUpdate, error) {
	result := deepcopy.Copy(repeat).(Repeat)

	eventValue, err := repeat.Event.VariableValue(score)

	if err != nil {
		return nil, err
	}

	result.Event = eventValue

	return result, nil
}
