package sdk

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/All-The-Vibes/ATV-StarterKit/pkg/installstate"
	"github.com/All-The-Vibes/ATV-StarterKit/pkg/monitor"
)

// IntelligenceOptions configures SDK behavior.
type IntelligenceOptions struct {
	MinInterval time.Duration // Minimum interval between queries (default: 30s)
	QueryBudget int           // Max queries per hour (default: 60)
}

// Intelligence is the GitHub Copilot SDK intelligence layer.
// It wraps the SDK client and manages rate limiting, auth detection,
// and graceful degradation to offline mode.
type Intelligence struct {
	watcher *monitor.Watcher
	opts    IntelligenceOptions

	mu          sync.RWMutex
	online      bool
	lastQuery   time.Time
	queriesUsed int
	queryWindow time.Time // start of current hour for budget tracking

	// Circuit breaker
	consecutiveFailures int
	lastFailure         time.Time
	circuitOpen         bool
	circuitOpenUntil    time.Time
}

const (
	defaultMinInterval  = 30 * time.Second
	defaultQueryBudget  = 60
	circuitBreakerLimit = 3
	circuitBreakerCool  = 5 * time.Minute
)

// NewIntelligence creates a new SDK intelligence layer.
func NewIntelligence(watcher *monitor.Watcher, opts IntelligenceOptions) *Intelligence {
	if opts.MinInterval == 0 {
		opts.MinInterval = defaultMinInterval
	}
	if opts.QueryBudget == 0 {
		opts.QueryBudget = defaultQueryBudget
	}
	return &Intelligence{
		watcher:     watcher,
		opts:        opts,
		queryWindow: time.Now().Truncate(time.Hour),
	}
}

// Start attempts to initialize the SDK client and authenticate.
// If auth is unavailable, the intelligence layer runs in offline mode.
func (i *Intelligence) Start(_ context.Context) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	// TODO: Replace with actual GitHub Copilot SDK initialization:
	//   client := copilot.NewClient(nil)
	//   if err := client.Start(); err != nil {
	//       i.online = false
	//       return nil  // offline is valid
	//   }
	//   i.online = true
	//   session, err := client.CreateSession(...)
	//
	// For now, start in offline mode until SDK is available.
	i.online = false
	return nil
}

// Stop cleans up the SDK session and client.
func (i *Intelligence) Stop() error {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.online = false
	return nil
}

// IsOnline reports whether the SDK is authenticated and available.
func (i *Intelligence) IsOnline() bool {
	i.mu.RLock()
	defer i.mu.RUnlock()
	return i.online
}

// Query requests context-aware recommendations from the SDK.
// Returns nil when offline or rate-limited.
func (i *Intelligence) Query(_ context.Context) ([]monitor.SDKRecommendation, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	if !i.online {
		return nil, nil
	}

	// Circuit breaker check
	if i.circuitOpen {
		if time.Now().Before(i.circuitOpenUntil) {
			return nil, nil
		}
		// Reset circuit breaker after cooldown
		i.circuitOpen = false
		i.consecutiveFailures = 0
	}

	// Rate limiting: minimum interval
	if time.Since(i.lastQuery) < i.opts.MinInterval {
		return nil, nil
	}

	// Rate limiting: hourly budget
	now := time.Now()
	currentWindow := now.Truncate(time.Hour)
	if currentWindow.After(i.queryWindow) {
		i.queryWindow = currentWindow
		i.queriesUsed = 0
	}
	if i.queriesUsed >= i.opts.QueryBudget {
		return nil, nil
	}

	// TODO: Replace with actual SDK query:
	//   result, err := i.session.Query(ctx, buildPrompt(i.watcher.State()))
	//   if err != nil {
	//       i.recordFailure()
	//       return nil, err
	//   }
	//   i.recordSuccess()
	//   return parseSDKRecommendations(result), nil

	i.lastQuery = now
	i.queriesUsed++
	return nil, nil
}

// Explain requests a detailed explanation for a specific recommendation.
func (i *Intelligence) Explain(_ context.Context, recID string) (string, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if !i.online {
		return "", fmt.Errorf("SDK offline — explanations unavailable")
	}

	// TODO: Replace with actual SDK explain query
	return fmt.Sprintf("Explanation for %s: SDK integration pending", recID), nil
}

// recordFailure tracks a failed SDK call for circuit breaker logic.
func (i *Intelligence) recordFailure() {
	i.consecutiveFailures++
	i.lastFailure = time.Now()

	if i.consecutiveFailures >= circuitBreakerLimit {
		i.circuitOpen = true
		i.circuitOpenUntil = time.Now().Add(circuitBreakerCool)
	}
}

// recordSuccess resets the circuit breaker failure counter.
func (i *Intelligence) recordSuccess() {
	i.consecutiveFailures = 0
}

// ToolState returns the current live state for SDK tool invocations.
func (i *Intelligence) ToolState() monitor.LiveState {
	return i.watcher.State()
}

// ToolManifest returns the install manifest for SDK tool invocations.
func (i *Intelligence) ToolManifest() (installstate.InstallManifest, error) {
	state := i.watcher.State()
	if !state.InstallSnapshot.HasManifest {
		return installstate.InstallManifest{}, fmt.Errorf("no manifest found")
	}
	return installstate.ReadManifest(i.watcher.Root())
}
