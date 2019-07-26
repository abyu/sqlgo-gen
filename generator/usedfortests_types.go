package generator

type AType struct {
	AField      int
	BField      string
	TaggedField string `key:"value" another:"value2"`
}
