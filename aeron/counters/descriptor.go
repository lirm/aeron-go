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
)

const (
	CncFile           string = "cnc.dat"
	CurrentCncVersion int32  = 5
)

type MetaDataFlyweight struct {
	flyweight.FWBase

	cncVersion flyweight.Int32Field

	driverBufLen     flyweight.Int32Field
	clientBufLen     flyweight.Int32Field
	metadataBuLen    flyweight.Int32Field
	valuesBufLen     flyweight.Int32Field
	clientLivenessTo flyweight.Int64Field
	errorLogLen      flyweight.Int32Field

	toDriverBuf  flyweight.RawDataField
	toClientsBuf flyweight.RawDataField
	metadataBuf  flyweight.RawDataField
	valuesBuf    flyweight.RawDataField
	errorBuf     flyweight.RawDataField
}

func (m *MetaDataFlyweight) Wrap(buf *atomic.Buffer, offset int) flyweight.Flyweight {
	pos := offset
	pos += m.cncVersion.Wrap(buf, pos)
	pos += m.driverBufLen.Wrap(buf, pos)
	pos += m.clientBufLen.Wrap(buf, pos)
	pos += m.metadataBuLen.Wrap(buf, pos)
	pos += m.valuesBufLen.Wrap(buf, pos)
	pos += m.clientLivenessTo.Wrap(buf, pos)
	pos += m.errorLogLen.Wrap(buf, pos)

	pos = int(util.AlignInt32(int32(pos), util.CacheLineLength*2))

	pos += m.toDriverBuf.Wrap(buf, pos, m.driverBufLen.Get())
	pos += m.toClientsBuf.Wrap(buf, pos, m.clientBufLen.Get())
	pos += m.metadataBuf.Wrap(buf, pos, m.metadataBuLen.Get())
	pos += m.valuesBuf.Wrap(buf, pos, m.valuesBufLen.Get())
	pos += m.errorBuf.Wrap(buf, pos, m.errorLogLen.Get())

	m.SetSize(pos - offset)
	return m
}

func CreateToDriverBuffer(cncFile *memmap.File) *atomic.Buffer {

	cncBuffer := atomic.MakeBuffer(cncFile.GetMemoryPtr(), cncFile.GetMemorySize())
	var meta MetaDataFlyweight
	meta.Wrap(cncBuffer, 0)

	return meta.toDriverBuf.Get()
}

func CreateToClientsBuffer(cncFile *memmap.File) *atomic.Buffer {

	cncBuffer := atomic.MakeBuffer(cncFile.GetMemoryPtr(), cncFile.GetMemorySize())
	var meta MetaDataFlyweight
	meta.Wrap(cncBuffer, 0)

	return meta.toClientsBuf.Get()
}

func CreateCounterMetadataBuffer(cncFile *memmap.File) *atomic.Buffer {

	cncBuffer := atomic.MakeBuffer(cncFile.GetMemoryPtr(), cncFile.GetMemorySize())
	var meta MetaDataFlyweight
	meta.Wrap(cncBuffer, 0)

	return meta.metadataBuf.Get()
}

func CreateCounterValuesBuffer(cncFile *memmap.File) *atomic.Buffer {

	cncBuffer := atomic.MakeBuffer(cncFile.GetMemoryPtr(), cncFile.GetMemorySize())
	var meta MetaDataFlyweight
	meta.Wrap(cncBuffer, 0)

	return meta.valuesBuf.Get()
}

func CreateErrorLogBuffer(cncFile *memmap.File) *atomic.Buffer {
	cncBuffer := atomic.MakeBuffer(cncFile.GetMemoryPtr(), cncFile.GetMemorySize())
	var meta MetaDataFlyweight
	meta.Wrap(cncBuffer, 0)

	return meta.errorBuf.Get()
}

func CncVersion(cncFile *memmap.File) int32 {
	cncBuffer := atomic.MakeBuffer(cncFile.GetMemoryPtr(), cncFile.GetMemorySize())
	var meta MetaDataFlyweight
	meta.Wrap(cncBuffer, 0)

	return meta.cncVersion.Get()
}

func ClientLivenessTimeout(cncFile *memmap.File) int64 {
	cncBuffer := atomic.MakeBuffer(cncFile.GetMemoryPtr(), cncFile.GetMemorySize())
	var meta MetaDataFlyweight
	meta.Wrap(cncBuffer, 0)

	return meta.clientLivenessTo.Get()
}
