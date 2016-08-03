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

	/** Add Publication */
	AddPublication int32 = 0x01
	/** Remove Publication */
	RemovePublication int32 = 0x02
	/** Add Subscriber */
	AddSubscription int32 = 0x04
	/** Remove Subscriber */
	RemoveSubscription int32 = 0x05
	/** Keepalive from Client */
	ClientKeepalive int32 = 0x06
)
