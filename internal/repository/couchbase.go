package repository

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"broozkan/postapi/internal/config"

	"github.com/couchbase/gocb/v2"
)

type Couchbase struct {
	Cluster        *gocb.Cluster
	PostBucket     *gocb.Bucket
	PostCollection string
}

const timeout = 20 * time.Second

func New(couchbaseConfig *config.Couchbase) (*Couchbase, error) {
	cluster, err := gocb.Connect(couchbaseConfig.URL, gocb.ClusterOptions{
		Username: couchbaseConfig.Username,
		Password: couchbaseConfig.Password,
	})
	if err != nil {
		return nil, err
	}

	//nolint:govet
	if err := cluster.WaitUntilReady(timeout, nil); err != nil {
		return nil, err
	}

	bucketMap := make(map[string]*gocb.Bucket)
	for _, bucketConfig := range couchbaseConfig.Buckets {
		bucket := cluster.Bucket(bucketConfig.Name)

		err = bucket.WaitUntilReady(timeout, nil)
		if err != nil {
			return nil, err
		}

		err = createScopes(cluster, bucket, bucketConfig.Scopes)
		if err != nil {
			return nil, err
		}

		if bucketConfig.CreatePrimaryIndex {
			err = createPrimaryIndex(cluster, bucket.Name(), "_default", "_default")
			if err != nil {
				return nil, err
			}
		}

		bucketMap[bucketConfig.Name] = bucket
	}

	c := &Couchbase{
		Cluster:        cluster,
		PostBucket:     bucketMap[couchbaseConfig.Buckets[0].Name],
		PostCollection: couchbaseConfig.Buckets[0].Scopes[0].Collections[0].Name,
	}

	err = c.createIndexIfNotExist("postIndex", "post", "id")
	if err != nil {
		return nil, err
	}
	return c, nil
}

func createScopes(cluster *gocb.Cluster, bucket *gocb.Bucket, scopes []config.ScopeConfig) error {
	existingScopes, err := bucket.Collections().GetAllScopes(nil)
	if err != nil {
		return err
	}

	existingScopesMap := make(map[string]bool, len(existingScopes))
	for _, scope := range existingScopes {
		existingScopesMap[scope.Name] = true
	}

	for _, scopeConf := range scopes {
		if scopeConf.Name != "" && !existingScopesMap[scopeConf.Name] {
			opts := &gocb.CreateScopeOptions{
				RetryStrategy: gocb.NewBestEffortRetryStrategy(nil),
			}

			if err := bucket.Collections().CreateScope(scopeConf.Name, opts); err != nil {
				return err
			}
		}
		if err := createCollections(cluster, bucket, scopeConf.Name, scopeConf.Collections); err != nil {
			return err
		}
	}
	return nil
}

func createCollections(cluster *gocb.Cluster, bucket *gocb.Bucket, scopeName string,
	collections []config.CollectionConfig) error {
	if scopeName == "" {
		scopeName = bucket.DefaultCollection().ScopeName()
	}

	for _, collection := range collections {
		spec := gocb.CollectionSpec{
			Name:      collection.Name,
			ScopeName: scopeName,
		}

		opts := &gocb.CreateCollectionOptions{
			RetryStrategy: gocb.NewBestEffortRetryStrategy(nil),
		}

		if err := bucket.Collections().CreateCollection(spec, opts); err != nil && !errors.Is(err,
			gocb.ErrCollectionExists) {
			return err
		}

		if collection.CreatePrimaryIndex {
			err := createPrimaryIndex(cluster, bucket.Name(), scopeName, collection.Name)
			if err != nil {
				return err
			}
		}

		for _, field := range collection.FieldIndexes {
			if err := createCollectionFieldIndex(
				cluster,
				bucket.Name(),
				scopeName,
				collection.Name,
				field,
			); err != nil {
				return err
			}
		}
	}
	return nil
}

func createCollectionFieldIndex(cluster *gocb.Cluster, bucket, scope, collection, field string) error {
	keyspaceName := fmt.Sprintf("default:%s.%s.%s", bucket, scope, collection)
	hasKeyspace, err := hasKeyspace(cluster, keyspaceName)
	if err != nil {
		return err
	}
	if !hasKeyspace {
		return fmt.Errorf("cannot create field index: could not find keyspace %s err: %v", keyspaceName, err)
	}

	fieldWithoutDots := strings.ReplaceAll(field, ".", "-")
	indexName := fmt.Sprintf(
		"%s_%s_%s_%s_FieldIndex",
		bucket,
		scope,
		collection,
		fieldWithoutDots,
	)

	namespace := fmt.Sprintf(
		"%s`.`%s`.`%s",
		bucket,
		scope,
		collection,
	)

	opts := &gocb.CreateQueryIndexOptions{
		IgnoreIfExists: true,
		RetryStrategy:  gocb.NewBestEffortRetryStrategy(nil),
	}

	fieldForNestedIndexes := strings.ReplaceAll(field, ".", "`.`")
	err = cluster.QueryIndexes().CreateIndex(
		namespace,
		indexName,
		[]string{fieldForNestedIndexes},
		opts,
	)
	if err != nil {
		return err
	}
	return nil
}

func hasKeyspace(cluster *gocb.Cluster, keyspaceName string) (bool, error) {
	type Keyspace struct {
		Path string `json:"path"`
	}
	type KeyspaceResult struct {
		Keyspaces Keyspace `json:"keyspaces"`
	}

	var keyspace KeyspaceResult
	const retryCount = 5
	for i := 0; i < retryCount; i++ {
		result, err := cluster.Query(fmt.Sprintf("SELECT * FROM system:keyspaces where keyspaces.`path`='%s'",
			keyspaceName), &gocb.QueryOptions{Adhoc: true})
		if err != nil {
			return false, err
		}

		if result.Next() {
			err = result.Row(&keyspace)
			if err != nil {
				return false, err
			}
		}
		if keyspace.Keyspaces.Path != "" {
			return true, nil
		}

		time.Sleep(1 * time.Second)
	}

	return false, nil
}

func (c *Couchbase) isIndexExists(index string) (bool, error) {
	res, err := c.Cluster.Query(fmt.Sprintf(`SELECT * FROM system:indexes where name='%s'`, index), &gocb.QueryOptions{Adhoc: true})
	if err != nil {
		return false, err
	}

	var indexResponse map[string]interface{}
	err = res.One(&indexResponse)

	if err != nil && err.Error() == "no result was available" {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return len(indexResponse) == 1, nil
}

func (c *Couchbase) createIndex(name, bucket, on string) error {
	_, err := c.Cluster.Query("CREATE INDEX "+name+" ON `"+bucket+"`(`"+on+"`)", nil)
	if err != nil {
		return err
	}
	return nil
}

func createPrimaryIndex(c *gocb.Cluster, bucket, scope, collection string) error {
	queryString := fmt.Sprintf("CREATE PRIMARY INDEX ON `default`:`%s`.`%s`.`%s`", bucket, scope, collection)
	_, err := c.Query(queryString, nil)
	if err != nil && strings.Contains(err.Error(), "already exists") {
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}

func (c *Couchbase) createIndexIfNotExist(name, bucket, on string) error {
	isExists, err := c.isIndexExists(name)
	if err != nil {
		return err
	}

	if isExists {
		fmt.Printf("index %s with name %s on %s bucket already exists, skipping...\n", on, name, bucket)
		return nil
	}

	return c.createIndex(name, bucket, on)
}
