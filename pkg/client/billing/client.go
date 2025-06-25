package billing

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"time"

	proto "Auth-service/pkg/proto/billing"
)

type BillingClient struct {
	conn   *grpc.ClientConn
	client proto.BillingServiceClient
}

func BillingAdapter(cfg *BillingClientConfig) (*BillingClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	maxMsgSize := 10500000

	dialOption := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(maxMsgSize),
			grpc.MaxCallSendMsgSize(maxMsgSize),
		),
	}

	address := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	conn, err := grpc.DialContext(ctx, address, dialOption...)
	if err != nil {
		log.Printf("Did not connect: %v", err)
		return nil, fmt.Errorf("failed to connect to %s: %v", address, err)
	}

	client := proto.NewBillingServiceClient(conn)

	return &BillingClient{
		conn,
		client,
	}, nil
}

func (s *BillingClient) CreateWallet(ctx context.Context, userID, currencyCode string) error {
	_, err := s.client.CreateWallet(ctx, &proto.CreateWalletRequest{
		UserId:       userID,
		CurrencyCode: currencyCode,
	})
	if err != nil {
		return err
	}
	return nil
}
