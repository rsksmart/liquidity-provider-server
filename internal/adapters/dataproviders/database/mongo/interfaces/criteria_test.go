package mongo_interfaces

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestCriteria_AddCondition(t *testing.T) {
	c := NewCriteria().
		AddCondition("field1", EQ, "value1").
		AddCondition("field2", GT, 10)

	assert.Len(t, c.Conditions, 2)
	assert.Equal(t, "field1", c.Conditions[0].Field)
	assert.Equal(t, EQ, c.Conditions[0].Operator)
	assert.Equal(t, "value1", c.Conditions[0].Value)
	assert.Equal(t, "field2", c.Conditions[1].Field)
	assert.Equal(t, GT, c.Conditions[1].Operator)
	assert.Equal(t, 10, c.Conditions[1].Value)
}

func TestCriteria_SetLogicalOperator(t *testing.T) {
	c := NewCriteria()
	c.SetLogicalOperator(OR)
	assert.Equal(t, OR, c.LogicalOperator)
}

func TestCriteria_SetSkip(t *testing.T) {
	c := NewCriteria()
	c.SetSkip(123)
	assert.Equal(t, int64(123), c.Skip)
}

func TestCriteria_SetLimit(t *testing.T) {
	c := NewCriteria()
	c.SetLimit(456)
	assert.Equal(t, int64(456), c.Limit)
}

func TestCriteria_SetSort(t *testing.T) {
	c := NewCriteria()
	c.SetSort("fieldA", 1).SetSort("fieldB", -1)
	assert.Equal(t, map[string]int{"fieldA": 1, "fieldB": -1}, c.Sort)
}

func TestCriteria_SetProjection(t *testing.T) {
	c := NewCriteria()
	c.SetProjection([]string{"fieldA", "fieldB"})
	assert.Equal(t, []string{"fieldA", "fieldB"}, c.ProjectFields)
}

func TestCriteria_ToBson(t *testing.T) {
	t.Run("No conditions returns empty bson", func(t *testing.T) {
		c := NewCriteria()
		b := c.ToBson()
		assert.Empty(t, b)
	})

	t.Run("Single EQ condition returns simple bson", func(t *testing.T) {
		c := NewCriteria().
			AddCondition("field1", EQ, 5)
		b := c.ToBson()
		assert.Equal(t, bson.M{"field1": 5}, b)
	})

	t.Run("Multiple conditions with AND logical operator", func(t *testing.T) {
		c := NewCriteria().
			AddCondition("field1", GT, 3).
			AddCondition("field1", LT, 10).
			AddCondition("field2", EQ, "value").
			SetLogicalOperator(AND)
		b := c.ToBson()
		expected := bson.M{
			"field1": bson.M{"$gt": 3, "$lt": 10},
			"field2": "value",
		}
		assert.Equal(t, expected, b)
	})

	t.Run("Multiple conditions with OR logical operator", func(t *testing.T) {
		c := NewCriteria().
			AddCondition("field1", GT, 3).
			AddCondition("field1", LT, 10).
			AddCondition("field2", EQ, "value").
			SetLogicalOperator(OR)
		b := c.ToBson()
		filters := []bson.M{
			{"field1": bson.M{"$gt": 3, "$lt": 10}},
			{"field2": "value"},
		}
		assert.Contains(t, b, string(OR))
		actualFilters, ok := b[string(OR)].([]bson.M)
		assert.True(t, ok)
		assert.ElementsMatch(t, filters, actualFilters)
	})
}

func TestCriteria_GetFindOptions(t *testing.T) {
	c := NewCriteria().
		SetSkip(10).
		SetLimit(20).
		SetSort("fieldA", 1).
		SetSort("fieldB", -1).
		SetProjection([]string{"fieldA", "fieldB"})
	opts := c.GetFindOptions()

	assert.Equal(t, int64(10), *opts.Skip)
	assert.Equal(t, int64(20), *opts.Limit)

	projDoc, ok := opts.Projection.(bson.M)
	assert.True(t, ok)
	expectedProj := bson.M{"fieldA": 1, "fieldB": 1}
	assert.Equal(t, expectedProj, projDoc)
}
