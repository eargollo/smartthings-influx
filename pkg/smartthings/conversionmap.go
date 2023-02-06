package smartthings

import (
	"fmt"
)

type ConversionMap map[string]map[string]float64

func ParseConversionMap(valuemap map[string]interface{}) (ConversionMap, error) {
	conversionMap := make(ConversionMap)
	for key, value := range valuemap {
		innermap, ok := value.(map[string]any)
		if !ok {
			return map[string]map[string]float64{}, fmt.Errorf("could not parse valuemap, it shold be in the format of metric to a map of values, got %v", value)
		}

		_, ok = conversionMap[key]
		if !ok {
			conversionMap[key] = make(map[string]float64)
		}

		for inkey, inval := range innermap {
			var inFloat float64
			_, ok := inval.(int)
			if ok {
				inFloat = float64(inval.(int))
			} else {
				inFloat, ok = inval.(float64)
				if !ok {
					return map[string]map[string]float64{}, fmt.Errorf("could not convert %v to a number for metric %s", inval, key)
				}
			}

			list := conversionMap[key]
			list[inkey] = inFloat
			conversionMap[key] = list
		}
	}
	return conversionMap, nil
}

func (c ConversionMap) ConvertValueToFloat(metric string, value any) (float64, error) {
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
			return 0, fmt.Errorf("there is no value map for metric '%s' and value '%s', can't convert", metric, stValue)
		}
		return metricMap[stValue], nil
	}
	return 0, nil
}
