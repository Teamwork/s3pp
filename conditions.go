package s3pp

import (
	"encoding/json"
	"strconv"
)

type Condition interface {
	Name() string
	Value() string
}

// Match returns a condition where the field must match value.
// Fields created from this Condition return value from Value.
func Match(field, value string) Condition {
	return matchCondition{field, value}
}

// StartsWith returns a condition where field must start with value.
// Fields created from this Condition return value from Value.
func StartsWith(field, value string) Condition {
	return startsWithCondition{field, value}
}

// Any returns a condition where field can have any content.
// Fields created from this Condition return an empty string from Value.
func Any(field string) Condition {
	return startsWithCondition{field, ""}
}

// Range returns a condition where field must be int eh range of min and max.
// Fields created from this Condition return min from Value.
func Range(field string, min, max int) Condition {
	return rangeCondition{field, min, max}
}

type startsWithCondition struct {
	name, value string
}

func (c startsWithCondition) Name() string {
	return c.name
}

func (c startsWithCondition) Value() string {
	return c.value
}

func (c startsWithCondition) MarshalJSON() ([]byte, error) {
	return json.Marshal([]string{"starts-with", "$" + c.name, c.value})
}

type matchCondition struct {
	name, value string
}

func (c matchCondition) Name() string {
	return c.name
}

func (c matchCondition) Value() string {
	return c.value
}

func (c matchCondition) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{c.name: c.value})
}

type rangeCondition struct {
	name     string
	min, max int
}

func (c rangeCondition) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{c.name, c.min, c.max})
}

func (c rangeCondition) Name() string {
	return c.name
}

func (c rangeCondition) Value() string {
	return strconv.Itoa(c.min)
}
