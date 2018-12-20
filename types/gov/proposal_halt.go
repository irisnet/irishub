package gov

var _ Proposal = (*HaltProposal)(nil)

type HaltProposal struct {
	TextProposal
}

//func (sp *HaltProposal) Execute(ctx sdk.Context, k Keeper) error {
//	logger := ctx.Logger().With("module", "x/gov")
//
//	if k.GetTerminatorHeight(ctx) == -1 {
//		k.SetTerminatorHeight(ctx, ctx.BlockHeight()+k.GetTerminatorPeriod(ctx))
//		logger.Info("Execute TerminatorProposal begin", "info", fmt.Sprintf("Terminator height:%d", k.GetTerminatorHeight(ctx)))
//	} else {
//		logger.Info("Terminator Period is in process.", "info", fmt.Sprintf("Terminator height:%d", k.GetTerminatorHeight(ctx)))
//
//	}
//	return nil
//}
