package models

type AnEntity struct {
	Id     int64  `table:"entiry" id:"id"`
	Field1 string `column:"field1"`
	Field2 string `column:"field2"`
}

type AnotherEntity struct {
	Id      int64  `table:"entiry" id:"id"`
	AField1 string `column:"afield1"`
	BField2 string `column:"bfield2"`
}
