// Code generated by falafel 0.9.1. DO NOT EDIT.
// source: walletunlocker.proto

// +build js

package lnrpc

import (
	"context"

	gateway "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

func RegisterWalletUnlockerJSONCallbacks(registry map[string]func(ctx context.Context,
	conn *grpc.ClientConn, reqJSON string, callback func(string, error))) {

	marshaler := &gateway.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames:   true,
			EmitUnpopulated: true,
		},
	}

	registry["lnrpc.WalletUnlocker.GenSeed"] = func(ctx context.Context,
		conn *grpc.ClientConn, reqJSON string, callback func(string, error)) {

		req := &GenSeedRequest{}
		err := marshaler.Unmarshal([]byte(reqJSON), req)
		if err != nil {
			callback("", err)
			return
		}

		client := NewWalletUnlockerClient(conn)
		resp, err := client.GenSeed(ctx, req)
		if err != nil {
			callback("", err)
			return
		}

		respBytes, err := marshaler.Marshal(resp)
		if err != nil {
			callback("", err)
			return
		}
		callback(string(respBytes), nil)
	}

	registry["lnrpc.WalletUnlocker.InitWallet"] = func(ctx context.Context,
		conn *grpc.ClientConn, reqJSON string, callback func(string, error)) {

		req := &InitWalletRequest{}
		err := marshaler.Unmarshal([]byte(reqJSON), req)
		if err != nil {
			callback("", err)
			return
		}

		client := NewWalletUnlockerClient(conn)
		resp, err := client.InitWallet(ctx, req)
		if err != nil {
			callback("", err)
			return
		}

		respBytes, err := marshaler.Marshal(resp)
		if err != nil {
			callback("", err)
			return
		}
		callback(string(respBytes), nil)
	}

	registry["lnrpc.WalletUnlocker.UnlockWallet"] = func(ctx context.Context,
		conn *grpc.ClientConn, reqJSON string, callback func(string, error)) {

		req := &UnlockWalletRequest{}
		err := marshaler.Unmarshal([]byte(reqJSON), req)
		if err != nil {
			callback("", err)
			return
		}

		client := NewWalletUnlockerClient(conn)
		resp, err := client.UnlockWallet(ctx, req)
		if err != nil {
			callback("", err)
			return
		}

		respBytes, err := marshaler.Marshal(resp)
		if err != nil {
			callback("", err)
			return
		}
		callback(string(respBytes), nil)
	}

	registry["lnrpc.WalletUnlocker.ChangePassword"] = func(ctx context.Context,
		conn *grpc.ClientConn, reqJSON string, callback func(string, error)) {

		req := &ChangePasswordRequest{}
		err := marshaler.Unmarshal([]byte(reqJSON), req)
		if err != nil {
			callback("", err)
			return
		}

		client := NewWalletUnlockerClient(conn)
		resp, err := client.ChangePassword(ctx, req)
		if err != nil {
			callback("", err)
			return
		}

		respBytes, err := marshaler.Marshal(resp)
		if err != nil {
			callback("", err)
			return
		}
		callback(string(respBytes), nil)
	}
}
