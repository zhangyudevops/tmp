// =================================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Static is the golang structure of table static for DAO operations like Where/Data.
type Static struct {
	g.Meta `orm:"table:static, do:true"`
	Id     interface{} //
	Name   interface{} // 应用名
	Path   interface{} // 在pod中的目录
}
