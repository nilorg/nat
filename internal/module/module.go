package module

import (
	"context"

	"github.com/nilorg/nat/internal/module/store"
)

// Init 初始化
func Init(ctx context.Context) {
	store.Init()
}
