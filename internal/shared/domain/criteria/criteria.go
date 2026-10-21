package criteria

import (
	"fmt"
	"regexp"
	"strings"
)

type Criteria struct {
	Filters []Filter
}

type Filter struct {
	Field    string
	Operator Operator
	Value    any
}

type Operator string

const (
	Equals Operator = "="
	AnyOf  Operator = "= ANY"
)

var fieldNamePattern = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

// ToPostgresql translates the criteria into a SQL WHERE fragment (without the
// "WHERE" keyword) and its positional args, numbering placeholders starting at
// startIndex. Filter.Field is validated against a strict identifier pattern to
// rule out SQL injection, since it is interpolated directly into the query.
func (c Criteria) ToPostgresql(startIndex int) (string, []any, error) {
	if len(c.Filters) == 0 {
		return "", nil, nil
	}

	conditions := make([]string, len(c.Filters))
	args := make([]any, len(c.Filters))
	for i, f := range c.Filters {
		if !fieldNamePattern.MatchString(f.Field) {
			return "", nil, fmt.Errorf("invalid filter field: %q", f.Field)
		}
		if f.Operator == AnyOf {
			conditions[i] = fmt.Sprintf("%s = ANY($%d)", f.Field, startIndex+i)
		} else {
			conditions[i] = fmt.Sprintf("%s %s $%d", f.Field, f.Operator, startIndex+i)
		}
		args[i] = f.Value
	}
	return strings.Join(conditions, " AND "), args, nil
}
