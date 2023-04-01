package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"testing"
	"time"

	"broozkan/postapi/handlers"
	"broozkan/postapi/internal/config"
	"broozkan/postapi/internal/models"
	"broozkan/postapi/internal/repository"
	"broozkan/postapi/internal/services"
	"broozkan/postapi/pkg/server"

	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"
)

func TestMain(m *testing.M) {
	flag.Parse()
	if testing.Short() {
		return
	}
	port, err := nat.NewPort("tcp", "8091")
	if err != nil {
		log.Fatal(err)
	}

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
	_, err = testcontainers.GenericContainer(
		ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		},
	)

	if err != nil {
		log.Fatal(err)
	}

	prepareDB()

	os.Exit(m.Run())
}

func TestIntegration(t *testing.T) {
	t.Run("test given unreachable couchbase url when call main it should give an error", func(t *testing.T) {
		err := os.Setenv("APP_ENV", "/test/invalid-config")
		assert.Nil(t, err)
		err = run()
		assert.NotNil(t, err)
	})

	conf, err := config.New("../.config", "/test/valid-config")
	assert.Nil(t, err)

	logger, _ := zap.NewDevelopment()

	s := initServer(logger, conf)

	time.Sleep(2 * time.Second)

	tests := []struct {
		Name                      string
		Setup                     func()
		Perform                   func() (interface{}, error)
		shouldCheckExpectedResult bool
		expectedResult            interface{}
		customAssertion           func(interface{})
		isErrorNil                bool
	}{
		{
			Name:  "given empty collection when posts/feed request arrived then it should return valid response",
			Setup: func() {},
			Perform: func() (interface{}, error) {
				var respBody []byte
				respBody, err = performRequest(http.MethodGet, "http://localhost:3000/posts/feed", http.NoBody)
				var listPostsResponse models.ListPostsResponse
				_ = json.Unmarshal(respBody, &listPostsResponse)
				return listPostsResponse, err
			},
			customAssertion:           nil,
			shouldCheckExpectedResult: true,
			expectedResult: models.ListPostsResponse{
				Page:       1,
				TotalPages: 0,
			},
			isErrorNil: true,
		},
		{
			Name:  "given create 25 normal with 10th is nsfw and 5 promoted post when /posts request arrived then it should return valid response",
			Setup: func() {},
			Perform: func() (interface{}, error) {
				for i := 1; i <= 25; i++ {
					nsfw := false
					if i == 10 {
						nsfw = true
					}
					body := &models.Post{
						Title:     fmt.Sprintf("Post %d", i),
						Author:    "t2_user123",
						Link:      fmt.Sprintf("https://example.com/post%d", i),
						Subreddit: "testsubreddit",
						Content:   "",
						Score:     i * 10,
						Promoted:  false,
						NSFW:      nsfw,
					}
					_, _ = performRequest(http.MethodPost, "http://localhost:3000/posts", body)
				}

				for i := 1; i <= 5; i++ {
					body := &models.Post{
						Title:     fmt.Sprintf("Post Promoted %d", i),
						Author:    "t2_user123",
						Link:      fmt.Sprintf("https://example.com/post_promoted%d", i),
						Subreddit: "testsubreddit_promoted",
						Content:   "",
						Score:     0,
						Promoted:  true,
						NSFW:      false,
					}
					_, _ = performRequest(http.MethodPost, "http://localhost:3000/posts", body)
				}
				time.Sleep(2 * time.Second)
				return nil, nil
			},
			shouldCheckExpectedResult: false,
			customAssertion:           nil,
			expectedResult:            nil,
			isErrorNil:                true,
		},
		{
			Name:  "given non-empty collection with ad placement have nsfw post when posts/feed request arrived then it should return valid response",
			Setup: func() {},
			Perform: func() (interface{}, error) {
				var respBody []byte
				respBody, err = performRequest(http.MethodGet, "http://localhost:3000/posts/feed", http.NoBody)
				var listPostsResponse models.ListPostsResponse
				_ = json.Unmarshal(respBody, &listPostsResponse)
				return listPostsResponse, err
			},
			customAssertion: func(data interface{}) {
				byteArr, _ := json.Marshal(data)
				var listPostsResponse models.ListPostsResponse
				_ = json.Unmarshal(byteArr, &listPostsResponse)
				assert.Equal(t, 26, len(listPostsResponse.Posts))
			},
			shouldCheckExpectedResult: false,
			expectedResult:            nil,
			isErrorNil:                true,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			test.Setup()
			var res interface{}
			res, err = test.Perform()
			if test.customAssertion != nil {
				test.customAssertion(res)
			}
			if test.shouldCheckExpectedResult {
				assert.Equal(t, test.expectedResult, res)
			}
			assert.Equal(t, test.isErrorNil, err == nil)
		})
	}

	err = Stop(s)
	assert.Nil(t, err)
	time.Sleep(2 * time.Second)

	conf2, err := config.New("../.config", "/test/ads-disabled-config")
	assert.Nil(t, err)

	_ = initServer(logger, conf2)
	time.Sleep(2 * time.Second)

	tests2 := []struct {
		Name                      string
		Setup                     func()
		Perform                   func() (interface{}, error)
		shouldCheckExpectedResult bool
		expectedResult            interface{}
		customAssertion           func(interface{})
		isErrorNil                bool
	}{
		{
			Name:  "given non-empty collection with ad placement have nsfw post when posts/feed request arrived then it should return valid response",
			Setup: func() {},
			Perform: func() (interface{}, error) {
				respBody, err := performRequest(http.MethodGet, "http://localhost:3000/posts/feed", http.NoBody)
				var listPostsResponse models.ListPostsResponse
				_ = json.Unmarshal(respBody, &listPostsResponse)
				return listPostsResponse, err
			},
			customAssertion: func(data interface{}) {
				byteArr, _ := json.Marshal(data)
				var listPostsResponse models.ListPostsResponse
				_ = json.Unmarshal(byteArr, &listPostsResponse)
				assert.Equal(t, 25, len(listPostsResponse.Posts))
			},
			shouldCheckExpectedResult: false,
			expectedResult:            nil,
			isErrorNil:                true,
		},
	}

	for _, test := range tests2 {
		t.Run(test.Name, func(t *testing.T) {
			test.Setup()
			res, err := test.Perform()
			if test.customAssertion != nil {
				test.customAssertion(res)
			}
			if test.shouldCheckExpectedResult {
				assert.Equal(t, test.expectedResult, res)
			}
			assert.Equal(t, test.isErrorNil, err == nil)
		})
	}
}

func initServer(logger *zap.Logger, conf *config.Config) *server.Server {
	couchbaseRepo, err := repository.NewPostRepository(&conf.Couchbase)
	if err != nil {
		log.Fatal("error while initializing couchbase", zap.Error(err))
	}

	postService := services.NewPostService(logger, conf, couchbaseRepo)

	postHandler := handlers.NewPostHandler(logger, conf, postService)

	serverHandlers := []server.Handler{
		postHandler,
	}

	s := server.New(logger, conf.Server, serverHandlers)

	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM)

	go s.Run()
	return &s
}

func Stop(s *server.Server) error {
	err := s.App.Shutdown()
	if err != nil {
		return errors.New("graceful shutdown failed")
	}
	return nil
}

func performRequest(method, url string, body interface{}) ([]byte, error) {
	var reqBody []byte
	if body != nil {
		var err error
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return respBody, nil
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
	createFinanceAPIBucketCommand := "/pools/default/buckets " +
		"-d flushEnabled=1 " +
		"-d name=post " +
		"-d ramQuotaMB=100 " +
		"-d replicaNumber=0 " +
		"-d evictionPolicy=fullEviction " +
		"-d bucketType=couchbase"
	runCurlCommand(newBaseURL, createFinanceAPIBucketCommand)

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
	command := fmt.Sprintf(baseURL, path)
	_, err := exec.Command("/bin/sh", "-c", command).CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
}

func isArmProcessor() bool {
	out, err := exec.Command("uname", "-m").Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(out)) == "arm64"
}
