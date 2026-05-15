package smartbox

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type stringOrInt int

func (v *stringOrInt) UnmarshalJSON(data []byte) error {
	var intVal int
	err := json.Unmarshal(data, &intVal)
	if err == nil {
		*v = stringOrInt(intVal)
		return nil
	}

	var strVal string
	err = json.Unmarshal(data, &strVal)
	if err == nil {
		intVal, err := strconv.Atoi(strVal)
		if err == nil {
			*v = stringOrInt(intVal)
			return nil
		}
	}

	return fmt.Errorf("invalid value for stringOrInt: %s", string(data))
}
