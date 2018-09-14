package parameter

var paramMapping = make(map[string]GovParameter)

func RegisterGovParamMapping(gps ...GovParameter) {
	for _, gp := range gps {
		if gp != nil {
			paramMapping[gp.GetStoreKey()] = gp
		}
	}
}
