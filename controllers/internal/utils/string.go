/*
Copyright 2021 Absa Group Limited

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

// Package utils provides common functionality to gslb controller
package utils

import (
	"encoding/json"
	"fmt"
)

// ToString converts type to formatted string. If value is struct, function returns formatted JSON. Function retrieves
// null for nil pointer references. Function doesn't return error. In case of marshal error it converts with %v formatter
// Only two possible errors can occur e.g.:
//	UnsupportedTypeError ToString(make(chan int));
//	UnsupportedValueError ToString(math.Inf(1));
//	In both cases function retrieves expected result. The pointer address in the first while "+Inf" in second
func ToString(v interface{}) string {
	value, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return fmt.Sprintf("%v", v)
	}
	return string(value)
}