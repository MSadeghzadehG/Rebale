package main

import "errors"

// in the name of God.

type Rebale interface {
	Connect(address string) error
	Ping() error
	Get(key string) (io.Reader, error)
	Set(key string, value io.Reader, length int) error
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

func (r *MyRebaleImpl) Get(key string) (io.Reader, error) {
	return nil, errors.New("implement me")
}

func (r *MyRebaleImpl) Set(key string, value io.Reader, length int) error {
	return errors.New("implement me")
}

func (r *MyRebaleImpl) Close() error {
	return errors.New("implement me")
}
