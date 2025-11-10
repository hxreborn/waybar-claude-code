package ccusage

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"
)

// TokenCounts represents token usage details
type TokenCounts struct {
	InputTokens              int `json:"inputTokens"`
	OutputTokens             int `json:"outputTokens"`
	CacheCreationInputTokens int `json:"cacheCreationInputTokens"`
	CacheReadInputTokens     int `json:"cacheReadInputTokens"`
}

// BurnRate represents usage rate metrics
type BurnRate struct {
	TokensPerMinute              float64 `json:"tokensPerMinute"`
	TokensPerMinuteForIndicator  float64 `json:"tokensPerMinuteForIndicator"`
	CostPerHour                  float64 `json:"costPerHour"`
}

// Projection represents projected usage until reset
type Projection struct {
	TotalTokens      int     `json:"totalTokens"`
	TotalCost        float64 `json:"totalCost"`
	RemainingMinutes int     `json:"remainingMinutes"`
}

// Block represents a single billing block from ccusage
type Block struct {
	ID              string      `json:"id"`
	StartTime       string      `json:"startTime"`
	EndTime         string      `json:"endTime"`
	ActualEndTime   string      `json:"actualEndTime"`
	IsActive        bool        `json:"isActive"`
	IsGap           bool        `json:"isGap"`
	Entries         int         `json:"entries"`
	TokenCounts     TokenCounts `json:"tokenCounts"`
	TotalTokens     int         `json:"totalTokens"`
	CostUSD         float64     `json:"costUSD"`
	Models          []string    `json:"models"`
	BurnRate        BurnRate    `json:"burnRate"`
	Projection      Projection  `json:"projection"`
}

// BlocksResponse represents the JSON response from ccusage blocks command
type BlocksResponse struct {
	Blocks []Block `json:"blocks"`
}

// BlocksData represents simplified data extracted from active block
type BlocksData struct {
	Entries         int
	TotalTokens     int
	InputTokens     int
	OutputTokens    int
	CostUSD         float64
	RemainingMinutes int
	CostPerHour     float64
}

// GetBlocks executes `npx ccusage@latest blocks --active --json --offline`
// and returns the parsed data. Returns zero values on error (graceful degradation for MVP).
func GetBlocks(ctx context.Context) (*BlocksData, error) {
	// Create command with timeout context
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "npx", "ccusage@latest", "blocks", "--active", "--json", "--offline")

	output, err := cmd.Output()
	if err != nil {
		// Return zero values on error (MVP: graceful degradation)
		return &BlocksData{}, fmt.Errorf("failed to execute ccusage: %w", err)
	}

	var response BlocksResponse
	if err := json.Unmarshal(output, &response); err != nil {
		return &BlocksData{}, fmt.Errorf("failed to parse ccusage output: %w", err)
	}

	// Extract active block (should be first one with --active flag)
	if len(response.Blocks) == 0 {
		return &BlocksData{}, fmt.Errorf("no active blocks found")
	}

	block := response.Blocks[0]

	return &BlocksData{
		Entries:          block.Entries,
		TotalTokens:      block.TotalTokens,
		InputTokens:      block.TokenCounts.InputTokens,
		OutputTokens:     block.TokenCounts.OutputTokens,
		CostUSD:          block.CostUSD,
		RemainingMinutes: block.Projection.RemainingMinutes,
		CostPerHour:      block.BurnRate.CostPerHour,
	}, nil
}
