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

// Project represents a project in Autotask
type Project struct {
	ID                    int64   `json:"id,omitempty"`
	ProjectName           string  `json:"projectName,omitempty"`
	Description           string  `json:"description,omitempty"`
	CompanyID             int64   `json:"companyID,omitempty"`
	Status                int     `json:"status,omitempty"`
	ProjectNumber         string  `json:"projectNumber,omitempty"`
	Type                  int     `json:"type,omitempty"`
	StartDate             string  `json:"startDate,omitempty"`
	EndDate               string  `json:"endDate,omitempty"`
	EstimatedHours        float64 `json:"estimatedHours,omitempty"`
	ProjectLeadResourceID int64   `json:"projectLeadResourceID,omitempty"`
	CompletedPercentage   float64 `json:"completedPercentage,omitempty"`
	DepartmentID          int64   `json:"departmentID,omitempty"`
	ContractID            int64   `json:"contractID,omitempty"`
	CreatorResourceID     int64   `json:"creatorResourceID,omitempty"`
	CreateDate            string  `json:"createDate,omitempty"`
	LastActivityDate      string  `json:"lastActivityDate,omitempty"`
}

// Task represents a task in Autotask
type Task struct {
	ID                 int64   `json:"id,omitempty"`
	TaskNumber         string  `json:"taskNumber,omitempty"`
	Title              string  `json:"title,omitempty"`
	Description        string  `json:"description,omitempty"`
	Status             int     `json:"status,omitempty"`
	Priority           int     `json:"priority,omitempty"`
	ProjectID          int64   `json:"projectID,omitempty"`
	AssignedResourceID int64   `json:"assignedResourceID,omitempty"`
	StartDate          string  `json:"startDate,omitempty"`
	EndDate            string  `json:"endDate,omitempty"`
	EstimatedHours     float64 `json:"estimatedHours,omitempty"`
	RemainingHours     float64 `json:"remainingHours,omitempty"`
	CompletedDate      string  `json:"completedDate,omitempty"`
	CreateDate         string  `json:"createDate,omitempty"`
	LastActivityDate   string  `json:"lastActivityDate,omitempty"`
	PhaseID            int64   `json:"phaseID,omitempty"`
	TaskType           int     `json:"taskType,omitempty"`
	CreatorResourceID  int64   `json:"creatorResourceID,omitempty"`
}

// TimeEntry represents a time entry in Autotask
type TimeEntry struct {
	ID               int64   `json:"id,omitempty"`
	ResourceID       int64   `json:"resourceID,omitempty"`
	TicketID         int64   `json:"ticketID,omitempty"`
	TaskID           int64   `json:"taskID,omitempty"`
	Type             int     `json:"type,omitempty"`
	DateWorked       string  `json:"dateWorked,omitempty"`
	StartDateTime    string  `json:"startDateTime,omitempty"`
	EndDateTime      string  `json:"endDateTime,omitempty"`
	HoursWorked      float64 `json:"hoursWorked,omitempty"`
	HoursToBill      float64 `json:"hoursToBill,omitempty"`
	SummaryNotes     string  `json:"summaryNotes,omitempty"`
	InternalNotes    string  `json:"internalNotes,omitempty"`
	NonBillable      bool    `json:"nonBillable,omitempty"`
	CreateDate       string  `json:"createDate,omitempty"`
	LastModifiedDate string  `json:"lastModifiedDate,omitempty"`
}

// Contract represents a contract in Autotask
type Contract struct {
	ID                      int64   `json:"id,omitempty"`
	ContractName            string  `json:"contractName,omitempty"`
	ContractNumber          string  `json:"contractNumber,omitempty"`
	CompanyID               int64   `json:"companyID,omitempty"`
	Status                  int     `json:"status,omitempty"`
	ServiceLevelAgreementID int64   `json:"serviceLevelAgreementID,omitempty"`
	StartDate               string  `json:"startDate,omitempty"`
	EndDate                 string  `json:"endDate,omitempty"`
	ContractType            int     `json:"contractType,omitempty"`
	IsDefaultContract       bool    `json:"isDefaultContract,omitempty"`
	SetupFee                float64 `json:"setupFee,omitempty"`
	EstimatedHours          float64 `json:"estimatedHours,omitempty"`
	CreatorResourceID       int64   `json:"creatorResourceID,omitempty"`
	CreateDate              string  `json:"createDate,omitempty"`
	LastActivityDate        string  `json:"lastActivityDate,omitempty"`
}

// ConfigurationItem represents a configuration item in Autotask
type ConfigurationItem struct {
	ID                    int64  `json:"id,omitempty"`
	CompanyID             int64  `json:"companyID,omitempty"`
	ConfigurationItemType int    `json:"configurationItemType,omitempty"`
	ReferenceTitle        string `json:"referenceTitle,omitempty"`
	ReferenceNumber       string `json:"referenceNumber,omitempty"`
	SerialNumber          string `json:"serialNumber,omitempty"`
	InstallDate           string `json:"installDate,omitempty"`
	ProductID             int64  `json:"productID,omitempty"`
	Status                int    `json:"status,omitempty"`
	Location              string `json:"location,omitempty"`
	Active                bool   `json:"active,omitempty"`
	CreateDate            string `json:"createDate,omitempty"`
	LastModifiedDate      string `json:"lastModifiedDate,omitempty"`
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

type ProjectResponse struct {
	Item Project `json:"item"`
}

type ProjectListResponse struct {
	Items       []Project   `json:"items"`
	PageDetails PageDetails `json:"pageDetails,omitempty"`
}

type TaskResponse struct {
	Item Task `json:"item"`
}

type TaskListResponse struct {
	Items       []Task      `json:"items"`
	PageDetails PageDetails `json:"pageDetails,omitempty"`
}

type TimeEntryResponse struct {
	Item TimeEntry `json:"item"`
}

type TimeEntryListResponse struct {
	Items       []TimeEntry `json:"items"`
	PageDetails PageDetails `json:"pageDetails,omitempty"`
}

type ContractResponse struct {
	Item Contract `json:"item"`
}

type ContractListResponse struct {
	Items       []Contract  `json:"items"`
	PageDetails PageDetails `json:"pageDetails,omitempty"`
}

type ConfigurationItemResponse struct {
	Item ConfigurationItem `json:"item"`
}

type ConfigurationItemListResponse struct {
	Items       []ConfigurationItem `json:"items"`
	PageDetails PageDetails         `json:"pageDetails,omitempty"`
}
