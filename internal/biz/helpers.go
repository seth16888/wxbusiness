package biz

import (
	"context"
	"fmt"

	"github.com/seth16888/wxbusiness/internal/data/entities"
)

func GetAppInfoFromCtx(ctx context.Context) (*entities.PlatformApp, error) {
  appInfo := ctx.Value("APP")
  if appInfo == nil {
    return nil, fmt.Errorf("appInfo is nil")
  }
  app, ok := appInfo.(*entities.PlatformApp)
  if !ok {
    return nil, fmt.Errorf("appInfo is not invalid")
  }

  return app, nil
}
