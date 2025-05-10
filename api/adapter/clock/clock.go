package clock

import (
	"time"

	"gitlab.com/digeon-inc/japan-association-for-clinical-engineers/e-privado/api/usecase/output_port"
)

type Clock struct{}

func New() output_port.Clock {
	return Clock{}
}

func (c Clock) Now() time.Time {
	return time.Now()
}
