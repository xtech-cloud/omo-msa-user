package cache

import "omo.msa.user/proxy/nosql"

type DatumInfo struct {
	Sex uint8
	UID string
	RealName string
	Phone string
	Job string
}

func (mine *DatumInfo)initInfo(db *nosql.Datum)  {
	mine.UID = db.UID.Hex()
	mine.RealName = db.Name
	mine.Phone = db.Phone
	mine.Job = db.Job
	mine.Sex = db.Sex
}
