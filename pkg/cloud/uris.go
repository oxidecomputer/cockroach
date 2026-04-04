// Copyright 2021 The Cockroach Authors.
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

package cloud

import (
	"net/url"
	"path"
	"strings"

	"github.com/cockroachdb/cockroach/pkg/roachpb"
)

const (
	// AuthParam is the query parameter for the cluster settings named
	// key in a URI.
	AuthParam = "AUTH"
	// AuthParamImplicit is the query parameter for the implicit authentication
	// mode in a URI.
	AuthParamImplicit = roachpb.ExternalStorageAuthImplicit
	// AuthParamSpecified is the query parameter for the specified authentication
	// mode in a URI.
	AuthParamSpecified = roachpb.ExternalStorageAuthSpecified
)

// GetPrefixBeforeWildcard gets the prefix of a path that does not contain glob-
// style matchers, up to the last path segment before the first one which has a
// glob character.
func GetPrefixBeforeWildcard(p string) string {
	globIndex := strings.IndexAny(p, "*?[")
	if globIndex < 0 {
		return p
	}
	return path.Dir(p[:globIndex])
}

// RedactKMSURI redacts the Master Key ID and the ExternalStorage secret
// credentials.
func RedactKMSURI(kmsURI string) (string, error) {
	sanitizedKMSURI, err := SanitizeExternalStorageURI(kmsURI, nil)
	if err != nil {
		return "", err
	}

	// Redact the path which contains the KMS Master Key identifier.
	uri, err := url.ParseRequestURI(sanitizedKMSURI)
	if err != nil {
		return "", err
	}
	uri.Path = "/redacted"
	return uri.String(), nil
}

// JoinPathPreservingTrailingSlash wraps path.Join but preserves the trailing
// slash if there was one in the suffix.
//
// This is particularly important when the joined path is used as a prefix for
// listing. When listing, the suffix *after the listed prefix* of each file name
// is what is returned and, importantly, what is used when grouping using a
// delimiter. E.g. when using `/` as a delimiter to find what might be called the
// immediate children in a directory, we pass that directory's path *with a
// trailing slash* as the prefix, so that the children do not start with a slash
// and get grouped into nothing. Thus it is important that if we use path.Join
// to construct the prefix, we always preserve the trailing slash.
func JoinPathPreservingTrailingSlash(prefix, suffix string) string {
	out := path.Join(prefix, suffix)
	// path.Clean removes trailing slashes, so put it back if needed.
	if strings.HasSuffix(suffix, "/") {
		out += "/"
	}
	return out
}
