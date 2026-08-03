package fxconfig

import (
	"context"

	"go.uber.org/fx"
	"google.golang.org/grpc"

	categoryv1 "github.com/Sokol111/ecommerce-category-query-service-api/gen/go/category_query/v1"
	"github.com/Sokol111/ecommerce-commons/pkg/core/config"
	grpcclient "github.com/Sokol111/ecommerce-commons/pkg/http/grpc/client"
)

func NewGrpcClientsModule() fx.Option {
	return fx.Module("category-query-grpc-clients",
		fx.Provide(func(loader *config.Loader) (grpcclient.Config, error) {
			return grpcclient.LoadConfig(loader, "category-query.grpc")
		}, fx.Private),
		fx.Provide(grpcclient.NewGrpcConnWithTokenSource, fx.Private),
		fx.Provide(func(conn *grpc.ClientConn) grpc.ClientConnInterface {
			return conn
		}, fx.Private),
		fx.Provide(categoryv1.NewCategoryQueryServiceClient),
		fx.Invoke(func(lc fx.Lifecycle, conn *grpc.ClientConn) {
			lc.Append(fx.Hook{
				OnStop: func(context.Context) error {
					return conn.Close()
				},
			})
		}),
	)
}
