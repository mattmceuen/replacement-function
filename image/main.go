// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

// Package main implements a validator function run by `kustomize config run`
package main

import (
	"fmt"
	"os"

	"sigs.k8s.io/kustomize/kyaml/kio"
	"sigs.k8s.io/kustomize/kyaml/yaml"

	"sigs.k8s.io/kustomize/api/k8sdeps/kunstruct"
	"sigs.k8s.io/kustomize/api/resmap"
	"sigs.k8s.io/kustomize/api/resource"
)

func main() {
	rw := &kio.ByteReadWriter{
		Reader:                os.Stdin,
		Writer:                os.Stdout,
		OmitReaderAnnotations: true,
		KeepReaderAnnotations: true,
	}
	p := kio.Pipeline{
		Inputs:  []kio.Reader{rw}, // read the inputs into a slice
		Filters: []kio.Filter{kubevalFilter{rw: rw}},
		Outputs: []kio.Writer{rw}, // copy the inputs to the output
	}
	if err := p.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

// kubevalFilter implements kio.Filter
type kubevalFilter struct {
	rw *kio.ByteReadWriter
}

func (f kubevalFilter) filterToPluginConfig() error {
	return KustomizePlugin.Config(nil, []byte(f.rw.FunctionConfig.MustString()))
}

// Filter checks each resource for validity, otherwise returning an error.
func (f kubevalFilter) Filter(in []*yaml.RNode) ([]*yaml.RNode, error) {
	err := f.filterToPluginConfig()
	if err != nil {
		return nil, err
	}

	resmap, err := rNodesToResMap(in)
	if err != nil {
		return nil, err
	}

	err = KustomizePlugin.Transform(resmap)
	if err != nil {
		return nil, err
	}

	rnodes, err := resMapToRNodes(resmap)
	if err != nil {
		return nil, err
	}

	return rnodes, nil
}

func rNodesToResMap(rnodes []*yaml.RNode) (resmap.ResMap, error) {
	resmap := resmap.New()
	factory := resource.NewFactory(kunstruct.NewKunstructuredFactoryImpl())
	for _, rnode := range rnodes {
		yaml, err := rnode.String()
		if err != nil {
			return nil, err
		}
		resource, err := factory.FromBytes([]byte(yaml))
		if err != nil {
			return nil, err
		}
		resmap.Append(resource)

	}
	return resmap, nil
}

func resMapToRNodes(resmap resmap.ResMap) ([]*yaml.RNode, error) {
	rnodes := []*yaml.RNode{}
	for _, resId := range resmap.AllIds() {
		resource, err := resmap.GetById(resId)
		if err != nil {
			return nil, err
		}

		yml, err := resource.AsYAML()
		if err != nil {
			return nil, err
		}

		rnode, err := yaml.Parse(string(yml))
		if err != nil {
			return nil, err
		}

		rnodes = append(rnodes, rnode)
	}
	return rnodes, nil
}
