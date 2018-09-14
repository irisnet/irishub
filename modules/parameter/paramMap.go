package parameter

var paramMapping = make(map[string]GovParameter)

func RegisterGovParamMapping(gp GovParameter) {

	if gp != nil {
		paramMapping[gp.GetStoreKey()] = gp
	}
}
