// Copyright (c) 2023 Benjamin Borbe All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package parse

import (
	"context"
	"fmt"
	"strconv"
)

func ParseFloat64(ctx context.Context, value interface{}) (float64, error) {
	switch v := value.(type) {
	case int:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	case string:
		return strconv.ParseFloat(v, 64)
	case fmt.Stringer:
		return ParseFloat64(ctx, v.String())
	default:
		return ParseFloat64(ctx, fmt.Sprintf("%v", value))
	}
}

func ParseFloat64Default(ctx context.Context, value interface{}, defaultValue float64) float64 {
	result, err := ParseFloat64(ctx, value)
	if err != nil {
		return defaultValue
	}
	return result
}
