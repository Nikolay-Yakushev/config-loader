package ports
import (
	"context"
)

type Config interface{
	GetValue(ctx context.Context)(string, error)
}