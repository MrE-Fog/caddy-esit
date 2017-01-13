// Copyright 2015-2016, Cyrill @ Schumacher.fm and the CoreStore contributors
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

package responseproxy_test

import (
	"bytes"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/SchumacherFM/caddyesi/responseproxy"
	"github.com/stretchr/testify/assert"
)

func TestWrapPipe(t *testing.T) {
	wOrg := httptest.NewRecorder()
	data := []byte(`Commander Data encrypts the computer with a fractal algorithm to protect it from the Borgs.`)

	buf := new(bytes.Buffer)
	pw := responseproxy.WrapPiped(buf, wOrg)

	n, err := pw.Write(data)
	if err := pw.Close(); err != nil {
		t.Fatal(err)
	}

	assert.NoError(t, err)
	assert.Exactly(t, len(data), n)
	assert.Exactly(t, 0, wOrg.Body.Len())

	assert.Exactly(t, len(data), buf.Len()) // race because of MutexBuffer

	if _, err := io.Copy(pw.Unwrap(), buf); err != nil {
		t.Fatal(err)
	}
	assert.Exactly(t, len(data), wOrg.Body.Len())
}
