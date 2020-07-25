package main

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
)

func getActionServices(files []*desc.FileDescriptor) ([]*desc.ServiceDescriptor, error) {
	var services []*desc.ServiceDescriptor

	for _, srv := range getServices(files) {
		if opts := srv.GetServiceOptions(); opts != nil {
			extensions, err := proto.ExtensionDescs(opts)
			if err != nil {
				return nil, fmt.Errorf("failed to examine service extensions: %w", err)
			}

			for _, ext := range extensions {
				if ext != nil && ext.Field == isActionSetExtensionFieldNumber {
					services = append(services, srv)
				}
			}
		}
	}

	return services, nil
}

func getServices(files []*desc.FileDescriptor) []*desc.ServiceDescriptor {
	var services []*desc.ServiceDescriptor

	for _, file := range files {
		services = append(services, file.GetServices()...)
	}

	return services
}
