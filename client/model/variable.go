package model

import (
	"fmt"

	"alda.io/client/json"
	log "alda.io/client/logging"
	"github.com/mohae/deepcopy"
)

// GetVariable returns the value of a variable, or an error if the variable in
// undefined.
func (score *Score) GetVariable(name string) ([]ScoreUpdate, error) {
	events, hit := score.Variables[name]

	if !hit {
		return nil, fmt.Errorf("undefined variable: %s", name)
	}

	return events, nil
}

// SetVariable defines the value of a variable.
func (score *Score) SetVariable(name string, value []ScoreUpdate) {
	log.Debug().
		Str("name", name).
		Interface("value", value).
		Msg("Setting variable.")
	score.Variables[name] = value
}

// A VariableDefinition stores a sequence of ScoreUpdates, using the provided
// variable name as a lookup key.
type VariableDefinition struct {
	SourceContext AldaSourceContext
	VariableName  string
	Events        []ScoreUpdate
}

// GetSourceContext implements HasSourceContext.GetSourceContext.
func (vd VariableDefinition) GetSourceContext() AldaSourceContext {
	return vd.SourceContext
}

// JSON implements RepresentableAsJSON.JSON.
func (vd VariableDefinition) JSON() *json.Container {
	events := json.Array()
	for _, event := range vd.Events {
		events.ArrayAppend(event.JSON())
	}

	return json.Object(
		"type", "variable-definition",
		"value", json.Object("events", events),
	)
}

// UpdateScore implements ScoreUpdate.UpdateScore by defining a variable.
func (vd VariableDefinition) UpdateScore(score *Score) error {
	eventValues := []ScoreUpdate{}

	for _, event := range vd.Events {
		eventValue, err := event.VariableValue(score)
		if err != nil {
			return err
		}

		eventValues = append(eventValues, eventValue)
	}

	score.SetVariable(vd.VariableName, eventValues)

	return nil
}

// DurationMs implements ScoreUpdate.DurationMs by returning 0, since a
// variable definition is conceptually instantaneous.
func (VariableDefinition) DurationMs(part *Part) float64 {
	return 0
}

// VariableValue implements ScoreUpdate.VariableValue by capturing the current
// value of each event in the definition.
func (vd VariableDefinition) VariableValue(score *Score) (ScoreUpdate, error) {
	result := deepcopy.Copy(vd).(VariableDefinition)
	result.Events = []ScoreUpdate{}

	for _, event := range vd.Events {
		eventValue, err := event.VariableValue(score)
		if err != nil {
			return nil, err
		}

		result.Events = append(result.Events, eventValue)
	}

	return result, nil
}

// A VariableReference dereferences a stored variable. A variable with the
// provided name is looked up, and assuming that it was previously defined, the
// corresponding sequence of events defined is used to update the score.
type VariableReference struct {
	SourceContext AldaSourceContext
	VariableName  string
}

// GetSourceContext implements HasSourceContext.GetSourceContext.
func (vr VariableReference) GetSourceContext() AldaSourceContext {
	return vr.SourceContext
}

// JSON implements RepresentableAsJSON.JSON.
func (vr VariableReference) JSON() *json.Container {
	return json.Object(
		"type", "variable-reference",
		"value", json.Object("name", vr.VariableName),
	)
}

// UpdateScore implements ScoreUpdate.UpdateScore by looking up a variable and
// (assuming it was previously defined) using the corresponding sequence of
// events to update the score.
func (vr VariableReference) UpdateScore(score *Score) error {
	events, err := score.GetVariable(vr.VariableName)
	if err != nil {
		return err
	}

	log.Debug().
		Str("name", vr.VariableName).
		Interface("events", events).
		Msg("Dereferenced variable.")

	for _, event := range events {
		if err := event.UpdateScore(score); err != nil {
			return err
		}
	}

	return nil
}

// DurationMs implements ScoreUpdate.DurationMs by looking up the sequence of
// events corresponding to the given variable name, and returning the sum of
// the durations of the events.
func (vr VariableReference) DurationMs(part *Part) float64 {
	events, err := part.score.GetVariable(vr.VariableName)
	if err != nil {
		// If the variable is undefined, an error will be thrown when we come back
		// through and look it up again for UpdateScore. So, we can safely ignore
		// the fact that the variable is undefined here and simply return 0.
		return 0
	}

	durationMs := 0.0

	for _, event := range events {
		durationMs += event.DurationMs(part)
	}

	return durationMs
}

// VariableValue implements ScoreUpdate.VariableValue by capturing the current
// value of the referenced variable.
func (vr VariableReference) VariableValue(score *Score) (ScoreUpdate, error) {
	events, err := score.GetVariable(vr.VariableName)
	if err != nil {
		return nil, err
	}

	return EventSequence{Events: events}, nil
}
