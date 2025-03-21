package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/grafana/tempo/cmd/tempo/app"
	"gopkg.in/yaml.v3"
)

const ManifestPath = "docs/sources/tempo/configuration/manifest.md"

const Cmd = "go run pkg/docsgen/generate_manifest.go"

var Manifest = fmt.Sprintf(`---
title: Manifest
description: This manifest lists of all Tempo options and their defaults.
weight: 110
---
[//]: # 'THIS FILE IS GENERATED AUTOMATICALLY BY %s'
[//]: # 'DO NOT EDIT THIS FILE DIRECTLY'

# Manifest

This document is a reference for all Tempo options and their defaults. If you are just getting
started with Tempo, refer to [Tempo examples](https://github.com/grafana/tempo/tree/main/example/docker-compose)
and other [configuration documentation](../). Most installations will require only setting 10 to 20 of these options.

## Complete configuration

`, Cmd)

func main() {
	newConfig := app.NewDefaultConfig()
	// Override values that depend on the host specifics
	const hostname = "hostname"
	newConfig.Distributor.DistributorRing.InstanceID = hostname
	newConfig.Compactor.ShardingRing.InstanceID = hostname
	newConfig.Ingester.LifecyclerConfig.ID = hostname
	newConfig.Ingester.LifecyclerConfig.InfNames = []string{"eth0"}
	newConfig.Generator.Ring.InstanceID = hostname
	newConfig.BackendWorker.Ring.InstanceID = hostname
	newConfig.Generator.InstanceID = hostname
	newConfig.BlockBuilder.InstanceID = hostname

	newConfigBytes, err := yaml.Marshal(newConfig)
	if err != nil {
		panic(err)
	}
	newManifest := Manifest + "```yaml\n" + string(newConfigBytes) + "```\n"

	err = os.WriteFile(ManifestPath, []byte(newManifest), 0o644)
	if err != nil {
		panic(err)
	}

	cmd := exec.Command("git", "diff", "--exit-code", ManifestPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatalf("The manifest with the default Tempo configuration has changed. Please run '%s' and commit the changes.", Cmd)
	}
}
