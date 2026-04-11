package grpc

import (
	"context"
	"fmt"

	pb "github.com/bookmarketing/marketing-saas/api-go/pkg/grpc"
	"github.com/bookmarketing/marketing-saas/api-go/internal/application/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

// GRPCServer represents the gRPC server
type GRPCServer struct {
	pb.UnimplementedCampaignServiceServer
	pb.UnimplementedABTestServiceServer
	pb.UnimplementedAnalyticsServiceServer
	campaignService  *service.CampaignService
	dashboardService *service.DashboardService
}

// NewGRPCServer creates a new gRPC server
func NewGRPCServer(campaignService *service.CampaignService, dashboardService *service.DashboardService) *GRPCServer {
	return &GRPCServer{
		campaignService:  campaignService,
		dashboardService: dashboardService,
	}
}

// Start starts the gRPC server
func (s *GRPCServer) Start(port string) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	
	pb.RegisterCampaignServiceServer(grpcServer, s)
	pb.RegisterABTestServiceServer(grpcServer, s)
	pb.RegisterAnalyticsServiceServer(grpcServer, s)
	
	reflection.Register(grpcServer)

	fmt.Printf("gRPC server started on port %s\n", port)
	return grpcServer.Serve(lis)
}

// GetCampaign implements CampaignService.GetCampaign
func (s *GRPCServer) GetCampaign(ctx context.Context, req *pb.GetCampaignRequest) (*pb.CampaignResponse, error) {
	campaign, err := s.campaignService.GetCampaign(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return campaignToProto(campaign), nil
}

// ListCampaigns implements CampaignService.ListCampaigns
func (s *GRPCServer) ListCampaigns(ctx context.Context, req *pb.ListCampaignsRequest) (*pb.ListCampaignsResponse, error) {
	campaigns, err := s.campaignService.GetAllCampaigns(ctx, "", int(req.Limit), int(req.Offset))
	if err != nil {
		return nil, err
	}

	response := &pb.ListCampaignsResponse{
		Total: int32(len(campaigns)),
	}

	for _, c := range campaigns {
		response.Campaigns = append(response.Campaigns, campaignToProto(c))
	}

	return response, nil
}

// CreateCampaign implements CampaignService.CreateCampaign
func (s *GRPCServer) CreateCampaign(ctx context.Context, req *pb.CreateCampaignRequest) (*pb.CampaignResponse, error) {
	// Implementation would convert proto to domain entity and call service
	return &pb.CampaignResponse{}, nil
}

// UpdateCampaign implements CampaignService.UpdateCampaign
func (s *GRPCServer) UpdateCampaign(ctx context.Context, req *pb.UpdateCampaignRequest) (*pb.CampaignResponse, error) {
	return &pb.CampaignResponse{}, nil
}

// DeleteCampaign implements CampaignService.DeleteCampaign
func (s *GRPCServer) DeleteCampaign(ctx context.Context, req *pb.DeleteCampaignRequest) (*pb.DeleteResponse, error) {
	err := s.campaignService.DeleteCampaign(ctx, req.Id)
	if err != nil {
		return &pb.DeleteResponse{Success: false, Message: err.Error()}, nil
	}
	return &pb.DeleteResponse{Success: true, Message: "Campaign deleted successfully"}, nil
}

// PauseCampaign implements CampaignService.PauseCampaign
func (s *GRPCServer) PauseCampaign(ctx context.Context, req *pb.PauseCampaignRequest) (*pb.OperationResponse, error) {
	err := s.campaignService.PauseCampaign(ctx, req.Id)
	if err != nil {
		return &pb.OperationResponse{Success: false, Message: err.Error()}, nil
	}
	return &pb.OperationResponse{Success: true, Message: "Campaign paused successfully"}, nil
}

// ActivateCampaign implements CampaignService.ActivateCampaign
func (s *GRPCServer) ActivateCampaign(ctx context.Context, req *pb.ActivateCampaignRequest) (*pb.OperationResponse, error) {
	err := s.campaignService.ActivateCampaign(ctx, req.Id)
	if err != nil {
		return &pb.OperationResponse{Success: false, Message: err.Error()}, nil
	}
	return &pb.OperationResponse{Success: true, Message: "Campaign activated successfully"}, nil
}

// GetDashboard implements AnalyticsService.GetDashboard
func (s *GRPCServer) GetDashboard(ctx context.Context, req *pb.GetDashboardRequest) (*pb.DashboardResponse, error) {
	dashboard, err := s.dashboardService.GetOverview(ctx, "")
	if err != nil {
		return nil, err
	}

	response := &pb.DashboardResponse{
		TotalCampaigns:  dashboard.TotalCampaigns,
		ActiveCampaigns: dashboard.ActiveCampaigns,
		TotalSpend:      dashboard.TotalSpend,
		TotalRevenue:    dashboard.TotalRevenue,
		Roi:             dashboard.ROI,
		ConversionRate:  dashboard.ConversionRate,
	}

	for _, channel := range dashboard.TopChannels {
		response.TopChannels = append(response.TopChannels, &pb.ChannelPerformance{
			Name:  channel.Name,
			Value: channel.Value,
		})
	}

	return response, nil
}

// GetMetrics implements AnalyticsService.GetMetrics
func (s *GRPCServer) GetMetrics(ctx context.Context, req *pb.GetMetricsRequest) (*pb.MetricsResponse, error) {
	// Implementation would fetch metrics
	return &pb.MetricsResponse{}, nil
}

// StreamMetrics implements AnalyticsService.StreamMetrics
func (s *GRPCServer) StreamMetrics(req *pb.StreamMetricsRequest, stream pb.AnalyticsService_StreamMetricsServer) error {
	// Implementation would stream metrics
	return nil
}

// Helper function to convert domain Campaign to proto
func campaignToProto(campaign interface{}) *pb.CampaignResponse {
	// This would properly convert the domain entity to proto
	// For now, return empty response
	return &pb.CampaignResponse{}
}
