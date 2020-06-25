package main

import (
	"context"
	"github.com/bmsandoval/vehicle-api/pkg/appcontext"
	"testing"
)

func BenchmarkAcquireServers(b *testing.B) {
	for i := 0; i < b.N; i++ {
		appCtx := appcontext.Context{
			GoContext: context.Background(),
		}

		_, _, _ = AcquireServers(appCtx, nil)
	}
}
