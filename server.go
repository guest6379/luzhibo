package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"luzhibo/api"
	"net/http"
	"os"
	"strconv"
	"time"
)

type checkRet struct {
	Pass bool
	Has  bool
	Live bool
	Err  bool
	Path string
}

type tasksRet struct {
	Tasks []*taskInfo
	Err   bool
	E     bool
}

type ajaxHandler struct{}

//ServeHTTP 实现接口
func (_ ajaxHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	switch r.Form.Get("act") {
	case "check":
		tr := checkRet{}
		url := r.Form.Get("url")
		oa := api.New(url)
		if oa == nil {
			tr.Pass = false
		} else {
			tr.Pass = true
			i, l, e := oa.GetRoomInfo()
			if e == nil {
				tr.Has = true
				tr.Path = fmt.Sprintf("[%s]%s_%s", oa.Site, i, time.Now().Format("20060102150405"))
				tr.Live = l
			} else {
				tr.Err = true
			}
		}
		j, _ := json.Marshal(tr)
		w.Write(j)
		return
	case "add":
		url, m, p, s := r.Form.Get("url"), r.Form.Get("m"), r.Form.Get("path"), r.Form.Get("run")
		mm, ss := m == "true", s == "true"
		if url != "" && p != "" {
			if addTaskEx(url, p, mm, ss) {
				w.Write([]byte("ok"))
				return
			}
		}
	case "del":
		i, d := r.Form.Get("id"), r.Form.Get("f")
		b := d == "true"
		c, e := strconv.Atoi(i)
		if e == nil {
			if delTask(c-1, b) {
				w.Write([]byte("ok"))
				return
			}
		}

	case "start":
		i := r.Form.Get("id")
		if startOrStopTask(i, true) {
			w.Write([]byte("ok"))
			return
		}
	case "stop":
		i := r.Form.Get("id")
		if startOrStopTask(i, false) {
			w.Write([]byte("ok"))
			return
		}
	case "tasks":
		list, o, e := getTaskInfoList()
		r := tasksRet{}
		r.Err = o
		r.Tasks = list
		r.E = e
		j, _ := json.Marshal(r)
		w.Write(j)
		return
	case "exist":
		p := r.Form.Get("path")
		if pp, _ := pathExist(p); pp {
			w.Write([]byte("exist"))
			return
		}
	case "get":
		i := r.Form.Get("id")
		ii, e := strconv.Atoi(i)
		if e == nil {
			inf, _ := getTaskInfo(ii - 1)
			pp := inf.Path
			if inf.M {
				pp += ".flv"
			}
			w.Header().Add("Content-Disposition", "attachment; filename=\""+pp+"\"")
			w.Header().Add("Content-Type", "video/x-flv")
			getAct(inf.Path, w)
		}
		return

	}
	w.Write([]byte(""))
}

type uiHandler struct{}

//ServeHTTP 实现接口
func (_ uiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(html))
}

func getFile(path string, w http.ResponseWriter) {
	f, e := os.Open(path)
	eof := false
	if e == nil {
		for {
			buf := make([]byte, bytes.MinRead)
			t, e := f.Read(buf)
			if e != nil {
				if e == io.EOF {
					eof = true
				} else {
					break
				}
			}
			_, e = w.Write(buf[:t])
			if e != nil || eof {
				break
			}
		}
	}
}

func getDir(path string, w http.ResponseWriter) {
	files, err := ioutil.ReadDir(path)
	if err == nil {
		for _, f := range files {
			if !f.IsDir() {
				p := path + "/" + f.Name()
				getFile(p, w)
			}
		}
	}
}

func getAct(path string, w http.ResponseWriter) {
	if pe, d := pathExist(path); pe {
		if d {
			getDir(path, w)
		} else {
			getFile(path, w)
		}
	} else {
		w.Write([]byte("no exist"))
	}

}

func startOrStopTask(i string, m bool) bool {
	c, e := strconv.Atoi(i)
	if e != nil {
		return false
	}
	c--
	if m {
		return startTask(c)
	}
	return stopTask(c)
}

func startServer(s string) {
	http.Handle("/", uiHandler{})
	http.Handle("/ajax", ajaxHandler{})
	http.ListenAndServe(s, nil)
}