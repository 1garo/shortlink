package service

import "fmt"

type ServiceError struct {
	Err  error
	Code int
}

func (s *ServiceError) Error() string {
	return fmt.Sprintf("status %d: err %v", s.Code, s.Err)
}
