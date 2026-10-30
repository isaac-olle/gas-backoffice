package vo

import "context"

type pictureContextKey struct{}

func WithPicture(ctx context.Context, picture string) context.Context {
	return context.WithValue(ctx, pictureContextKey{}, picture)
}

func PictureFromContext(ctx context.Context) (string, bool) {
	picture, ok := ctx.Value(pictureContextKey{}).(string)
	return picture, ok && picture != ""
}
