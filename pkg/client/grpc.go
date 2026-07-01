package client

import (
	"github.com/knadh/koanf/v2"
	"go.uber.org/fx"

	categoryv1 "github.com/Sokol111/ecommerce-category-query-service-api/gen/connect/category_query/v1"
	grpcclient "github.com/Sokol111/ecommerce-commons/pkg/grpc/client"
)

// Module wires a native gRPC client for CategoryQueryService.
// Configuration is read from koanf under key "category-query.grpc".
func Module() fx.Option {
	return fx.Module("category-query-grpc-client",
		fx.Provide(func(k *koanf.Koanf) (grpcclient.Config, error) {
			return grpcclient.LoadConfig(k, "category-query.grpc")
		}, fx.Private),
		fx.Provide(grpcclient.NewConn, fx.Private),
		fx.Provide(categoryv1.NewCategoryQueryServiceClient),
	)
}
