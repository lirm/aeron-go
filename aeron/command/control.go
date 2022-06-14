/*
Copyright 2016 Stanislav Liberman
Copyright (C) 2022 Talos, Inc.

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

package command

const (
	// Clients to Media Driver

	// AddPublication command from client
	AddPublication int32 = 0x01
	// RemovePublication command from client
	RemovePublication int32 = 0x02
	// AddExclusivePublication from client
	AddExclusivePublication int32 = 0x03
	// AddSubscription command from client
	AddSubscription int32 = 0x04
	// RemoveSubscription command from client
	RemoveSubscription int32 = 0x05
	// ClientKeepalive message from client */
	ClientKeepalive int32 = 0x06
	// AddDestination to add an MDC
	AddDestination = 0x07
	// RemoveDestination to remove an MDC
	RemoveDestination = 0x08
	// AddCounter command from client
	AddCounter = 0x09
	// RemoveCounter command from client
	RemoveCounter = 0x0A
	// ClientClose command from client
	ClientClose = 0x0B
	// RecordingPosition is rhe position a recording has reached when being archived
	RecordingPosition = 0x64 // 100
)
