package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

// 用户端网关调用 ai-mate-server 的 gRPC 连接。
type Client struct {
	// ai-mate-server gRPC 地址。
	target string
	// 单次 RPC 默认超时时间。
	timeout time.Duration
	// 底层 gRPC 连接。
	conn *grpc.ClientConn
}

// 创建 ai-mate-server gRPC 连接，连接失败时阻止服务启动。
func New(target string, timeout time.Duration) (*Client, error) {
	if strings.TrimSpace(target) == "" {
		return nil, fmt.Errorf("core grpc target is required")
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("create core grpc client %s: %w", target, err)
	}
	if err := waitReady(ctx, conn); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("connect core grpc %s: %w", target, err)
	}
	return &Client{target: target, timeout: timeout, conn: conn}, nil
}

func waitReady(ctx context.Context, conn *grpc.ClientConn) error {
	conn.Connect()
	for {
		state := conn.GetState()
		if state == connectivity.Ready {
			return nil
		}
		if !conn.WaitForStateChange(ctx, state) {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			return fmt.Errorf("wait for grpc ready state failed")
		}
		conn.Connect()
	}
}

// 提供具体业务 client 初始化 proto stub 所需的底层连接。
func (c *Client) Conn() *grpc.ClientConn {
	return c.conn
}

// 单次 Core RPC 默认超时时间。
func (c *Client) Timeout() time.Duration {
	return c.timeout
}

// 关闭底层 gRPC 连接。
func (c *Client) Close() error {
	return c.conn.Close()
}
