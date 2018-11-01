// Copyright 2017 Xiaomi, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	log "github.com/Sirupsen/logrus"
	"github.com/sethvargo/go-password/password"
	"github.com/spf13/viper"
	"github.com/toolkits/str"
)

func HashIt(passwd string) (hashed string) {
	salt := viper.GetString("salt")
	log.Debugf("salf is %v", salt)
	if salt == "" {
		log.Error("salt is empty, please check your app.ini")
	}
	hashed = str.Md5Encode(salt + passwd)
	return
}

func GeneratePass(length int) string {
	res, err := password.Generate(length, 5, 0, false, true)
	if err != nil {
		log.Fatal(err)
	}
	return res
}
