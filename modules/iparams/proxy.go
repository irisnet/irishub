package iparams

const (
	Gov    = "gov"
	Global = "global"
)

// Getter returns readonly struct,default get from Global
func (k Keeper) GlobalGetter() GlobalGetter {
	return GlobalGetter{k}
}

// Getter returns readonly struct,default get from Global
func (k Keeper) GovGetter() GovGetter {
	return GovGetter{k}
}

func (k Keeper) GlobalSetter() GlobalSetter {
	return GlobalSetter{GlobalGetter{k}}
}

func (k Keeper) GovSetter() GovSetter {
	return GovSetter{GovGetter{k}}
}
