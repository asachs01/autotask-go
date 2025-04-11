package autotask

import (
	"encoding/json"
	"testing"
)

func TestEntityQueryParamsBuildQueryString(t *testing.T) {
	// Test simple filter
	filter := ParseFilterString("name='Test'")
	params := NewEntityQueryParams(filter)
	queryString := params.BuildQueryString()

	// Verify query string contains the expected filter
	var queryParams map[string]interface{}
	err := json.Unmarshal([]byte(queryString), &queryParams)
	AssertNil(t, err, "error unmarshaling query string should be nil")

	filterJson, ok := queryParams["filter"]
	AssertTrue(t, ok, "query params should contain filter")
	AssertNotNil(t, filterJson, "filter should not be nil")
}

func TestNewQueryFilter(t *testing.T) {
	// Test simple filter
	filter := NewQueryFilter("name", "eq", "Test")
	AssertEqual(t, "name", filter.Field, "field should match")
	AssertEqual(t, QueryOperator("eq"), filter.Operator, "operator should match")
	AssertEqual(t, "Test", filter.Value, "value should match")
}

func TestNewAndFilterGroup(t *testing.T) {
	// Create filters
	filter1 := NewQueryFilter("name", "eq", "Test")
	filter2 := NewQueryFilter("active", "eq", true)

	// Create AND filter group
	group := NewAndFilterGroup(filter1, filter2)

	// Verify group
	AssertEqual(t, LogicalOperator("and"), group.Operator, "operator should match")
	AssertLen(t, group.Items, 2, "items length should match")
}

func TestNewOrFilterGroup(t *testing.T) {
	// Create filters
	filter1 := NewQueryFilter("name", "eq", "Test")
	filter2 := NewQueryFilter("name", "eq", "Test2")

	// Create OR filter group
	group := NewOrFilterGroup(filter1, filter2)

	// Verify group
	AssertEqual(t, LogicalOperator("or"), group.Operator, "operator should match")
	AssertLen(t, group.Items, 2, "items length should match")
}

func TestNewEntityQueryParams(t *testing.T) {
	// Create query params
	params := NewEntityQueryParams(nil)

	// Set max records explicitly
	params.MaxRecords = 500

	// Verify default values
	AssertEqual(t, 500, params.MaxRecords, "max records should match")

	// Test with fields
	params = params.WithFields("id", "name")
	AssertLen(t, params.Fields, 2, "fields length should match")
	AssertEqual(t, "id", params.Fields[0], "first field should match")
	AssertEqual(t, "name", params.Fields[1], "second field should match")

	// Test with include fields
	params = params.WithIncludeFields("description")
	AssertLen(t, params.IncludeFields, 1, "include fields length should match")
	AssertEqual(t, "description", params.IncludeFields[0], "include field should match")

	// Test with exclude fields
	params = params.WithExcludeFields("createdDate")
	AssertLen(t, params.ExcludeFields, 1, "exclude fields length should match")
	AssertEqual(t, "createdDate", params.ExcludeFields[0], "exclude field should match")

	// Test with max records
	params = params.WithMaxRecords(100)
	AssertEqual(t, 100, params.MaxRecords, "max records should match")

	// Test with page
	params = params.WithPage(2)
	AssertEqual(t, 2, params.Page, "page should match")
}

func TestParseFilterString(t *testing.T) {
	// Test simple filter
	filter := "name='Test'"
	result := ParseFilterString(filter)

	// Verify result is a QueryFilter
	queryFilter, ok := result.(QueryFilter)
	AssertTrue(t, ok, "result should be of type QueryFilter")
	AssertEqual(t, "name", queryFilter.Field, "field should match")
	AssertEqual(t, QueryOperator("eq"), queryFilter.Operator, "operator should match")
	AssertEqual(t, "'Test'", queryFilter.Value, "value should match")

	// Test filter with AND
	filter = "name='Test' AND active=true"
	result = ParseFilterString(filter)

	// Verify result is a FilterGroup
	andGroup, ok := result.(FilterGroup)
	AssertTrue(t, ok, "result should be of type FilterGroup")
	AssertEqual(t, LogicalOperator("and"), andGroup.Operator, "operator should match")
	AssertLen(t, andGroup.Items, 2, "items length should match")

	// Test filter with OR
	filter = "name='Test' OR active=true"
	result = ParseFilterString(filter)

	// Verify result is a FilterGroup
	orGroup, ok := result.(FilterGroup)
	AssertTrue(t, ok, "result should be of type FilterGroup")
	AssertEqual(t, LogicalOperator("or"), orGroup.Operator, "operator should match")
	AssertLen(t, orGroup.Items, 2, "items length should match")

	// Test filter with parentheses
	filter = "(name='Test' OR name='Test2') AND active=true"
	result = ParseFilterString(filter)

	// Verify result is a FilterGroup
	complexGroup, ok := result.(FilterGroup)
	AssertTrue(t, ok, "result should be of type FilterGroup")
	AssertEqual(t, LogicalOperator("and"), complexGroup.Operator, "operator should match")
	AssertLen(t, complexGroup.Items, 2, "items length should match")

	// Verify first item is a FilterGroup
	nestedGroup, ok := complexGroup.Items[0].(FilterGroup)
	AssertTrue(t, ok, "first item should be of type FilterGroup")
	AssertEqual(t, LogicalOperator("or"), nestedGroup.Operator, "nested operator should match")
	AssertLen(t, nestedGroup.Items, 2, "nested items length should match")
}

func TestEntityQueryParamsWithMethods(t *testing.T) {
	// Create a new EntityQueryParams
	filter := NewQueryFilter("Status", OperatorEquals, 1)
	params := NewEntityQueryParams(filter)

	// Test WithFields
	fields := []string{"id", "name", "status"}
	params = params.WithFields(fields...)
	AssertEqual(t, len(fields), len(params.Fields), "fields length should match")
	for i, field := range fields {
		AssertEqual(t, field, params.Fields[i], "field should match")
	}

	// Test WithIncludeFields
	includeFields := []string{"assignedTo", "createdBy"}
	params = params.WithIncludeFields(includeFields...)
	AssertEqual(t, len(includeFields), len(params.IncludeFields), "include fields length should match")
	for i, field := range includeFields {
		AssertEqual(t, field, params.IncludeFields[i], "include field should match")
	}

	// Test WithExcludeFields
	excludeFields := []string{"description", "notes"}
	params = params.WithExcludeFields(excludeFields...)
	AssertEqual(t, len(excludeFields), len(params.ExcludeFields), "exclude fields length should match")
	for i, field := range excludeFields {
		AssertEqual(t, field, params.ExcludeFields[i], "exclude field should match")
	}

	// Test WithMaxRecords
	params = params.WithMaxRecords(100)
	AssertEqual(t, 100, params.MaxRecords, "max records should match")

	// Test WithPage
	params = params.WithPage(2)
	AssertEqual(t, 2, params.Page, "page should match")

	// Verify filter is still present
	AssertEqual(t, 1, len(params.Filter), "filter should have 1 item")

	// Verify filter is the same as the one we created
	qf, ok := params.Filter[0].(QueryFilter)
	AssertTrue(t, ok, "filter should be a QueryFilter")
	AssertEqual(t, "Status", qf.Field, "filter field should match")
	AssertEqual(t, OperatorEquals, qf.Operator, "filter operator should match")
	AssertEqual(t, 1, qf.Value, "filter value should match")
}

func TestParseSimpleFilter(t *testing.T) {
	// Test not equals operator
	filter := parseSimpleFilter("name != 'Test'")
	AssertEqual(t, "name", filter.Field, "field should match")
	AssertEqual(t, QueryOperator("noteq"), filter.Operator, "operator should match")
	AssertEqual(t, "'Test'", filter.Value, "value should match")

	// Test contains operator
	filter = parseSimpleFilter("name contains 'test'")
	AssertEqual(t, "name", filter.Field, "field should match")
	AssertEqual(t, QueryOperator("contains"), filter.Operator, "operator should match")
	AssertEqual(t, "test", filter.Value, "value should match")

	// Test greater than operator
	filter = parseSimpleFilter("count > 10")
	AssertEqual(t, "count", filter.Field, "field should match")
	AssertEqual(t, QueryOperator("greaterThan"), filter.Operator, "operator should match")
	AssertEqual(t, "10", filter.Value, "value should match")

	// Test less than operator
	filter = parseSimpleFilter("count < 10")
	AssertEqual(t, "count", filter.Field, "field should match")
	AssertEqual(t, QueryOperator("lessThan"), filter.Operator, "operator should match")
	AssertEqual(t, "10", filter.Value, "value should match")

	// Test equals operator (default)
	filter = parseSimpleFilter("name = 'Test'")
	AssertEqual(t, "name", filter.Field, "field should match")
	AssertEqual(t, QueryOperator("eq"), filter.Operator, "operator should match")
	AssertEqual(t, "'Test'", filter.Value, "value should match")

	// Test invalid filter
	filter = parseSimpleFilter("invalid filter")
	AssertEqual(t, "", filter.Field, "field should be empty")
	AssertEqual(t, QueryOperator(""), filter.Operator, "operator should be empty")
	AssertNil(t, filter.Value, "value should be nil")
}

func TestParseValueWithOperator(t *testing.T) {
	// Test boolean true
	filter := parseValueWithOperator("active", OperatorEquals, "true")
	AssertEqual(t, "active", filter.Field, "field should match")
	AssertEqual(t, QueryOperator("eq"), filter.Operator, "operator should match")
	AssertEqual(t, true, filter.Value, "value should be true")

	// Test boolean false
	filter = parseValueWithOperator("active", OperatorEquals, "false")
	AssertEqual(t, "active", filter.Field, "field should match")
	AssertEqual(t, QueryOperator("eq"), filter.Operator, "operator should match")
	AssertEqual(t, false, filter.Value, "value should be false")

	// Test string with quotes
	filter = parseValueWithOperator("name", OperatorEquals, "'Test'")
	AssertEqual(t, "name", filter.Field, "field should match")
	AssertEqual(t, QueryOperator("eq"), filter.Operator, "operator should match")
	AssertEqual(t, "'Test'", filter.Value, "value should match")

	// Test string with double quotes
	filter = parseValueWithOperator("name", OperatorEquals, "\"Test\"")
	AssertEqual(t, "name", filter.Field, "field should match")
	AssertEqual(t, QueryOperator("eq"), filter.Operator, "operator should match")
	AssertEqual(t, "\"Test\"", filter.Value, "value should match")

	// Test string without quotes
	filter = parseValueWithOperator("name", OperatorEquals, "Test")
	AssertEqual(t, "name", filter.Field, "field should match")
	AssertEqual(t, QueryOperator("eq"), filter.Operator, "operator should match")
	AssertEqual(t, "Test", filter.Value, "value should match")
}

func TestParseInt(t *testing.T) {
	// Test valid integer
	result := parseInt("123")
	AssertEqual(t, int64(123), result, "result should be 123")

	// Test negative integer
	result = parseInt("-123")
	AssertEqual(t, int64(-123), result, "result should be -123")

	// Test zero
	result = parseInt("0")
	AssertEqual(t, int64(0), result, "result should be 0")

	// Test invalid integer (should return 0)
	result = parseInt("not a number")
	AssertEqual(t, int64(0), result, "result should be 0 for invalid input")
}

func TestParseFloat(t *testing.T) {
	// Test valid float
	result := parseFloat("123.45")
	AssertEqual(t, float64(123.45), result, "result should be 123.45")

	// Test negative float
	result = parseFloat("-123.45")
	AssertEqual(t, float64(-123.45), result, "result should be -123.45")

	// Test integer as float
	result = parseFloat("123")
	AssertEqual(t, float64(123), result, "result should be 123.0")

	// Test zero
	result = parseFloat("0")
	AssertEqual(t, float64(0), result, "result should be 0.0")

	// Test invalid float (should return 0)
	result = parseFloat("not a number")
	AssertEqual(t, float64(0), result, "result should be 0.0 for invalid input")
}
