package sdk

import (
	"context"
	"testing"
	"time"

	"github.com/All-The-Vibes/ATV-StarterKit/pkg/monitor"
)

func newTestIntelligence(t *testing.T) *Intelligence {
	t.Helper()
	root := t.TempDir()
	w, err := monitor.NewWatcher(root, monitor.WatcherOptions{})
	if err != nil {
		t.Fatal(err)
	}
	return NewIntelligence(w, IntelligenceOptions{})
}

func TestNewIntelligence_Defaults(t *testing.T) {
	i := newTestIntelligence(t)
	if i.opts.MinInterval != defaultMinInterval {
		t.Errorf("MinInterval = %v, want %v", i.opts.MinInterval, defaultMinInterval)
	}
	if i.opts.QueryBudget != defaultQueryBudget {
		t.Errorf("QueryBudget = %d, want %d", i.opts.QueryBudget, defaultQueryBudget)
	}
}

func TestNewIntelligence_CustomOpts(t *testing.T) {
	root := t.TempDir()
	w, err := monitor.NewWatcher(root, monitor.WatcherOptions{})
	if err != nil {
		t.Fatal(err)
	}
	i := NewIntelligence(w, IntelligenceOptions{
		MinInterval: 10 * time.Second,
		QueryBudget: 100,
	})
	if i.opts.MinInterval != 10*time.Second {
		t.Errorf("MinInterval = %v, want 10s", i.opts.MinInterval)
	}
	if i.opts.QueryBudget != 100 {
		t.Errorf("QueryBudget = %d, want 100", i.opts.QueryBudget)
	}
}

func TestStart_OfflineByDefault(t *testing.T) {
	i := newTestIntelligence(t)
	if err := i.Start(context.Background()); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	if i.IsOnline() {
		t.Error("expected offline mode (SDK not available)")
	}
}

func TestStop(t *testing.T) {
	i := newTestIntelligence(t)
	if err := i.Start(context.Background()); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	if err := i.Stop(); err != nil {
		t.Fatalf("Stop() error = %v", err)
	}
	if i.IsOnline() {
		t.Error("expected offline after Stop()")
	}
}

func TestQuery_OfflineReturnsNil(t *testing.T) {
	i := newTestIntelligence(t)
	if err := i.Start(context.Background()); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	recs, err := i.Query(context.Background())
	if err != nil {
		t.Fatalf("Query() error = %v", err)
	}
	if recs != nil {
		t.Error("expected nil recommendations when offline")
	}
}

func TestExplain_OfflineReturnsError(t *testing.T) {
	i := newTestIntelligence(t)
	if err := i.Start(context.Background()); err != nil {
		t.Fatalf("Start() error = %v", err)
	}
	_, err := i.Explain(context.Background(), "test-rec")
	if err == nil {
		t.Error("expected error when explaining offline")
	}
}

func TestCircuitBreaker(t *testing.T) {
	i := newTestIntelligence(t)
	i.online = true

	// Simulate failures
	for range circuitBreakerLimit {
		i.recordFailure()
	}
	if !i.circuitOpen {
		t.Error("circuit breaker should be open after limit failures")
	}

	// Query should return nil while circuit is open
	recs, err := i.Query(context.Background())
	if err != nil {
		t.Fatalf("Query() error = %v", err)
	}
	if recs != nil {
		t.Error("expected nil when circuit breaker is open")
	}
}

func TestCircuitBreaker_Reset(t *testing.T) {
	i := newTestIntelligence(t)
	i.online = true

	for range circuitBreakerLimit {
		i.recordFailure()
	}
	// Reset via success
	i.recordSuccess()
	if i.consecutiveFailures != 0 {
		t.Errorf("consecutive failures = %d after recordSuccess, want 0", i.consecutiveFailures)
	}
}

func TestRateLimit_MinInterval(t *testing.T) {
	i := newTestIntelligence(t)
	i.online = true
	i.opts.MinInterval = time.Hour // very long interval

	// First query succeeds
	i.lastQuery = time.Time{} // reset
	_, err := i.Query(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	// Second query should be rate-limited
	recs, err := i.Query(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if recs != nil {
		t.Error("expected nil from rate-limited query")
	}
}

func TestRateLimit_Budget(t *testing.T) {
	i := newTestIntelligence(t)
	i.online = true
	i.opts.MinInterval = 0 // disable interval limiting
	i.opts.QueryBudget = 2

	// First two queries consume budget
	if _, err := i.Query(context.Background()); err != nil {
		t.Fatal(err)
	}
	if _, err := i.Query(context.Background()); err != nil {
		t.Fatal(err)
	}

	// Third should be budget-limited
	recs, err := i.Query(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if recs != nil {
		t.Error("expected nil from budget-exhausted query")
	}
}

func TestToolState(t *testing.T) {
	i := newTestIntelligence(t)
	// ToolState returns current live state (SchemaVersion only set after Start)
	state := i.ToolState()
	// Should return zero-value state since watcher not started
	if state.Brainstorms != nil {
		t.Error("expected nil brainstorms from unstarted watcher")
	}
}
