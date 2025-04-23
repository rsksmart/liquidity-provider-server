package mongo_interfaces

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Operator string

const (
	EQ  Operator = "$eq"
	NE  Operator = "$ne"
	GT  Operator = "$gt"
	GTE Operator = "$gte"
	LT  Operator = "$lt"
	LTE Operator = "$lte"
	IN  Operator = "$in"
)

type LogicalOperator string

const (
	AND LogicalOperator = "$and"
	OR  LogicalOperator = "$or"
)

type Condition struct {
	Field    string
	Operator Operator
	Value    interface{}
}

type Criteria struct {
	Conditions      []Condition
	LogicalOperator LogicalOperator
	Skip            int64
	Limit           int64
	Sort            map[string]int
	ProjectFields   []string
}

func NewCriteria() *Criteria {
	return &Criteria{
		Conditions:      make([]Condition, 0),
		LogicalOperator: AND,
		Skip:            0,
		Limit:           0,
		Sort:            make(map[string]int),
	}
}
func (c *Criteria) AddCondition(field string, operator Operator, value interface{}) *Criteria {
	c.Conditions = append(c.Conditions, Condition{Field: field, Operator: operator, Value: value})
	return c
}

func (c *Criteria) SetLogicalOperator(op LogicalOperator) *Criteria {
	c.LogicalOperator = op
	return c
}

func (c *Criteria) SetSkip(skip int64) *Criteria {
	c.Skip = skip
	return c
}

func (c *Criteria) SetLimit(limit int64) *Criteria {
	c.Limit = limit
	return c
}

func (c *Criteria) SetSort(field string, order int) *Criteria {
	c.Sort[field] = order
	return c
}

func (c *Criteria) SetProjection(fields []string) *Criteria {
	c.ProjectFields = fields
	return c
}

func (c *Criteria) buildFieldConditions() map[string]bson.M {
	fieldConditions := make(map[string]bson.M)
	for _, cond := range c.Conditions {
		if cur, ok := fieldConditions[cond.Field]; ok {
			cur[string(cond.Operator)] = cond.Value
			fieldConditions[cond.Field] = cur
		} else {
			if cond.Operator == EQ {
				fieldConditions[cond.Field] = bson.M{"$eq": cond.Value}
			} else {
				fieldConditions[cond.Field] = bson.M{string(cond.Operator): cond.Value}
			}
		}
	}
	return fieldConditions
}

func (c *Criteria) fieldConditionsToAndBson(fieldConditions map[string]bson.M) bson.M {
	result := bson.M{}
	for field, conds := range fieldConditions {
		if eqVal, hasEq := conds["$eq"]; hasEq && len(conds) == 1 {
			result[field] = eqVal
		} else {
			result[field] = conds
		}
	}
	return result
}

func (c *Criteria) fieldConditionsToOrBson(fieldConditions map[string]bson.M) bson.M {
	filters := make([]bson.M, 0, len(fieldConditions))
	for field, conds := range fieldConditions {
		if eqVal, hasEq := conds["$eq"]; hasEq && len(conds) == 1 {
			filters = append(filters, bson.M{field: eqVal})
		} else {
			filters = append(filters, bson.M{field: conds})
		}
	}
	return bson.M{string(c.LogicalOperator): filters}
}

func (c *Criteria) ToBson() bson.M {
	if len(c.Conditions) == 0 {
		return bson.M{}
	}
	if len(c.Conditions) == 1 && c.Conditions[0].Operator == EQ {
		return bson.M{c.Conditions[0].Field: c.Conditions[0].Value}
	}

	fieldConditions := c.buildFieldConditions()

	if c.LogicalOperator == AND {
		return c.fieldConditionsToAndBson(fieldConditions)
	}

	return c.fieldConditionsToOrBson(fieldConditions)
}

func (c *Criteria) GetFindOptions() *options.FindOptions {
	opts := options.Find()

	if c.Skip > 0 {
		opts.SetSkip(c.Skip)
	}

	if c.Limit > 0 {
		opts.SetLimit(c.Limit)
	}

	if len(c.Sort) > 0 {
		sort := bson.D{}
		for field, order := range c.Sort {
			sort = append(sort, bson.E{Key: field, Value: order})
		}
		opts.SetSort(sort)
	}

	if len(c.ProjectFields) > 0 {
		projection := bson.M{}
		for _, field := range c.ProjectFields {
			projection[field] = 1
		}
		opts.SetProjection(projection)
	}

	return opts
}
