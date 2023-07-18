// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package dao

import (
	"pack/internal/dao/internal"
)

// internalStaticDao is internal type for wrapping internal DAO implements.
type internalStaticDao = *internal.StaticDao

// staticDao is the data access object for table static.
// You can define custom methods on it to extend its functionality as you wish.
type staticDao struct {
	internalStaticDao
}

var (
	// Static is globally public accessible object for table static operations.
	Static = staticDao{
		internal.NewStaticDao(),
	}
)

// Fill with you ideas below.