package firestorex

import (
	"context"
	"fmt"
	"reflect"

	"cloud.google.com/go/firestore"
	"golang.org/x/sync/errgroup"
)

type settings struct {
	maxConcurrecy  int
	batchChunkSize int
}

type Option func(*settings)

func defaultSettings() *settings {
	s := new(settings)
	s.maxConcurrecy = 10
	s.batchChunkSize = 500
	return s
}

func MaxConcurrency(max int) Option {
	return func(s *settings) {
		s.maxConcurrecy = max
	}
}

func BatchChunkSize(chunkSize int) Option {
	return func(s *settings) {
		s.batchChunkSize = chunkSize
	}
}

func BatchWrite(ctx context.Context, client *firestore.Client, collection string, v interface{}, opts ...Option) error {
	s := defaultSettings()
	for _, opt := range opts {
		opt(s)
	}

	data := reflect.ValueOf(v)
	if data.Kind() != reflect.Slice {
		return fmt.Errorf("v is expected to be a slice")
	}

	colRef := client.Collection(collection)
	g, ctx := errgroup.WithContext(ctx)

	batches := func(ctx context.Context, client *firestore.Client, v reflect.Value) error {
		throttle := make(chan struct{}, s.maxConcurrecy)
		max := v.Len()
		for i := 0; i < max; i += s.batchChunkSize {
			end := i + s.batchChunkSize
			if end > max {
				end = max
			}
			d := v.Slice(i, end)
			g.Go(func() error {
				defer func() { <-throttle }()
				throttle <- struct{}{}
				batch := client.Batch()
				for i := 0; i < d.Len(); i++ {
					batch.Set(colRef.NewDoc(), d.Index(i).Interface())
				}
				_, err := batch.Commit(ctx)
				return err
			})
		}

		return g.Wait()
	}

	return batches(ctx, client, data)
}
