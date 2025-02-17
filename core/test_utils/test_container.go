package test_utils

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/pivonroll/EventStore-Client-Go/core/connection"
	"github.com/stretchr/testify/require"
)

type EventStoreEnvironmentVariable string

const (
	EVENTSTORE_DOCKER_REPOSITORY_ENV    EventStoreEnvironmentVariable = "EVENTSTORE_DOCKER_REPOSITORY"
	EVENTSTORE_DOCKER_TAG_ENV           EventStoreEnvironmentVariable = "EVENTSTORE_DOCKER_TAG"
	EVENTSTORE_DOCKER_PORT_ENV          EventStoreEnvironmentVariable = "EVENTSTORE_DOCKER_PORT"
	EVENTSTORE_MAX_APPEND_SIZE_IN_BYTES EventStoreEnvironmentVariable = "EVENTSTORE_MAX_APPEND_SIZE"
)

const (
	DEFAULT_EVENTSTORE_DOCKER_REPOSITORY = "docker.pkg.github.com/eventstore/eventstore-client-grpc-testdata/eventstore-client-grpc-testdata"
	DEFAULT_EVENTSTORE_DOCKER_TAG        = "21.6.0-buster-slim"
	DEFAULT_EVENTSTORE_DOCKER_PORT       = "2114"
)

type CloseFunc func()

func InitializeGrpcClientWithPrePopulatedDatabase(t *testing.T) (connection.GrpcClient, CloseFunc) {
	return InitializeGrpcClient(t, map[string]string{
		"EVENTSTORE_DB":     "/data/integration-tests",
		"EVENTSTORE_MEM_DB": "false",
	})
}

func InitializeGrpcClient(t *testing.T,
	environmentVariableOverrides map[string]string) (connection.GrpcClient, CloseFunc) {
	container := StartEventStoreInDockerContainer(environmentVariableOverrides)
	clientInstance := createGrpcClientConnectedToContainer(t, container, doNotUseTLS)
	closeFunc := func() {
		container.Close()
		clientInstance.Close()
	}
	return clientInstance, closeFunc
}

func InitializeGrpcClientWithCredentials(t *testing.T,
	username string,
	password string,
	environmentVariableOverrides map[string]string) (connection.GrpcClient, CloseFunc) {
	container := StartEventStoreInDockerContainer(environmentVariableOverrides)
	clientInstance := createGrpcClientConnectedToContainerWithCredentials(t, username, password, container, doNotUseTLS)
	closeFunc := func() {
		container.Close()
		clientInstance.Close()
	}
	return clientInstance, closeFunc
}

func InitializeGrpcClientWithTLS(t *testing.T,
	environmentVariableOverrides map[string]string) (connection.GrpcClient, CloseFunc) {
	container := StartEventStoreInDockerContainer(environmentVariableOverrides)
	clientInstance := createGrpcClientConnectedToContainer(t, container, useTLS)
	closeFunc := func() {
		container.Close()
		clientInstance.Close()
	}
	return clientInstance, closeFunc
}

func InitializeContainerAndGrpcClient(t *testing.T,
	environmentVariableOverrides map[string]string) (*Container, connection.GrpcClient) {
	container := StartEventStoreInDockerContainer(environmentVariableOverrides)
	clientInstance := createGrpcClientConnectedToContainer(t, container, doNotUseTLS)

	return container, clientInstance
}

func StartEventStoreInDockerContainer(environmentVariableOverrides map[string]string) *Container {
	envVariables := readOsEnvironmentVariables(environmentVariableOverrides)
	dockerRunOptions := &dockertest.RunOptions{
		Repository:   envVariables[string(EVENTSTORE_DOCKER_REPOSITORY_ENV)],
		Tag:          envVariables[string(EVENTSTORE_DOCKER_TAG_ENV)],
		ExposedPorts: []string{envVariables[string(EVENTSTORE_DOCKER_PORT_ENV)]},
	}

	otherVariables := envVariables
	delete(otherVariables, string(EVENTSTORE_DOCKER_REPOSITORY_ENV))
	delete(otherVariables, string(EVENTSTORE_DOCKER_TAG_ENV))
	delete(otherVariables, string(EVENTSTORE_DOCKER_PORT_ENV))

	var env []string

	for key, value := range otherVariables {
		env = append(env, fmt.Sprintf("%s=%s", key, value))
	}

	dockerRunOptions.Env = env

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker. Reason: %v", err)
	}
	return getEventStoreDockerContainer(pool, dockerRunOptions)
}

func GetEnableAdminUiWithAtomPub() map[string]string {
	return map[string]string{
		"EVENTSTORE_ENABLE_EXTERNAL_TCP":       "true",
		"EVENTSTORE_ENABLE_ATOM_PUB_OVER_HTTP": "true",
	}
}

// Container ...
type Container struct {
	Endpoint string
	Resource *dockertest.Resource
}

func (container *Container) Close() {
	err := container.Resource.Close()
	if err != nil {
		panic(err)
	}
}

func getEventStoreDockerContainer(pool *dockertest.Pool, options *dockertest.RunOptions) *Container {
	err := setTLSContext(options)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nStarting docker container...")

	resource, err := pool.RunWithOptions(options)
	if err != nil {
		log.Fatalf("Could not start resource. Reason: %v", err)
	}

	fmt.Printf("Started container with id: %v, name: %s\n",
		resource.Container.ID,
		resource.Container.Name)

	endpoint := fmt.Sprintf("localhost:%s", resource.GetPort("2113/tcp"))

	// Disable certificate verification
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	err = pool.Retry(func() error {
		if resource != nil && resource.Container != nil {
			containerInfo, containerError := pool.Client.InspectContainer(resource.Container.ID)
			if containerError == nil && containerInfo.State.Running == false {
				return fmt.Errorf("unexpected exit of container check the container logs for more information, container ID: %v", resource.Container.ID)
			}
		}

		healthCheckEndpoint := fmt.Sprintf("https://%s", endpoint)
		_, err := http.Get(healthCheckEndpoint)
		return err
	})

	if err != nil {
		log.Printf("HealthCheck failed. Reason: %v\n", err)

		closeErr := resource.Close()

		if closeErr != nil {
			log.Fatalf("Failed to close docker resource. Reason: %v", err)
		}
		log.Fatalln("Stopping docker resource")
	}

	return &Container{
		Endpoint: endpoint,
		Resource: resource,
	}
}

func setTLSContext(options *dockertest.RunOptions) error {
	err := verifyCertificatesExist()
	if err != nil {
		return err
	}

	options.Env = append(options.Env, []string{
		"EVENTSTORE_CERTIFICATE_FILE=/etc/eventstore/certs/node/node.crt",
		"EVENTSTORE_CERTIFICATE_PRIVATE_KEY_FILE=/etc/eventstore/certs/node/node.key",
		"EVENTSTORE_TRUSTED_ROOT_CERTIFICATES_PATH=/etc/eventstore/certs/ca",
	}...)

	certsDir, err := getCertificatesDir()
	if err != nil {
		return err
	}
	options.Mounts = []string{
		fmt.Sprintf("%v:/etc/eventstore/certs", certsDir),
	}
	return nil
}

func verifyCertificatesExist() error {
	certs := []string{
		path.Join("ca", "ca.crt"),
		path.Join("ca", "ca.key"),
		path.Join("node", "node.crt"),
		path.Join("node", "node.key"),
	}

	certsDir, err := getCertificatesDir()
	if err != nil {
		return err
	}

	for _, f := range certs {
		if _, err := os.Stat(path.Join(certsDir, f)); os.IsNotExist(err) {
			return fmt.Errorf("could not locate the certificates needed to run EventStoreDB and the tests. Please run 'docker-compose up' for generating the certificates")
		}
	}
	return nil
}

func getCertificatesDir() (string, error) {
	rootDir, err := getRootDir()
	if err != nil {
		return "", err
	}
	return path.Join(rootDir, "certs"), nil
}

func getRootDir() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	currentDir = strings.Replace(currentDir, "\\", "/", -1)
	finalPath := path.Clean(path.Join(currentDir, "../"))

	for {
		if finalPath == "/" {
			break
		}
		_, stdErr := os.Stat(path.Join(finalPath, "go.mod"))
		if os.IsNotExist(stdErr) {
			finalPath = path.Clean(path.Join(finalPath, "../"))
		} else {
			break
		}
	}

	return finalPath, nil
}

type useTLSType bool

const (
	useTLS      useTLSType = true
	doNotUseTLS useTLSType = false
)

func createClientURI(username string, password string, containerEndpoint string, shouldUseTLS useTLSType) string {
	clientURI := fmt.Sprintf("esdb://%s:%s@%s", username, password, containerEndpoint)

	if shouldUseTLS == doNotUseTLS {
		clientURI += "?tlsverifycert=false"
	}

	return clientURI
}

func createGrpcClientConnectedToContainer(t *testing.T,
	container *Container, shouldUseTLS useTLSType) connection.GrpcClient {
	return createGrpcClientConnectedToContainerWithCredentials(t, "admin", "changeit", container, shouldUseTLS)
}

func createGrpcClientConnectedToContainerWithCredentials(t *testing.T,
	username string,
	password string,
	container *Container, shouldUseTLS useTLSType) connection.GrpcClient {
	clientURI := createClientURI(username, password, container.Endpoint, shouldUseTLS)
	fmt.Println("Starting grpc client at:", clientURI)
	config, err := connection.ParseConnectionString(clientURI)
	require.NoError(t, err)

	grpcClient := connection.NewGrpcClient(*config)

	return grpcClient
}

func readOsEnvironmentVariables(override map[string]string) map[string]string {
	variables := map[string]string{}
	variables[string(EVENTSTORE_DOCKER_REPOSITORY_ENV)] = DEFAULT_EVENTSTORE_DOCKER_REPOSITORY
	variables[string(EVENTSTORE_DOCKER_TAG_ENV)] = DEFAULT_EVENTSTORE_DOCKER_TAG
	variables[string(EVENTSTORE_DOCKER_PORT_ENV)] = DEFAULT_EVENTSTORE_DOCKER_PORT

	if value, exists := os.LookupEnv(string(EVENTSTORE_DOCKER_REPOSITORY_ENV)); exists {
		variables[string(EVENTSTORE_DOCKER_REPOSITORY_ENV)] = value
	}

	if value, exists := os.LookupEnv(string(EVENTSTORE_DOCKER_TAG_ENV)); exists {
		variables[string(EVENTSTORE_DOCKER_TAG_ENV)] = value
	}

	if value, exists := os.LookupEnv(string(EVENTSTORE_DOCKER_PORT_ENV)); exists {
		variables[string(EVENTSTORE_DOCKER_PORT_ENV)] = value
	}

	for key, value := range override {
		variables[key] = value
	}

	printAllEnvironmentVariables(variables)
	return variables
}

func printAllEnvironmentVariables(variables map[string]string) {
	for key, value := range variables {
		fmt.Printf("%s=%s\n", key, value)
	}
}
