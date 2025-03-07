package autotask

import (
	"encoding/json"
	"fmt"
	"strings"
)

// QueryOperator represents the type of query operation
type QueryOperator string

const (
	OperatorEquals         QueryOperator = "eq"
	OperatorNotEquals      QueryOperator = "noteq"
	OperatorBeginsWith     QueryOperator = "beginsWith"
	OperatorEndsWith       QueryOperator = "endsWith"
	OperatorContains       QueryOperator = "contains"
	OperatorNotContains    QueryOperator = "notContains"
	OperatorGreaterThan    QueryOperator = "greaterThan"
	OperatorLessThan       QueryOperator = "lessThan"
	OperatorGreaterOrEqual QueryOperator = "greaterOrEqual"
	OperatorLessOrEqual    QueryOperator = "lessOrEqual"
	OperatorIn             QueryOperator = "in"
	OperatorNotIn          QueryOperator = "notIn"
	OperatorIsNull         QueryOperator = "isNull"
	OperatorIsNotNull      QueryOperator = "isNotNull"
)

// LogicalOperator represents the type of logical operation (AND, OR)
type LogicalOperator string

const (
	LogicalOperatorAnd LogicalOperator = "and"
	LogicalOperatorOr  LogicalOperator = "or"
)

// QueryFilter represents a single filter condition
type QueryFilter struct {
	Field    string        `json:"field,omitempty"`
	Operator QueryOperator `json:"op,omitempty"`
	Value    interface{}   `json:"value,omitempty"`
}

// FilterGroup represents a group of filters with a logical operator
type FilterGroup struct {
	Operator LogicalOperator `json:"op"`
	Items    []interface{}   `json:"items"` // Can contain QueryFilter or nested FilterGroup
}

// EntityQueryParams represents the parameters for a query request
type EntityQueryParams struct {
	Filter        interface{} `json:"filter,omitempty"` // Can be QueryFilter, FilterGroup, or []QueryFilter
	Fields        []string    `json:"fields,omitempty"`
	MaxRecords    int         `json:"maxRecords,omitempty"`
	IncludeFields []string    `json:"includeFields,omitempty"`
	ExcludeFields []string    `json:"excludeFields,omitempty"`
}

// BuildQueryString builds the query string for the request
func (p *EntityQueryParams) BuildQueryString() string {
	if p == nil {
		return ""
	}

	jsonBytes, err := json.Marshal(p)
	if err != nil {
		return ""
	}

	return string(jsonBytes)
}

// NewQueryFilter creates a new query filter with the given parameters
func NewQueryFilter(field string, operator QueryOperator, value interface{}) QueryFilter {
	return QueryFilter{
		Field:    field,
		Operator: operator,
		Value:    value,
	}
}

// NewAndFilterGroup creates a new filter group with AND logic
func NewAndFilterGroup(items ...interface{}) FilterGroup {
	return FilterGroup{
		Operator: LogicalOperatorAnd,
		Items:    items,
	}
}

// NewOrFilterGroup creates a new filter group with OR logic
func NewOrFilterGroup(items ...interface{}) FilterGroup {
	return FilterGroup{
		Operator: LogicalOperatorOr,
		Items:    items,
	}
}

// NewEntityQueryParams creates a new query parameters object with the given filter
func NewEntityQueryParams(filter interface{}) *EntityQueryParams {
	return &EntityQueryParams{
		Filter: filter,
	}
}

// WithFields adds field selection to the query parameters
func (p *EntityQueryParams) WithFields(fields ...string) *EntityQueryParams {
	p.Fields = fields
	return p
}

// WithIncludeFields adds fields to include in the response
func (p *EntityQueryParams) WithIncludeFields(fields ...string) *EntityQueryParams {
	p.IncludeFields = fields
	return p
}

// WithExcludeFields adds fields to exclude from the response
func (p *EntityQueryParams) WithExcludeFields(fields ...string) *EntityQueryParams {
	p.ExcludeFields = fields
	return p
}

// WithMaxRecords sets the maximum number of records to return
func (p *EntityQueryParams) WithMaxRecords(max int) *EntityQueryParams {
	p.MaxRecords = max
	return p
}

// ParseFilterString parses a filter string into a query filter or filter group
// Supports basic conditions and complex conditions with AND/OR operators
// Examples:
//   - "Status=1"
//   - "Status!=5"
//   - "Name contains 'test'"
//   - "Status=1 AND AssignedTo=123"
//   - "Status=1 OR Status=2"
//   - "Status=1 AND (AssignedTo=123 OR AssignedTo=456)"
func ParseFilterString(filterStr string) interface{} {
	if filterStr == "" {
		return nil
	}

	// Check if the filter contains logical operators
	if strings.Contains(strings.ToUpper(filterStr), " AND ") || strings.Contains(strings.ToUpper(filterStr), " OR ") {
		return parseComplexFilter(filterStr)
	}

	// Handle simple filter
	return parseSimpleFilter(filterStr)
}

// parseSimpleFilter parses a simple filter string (no logical operators)
func parseSimpleFilter(filterStr string) QueryFilter {
	// Handle not equals
	if strings.Contains(filterStr, "!=") {
		parts := strings.Split(filterStr, "!=")
		if len(parts) == 2 {
			field := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			return parseValueWithOperator(field, OperatorNotEquals, value)
		}
	}

	// Handle contains
	if strings.Contains(strings.ToLower(filterStr), " contains ") {
		parts := strings.Split(strings.ToLower(filterStr), " contains ")
		if len(parts) == 2 {
			field := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			// Remove quotes if present
			value = strings.Trim(value, "'\"")
			return NewQueryFilter(field, OperatorContains, value)
		}
	}

	// Handle greater than
	if strings.Contains(filterStr, ">") {
		parts := strings.Split(filterStr, ">")
		if len(parts) == 2 {
			field := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			return parseValueWithOperator(field, OperatorGreaterThan, value)
		}
	}

	// Handle less than
	if strings.Contains(filterStr, "<") {
		parts := strings.Split(filterStr, "<")
		if len(parts) == 2 {
			field := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			return parseValueWithOperator(field, OperatorLessThan, value)
		}
	}

	// Handle equals (default)
	parts := strings.Split(filterStr, "=")
	if len(parts) == 2 {
		field := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		return parseValueWithOperator(field, OperatorEquals, value)
	}

	// If we can't parse it, return a dummy filter
	return QueryFilter{}
}

// parseValueWithOperator parses a value and applies the appropriate type conversion
func parseValueWithOperator(field string, operator QueryOperator, valueStr string) QueryFilter {
	// Handle boolean values
	if valueStr == "true" {
		return NewQueryFilter(field, operator, true)
	} else if valueStr == "false" {
		return NewQueryFilter(field, operator, false)
	}

	// Handle numeric values
	if isNumeric(valueStr) {
		// Try to parse as int first
		if !strings.Contains(valueStr, ".") {
			return NewQueryFilter(field, operator, parseInt(valueStr))
		}
		// If it has a decimal point, parse as float
		return NewQueryFilter(field, operator, parseFloat(valueStr))
	}

	// Default to string
	return NewQueryFilter(field, operator, valueStr)
}

// parseComplexFilter parses a complex filter string with logical operators
func parseComplexFilter(filterStr string) interface{} {
	// Handle parentheses for nested conditions
	if strings.Contains(filterStr, "(") && strings.Contains(filterStr, ")") {
		// This is a complex case with nested conditions
		// For now, we'll implement a simple version that handles one level of nesting
		return parseNestedFilter(filterStr)
	}

	// Check if it's an AND condition
	if strings.Contains(strings.ToUpper(filterStr), " AND ") {
		parts := strings.Split(strings.ToUpper(filterStr), " AND ")
		filters := make([]interface{}, len(parts))
		for i, part := range parts {
			filters[i] = parseSimpleFilter(strings.TrimSpace(part))
		}
		return NewAndFilterGroup(filters...)
	}

	// Check if it's an OR condition
	if strings.Contains(strings.ToUpper(filterStr), " OR ") {
		parts := strings.Split(strings.ToUpper(filterStr), " OR ")
		filters := make([]interface{}, len(parts))
		for i, part := range parts {
			filters[i] = parseSimpleFilter(strings.TrimSpace(part))
		}
		return NewOrFilterGroup(filters...)
	}

	// If we can't parse it as a complex filter, try as a simple filter
	return parseSimpleFilter(filterStr)
}

// parseNestedFilter parses a filter string with nested conditions
func parseNestedFilter(filterStr string) interface{} {
	// This is a simplified implementation that handles basic nesting
	// A full implementation would need a proper parser

	// Check if it's an AND condition with nested OR
	if strings.Contains(strings.ToUpper(filterStr), " AND ") {
		parts := strings.Split(strings.ToUpper(filterStr), " AND ")
		filters := make([]interface{}, 0, len(parts))

		for _, part := range parts {
			part = strings.TrimSpace(part)
			if strings.HasPrefix(part, "(") && strings.HasSuffix(part, ")") {
				// This is a nested condition
				nestedPart := part[1 : len(part)-1] // Remove parentheses
				if strings.Contains(strings.ToUpper(nestedPart), " OR ") {
					// Parse as OR group
					nestedParts := strings.Split(strings.ToUpper(nestedPart), " OR ")
					nestedFilters := make([]interface{}, len(nestedParts))
					for i, np := range nestedParts {
						nestedFilters[i] = parseSimpleFilter(strings.TrimSpace(np))
					}
					filters = append(filters, NewOrFilterGroup(nestedFilters...))
				} else {
					// Parse as simple filter
					filters = append(filters, parseSimpleFilter(nestedPart))
				}
			} else {
				// Parse as simple filter
				filters = append(filters, parseSimpleFilter(part))
			}
		}

		return NewAndFilterGroup(filters...)
	}

	// Check if it's an OR condition with nested AND
	if strings.Contains(strings.ToUpper(filterStr), " OR ") {
		parts := strings.Split(strings.ToUpper(filterStr), " OR ")
		filters := make([]interface{}, 0, len(parts))

		for _, part := range parts {
			part = strings.TrimSpace(part)
			if strings.HasPrefix(part, "(") && strings.HasSuffix(part, ")") {
				// This is a nested condition
				nestedPart := part[1 : len(part)-1] // Remove parentheses
				if strings.Contains(strings.ToUpper(nestedPart), " AND ") {
					// Parse as AND group
					nestedParts := strings.Split(strings.ToUpper(nestedPart), " AND ")
					nestedFilters := make([]interface{}, len(nestedParts))
					for i, np := range nestedParts {
						nestedFilters[i] = parseSimpleFilter(strings.TrimSpace(np))
					}
					filters = append(filters, NewAndFilterGroup(nestedFilters...))
				} else {
					// Parse as simple filter
					filters = append(filters, parseSimpleFilter(nestedPart))
				}
			} else {
				// Parse as simple filter
				filters = append(filters, parseSimpleFilter(part))
			}
		}

		return NewOrFilterGroup(filters...)
	}

	// If we can't parse it as a nested filter, try as a simple filter
	return parseSimpleFilter(filterStr)
}

// Helper functions for type conversion
func isNumeric(s string) bool {
	// Check if the string is a valid number
	_, err := fmt.Sscanf(s, "%f", &struct{}{})
	return err == nil
}

func parseInt(s string) int64 {
	var i int64
	fmt.Sscanf(s, "%d", &i)
	return i
}

func parseFloat(s string) float64 {
	var f float64
	fmt.Sscanf(s, "%f", &f)
	return f
}
