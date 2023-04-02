package config

type (
	Couchbase struct {
		URL      string
		Username string
		Password string
		Buckets  []BucketConfig
	}

	BucketConfig struct {
		Name               string
		CreatePrimaryIndex bool
		Scopes             []ScopeConfig
	}

	ScopeConfig struct {
		Name        string
		Collections []CollectionConfig
	}

	CollectionConfig struct {
		Name               string
		CreatePrimaryIndex bool
		FieldIndexes       []string
	}
)
