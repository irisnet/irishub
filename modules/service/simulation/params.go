package simulation

// Default simulation operation weights for messages
const (
	DefaultWeightMsgDefineService         int = 100
	DefaultWeightMsgBindService           int = 80
	DefaultWeightMsgUpdateServiceBinding  int = 50
	DefaultWeightMsgSetWithdrawAddress    int = 20
	DefaultWeightMsgDisableServiceBinding int = 20
	DefaultWeightMsgEnableServiceBinding  int = 20
	DefaultWeightMsgRefundServiceDeposit  int = 20
	DefaultWeightMsgCallService           int = 100
	DefaultWeightMsgRespondService        int = 100
	DefaultWeightMsgStartRequestContext   int = 100
	DefaultWeightMsgPauseRequestContext   int = 100
	DefaultWeightMsgKillRequestContext    int = 100
	DefaultWeightMsgUpdateRequestContext  int = 100
	DefaultWeightMsgWithdrawEarnedFees    int = 20
)
