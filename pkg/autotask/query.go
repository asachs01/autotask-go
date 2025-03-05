package autotask

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

// QueryFilter represents a single filter condition
type QueryFilter struct {
	Field    string        `json:"field"`
	Operator QueryOperator `json:"op"`
	Value    interface{}   `json:"value,omitempty"`
}

// EntityQueryParams represents the parameters for a query request
type EntityQueryParams struct {
	Filter     []QueryFilter `json:"filter,omitempty"`
	Fields     []string      `json:"fields,omitempty"`
	MaxRecords int           `json:"maxRecords,omitempty"`
}

// BuildQueryString builds the query string for the request
func (p *EntityQueryParams) BuildQueryString() string {
	if p == nil {
		return ""
	}
	return ""
}

// NewQueryFilter creates a new query filter with the given parameters
func NewQueryFilter(field string, operator QueryOperator, value interface{}) QueryFilter {
	return QueryFilter{
		Field:    field,
		Operator: operator,
		Value:    value,
	}
}

// NewEntityQueryParams creates a new query parameters object with the given filters
func NewEntityQueryParams(filters ...QueryFilter) *EntityQueryParams {
	return &EntityQueryParams{
		Filter: filters,
	}
}

// WithFields adds field selection to the query parameters
func (p *EntityQueryParams) WithFields(fields ...string) *EntityQueryParams {
	p.Fields = fields
	return p
}

// WithMaxRecords sets the maximum number of records to return
func (p *EntityQueryParams) WithMaxRecords(max int) *EntityQueryParams {
	p.MaxRecords = max
	return p
}
