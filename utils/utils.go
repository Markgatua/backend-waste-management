package utils

 import (
 	"fmt"
 	"ttnmwastemanagementsystem/gen"
 )

 func GetNextTableID(table string) int32{
 	query := fmt.Sprint("SELECT nextval ('",table,"_id_seq')+1")
 	var nextVal int32
 	gen.REPO.DB.Get(&nextVal,query)
 	return nextVal
 }