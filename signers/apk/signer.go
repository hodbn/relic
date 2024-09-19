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

package apk

import (
	"errors"
	"io"

	"github.com/mind-security/relic/v8/lib/certloader"
	"github.com/mind-security/relic/v8/lib/magic"
	"github.com/mind-security/relic/v8/signers"
	"github.com/mind-security/relic/v8/signers/zipbased"
)

// Sign Android packages

var ApkSigner = &signers.Signer{
	Name:      "apk",
	Magic:     magic.FileTypeAPK,
	CertTypes: signers.CertTypeX509,
	Transform: zipbased.Transform,
	Sign:      sign,
	Verify:    verify,
}

const (
	sigMagic = "APK Sig Block 42"
	sigApkV2 = 0x7109871a
)

var (
	errMalformed = errors.New("malformed APK signing block")
	errTruncated = errors.New("truncated APK signing block sequence")
)

func init() {
	signers.Register(ApkSigner)
}

func sign(r io.Reader, cert *certloader.Certificate, opts signers.SignOpts) ([]byte, error) {
	digest, err := digestApkStream(r, opts.Hash)
	if err != nil {
		return nil, err
	}
	patchset, err := digest.Sign(cert)
	if err != nil {
		return nil, err
	}
	return opts.SetBinPatch(patchset)
}
