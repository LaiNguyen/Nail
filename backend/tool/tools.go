package tool

import (
	uuid "github.com/nu7hatch/gouuid"
)

func UUID() string {
	u, _ := uuid.NewV4()
	return u.String()
}
