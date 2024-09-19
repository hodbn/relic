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

package cat

// Sign Microsoft security catalog files

import (
	"errors"
	"io"
	"io/ioutil"

	"github.com/mind-security/relic/v8/lib/authenticode"
	"github.com/mind-security/relic/v8/lib/certloader"
	"github.com/mind-security/relic/v8/lib/magic"
	"github.com/mind-security/relic/v8/lib/pkcs7"
	"github.com/mind-security/relic/v8/lib/pkcs9"
	"github.com/mind-security/relic/v8/signers"
	"github.com/mind-security/relic/v8/signers/pkcs"
)

var CatSigner = &signers.Signer{
	Name:      "cat",
	Magic:     magic.FileTypeCAT,
	CertTypes: signers.CertTypeX509,
	Sign:      sign,
	Verify:    pkcs.Verify,
}

func init() {
	signers.Register(CatSigner)
}

func sign(r io.Reader, cert *certloader.Certificate, opts signers.SignOpts) ([]byte, error) {
	blob, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	oldpsd, err := pkcs7.Unmarshal(blob)
	if err != nil {
		return nil, err
	}
	if !oldpsd.Content.ContentInfo.ContentType.Equal(authenticode.OidCertTrustList) {
		return nil, errors.New("not a security catalog")
	}
	sig := pkcs7.NewBuilder(cert.Signer(), cert.Chain(), opts.Hash)
	if err := sig.SetContentInfo(oldpsd.Content.ContentInfo); err != nil {
		return nil, err
	}
	newpsd, err := sig.Sign()
	if err != nil {
		return nil, err
	}
	ts, err := pkcs9.TimestampAndMarshal(opts.Context(), newpsd, cert.Timestamper, true)
	if err != nil {
		return nil, err
	}
	return opts.SetPkcs7(ts)
}
