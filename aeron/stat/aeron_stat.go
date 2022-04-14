/*
Copyright 2016-2018 Stanislav Liberman

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

package main

import (
	"fmt"
	"github.com/corymonroe-coinbase/aeron-go/aeron"
	"github.com/corymonroe-coinbase/aeron-go/aeron/counters"
	"time"
)

var ansiCls = "\u001b[2J"
var ansiHome = "\u001b[H"

func main() {
	ctx := aeron.NewContext()
	counterFile, _, _ := counters.MapFile(ctx.CncFileName())

	reader := counters.NewReader(counterFile.ValuesBuf.Get(), counterFile.MetaDataBuf.Get())

	for {
		fmt.Print(ansiCls + ansiHome)
		reader.Scan(func(counter counters.Counter) {
			fmt.Printf("%3d: %20d - %s\n", counter.Id, counter.Value, counter.Label)
		})
		time.Sleep(time.Second)
	}
}
