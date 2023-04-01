package repository_test

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"testing"
	"time"

	"broozkan/postapi/internal/config"
	"broozkan/postapi/internal/repository"

	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"
)

type CouchbaseTestSuite struct {
	suite.Suite
	container       testcontainers.Container
	couchbaseConfig *config.Couchbase
}

const DBUsername = "admin"
const DBPassword = "password"

func (s *CouchbaseTestSuite) SetupSuite() {
	port, err := nat.NewPort("tcp", "8091")
	assert.Nil(s.T(), err)

	image := "couchbase/server:7.0.2"

	if isArmProcessor() {
		image = "couchbase/server:community-aarch64"
	}

	req := testcontainers.ContainerRequest{
		Image: image,
		ExposedPorts: []string{
			"8091:8091/tcp",
			"8092:8092/tcp",
			"8093:8093/tcp",
			"8094:8094/tcp",
			"11207:11207/tcp",
			"11210:11210/tcp",
			"11211:11211/tcp",
		},
		WaitingFor: wait.ForListeningPort(port),
	}

	ctx := context.Background()
	s.container, err = testcontainers.GenericContainer(
		ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})

	assert.Nil(s.T(), err)
	prepareDB()

	s.couchbaseConfig = &config.Couchbase{
		URL:      "localhost",
		Username: DBUsername,
		Password: DBPassword,
		Buckets: []config.BucketConfig{
			{
				Name:               "post",
				CreatePrimaryIndex: false,
				Scopes: []config.ScopeConfig{
					{
						Name: "_default",
						Collections: []config.CollectionConfig{
							{
								Name:               "posts",
								CreatePrimaryIndex: true,
							},
						},
					},
				},
			},
		},
	}
	time.Sleep(15 * time.Second)
	assert.Nil(s.T(), err)
}

func prepareDB() {
	const baseURL = "curl -s -u Administrator:password -X POST http://localhost:8091%s"
	const initializeNodeCommand = "/nodes/self/controller/settings " +
		"-d path=%2Fopt%2Fcouchbase%2Fvar%2Flib%2Fcouchbase%2Fdata " +
		"-d index_path=%2Fopt%2Fcouchbase%2Fvar%2Flib%2Fcouchbase%2Fadata " +
		"-d cbas_path=%2Fopt%2Fcouchbase%2Fvar%2Flib%2Fcouchbase%2Fedata " +
		"-d eventing_path=%2Fopt%2Fcouchbase%2Fvar%2Flib%2Fcouchbase%2Fidata"
	runCurlCommand(baseURL, initializeNodeCommand)

	const startServicesCommand = "/node/controller/setupServices -d services=kv%2Cindex%2Cn1ql%2Cfts"
	runCurlCommand(baseURL, startServicesCommand)

	const setBucketMemoryQuotaCommand = "/pools/default -d memoryQuota=256 -d indexMemoryQuota=256 -d ftsMemoryQuota=256"
	runCurlCommand(baseURL, setBucketMemoryQuotaCommand)

	const createAdminUserCommand = "/settings/web -d port=8091 -d username=admin -d password=password"
	runCurlCommand(baseURL, createAdminUserCommand)

	const newBaseURL = "curl -s -u admin:password -X POST http://localhost:8091%s"
	createPostAPIBucketCommand := "/pools/default/buckets " +
		"-d flushEnabled=1 " +
		"-d name=post " +
		"-d ramQuotaMB=100 " +
		"-d replicaNumber=0 " +
		"-d evictionPolicy=fullEviction " +
		"-d bucketType=couchbase"
	runCurlCommand(newBaseURL, createPostAPIBucketCommand)

	createTestBucketCommand := "/pools/default/buckets " +
		"-d flushEnabled=1 " +
		"-d name=test " +
		"-d ramQuotaMB=100 " +
		"-d replicaNumber=0 " +
		"-d evictionPolicy=fullEviction " +
		"-d bucketType=couchbase"
	runCurlCommand(newBaseURL, createTestBucketCommand)

	setIndexStorageSetting := "/settings/indexes " +
		"-d indexerThreads=0 " +
		"-d logLevel=info " +
		"-d maxRollbackPoints=5 " +
		"-d memorySnapshotInterval=200 " +
		"-d stableSnapshotInterval=5000 " +
		"-d storageMode=forestdb"
	runCurlCommand(newBaseURL, setIndexStorageSetting)
}

func runCurlCommand(baseURL, path string) {
	logger := zap.NewExample()
	command := fmt.Sprintf(baseURL, path)
	output, err := exec.Command("/bin/sh", "-c", command).CombinedOutput()
	if err != nil {
		logger.Debug(string(output))
		log.Fatal(err)
	}
}

func (s *CouchbaseTestSuite) TeardownSuite() {
	err := s.container.Terminate(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}

func TestCouchbase(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	suite.Run(t, new(CouchbaseTestSuite))
}

func (s *CouchbaseTestSuite) TestCouchbase() {
	s.Run("given valid db configuration and bucket when connecting then it should succeed", func() {
		cbConfig := &config.Couchbase{
			URL:      "localhost",
			Username: DBUsername,
			Password: DBPassword,
			Buckets: []config.BucketConfig{
				{
					Name:               "post",
					CreatePrimaryIndex: false,
					Scopes: []config.ScopeConfig{
						{
							Name: "",
							Collections: []config.CollectionConfig{
								{
									Name:               "posts",
									CreatePrimaryIndex: true,
								},
							},
						},
					},
				},
			},
		}
		cb, err := repository.New(cbConfig)

		assert.Nil(s.T(), err)
		assert.NotNil(s.T(), cb)

		err = cb.Cluster.Close(nil)
		assert.Nil(s.T(), err)
	})

	s.Run("given invalid db configuration and bucket when connecting then it should fail", func() {
		cbConfig := config.Couchbase{
			URL:      "localhost",
			Username: "test",
			Password: "testPassword",
			Buckets:  []config.BucketConfig{},
		}

		cb, err := repository.New(&cbConfig)

		assert.Nil(s.T(), cb)
		assert.NotNil(s.T(), err)
	})

	s.Run("given invalid bucket and bucket when connecting then it should fail", func() {
		cbConfig := &config.Couchbase{
			URL:      "localhost",
			Username: "admin",
			Password: "password",
			Buckets: []config.BucketConfig{
				{
					Name: "invalidBucketName",
				},
			},
		}

		cb, err := repository.New(cbConfig)

		assert.NotNil(s.T(), err)
		assert.Nil(s.T(), cb)
	})
}

func isArmProcessor() bool {
	out, err := exec.Command("uname", "-m").Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(out)) == "arm64"
}
