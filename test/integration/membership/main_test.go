//go:build integration

package membership_integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	tcnetwork "github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/medincident/medincident-command-service/internal/services/membership"
	"github.com/medincident/medincident-command-service/internal/services/zitadel"
)

const (
	zitadelImage     = "ghcr.io/zitadel/zitadel:v4.13.1"
	zitadelMasterKey = "MasterkeyMustBeExact32CharsXXXXX"
)

var (
	testDB     *gorm.DB
	testLogger = zerolog.Nop()

	empSvc *membership.EmployeeService

	zitadelBaseURL string
	zitadelPAT     string

	// Populated in setup after Zitadel boots.
	testUserAliceID string
	testUserBobID   string
	testUserCarolID string

	cleanupFuncs []func()
)

func TestMain(m *testing.M) {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Minute)

	if err := setup(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "setup failed: %v\n", err)
		cancel()
		runCleanup()
		os.Exit(1)
	}
	cancel()

	code := m.Run()
	runCleanup()
	os.Exit(code)
}

func runCleanup() {
	for i := len(cleanupFuncs) - 1; i >= 0; i-- {
		cleanupFuncs[i]()
	}
}

func setup(ctx context.Context) error {
	// 1. Docker network for Zitadel + its own Postgres.
	dockerNet, err := tcnetwork.New(ctx)
	if err != nil {
		return fmt.Errorf("create docker network: %w", err)
	}
	cleanupFuncs = append(cleanupFuncs, func() {
		_ = dockerNet.Remove(context.Background())
	})

	// 2. Postgres dedicated to Zitadel (inside the docker network, alias "db").
	zpg, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:16-alpine",
			ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_DB":       "zitadel",
				"POSTGRES_USER":     "postgres",
				"POSTGRES_PASSWORD": "postgres",
			},
			Networks: []string{dockerNet.Name},
			NetworkAliases: map[string][]string{
				dockerNet.Name: {"db"},
			},
			WaitingFor: wait.ForAll(
				wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
				wait.ForListeningPort("5432/tcp"),
			).WithDeadline(60 * time.Second),
		},
		Started: true,
	})
	if err != nil {
		return fmt.Errorf("start zitadel-postgres: %w", err)
	}
	cleanupFuncs = append(cleanupFuncs, func() { _ = zpg.Terminate(context.Background()) })

	// 3. Command-service's own Postgres, exposed to host, migrations run against it.
	appPG, err := tcpostgres.Run(ctx,
		"postgres:17-alpine",
		tcpostgres.WithDatabase("medincident"),
		tcpostgres.WithUsername("test"),
		tcpostgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	if err != nil {
		return fmt.Errorf("start app-postgres: %w", err)
	}
	cleanupFuncs = append(cleanupFuncs, func() { _ = appPG.Terminate(context.Background()) })

	appDSN, err := appPG.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return fmt.Errorf("get app DSN: %w", err)
	}
	if err := runMigrations(appDSN); err != nil {
		return fmt.Errorf("run migrations: %w", err)
	}
	testDB, err = gorm.Open(postgres.Open(appDSN), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		return fmt.Errorf("open gorm: %w", err)
	}

	// 4. Zitadel container (attached to the same docker network so it can
	//    reach Postgres via the "db" alias).
	zc, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: zitadelImage,
			Cmd: []string{
				"start-from-init",
				"--masterkey", zitadelMasterKey,
				"--tlsMode", "disabled",
				"--config", "/zitadel-config.yaml",
				"--steps", "/zitadel-steps.yaml",
			},
			ExposedPorts: []string{"8080/tcp"},
			Env: map[string]string{
				"ZITADEL_EXTERNALSECURE": "false",
				"ZITADEL_EXTERNALDOMAIN": "localhost",
			},
			Networks: []string{dockerNet.Name},
			Files: []testcontainers.ContainerFile{
				{HostFilePath: "data/zitadel-config.yaml", ContainerFilePath: "/zitadel-config.yaml", FileMode: 0o644},
				{HostFilePath: "data/zitadel-steps.yaml", ContainerFilePath: "/zitadel-steps.yaml", FileMode: 0o644},
			},
			WaitingFor: wait.ForHTTP("/debug/healthz").
				WithPort("8080/tcp").
				WithStartupTimeout(360 * time.Second).
				WithPollInterval(5 * time.Second),
		},
		Started: true,
	})
	if err != nil {
		return fmt.Errorf("start zitadel: %w", err)
	}
	cleanupFuncs = append(cleanupFuncs, func() { _ = zc.Terminate(context.Background()) })

	mappedPort, err := zc.MappedPort(ctx, "8080/tcp")
	if err != nil {
		return fmt.Errorf("get zitadel mapped port: %w", err)
	}
	host, err := zc.Host(ctx)
	if err != nil {
		return fmt.Errorf("get zitadel host: %w", err)
	}
	zitadelBaseURL = fmt.Sprintf("http://%s:%s", host, mappedPort.Port())

	// 5. Extract PAT from Zitadel logs.
	zitadelPAT, err = extractPAT(ctx, zc)
	if err != nil {
		return fmt.Errorf("extract PAT: %w", err)
	}
	fmt.Printf("Zitadel ready at %s (PAT length: %d)\n", zitadelBaseURL, len(zitadelPAT))

	// 6. Create test human users via Zitadel HTTP API.
	testUserAliceID, err = createHumanUser("Alice", "Andersson", "alice@test.local")
	if err != nil {
		return fmt.Errorf("create alice: %w", err)
	}
	testUserBobID, err = createHumanUser("Bob", "Baker", "bob@test.local")
	if err != nil {
		return fmt.Errorf("create bob: %w", err)
	}
	testUserCarolID, err = createHumanUser("Carol", "Carter", "carol@test.local")
	if err != nil {
		return fmt.Errorf("create carol: %w", err)
	}
	fmt.Printf("Test users: alice=%s bob=%s carol=%s\n", testUserAliceID, testUserBobID, testUserCarolID)

	// 7. Wire the zitadel.Service (PAT auth, test-only).
	zsvc, err := zitadel.NewServiceFromPAT(ctx, zitadelBaseURL, zitadelPAT)
	if err != nil {
		return fmt.Errorf("build zitadel service: %w", err)
	}

	empSvc = membership.NewEmployeeService(testDB, zsvc, &testLogger)
	return nil
}

// runMigrations invokes dbmate via go tool against the given DSN.
func runMigrations(dsn string) error {
	cmd := exec.Command("go", "tool", "dbmate",
		"--url", dsn,
		"--migrations-dir", "../../../db/migrations",
		"--no-dump-schema",
		"up",
	)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// extractPAT reads the Zitadel container's logs and pulls out the PAT
// printed on first boot (after the "serviceaccount" marker). Tolerant
// of docker-multiplexed log prefixes.
func extractPAT(ctx context.Context, c testcontainers.Container) (string, error) {
	reader, err := c.Logs(ctx)
	if err != nil {
		return "", fmt.Errorf("get logs: %w", err)
	}
	defer reader.Close()

	raw, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("read logs: %w", err)
	}
	// Strip docker multiplexing headers (non-printable chars).
	clean := strings.Map(func(r rune) rune {
		if r >= 32 || r == '\n' || r == '\r' {
			return r
		}
		return -1
	}, string(raw))

	lines := strings.Split(clean, "\n")
	for i, line := range lines {
		if strings.Contains(strings.ToLower(line), "serviceaccount") {
			for j := i + 1; j < len(lines); j++ {
				candidate := strings.TrimSpace(lines[j])
				if candidate != "" && len(candidate) > 10 {
					return candidate, nil
				}
			}
		}
	}
	return "", fmt.Errorf("PAT not found in container logs (%d lines)", len(lines))
}

// zitadelAPI POSTs a JSON body to a Zitadel API path with the
// extracted PAT as a bearer token.
func zitadelAPI(path string, body any) (map[string]any, error) {
	var r io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		r = bytes.NewReader(data)
	}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, zitadelBaseURL+path, r)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+zitadelPAT)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result map[string]any
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("status %d: %s", resp.StatusCode, respBody)
	}
	if errCode, ok := result["code"]; ok {
		return result, fmt.Errorf("zitadel API error (code=%v): %s", errCode, respBody)
	}
	return result, nil
}

// createHumanUser calls the User v2 AddHumanUser endpoint and returns
// the created user's Zitadel user ID.
func createHumanUser(givenName, familyName, email string) (string, error) {
	resp, err := zitadelAPI("/zitadel.user.v2.UserService/AddHumanUser", map[string]any{
		"username": fmt.Sprintf("%s-%d", strings.ToLower(givenName), time.Now().UnixNano()),
		"profile": map[string]any{
			"givenName":  givenName,
			"familyName": familyName,
		},
		"email": map[string]any{
			"email":      email,
			"isVerified": true,
		},
	})
	if err != nil {
		return "", err
	}
	if id, ok := resp["userId"].(string); ok && id != "" {
		return id, nil
	}
	return "", fmt.Errorf("no userId in response: %v", resp)
}
