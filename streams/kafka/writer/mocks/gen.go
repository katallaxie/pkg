package mocks

//go:generate go run github.com/golang/mock/mockgen -destination=mocks.go -package mocks github.com/katallaxie/pkg/streams/kafka/writer Writer
