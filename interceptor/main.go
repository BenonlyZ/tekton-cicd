/*
Copyright 2020 The Tekton Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/tektoncd/triggers/pkg/apis/triggers/v1alpha1"
	"github.com/xanzy/go-gitlab"
	"google.golang.org/grpc/codes"
)

var count int

type Result struct {
	Body string `json:"body"`
}

func main() {

	http.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		count++

		d, err := ioutil.ReadAll(r.Body)
		if err != nil {
			respondErr(w, 500, codes.Unknown, err)
			return
		}

		if len(d) != 0 {
			fmt.Println("请求方法为:", r.Method)
			fmt.Printf("统计次数:%d\n", count)
			var e Result
			p := gitlab.PushEvent{}
			var err error
			err = json.Unmarshal(d, &e)
			if err != nil {
				respondErr(w, 500, codes.Unknown, err)
				return
			}
			fmt.Println("事件信息: ", e.Body)

			err = json.Unmarshal([]byte(e.Body), &p)
			if err != nil {
				respondErr(w, 500, codes.Unknown, err)
				return
			}

			//添加附加信息
			pipe := make(map[string]interface{})
			pipe["extra"] = "extraParams"
			//拦截器返回固定格式
			iRes := &v1alpha1.InterceptorResponse{
				Extensions: pipe,
				Continue:   true,
				Status: v1alpha1.Status{
					Code:    codes.OK,
					Message: strconv.Itoa(count) + ": " + string(d),
				},
			}

			b, err := json.Marshal(iRes)
			if err != nil {
				respondErr(w, 500, codes.Unknown, err)
				return
			}
			_, err = w.Write(b)
			if err != nil {
				respondErr(w, 500, codes.Unknown, err)
				return
			}
		}

	})
	http.HandleFunc("/unready", func(w http.ResponseWriter, r *http.Request) {
		d, err := ioutil.ReadAll(r.Body)
		if err != nil {
			respondErr(w, 500, codes.Unknown, err)
			return
		}

		fmt.Println(string(d))
		respondErr(w, 500, codes.Unknown, errors.New("unready err"))
		fmt.Println("返回码: 500")
	})

	http.ListenAndServe(":8080", nil)
}

func respondErr(w http.ResponseWriter, statusCode int, code codes.Code, err error) {
	log.Printf("respond %d: %s", statusCode, err)

	b, err := json.Marshal(v1alpha1.InterceptorResponse{
		Continue: false,
		Status: v1alpha1.Status{
			Code:    code,
			Message: err.Error(),
		},
	})
	if err != nil {
		log.Printf("marshall error response: %v", err)
	}

	w.WriteHeader(statusCode)
	_, err = w.Write(b)
	if err != nil {
		log.Printf("write error response: %v", err)
	}
}
