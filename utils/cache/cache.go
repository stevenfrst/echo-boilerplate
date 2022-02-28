package cacheHelper

import (
	"github.com/coocood/freecache"
)

const cacheSize = 100 * 1024
const CacheExpire = 60

var AllUserKey = []byte("GetUser")
var AllVerified = []byte("GetVerifiedUser")
var AllNotVerified = []byte("GetNotVerifiedUser")

var Cache = freecache.NewCache(cacheSize)
