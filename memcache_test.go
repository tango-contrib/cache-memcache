// Copyright 2014 The Macaron Authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package cache

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Unknwon/com"
	"github.com/lunny/tango"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/tango-contrib/cache"
)

type CacheAction struct {
	cache.Cache
}

func (c *CacheAction) Get() {
	So(c.Put("uname", "unknwon", 1), ShouldBeNil)
	So(c.Put("uname2", "unknwon2", 1), ShouldBeNil)
	So(c.IsExist("uname"), ShouldBeTrue)

	So(c.Cache.Get("404"), ShouldBeNil)
	So(c.Cache.Get("uname").(string), ShouldEqual, "unknwon")

	time.Sleep(1 * time.Second)
	So(c.Cache.Get("uname"), ShouldBeNil)
	time.Sleep(1 * time.Second)
	So(c.Cache.Get("uname2"), ShouldBeNil)

	So(c.Put("uname", "unknwon", 0), ShouldBeNil)
	So(c.Delete("uname"), ShouldBeNil)
	So(c.Cache.Get("uname"), ShouldBeNil)

	So(c.Put("uname", "unknwon", 0), ShouldBeNil)
	So(c.Flush(), ShouldBeNil)
	So(c.Cache.Get("uname"), ShouldBeNil)
}

type Cache2Action struct {
	cache.Cache
}

func (c *Cache2Action) Get() {
	So(c.Incr("404"), ShouldNotBeNil)
	So(c.Decr("404"), ShouldNotBeNil)

	So(c.Put("int", 0, 0), ShouldBeNil)
	So(c.Put("int64", int64(0), 0), ShouldBeNil)
	So(c.Put("string", "hi", 0), ShouldBeNil)

	So(c.Incr("int"), ShouldBeNil)
	So(c.Incr("int64"), ShouldBeNil)

	So(c.Decr("int"), ShouldBeNil)
	So(c.Decr("int64"), ShouldBeNil)

	So(c.Incr("string"), ShouldNotBeNil)
	So(c.Decr("string"), ShouldNotBeNil)

	So(com.StrTo(c.Cache.Get("int").(string)).MustInt(), ShouldEqual, 0)
	So(com.StrTo(c.Cache.Get("int64").(string)).MustInt64(), ShouldEqual, 0)

	So(c.Flush(), ShouldBeNil)
}

func Test_MemcacheCacher(t *testing.T) {
	Convey("Test memcache cache adapter", t, func() {
		opt := cache.Options{
			Adapter:       "memcache",
			AdapterConfig: "127.0.0.1:9090",
		}

		Convey("Basic operations", func() {
			t := tango.New()
			t.Use(cache.New(opt))

			t.Get("/", new(CacheAction))

			resp := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "/", nil)
			So(err, ShouldBeNil)
			t.ServeHTTP(resp, req)
		})

		Convey("Increase and decrease operations", func() {
			t := tango.New()
			t.Use(cache.New(opt))
			t.Get("/", new(Cache2Action))

			resp := httptest.NewRecorder()
			req, err := http.NewRequest("GET", "/", nil)
			So(err, ShouldBeNil)
			t.ServeHTTP(resp, req)
		})
	})
}
