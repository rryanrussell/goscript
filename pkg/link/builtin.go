package link

func GetBuiltInSymbols() map[string]*Symbol {
	return map[string]*Symbol{
		"float64": {
			Type: "float64",
			Kind: BuiltIn,
		},
	}
}
