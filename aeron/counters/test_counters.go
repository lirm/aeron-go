// Copyright 2022 Steven Stern
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package counters

import (
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/util"
)

// InitAndWrapMetaData is a test method to inject the desired values into the buffer, then wrap it as a
// MetaDataFlyweight.
func InitAndWrapMetaData(buf *atomic.Buffer, offset int, toDriverBufLen int32, toClientBufLen int32,
	metadataBufLen int32, valuesBufLen int32, errorLogLen int32) *MetaDataFlyweight {
	buf.PutInt32(int32(offset), CurrentCncVersion)
	buf.PutInt32(int32(offset)+util.SizeOfInt32, toDriverBufLen)
	buf.PutInt32(int32(offset)+2*util.SizeOfInt32, toClientBufLen)
	buf.PutInt32(int32(offset)+3*util.SizeOfInt32, metadataBufLen)
	buf.PutInt32(int32(offset)+4*util.SizeOfInt32, valuesBufLen)
	buf.PutInt32(int32(offset)+5*util.SizeOfInt32, errorLogLen)
	var ret = new(MetaDataFlyweight)
	ret.Wrap(buf, offset)
	return ret
}
