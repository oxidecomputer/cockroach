// Copyright 2018 The Cockroach Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package aws

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/cockroachdb/cockroach/pkg/util/timeutil"
	"github.com/cockroachdb/errors"
	"github.com/cockroachdb/errors/oserror"
)

const sshPublicKeyFile = "${HOME}/.ssh/id_rsa.pub"

// sshKeyExists checks to see if there is a an SSH key with the given name in the given region.
func (p *Provider) sshKeyExists(keyName, region string) (bool, error) {
	var data struct {
		KeyPairs []struct {
			KeyName string
		}
	}
	args := []string{
		"ec2", "describe-key-pairs",
		"--region", region,
	}
	err := p.runJSONCommand(args, &data)
	if err != nil {
		return false, err
	}
	for _, keyPair := range data.KeyPairs {
		if keyPair.KeyName == keyName {
			return true, nil
		}
	}
	return false, nil
}

// sshKeyImport takes the user's local, public SSH key and imports it into the ec2 region so that
// we can create new hosts with it.
func (p *Provider) sshKeyImport(keyName, region string) error {
	_, err := os.Stat(os.ExpandEnv(sshPublicKeyFile))
	if err != nil {
		if oserror.IsNotExist(err) {
			return errors.Wrapf(err, "please run ssh-keygen externally to create your %s file", sshPublicKeyFile)
		}
		return err
	}

	var data struct {
		KeyName string
	}
	_ = data.KeyName // silence unused warning

	user, err := p.FindActiveAccount()
	if err != nil {
		return err
	}

	timestamp := timeutil.Now()
	createdAt := timestamp.Format(time.RFC3339)

	IAMUserNameTag := fmt.Sprintf("{Key=IAMUserName,Value=%s}", user)
	createdAtTag := fmt.Sprintf("{Key=CreatedAt,Value=%s}", createdAt)
	tagSpecs := fmt.Sprintf("ResourceType=key-pair,Tags=[%s, %s]", IAMUserNameTag, createdAtTag)

	args := []string{
		"ec2", "import-key-pair",
		"--region", region,
		"--key-name", keyName,
		"--public-key-material", fmt.Sprintf("fileb://%s", sshPublicKeyFile),
		"--tag-specifications", tagSpecs,
	}
	err = p.runJSONCommand(args, &data)
	// If two roachprod instances run at the same time with the same key, they may
	// race to upload the key pair.
	if err == nil || strings.Contains(err.Error(), "InvalidKeyPair.Duplicate") {
		return nil
	}
	return err
}

// sshKeyName computes the name of the ec2 ssh key that we'll store the local user's public key in
func (p *Provider) sshKeyName() (string, error) {
	user, err := p.FindActiveAccount()
	if err != nil {
		return "", err
	}

	keyBytes, err := ioutil.ReadFile(os.ExpandEnv(sshPublicKeyFile))
	if err != nil {
		if oserror.IsNotExist(err) {
			return "", errors.Wrapf(err, "please run ssh-keygen externally to create your %s file", sshPublicKeyFile)
		}
		return "", err
	}

	hash := sha1.New()
	if _, err := hash.Write(keyBytes); err != nil {
		return "", err
	}
	hashBytes := hash.Sum(nil)
	hashText := base64.URLEncoding.EncodeToString(hashBytes)

	return fmt.Sprintf("%s-%s", user, hashText), nil
}
