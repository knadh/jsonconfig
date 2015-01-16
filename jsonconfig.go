/*
	Kailash Nadh,
	http://nadh.in/code/jsonconfig
	Jan 2015

	A super tiny (pseudo) JSON configuration parser for Go with
	comment support
*/

package jsonconfig

import (
	"regexp"
	"encoding/json"
	"io/ioutil"
	"errors"
)

func Load(filename string, config interface{})  error {
	// read file
	data, err := ioutil.ReadFile(filename)
	if(err != nil) {
		return errors.New("Error reading file")
	}

	// regex monstrosity because of the lack of lookbehinds/aheads
	// stand alone comments
	r1, _ := regexp.Compile(`(?m)^(\s+)?//(.*)$`)

	// numbers and boolean
	r2, _ := regexp.Compile(`(?m)"(.+?)":(\s+)?([0-9\.\-]+|true|false|null)(\s+)?,(\s+)?//(.*)$`)

	// strings
	r3, _ := regexp.Compile(`(?m)"(.+?)":(\s+)?"(.+?)"(\s+)?,(\s+)?//(.*)$`)

	// arrays and objects
	r4, _ := regexp.Compile(`(?m)"(.+?)":(\s+)?([\{\[])(.+?)([\}\]])(\s+)?,(\s+)?//(.*)$`)
	

	res := r1.ReplaceAllString(string(data), "")
	res = r2.ReplaceAllString(res, `"$1": $3,`)
	res = r3.ReplaceAllString(res, `"$1": "$3",`)
	res = r4.ReplaceAllString(res, `"$1": $3$4$5,`)

	// decode json
	if err := json.Unmarshal([]byte(res), &config); err != nil {
		return err
	}

	return nil
}
