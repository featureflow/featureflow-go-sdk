package featureflow

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/cucumber/godog"
)

// applicationTagContextType captures the headers of the requests a client makes
// against a local httptest server — no live network calls. The initial features poll
// happens synchronously inside Client(), but register is a fire-and-forget goroutine,
// so both captures go through buffered channels and the Then steps wait with a timeout.
type applicationTagContextType struct {
	server          *httptest.Server
	featuresHeaders chan http.Header
	registerHeaders chan http.Header
}

var applicationTagContext applicationTagContextType

// anEndpointCapturingRequestHeaders serves both the "features endpoint" and the
// "events endpoint" Givens: one server answers every SDK path and records the
// request headers per concern.
func anEndpointCapturingRequestHeaders() error {
	applicationTagContext.featuresHeaders = make(chan http.Header, 4)
	applicationTagContext.registerHeaders = make(chan http.Header, 4)
	applicationTagContext.server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/sdk/v1/features":
			select {
			case applicationTagContext.featuresHeaders <- r.Header.Clone():
			default:
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte("{}"))
		case "/api/sdk/v1/register":
			select {
			case applicationTagContext.registerHeaders <- r.Header.Clone():
			default:
			}
		default:
			w.WriteHeader(http.StatusOK)
		}
	}))
	return nil
}

func aClientIsInitialisedAgainstItWithApplication(application string) error {
	_, err := Client("sdk-srv-env-test", Config{
		BaseURL:     applicationTagContext.server.URL,
		Application: application,
	})
	return err
}

func aClientWithARegisteredFeatureIsInitialisedAgainstItWithApplication(application string) error {
	_, err := Client("sdk-srv-env-test", Config{
		BaseURL:     applicationTagContext.server.URL,
		Application: application,
		WithFeatures: []FeatureRegistration{
			WithFeature("example-feature", "off").Build(),
		},
	})
	return err
}

func aClientIsInitialisedAgainstItWithNoApplication() error {
	return aClientIsInitialisedAgainstItWithApplication("")
}

// capturedHeaders waits for the capture because register (and any event send) is a
// fire-and-forget goroutine — the client constructor returns before the request lands.
func capturedHeaders(capture chan http.Header, request string) (http.Header, error) {
	select {
	case headers := <-capture:
		return headers, nil
	case <-time.After(2 * time.Second):
		return nil, fmt.Errorf("no %s request was captured within 2s", request)
	}
}

func assertApplicationHeader(capture chan http.Header, request, expected string) error {
	headers, err := capturedHeaders(capture, request)
	if err != nil {
		return err
	}
	if got := headers.Get("X-Featureflow-Application"); got != expected {
		return fmt.Errorf("expected the %s request to carry X-Featureflow-Application %q, got %q", request, expected, got)
	}
	return nil
}

func theCapturedFeaturesRequestHasApplication(expected string) error {
	return assertApplicationHeader(applicationTagContext.featuresHeaders, "features", expected)
}

func theCapturedRegisterRequestHasApplication(expected string) error {
	return assertApplicationHeader(applicationTagContext.registerHeaders, "register", expected)
}

func theCapturedFeaturesRequestHasNoApplicationHeader() error {
	headers, err := capturedHeaders(applicationTagContext.featuresHeaders, "features")
	if err != nil {
		return err
	}
	if got, present := headers["X-Featureflow-Application"]; present {
		return fmt.Errorf("expected no X-Featureflow-Application header, got %q", got)
	}
	return nil
}

func ApplicationTagFeatureContext(ctx *godog.ScenarioContext) {
	ctx.Step(`^a features endpoint capturing request headers$`, anEndpointCapturingRequestHeaders)
	ctx.Step(`^an events endpoint capturing request headers$`, anEndpointCapturingRequestHeaders)
	ctx.Step(`^a client is initialised against it with application "([^"]*)"$`, aClientIsInitialisedAgainstItWithApplication)
	ctx.Step(`^a client with a registered feature is initialised against it with application "([^"]*)"$`, aClientWithARegisteredFeatureIsInitialisedAgainstItWithApplication)
	ctx.Step(`^a client is initialised against it with no application$`, aClientIsInitialisedAgainstItWithNoApplication)
	ctx.Step(`^the captured features request has X-Featureflow-Application "([^"]*)"$`, theCapturedFeaturesRequestHasApplication)
	ctx.Step(`^the captured register request has X-Featureflow-Application "([^"]*)"$`, theCapturedRegisterRequestHasApplication)
	ctx.Step(`^the captured features request has no X-Featureflow-Application header$`, theCapturedFeaturesRequestHasNoApplicationHeader)

	ctx.Before(func(c context.Context, sc *godog.Scenario) (context.Context, error) {
		// A FEATUREFLOW_APPLICATION set on the host must not leak into the
		// no-application scenario.
		os.Unsetenv("FEATUREFLOW_APPLICATION")
		applicationTagContext = applicationTagContextType{}
		return c, nil
	})
	ctx.After(func(c context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		if applicationTagContext.server != nil {
			applicationTagContext.server.Close()
		}
		return c, nil
	})
}

// TestApplicationEnvironmentVariable covers the config-vs-environment precedence the
// shared scenarios leave to each SDK's own config idiom: FEATUREFLOW_APPLICATION
// applies when Config.Application is unset, and a value set in code wins over it.
func TestApplicationEnvironmentVariable(t *testing.T) {
	captured := make(chan http.Header, 4)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/sdk/v1/features" {
			select {
			case captured <- r.Header.Clone():
			default:
			}
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("{}"))
	}))
	defer server.Close()

	t.Setenv("FEATUREFLOW_APPLICATION", "env-app")

	assertNextApplication := func(expected string) {
		t.Helper()
		select {
		case headers := <-captured:
			if got := headers.Get("X-Featureflow-Application"); got != expected {
				t.Fatalf("expected X-Featureflow-Application %q, got %q", expected, got)
			}
		case <-time.After(2 * time.Second):
			t.Fatal("no features request was captured within 2s")
		}
	}

	if _, err := Client("sdk-srv-env-test", Config{BaseURL: server.URL}); err != nil {
		t.Fatal(err)
	}
	assertNextApplication("env-app")

	if _, err := Client("sdk-srv-env-test", Config{BaseURL: server.URL, Application: "code-app"}); err != nil {
		t.Fatal(err)
	}
	assertNextApplication("code-app")
}
