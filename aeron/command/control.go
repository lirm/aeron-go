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

package command

const (
	// Clients to Media Driver

	// AddPublication command from client
	AddPublication int32 = 0x01
	// RemovePublication command from client
	RemovePublication int32 = 0x02
	// Add Exclusive Publication
	AddExclusivePublication int32 = 0x03
	// AddSubscription command from client
	AddSubscription int32 = 0x04
	// RemoveSubscription command from client
	RemoveSubscription int32 = 0x05
	// ClientKeepalive message from client */
	ClientKeepalive int32 = 0x06

	AddDestination = 0x07
	// Remove Destination
	RemoveDestination = 0x08
	// Add Counter
	AddCounter = 0x09
	// Remove Counter
	RemoveCounter = 0x0A
	// Client Close
	ClientClose = 0x0B
)
