package ccusage

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
)

type TokenCounts struct {
	InputTokens              int `json:"inputTokens"`
	OutputTokens             int `json:"outputTokens"`
	CacheCreationInputTokens int `json:"cacheCreationInputTokens"`
	CacheReadInputTokens     int `json:"cacheReadInputTokens"`
}

type BurnRate struct {
	TokensPerMinute             float64 `json:"tokensPerMinute"`
	TokensPerMinuteForIndicator float64 `json:"tokensPerMinuteForIndicator"`
	CostPerHour                 float64 `json:"costPerHour"`
}

type Projection struct {
	TotalTokens      int     `json:"totalTokens"`
	TotalCost        float64 `json:"totalCost"`
	RemainingMinutes int     `json:"remainingMinutes"`
}

type Block struct {
	ID            string      `json:"id"`
	StartTime     string      `json:"startTime"`
	EndTime       string      `json:"endTime"`
	ActualEndTime string      `json:"actualEndTime"`
	IsActive      bool        `json:"isActive"`
	IsGap         bool        `json:"isGap"`
	Entries       int         `json:"entries"`
	TokenCounts   TokenCounts `json:"tokenCounts"`
	TotalTokens   int         `json:"totalTokens"`
	CostUSD       float64     `json:"costUSD"`
	Models        []string    `json:"models"`
	BurnRate      BurnRate    `json:"burnRate"`
	Projection    Projection  `json:"projection"`
}

type BlocksResponse struct {
	Blocks []Block `json:"blocks"`
}

type BlocksData struct {
	Entries                  int
	TotalTokens              int
	InputTokens              int
	OutputTokens             int
	CacheCreationInputTokens int
	CacheReadInputTokens     int
	CostUSD                  float64
	RemainingMinutes         int
	CostPerHour              float64
}

func GetBlocks(ctx context.Context) (*BlocksData, error) {
	cmd := exec.CommandContext(ctx, "npx", "ccusage@latest", "blocks", "--active", "--json", "--offline")

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("execute ccusage (npx ccusage@latest blocks --active --json --offline): %w", err)
	}

	var response BlocksResponse
	if err := json.Unmarshal(output, &response); err != nil {
		return nil, fmt.Errorf("parse ccusage json output: %w", err)
	}

	if len(response.Blocks) == 0 {
		return nil, fmt.Errorf("no active usage blocks found in ccusage response")
	}

	block := response.Blocks[0]

	return &BlocksData{
		Entries:                  block.Entries,
		TotalTokens:              block.TotalTokens,
		InputTokens:              block.TokenCounts.InputTokens,
		OutputTokens:             block.TokenCounts.OutputTokens,
		CacheCreationInputTokens: block.TokenCounts.CacheCreationInputTokens,
		CacheReadInputTokens:     block.TokenCounts.CacheReadInputTokens,
		CostUSD:                  block.CostUSD,
		RemainingMinutes:         block.Projection.RemainingMinutes,
		CostPerHour:              block.BurnRate.CostPerHour,
	}, nil
}
