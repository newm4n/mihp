package helper

import "testing"

func TestNewProperties(t *testing.T) {
	prop := NewProperties()
	prop["abc"] = "aesfasef"
	prop["abc.def"] = "aesfasef"
	prop["abc.def.ghi"] = "aesfasef"
	prop["abc.cde"] = "aesfasef"
	prop["abc.cde.fgh"] = "aesfasef"
	prop["abc.aiu"] = "aesfasef"
	prop["abc.aiu.asef"] = "aesfasef"
	t.Log(prop.String())
}
