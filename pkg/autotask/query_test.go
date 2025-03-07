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
	filterMap, ok := params.Filter[0].(map[string]interface{})
	AssertTrue(t, ok, "filter should be a map")
	AssertEqual(t, "Status", filterMap["field"], "filter field should match")
	AssertEqual(t, "eq", filterMap["op"], "filter operator should match")
	AssertEqual(t, float64(1), filterMap["value"], "filter value should match")
}
