/*
Copyright 2016 Stanislav Liberman
Copyright 2020 Evan Wies

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

package util

import "fmt"

func SemanticVersionCompose(major uint8, minor uint8, patch uint8) uint32 {
	return (uint32(major) << 16) | (uint32(minor) << 8) | uint32(patch)
}

func SemanticVersionMajor(version uint32) uint8 {
	return uint8((version >> 16) & 0xFF)
}

func SemanticVersionMinor(version uint32) uint8 {
	return uint8((version >> 8) & 0xFF)
}

func SemanticVersionPatch(version uint32) uint8 {
	return uint8(version & 0xFF)
}

func SemanticVersionToString(version uint32) string {
	return fmt.Sprintf("%d.%d.%d",
		SemanticVersionMajor(version),
		SemanticVersionMinor(version),
		SemanticVersionPatch(version))
}
