package monitor

import (
	"fmt"
	"strings"
)

type ConversionMap map[string]map[string]float64

func (c ConversionMap) Convert(metric string, value any) (float64, error) {
	_, ok := value.(float64)
	if ok {
		return value.(float64), nil
	}

	_, ok = value.(string)
	if ok {
		stValue := value.(string)
		// Check if there is a map for metric
		metricMap, ok := c[metric]
		if !ok {
			// Try lowercase since yaml files ignore case on tag
			metricMap, ok = c[strings.ToLower(metric)]
			if !ok {
				return 0, fmt.Errorf("there is no value map for metric '%s' and value '%s', can't convert", metric, stValue)
			}
		}

		result, ok := metricMap[stValue]
		if !ok {
			result, ok = metricMap[strings.ToLower(stValue)]
			if !ok {
				return 0, fmt.Errorf("there is no value map for metric '%s' and value '%v', can't convert", metric, value)
			}
		}

		return result, nil
	}
	return 0, fmt.Errorf("there is no value map for metric '%s' and value '%v', can't convert", metric, value)
}
