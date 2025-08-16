package core

import "context"

type Sender interface {
	Process(context.Context, []byte) error
}
