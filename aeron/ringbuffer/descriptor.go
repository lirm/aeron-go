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

package rb

import (
	"github.com/lirm/aeron-go/aeron/util"
)

var Descriptor = struct {
	TAIL_POSITION_OFFSET       int32
	HEAD_CACHE_POSITION_OFFSET int32
	HEAD_POSITION_OFFSET       int32
	CORRELATION_COUNTER_OFFSET int32
	CONSUMER_HEARTBEAT_OFFSET  int32
	TRAILER_LENGTH             int32
}{
	util.CACHE_LINE_LENGTH * 2,
	util.CACHE_LINE_LENGTH * 4,
	util.CACHE_LINE_LENGTH * 6,
	util.CACHE_LINE_LENGTH * 8,
	util.CACHE_LINE_LENGTH * 10,
	util.CACHE_LINE_LENGTH * 12,
}
