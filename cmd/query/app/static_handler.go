// Copyright (c) 2017 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package app

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

const (
	defaultStaticAssetsRoot = "jaeger-ui-build/build/"
)

var (
	staticRootFiles = []string{"favicon.ico"}
)

// StaticAssetsHandler handles static assets
type StaticAssetsHandler struct {
	staticAssetsRoot string
}

// NewStaticAssetsHandler returns a StaticAssetsHandler
func NewStaticAssetsHandler(staticAssetsRoot string) (*StaticAssetsHandler, error) {
	return newStaticAssetsHandler(staticAssetsRoot, defaultStaticAssetsRoot)
}

func newStaticAssetsHandler(staticAssetsRoot, defaultStaticAssetsRoot string) (*StaticAssetsHandler, error) {
	if staticAssetsRoot == "" {
		staticAssetsRoot = defaultStaticAssetsRoot
	}
	if !strings.HasSuffix(staticAssetsRoot, "/") {
		staticAssetsRoot = staticAssetsRoot + "/"
	}
	if _, err := ioutil.ReadFile(staticAssetsRoot + "index.html"); err != nil {
		return nil, errors.Wrap(err, "Cannot read index.html, static-files argument may be incorrect")
	}
	return &StaticAssetsHandler{staticAssetsRoot: staticAssetsRoot}, nil
}

// RegisterRoutes registers routes for this handler on the given router
func (sH *StaticAssetsHandler) RegisterRoutes(router *mux.Router) {
	router.PathPrefix("/static").Handler(http.FileServer(http.Dir(sH.staticAssetsRoot)))
	for _, file := range staticRootFiles {
		router.Path("/" + file).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, sH.staticAssetsRoot+file)
		})
	}
	router.NotFoundHandler = http.HandlerFunc(sH.notFound)
}

func (sH *StaticAssetsHandler) notFound(w http.ResponseWriter, r *http.Request) {
	// don't allow returning "304 Not Modified" for index.html because
	// the cached versions might have the wrong filenames for javascript assets
	delete(r.Header, "If-Modified-Since")
	http.ServeFile(w, r, sH.staticAssetsRoot+"index.html")
}
