/*
Copyright 2016 Stanislav Liberman

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package counters

import (
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/flyweight"
	"github.com/lirm/aeron-go/aeron/util"
	"github.com/lirm/aeron-go/aeron/util/memmap"
	"github.com/op/go-logging"
	"log"
)

var logger = logging.MustGetLogger("counters")

const (
	CncFile           string = "cnc.dat"
	CurrentCncVersion int32  = 5
)

type MetaDataFlyweight struct {
	flyweight.FWBase

	CncVersion flyweight.Int32Field

	ToDriverBufLen   flyweight.Int32Field
	ToClientBufLen   flyweight.Int32Field
	metadataBuLen    flyweight.Int32Field
	valuesBufLen     flyweight.Int32Field
	ClientLivenessTo flyweight.Int64Field
	errorLogLen      flyweight.Int32Field

	ToDriverBuf  flyweight.RawDataField
	ToClientsBuf flyweight.RawDataField
	MetaDataBuf  flyweight.RawDataField
	ValuesBuf    flyweight.RawDataField
	ErrorBuf     flyweight.RawDataField
}

func (m *MetaDataFlyweight) Wrap(buf *atomic.Buffer, offset int) flyweight.Flyweight {
	pos := offset
	pos += m.CncVersion.Wrap(buf, pos)
	pos += m.ToDriverBufLen.Wrap(buf, pos)
	pos += m.ToClientBufLen.Wrap(buf, pos)
	pos += m.metadataBuLen.Wrap(buf, pos)
	pos += m.valuesBufLen.Wrap(buf, pos)
	pos += m.ClientLivenessTo.Wrap(buf, pos)
	pos += m.errorLogLen.Wrap(buf, pos)

	pos = int(util.AlignInt32(int32(pos), util.CacheLineLength*2))

	pos += m.ToDriverBuf.Wrap(buf, pos, m.ToDriverBufLen.Get())
	pos += m.ToClientsBuf.Wrap(buf, pos, m.ToClientBufLen.Get())
	pos += m.MetaDataBuf.Wrap(buf, pos, m.metadataBuLen.Get())
	pos += m.ValuesBuf.Wrap(buf, pos, m.valuesBufLen.Get())
	pos += m.ErrorBuf.Wrap(buf, pos, m.errorLogLen.Get())

	m.SetSize(pos - offset)
	return m
}

func MapFile(filename string) (*MetaDataFlyweight, *memmap.File) {

	logger.Debugf("Trying to map file: %s", filename)
	cncMap, err := memmap.MapExisting(filename, 0, 0)
	if err != nil {
		log.Fatalf("Failed to map the file %s with %s", filename, err.Error())
	}

	cncBuffer := atomic.MakeBuffer(cncMap.GetMemoryPtr(), cncMap.GetMemorySize())
	var meta MetaDataFlyweight
	meta.Wrap(cncBuffer, 0)

	cncVer := meta.CncVersion.Get()
	logger.Debugf("Mapped %s for ver %d", filename, cncVer)

	if CurrentCncVersion != cncVer {
		log.Fatalf("aeron cnc file version not understood: version=%d", cncVer)
	}

	return &meta, cncMap
}
