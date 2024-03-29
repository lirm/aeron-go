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

package aeron

import (
	"github.com/lirm/aeron-go/aeron/atomic"
	"github.com/lirm/aeron-go/aeron/counters"
)

type Counter struct {
	registrationId  int64
	clientConductor *ClientConductor
	atomicCounter   *counters.AtomicCounter
	counterId       int32
	isClosed        atomic.Bool
}

func NewCounter(
	registrationId int64,
	clientConductor *ClientConductor,
	counterId int32) (*Counter, error) {

	counter := new(Counter)
	counter.registrationId = registrationId
	counter.clientConductor = clientConductor
	atomicCounter, err := counters.NewAtomicCounter(clientConductor.CounterReader(), counterId)
	if err != nil {
		return nil, err
	}
	counter.atomicCounter = atomicCounter
	counter.counterId = counterId

	return counter, nil
}

func (c *Counter) RegistrationId() int64 {
	return c.registrationId
}

func (c *Counter) Id() int32 {
	return c.counterId
}

func (c *Counter) Counter() *counters.AtomicCounter {
	return c.atomicCounter
}

func (c *Counter) Close() error {
	if c.isClosed.CompareAndSet(false, true) {
		if c.clientConductor != nil {
			return c.clientConductor.releaseCounter(*c)
		}
	}
	return nil
}

func (c *Counter) IsClosed() bool {
	return c.isClosed.Get()
}
