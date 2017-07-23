package poloniexapi

import (
	"strconv"
)

func interfaceTo2FloatArray(in interface{}) ([][2]float64, error) {
	var err error
	out := make([][2]float64, 0)

	for _, v := range in.([]interface{}) {
		var subout [2]float64

		subout[0], err = strconv.ParseFloat(v.([]interface{})[0].(string), 64)
		if err != nil {
			return out, err
		}
		subout[1] = v.([]interface{})[1].(float64)
		if err != nil {
			return out, err
		}

		out = append(out, subout)
	}

	return out, nil
}
