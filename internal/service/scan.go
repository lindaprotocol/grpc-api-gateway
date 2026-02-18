package service

import (
    "context"
    
    "github.com/lindaprotocol/grpc-gateway/pkg/api/protocol"
    "google.golang.org/grpc"
    "gorm.io/gorm"
)

type ScanService struct {
    protocol.UnimplementedScanServiceServer
    db          *gorm.DB
    walletClient protocol.WalletClient
}

func NewScanService(conn *grpc.ClientConn, db *gorm.DB) *ScanService {
    return &ScanService{
        walletClient: protocol.NewWalletClient(conn),
        db:          db,
    }
}

// Implement your scan service methods here
func (s *ScanService) GetHomepageBundle(ctx context.Context, req *protocol.EmptyMessage) (*protocol.HomepageBundle, error) {
    // Your implementation
    return &protocol.HomepageBundle{}, nil
}
