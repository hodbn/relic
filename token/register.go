//
// Copyright (c) SAS Institute Inc.
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
//

package token

import (
	"io"

	"github.com/mind-security/relic/v8/config"
	"github.com/mind-security/relic/v8/lib/passprompt"
)

type (
	OpenFunc func(cfg *config.Config, tokenName string, prompt passprompt.PasswordGetter) (Token, error)
	ListFunc func(provider string, dest io.Writer) error
)

var (
	Openers = make(map[string]OpenFunc)
	Listers = make(map[string]ListFunc)
)
