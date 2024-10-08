// Copyright (c) 2024 Intel Corporation
//
// SPDX-License-Identifier: Apache-2.0
package api

import (
	"context"
	"fmt"
	"net"

	cdhapi "github.com/confidential-containers/guest-components/confidential-data-hub/golang/pkg/api/cdhapi"
	"github.com/containerd/ttrpc"
)

const (
	CDHTtrpcSocket     = "/run/confidential-containers/cdh.sock"
	SealedSecretPrefix = "sealed."
)

type cdhTtrpcClient struct {
	conn               net.Conn
	sealedSecretClient cdhapi.SealedSecretServiceService
	secureMountClient  cdhapi.SecureMountServiceService
}

func CreateCDHTtrpcClient(sockAddress string) (*cdhTtrpcClient, error) {
	conn, err := net.Dial("unix", sockAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to cdh sock %q: %w", sockAddress, err)
	}

	ttrpcClient := ttrpc.NewClient(conn)
	sealedSecretClient := cdhapi.NewSealedSecretServiceClient(ttrpcClient)
	securelMountClient := cdhapi.NewSecureMountServiceClient(ttrpcClient)

	c := &cdhTtrpcClient{
		conn:               conn,
		sealedSecretClient: sealedSecretClient,
		secureMountClient:  securelMountClient,
	}
	return c, nil
}

func (c *cdhTtrpcClient) Close() error {
	return c.conn.Close()
}

func (c *cdhTtrpcClient) UnsealSecret(ctx context.Context, secret string) (string, error) {
	input := cdhapi.UnsealSecretInput{Secret: []byte(secret)}
	output, err := c.sealedSecretClient.UnsealSecret(ctx, &input)
	if err != nil {
		return "", fmt.Errorf("failed to unseal secret: %w", err)
	}

	return string(output.GetPlaintext()[:]), nil
}

func (c *cdhTtrpcClient) SecureMount(ctx context.Context, volume_type string, options map[string]string, flags []string, mountpoint string) (string, error) {
	input := cdhapi.SecureMountRequest{VolumeType: volume_type, Options: options, Flags: flags, MountPoint: mountpoint}
	output, err := c.secureMountClient.SecureMount(ctx, &input)
	if err != nil {
		return "", fmt.Errorf("failed to unseal secret: %w", err)
	}

	return output.GetMountPath(), nil
}
