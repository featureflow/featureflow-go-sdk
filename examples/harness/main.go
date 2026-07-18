// Local harness: net/http app using the SDK from this repo.
//
//	FEATUREFLOW_SERVER_KEY=srv-env-xxx go run ./examples/harness
//	FEATUREFLOW_BASE_URL=https://beta.featureflow-staging.com FEATUREFLOW_SERVER_KEY=srv-env-xxx go run ./examples/harness
//
// Open http://127.0.0.1:3456/ — or use the JSON API below.
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/featureflow/featureflow-go-sdk/featureflow"
)

// buildUser mirrors what an application would do: attach attributes alongside
// the userId so audience rules (e.g. role equals "pvt_tester") can match.
func buildUser(userId, role, tier string) *featureflow.User {
	if role == "" {
		role = "pvt_tester"
	}
	if tier == "" {
		tier = "gold"
	}
	user, _ := featureflow.NewUserBuilder(userId).
		WithAttribute("role", role).
		WithAttribute("tier", tier).
		Build()
	return user
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func main() {
	apiKey := os.Getenv("FEATUREFLOW_SERVER_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "Set FEATUREFLOW_SERVER_KEY to your Featureflow server API key.")
		fmt.Fprintln(os.Stderr, "Example:")
		fmt.Fprintln(os.Stderr, "  FEATUREFLOW_SERVER_KEY=srv-env-xxxxxxxxx go run ./examples/harness")
		os.Exit(1)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3456"
	}

	// The harness disables the events API by default so it doesn't record
	// evaluation stats; set FEATUREFLOW_DISABLE_EVENTS=false to send events
	// like a full production client.
	disableEvents := os.Getenv("FEATUREFLOW_DISABLE_EVENTS") != "false"

	config := featureflow.Config{
		BaseURL:       os.Getenv("FEATUREFLOW_BASE_URL"), // empty = SDK default https://app.featureflow.io
		DisableEvents: disableEvents,
		WithFeatures: []featureflow.FeatureRegistration{
			featureflow.WithFeature("harness-demo", "off").Build(),
		},
	}

	client, err := featureflow.Client(apiKey, config)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to initialize the Featureflow client:", err)
		os.Exit(1)
	}

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	http.HandleFunc("/api/config", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"baseUrl":       client.Config.BaseURL,
			"disableEvents": client.Config.DisableEvents,
		})
	})

	http.HandleFunc("/api/evaluate", func(w http.ResponseWriter, r *http.Request) {
		featureKey := r.URL.Query().Get("feature")
		if featureKey == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": `Query parameter "feature" is required`})
			return
		}
		userId := r.URL.Query().Get("userId")
		if userId == "" {
			userId = "anonymous"
		}

		user := buildUser(userId, r.URL.Query().Get("role"), r.URL.Query().Get("tier"))
		evaluated := client.Evaluate(featureKey, user)
		writeJSON(w, http.StatusOK, map[string]interface{}{
			"feature": featureKey,
			"user":    user,
			"value":   evaluated.Value(),
			"isOn":    evaluated.IsOn(),
			"isOff":   evaluated.IsOff(),
		})
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "Not found"})
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, indexHTML)
	})

	fmt.Printf("Featureflow Go SDK harness listening on http://127.0.0.1:%s/\n", port)
	fmt.Println("  GET  /health  GET  /api/config")
	fmt.Println("  GET  /api/evaluate?feature=<key>&userId=<id>")
	fmt.Println("  effective baseUrl:", client.Config.BaseURL)
	if disableEvents {
		fmt.Println("  (harness disables events HTTP by default; set FEATUREFLOW_DISABLE_EVENTS=false to enable)")
	}

	if err := http.ListenAndServe("127.0.0.1:"+port, nil); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

const indexHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <title>Featureflow Go SDK — harness</title>
  <style>
    body { font-family: system-ui, sans-serif; max-width: 42rem; margin: 2rem auto; padding: 0 1rem; }
    label { display: block; margin-top: 0.75rem; }
    input { width: 100%; max-width: 24rem; padding: 0.25rem 0.5rem; }
    pre { background: #f4f4f4; padding: 0.75rem; overflow: auto; }
    .err { color: #c00; }
  </style>
</head>
<body>
  <h1>Featureflow Go SDK harness</h1>
  <p>Uses the SDK from this repo via <code>go run ./examples/harness</code>. Set <code>FEATUREFLOW_SERVER_KEY</code>; use <code>FEATUREFLOW_BASE_URL</code> for staging.</p>
  <p>SDK baseUrl: <strong id="baseUrl">…</strong></p>
  <label>Feature key <input id="feature" type="text" value="harness-demo" placeholder="my-feature-key" /></label>
  <label>User id <input id="userId" type="text" value="harness-user" placeholder="user id" /></label>
  <label>Role attribute <input id="role" type="text" value="pvt_tester" placeholder="role" /></label>
  <label>Tier attribute <input id="tier" type="text" value="gold" placeholder="tier" /></label>
  <p><button type="button" id="go">Evaluate</button></p>
  <h2>Result</h2>
  <pre id="out">—</pre>
  <h2>User</h2>
  <pre id="userJson">—</pre>
  <h2>Raw JSON</h2>
  <pre id="raw">—</pre>
  <script>
  (function () {
    function get(path) { return fetch(path).then(function (r) { return r.json(); }); }
    get('/api/config').then(function (c) {
      document.getElementById('baseUrl').textContent = c.baseUrl + (c.disableEvents ? ' (events disabled)' : '');
    }).catch(function () {
      document.getElementById('baseUrl').textContent = '(unavailable)';
    });
    var outEl = document.getElementById('out');
    var userJsonEl = document.getElementById('userJson');
    var rawEl = document.getElementById('raw');
    document.getElementById('go').onclick = function () {
      var q = '/api/evaluate?feature=' + encodeURIComponent(document.getElementById('feature').value.trim()) +
        '&userId=' + encodeURIComponent(document.getElementById('userId').value.trim() || 'anonymous') +
        '&role=' + encodeURIComponent(document.getElementById('role').value.trim()) +
        '&tier=' + encodeURIComponent(document.getElementById('tier').value.trim());
      get(q).then(function (d) {
        outEl.textContent = 'value: ' + d.value + '\n' + 'isOn: ' + d.isOn + ', isOff: ' + d.isOff;
        outEl.className = '';
        userJsonEl.textContent = JSON.stringify(d.user, null, 2);
        rawEl.textContent = JSON.stringify(d, null, 2);
      }).catch(function () {
        outEl.textContent = 'Request failed';
        outEl.className = 'err';
      });
    };
  })();
  </script>
</body>
</html>
`
