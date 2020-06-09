package upload

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

func Test_parse(t *testing.T) {
	type args struct {
		usn string
	}
	tests := []struct {
		name     string
		args     args
		wantName string
		wantCfg  Config
		wantErr  bool
	}{
		{
			"ali",
			args{"ali"},
			"ali",
			Config{},
			false,
		},
		{
			"ali-cfg",
			args{"ali:uid=dxkite&pwd=dxkite"},
			"ali",
			Config{
				"uid": "dxkite",
				"pwd": "dxkite",
			},
			false,
		},
		{
			"ali-empty-cfg",
			args{"ali:"},
			"ali",
			Config{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotName, gotCfg, err := parse(tt.args.usn)
			if (err != nil) != tt.wantErr {
				t.Errorf("parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotName != tt.wantName {
				t.Errorf("parse() gotName = %v, want %v", gotName, tt.wantName)
			}
			if !reflect.DeepEqual(gotCfg, tt.wantCfg) {
				t.Errorf("parse() gotCfg = %v, want %v", gotCfg, tt.wantCfg)
			}
		})
	}
}

func uploadTest(t *testing.T, usn string) {
	if data, err := ioutil.ReadFile("./test/1.png"); err == nil {
		res, er := Upload(usn, &FileObject{
			Name: "cdn.png",
			Data: data,
		})
		if er != nil {
			t.Error(er)
		} else {
			r, er := http.Get(res.Url)
			if er != nil {
				t.Error(er)
			}
			buf := &bytes.Buffer{}
			if _, err := io.Copy(buf, r.Body); err != nil {
				t.Error(err)
			}
			if bytes.Compare(data, buf.Bytes()) != 0 {
				fmt.Println("uploaded but not raw data", res.Url)
			}
			fmt.Println("uploaded", res.Url)
		}
	} else {
		t.Error(err)
	}
}
