package mySQL

import (
	"io/ioutil"

	Util "github.com/Nimajjj/Tidder/go/utility"
)

func DefaultPP() string {
	content, err := ioutil.ReadFile("./images/profiles/default_pp.txt")
	if err != nil {
		Util.Error(err)
	}

	return string(content)
}
