package parameter

func RegisterGovParamMapping(gps ...GovParameter) {
	for _, gp := range gps {
		if gp != nil {
			ParamMapping[gp.GetStoreKey()] = gp
		}
	}
}