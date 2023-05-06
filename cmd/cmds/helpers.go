package cmds

import (
	"context"

	"github.com/pkg/errors"
	"github.com/trevorstarick/qobuz-sync/client"
)

func GetClientFromContext(ctx context.Context) (*client.Client, error) {
	switch t := ctx.Value(client.Key{}).(type) {
	case *client.Client:
		return t, nil
	case nil:
		return nil, errors.New("client is nil")
	default:
		return nil, errors.New("client is not a *Client")
	}
}
