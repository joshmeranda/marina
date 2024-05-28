package storage

import (
	"context"
	"fmt"

	clientv3 "go.etcd.io/etcd/client/v3"
)

type Marshaller func(v any) ([]byte, error)
type Unmarshaller func([]byte, any) error

type etcdStore[V any] struct {
	kv clientv3.KV

	marshaller   Marshaller
	unmarshaller Unmarshaller
}

func NewEtcdStore[V any](kv clientv3.KV, marshaller Marshaller, unmarshaller Unmarshaller) KeyValueStore[string, V] {
	return &etcdStore[V]{
		kv: kv,

		marshaller:   marshaller,
		unmarshaller: unmarshaller,
	}
}

func (s *etcdStore[V]) Get(ctx context.Context, key string) (V, error) {
	var value V

	resp, err := s.kv.Get(ctx, key)
	if err != nil {
		return value, err
	}

	if resp.Count != 1 {
		return value, fmt.Errorf("expected 1 key in response but found %d", resp.Count)
	}

	if n := len(resp.Kvs); n != 1 {
		return value, fmt.Errorf("expected 1 key in response but found %d", len(resp.Kvs))
	}

	data := []byte(resp.Kvs[0].Value)
	if err := s.unmarshaller(data, &value); err != nil {
		return value, fmt.Errorf("error unmarshalling response from etcd: %w", err)
	}

	return value, nil
}

func (s *etcdStore[V]) Set(ctx context.Context, key string, value V) error {
	data, err := s.marshaller(value)
	if err != nil {
		return fmt.Errorf("error marshalling value: %w", err)
	}

	_, err = s.kv.Put(ctx, key, string(data))
	if err != nil {
		return fmt.Errorf("error putting value to etcd: %w", err)
	}

	return nil
}

func (s *etcdStore[V]) Delete(ctx context.Context, key string) error {
	_, err := s.kv.Delete(ctx, key)
	if err != nil {
		return fmt.Errorf("error deleting value from etcd: %w", err)
	}

	return nil
}
