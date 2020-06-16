// Copyright 2020 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package data

const (
	TypeBool   = "bool"
	TypeInt    = "int"
	TypeUint   = "uint" // uint is a shortcut for int starts from zero rather than a representation of uint type
	TypeFloat  = "float"
	TypeString = "string"
	TypeSize   = "size"
	TypeTime   = "time"
	TypeList   = "list"
	TypeChoice = "choice"
	TypeMap    = "map"
	TypeStruct = "struct" // struct is an alias of map
)
