package slashing

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type GenesisState struct {
	Params          Params                          `json:"params"`
	SigningInfos    map[string]ValidatorSigningInfo `json:"signing_infos"`
	MissedBlocks    map[string][]MissedBlock        `json:"missed_blocks"`
	SlashingPeriods []ValidatorSlashingPeriod       `json:"slashing_periods"`
}

type Params struct {
	MaxEvidenceAge          int64         `json:"max_evidence_age"`
	SignedBlocksWindow      int64         `json:"signed_blocks_window"`
	MinSignedPerWindow      sdk.Dec       `json:"min_signed_per_window"`
	DoubleSignJailDuration  time.Duration `json:"double_sign_jail_duration"`
	DowntimeJailDuration    time.Duration `json:"downtime_jail_duration"`
	CensorshipJailDuration  time.Duration `json:"censorship_jail_duration"`
	SlashFractionDoubleSign sdk.Dec       `json:"slash_fraction_double_sign"`
	SlashFractionDowntime   sdk.Dec       `json:"slash_fraction_downtime"`
	SlashFractionCensorship sdk.Dec       `json:"slash_fraction_censorship"`
}

type MissedBlock struct {
	Index  int64 `json:"index"`
	Missed bool  `json:"missed"`
}

type ValidatorSigningInfo struct {
	StartHeight         int64     `json:"start_height"`          // height at which validator was first a candidate OR was unjailed
	IndexOffset         int64     `json:"index_offset"`          // index offset into signed block bit array
	JailedUntil         time.Time `json:"jailed_until"`          // timestamp validator cannot be unjailed until
	MissedBlocksCounter int64     `json:"missed_blocks_counter"` // missed blocks counter (to avoid scanning the array every time)
}

type ValidatorSlashingPeriod struct {
	ValidatorAddr sdk.ConsAddress `json:"validator_addr"` // validator which this slashing period is for
	StartHeight   int64           `json:"start_height"`   // starting height of the slashing period
	EndHeight     int64           `json:"end_height"`     // ending height of the slashing period, or sentinel value of 0 for in-progress
	SlashedSoFar  sdk.Dec         `json:"slashed_so_far"` // fraction of validator stake slashed so far in this slashing period
}
