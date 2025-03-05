package autotask

// Company represents an Autotask company
type Company struct {
	ID                      int64  `json:"id,omitempty"`
	CompanyName             string `json:"companyName,omitempty"`
	CompanyNumber           string `json:"companyNumber,omitempty"`
	Phone                   string `json:"phone,omitempty"`
	WebAddress              string `json:"webAddress,omitempty"`
	Active                  bool   `json:"active,omitempty"`
	Address1                string `json:"address1,omitempty"`
	Address2                string `json:"address2,omitempty"`
	City                    string `json:"city,omitempty"`
	State                   string `json:"state,omitempty"`
	PostalCode              string `json:"postalCode,omitempty"`
	Country                 string `json:"country,omitempty"`
	TerritoryID             int64  `json:"territoryID,omitempty"`
	AccountNumber           string `json:"accountNumber,omitempty"`
	TaxRegionID             int64  `json:"taxRegionID,omitempty"`
	ParentCompanyID         int64  `json:"parentCompanyID,omitempty"`
	CompanyType             int    `json:"companyType,omitempty"`
	BillToCompanyID         int64  `json:"billToCompanyID,omitempty"`
	BillToAddress1          string `json:"billToAddress1,omitempty"`
	BillToAddress2          string `json:"billToAddress2,omitempty"`
	BillToCity              string `json:"billToCity,omitempty"`
	BillToState             string `json:"billToState,omitempty"`
	BillToZipCode           string `json:"billToZipCode,omitempty"`
	BillToCountryID         int64  `json:"billToCountryID,omitempty"`
	BillToAttention         string `json:"billToAttention,omitempty"`
	BillToAddressToUse      int    `json:"billToAddressToUse,omitempty"`
	InvoiceMethod           int    `json:"invoiceMethod,omitempty"`
	InvoiceNonContractItems bool   `json:"invoiceNonContractItems,omitempty"`
	InvoiceTemplateID       int64  `json:"invoiceTemplateID,omitempty"`
	QuoteTemplateID         int64  `json:"quoteTemplateID,omitempty"`
	TaxID                   string `json:"taxID,omitempty"`
	TaxExempt               bool   `json:"taxExempt,omitempty"`
	CreatedDate             string `json:"createdDate,omitempty"`
	LastActivityDate        string `json:"lastActivityDate,omitempty"`
	DateStamp               string `json:"dateStamp,omitempty"`
}

// Ticket represents an Autotask ticket
type Ticket struct {
	ID                      int64  `json:"id,omitempty"`
	TicketNumber            string `json:"ticketNumber,omitempty"`
	Title                   string `json:"title,omitempty"`
	Description             string `json:"description,omitempty"`
	Status                  int    `json:"status,omitempty"`
	Priority                int    `json:"priority,omitempty"`
	DueDateTime             string `json:"dueDateTime,omitempty"`
	CreateDate              string `json:"createDate,omitempty"`
	LastActivityDate        string `json:"lastActivityDate,omitempty"`
	CompanyID               int64  `json:"companyID,omitempty"`
	ContactID               int64  `json:"contactID,omitempty"`
	AccountID               int64  `json:"accountID,omitempty"`
	QueueID                 int64  `json:"queueID,omitempty"`
	AssignedResourceID      int64  `json:"assignedResourceID,omitempty"`
	AssignedResourceRoleID  int64  `json:"assignedResourceRoleID,omitempty"`
	TicketType              int    `json:"ticketType,omitempty"`
	IssueType               int    `json:"issueType,omitempty"`
	SubIssueType            int    `json:"subIssueType,omitempty"`
	ServiceLevelAgreementID int64  `json:"serviceLevelAgreementID,omitempty"`
	Source                  int    `json:"source,omitempty"`
	CreatorResourceID       int64  `json:"creatorResourceID,omitempty"`
	CompletedDate           string `json:"completedDate,omitempty"`
}

// Contact represents an Autotask contact
type Contact struct {
	ID               int64  `json:"id,omitempty"`
	FirstName        string `json:"firstName,omitempty"`
	LastName         string `json:"lastName,omitempty"`
	CompanyID        int64  `json:"companyID,omitempty"`
	Email            string `json:"emailAddress,omitempty"`
	Phone            string `json:"phone,omitempty"`
	MobilePhone      string `json:"mobilePhone,omitempty"`
	Title            string `json:"title,omitempty"`
	Active           bool   `json:"active,omitempty"`
	Address1         string `json:"address1,omitempty"`
	Address2         string `json:"address2,omitempty"`
	City             string `json:"city,omitempty"`
	State            string `json:"state,omitempty"`
	PostalCode       string `json:"postalCode,omitempty"`
	Country          string `json:"country,omitempty"`
	PrimaryContact   bool   `json:"isPrimaryContact,omitempty"`
	LastActivityDate string `json:"lastActivityDate,omitempty"`
	CreatedDate      string `json:"createDate,omitempty"`
}

// Resource represents a resource in Autotask
type Resource struct {
	ID        int64  `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Active    bool   `json:"active"`
}

// Response types
type CompanyResponse struct {
	Item Company `json:"item"`
}

type CompanyListResponse struct {
	Items       []Company   `json:"items"`
	PageDetails PageDetails `json:"pageDetails,omitempty"`
}

type TicketResponse struct {
	Item Ticket `json:"item"`
}

type TicketListResponse struct {
	Items       []Ticket    `json:"items"`
	PageDetails PageDetails `json:"pageDetails,omitempty"`
}

type ContactResponse struct {
	Item Contact `json:"item"`
}

type ContactListResponse struct {
	Items       []Contact   `json:"items"`
	PageDetails PageDetails `json:"pageDetails,omitempty"`
}
