package cache

import "omo.msa.user/proxy/nosql"

type DatumInfo struct {
	UID string
	Job string
}

func (mine *DatumInfo) initInfo(db *nosql.Datum) {
	mine.UID = db.UID.Hex()
	mine.Job = db.Job
}
