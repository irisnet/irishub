package parameter

var ParamMapping = make(map[string]GovParameter)

func RegisterGovParamMapping(gps ...GovParameter) {
	for _, gp := range gps {
		if gp != nil {
			ParamMapping[gp.GetStoreKey()] = gp
		}
	}
}
