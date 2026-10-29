package usecase

import "context"

type actorKey struct{}

// Actor adalah user yang sedang melakukan request (diisi dari JWT oleh handler).
type Actor struct {
	UserID  string
	RoleIDs []string
}

// WithActor menyisipkan actor ke context agar usecase bisa memeriksa
// siapa yang melakukan aksi (mis. mencegah eskalasi role).
func WithActor(ctx context.Context, actor Actor) context.Context {
	return context.WithValue(ctx, actorKey{}, actor)
}

// ActorFrom mengambil actor dari context.
func ActorFrom(ctx context.Context) (Actor, bool) {
	a, ok := ctx.Value(actorKey{}).(Actor)
	return a, ok
}
