// =================================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Upload is the golang structure of table upload for DAO operations like Where/Data.
type Upload struct {
	g.Meta `orm:"table:upload, do:true"`
	Id     interface{} //
	Name   interface{} // 升级包名
	Md5    interface{} // md5值
	UpTime *gtime.Time // 上传时间
}