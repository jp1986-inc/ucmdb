package rest

// Authentication structure

type RequestToken struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	ClientContext string `json:"clientContext"`
}

type Token struct {
	Token string `json:"token"`
}

// TopologyQuery structure is used to query UCMDB about data model
type AttributeConditions struct {
	Attribute string      `json:"attribute"`
	Operator  string      `json:"operator"`
	Value     interface{} `json:"value"`
}

type LinkCondition struct {
	LinkIdentifier string `json:"linkIdentifier"`
	MinCardinality string `json:"minCardinality"`
	MaxCardinality string `json:"maxCardinality"`
}

type Node struct {
	Type                string                `json:"type"`
	QueryIdentifier     string                `json:"queryIdentifier"`
	Visible             bool                  `json:"visible"`
	IncludeSubtypes     bool                  `json:"includeSubtypes"`
	Layout              []string              `json:"layout"`
	AttributeConditions []AttributeConditions `json:"attributeConditions"`
	LinkConditions      []LinkCondition       `json:"linkConditions"`
}

type Relation struct {
	Type            string   `json:"type"`
	QueryIdentifier string   `json:"queryIdentifier"`
	Visible         bool     `json:"visible"`
	IncludeSubtypes bool     `json:"includeSubtypes"`
	Layout          []string `json:"layout"`
	From            string   `json:"from"`
	To              string   `json:"to"`
}

type TopologyQuery struct {
	Nodes     []Node     `json:"nodes"`
	Relations []Relation `json:"relations"`
}

// TopologyData structure is used (1) to create topology tree and (2) as return type for Topology Query
// DataModelCreated is used as return type for createing topology tree in UCMDB

type DataInConfigurationItem struct {
	UcmdbId    string                 `json:"ucmdbId"`
	GlobalId   string                 `json:"globalId"`
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties"`
}

type DataInRelation struct {
	UcmdbId    string                 `json:"ucmdbId"`
	GlobalId   string                 `json:"globalId"`
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties"`
	End1Id     string                 `json:"end1Id"`
	End2Id     string                 `json:"end2Id"`
}

// structure of toplogy data to create topology tree in ucmdb
type TopologyData struct {
	CIS       []DataInConfigurationItem `json:"cis"`
	Relations []DataInRelation          `json:"relations"`
}

type DataModelChange struct {
	AddedCis   []string `json:"addedCis"`
	RemovedCis []string `json:"removedCis"`
	UpdatedCis []string `json:"updatedCis"`
	IgnoredCis []string `json:"ignoredCis"`
}
