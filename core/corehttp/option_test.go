
//此源码被清华学神尹成大魔王专业翻译分析并修改
//尹成QQ77025077
//尹成微信18510341407
//尹成所在QQ群721929980
//尹成邮箱 yinc13@mails.tsinghua.edu.cn
//尹成毕业于清华大学,微软区块链领域全球最有价值专家
//https://mvp.microsoft.com/zh-cn/PublicProfile/4033620
package corehttp

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	version "github.com/ipfs/go-ipfs"
)

type testcasecheckversion struct {
	userAgent    string
	uri          string
	shouldHandle bool
	responseBody string
	responseCode int
}

func (tc testcasecheckversion) body() string {
	if !tc.shouldHandle && tc.responseBody == "" {
		return fmt.Sprintf("%s (%s != %s)\n", errAPIVersionMismatch, version.ApiVersion, tc.userAgent)
	}

	return tc.responseBody
}

func TestCheckVersionOption(t *testing.T) {
	tcs := []testcasecheckversion{
		{"/go-ipfs/0.1/", APIPath + "/test/", false, "", http.StatusBadRequest},
		{"/go-ipfs/0.1/", APIPath + "/version", true, "check!", http.StatusOK},
		{version.ApiVersion, APIPath + "/test", true, "check!", http.StatusOK},
		{"Mozilla Firefox/no go-ipfs node", APIPath + "/test", true, "check!", http.StatusOK},
		{"/go-ipfs/0.1/", "/webui", true, "check!", http.StatusOK},
	}

	for _, tc := range tcs {
		t.Logf("%#v", tc)
		r := httptest.NewRequest("POST", tc.uri, nil)
r.Header.Add("User-Agent", tc.userAgent) //旧版本，应该失败

		called := false
		root := http.NewServeMux()
		mux, err := CheckVersionOption()(nil, nil, root)
		if err != nil {
			t.Fatal(err)
		}

		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			called = true
			if !tc.shouldHandle {
				t.Error("handler was called even though version didn't match")
			} else {
				io.WriteString(w, "check!")
			}
		})

		w := httptest.NewRecorder()

		root.ServeHTTP(w, r)

		if tc.shouldHandle && !called {
			t.Error("handler wasn't called even though it should have")
		}

		if w.Code != tc.responseCode {
			t.Errorf("expected code %d but got %d", tc.responseCode, w.Code)
		}

		if w.Body.String() != tc.body() {
			t.Errorf("expected error message %q, got %q", tc.body(), w.Body.String())
		}
	}
}