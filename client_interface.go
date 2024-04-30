package main

import "errors"

// in the name of God.

type Rebale interface {
	Connect(address string) error
	Ping() error
	Get(key string) (interface{}, error)
	Set(key string, value interface{}, length int) error
	Close() error
}

type MyRebaleImpl struct {
	// TODO fill this with your required fileds
}

func (r *MyRebaleImpl) Connect(address string) error {
	return errors.New("implement me")
}

func (r *MyRebaleImpl) Ping() error {
	return errors.New("implement me")
}

func (r *MyRebaleImpl) Get(key string) (interface{}, error) {
	return nil, errors.New("implement me")
}

func (r *MyRebaleImpl) Set(key string, value interface{}, length int) error {
	return errors.New("implement me")
}

func (r *MyRebaleImpl) Close() error {
	return errors.New("implement me")
}
