package pubsub

import (
    "context"
    "errors"
    "sync"
)

// MsgHandler is a callback function that processes messages delivered to subscribers.
type MsgHandler func(msg interface{})

// Subscription represents a subscription to a topic.
type Subscription interface {
    // Unsubscribe cancels the subscription.
    Unsubscribe()
}

// PubSub represents a publish-subscribe system.
type PubSub interface {
    // Subscribe creates a new subscription to the given topic.
    Subscribe(topic string, cb MsgHandler) (Subscription, error)

    // Publish sends a message to all subscribers of the given topic.
    Publish(topic string, msg interface{}) error

    // Close shuts down the pub-sub system.
    // It respects the provided context, so if it's canceled, Close returns immediately.
    Close(ctx context.Context) error
}

type pubSub struct {
    mu            sync.Mutex
    subscriptions map[string][]*subscription
    closed        bool // Добавляем поле для отслеживания закрытия PubSub
    wg            sync.WaitGroup // Добавляем WaitGroup для ожидания завершения горутин
}

type subscription struct {
    ctx        context.Context
    cancelFunc context.CancelFunc
    messages   chan interface{}
}

func (s *subscription) Unsubscribe() {
    s.cancelFunc()
}

func (ps *pubSub) Subscribe(topic string, cb MsgHandler) (Subscription, error) {
    ps.mu.Lock()
    defer ps.mu.Unlock()

    if ps.closed { // Проверяем, закрыт ли PubSub
        return nil, errors.New("pubsub is closed")
    }

    ctx, cancel := context.WithCancel(context.Background())
    sub := &subscription{
        ctx:        ctx,
        cancelFunc: cancel,
        messages:   make(chan interface{}),
    }

    if _, ok := ps.subscriptions[topic]; !ok {
        ps.subscriptions[topic] = []*subscription{sub}
    } else {
        ps.subscriptions[topic] = append(ps.subscriptions[topic], sub)
    }

    // Добавляем WaitGroup для синхронизации завершения горутин
    ps.wg.Add(1)
    go func() {
        defer ps.wg.Done() // Уменьшаем счетчик при завершении горутины
        defer close(sub.messages) // Закрываем канал сообщений при завершении горутины
        for {
            select {
            case <-ctx.Done():
                return
            case msg, ok := <-sub.messages:
                if !ok {
                    return
                }
                cb(msg)
            }
        }
    }()

    return sub, nil
}


func (ps *pubSub) Publish(topic string, msg interface{}) error {
    ps.mu.Lock()
    defer ps.mu.Unlock()

    if ps.closed { // Проверяем, закрыт ли PubSub
        return errors.New("pubsub is closed")
    }

    if subs, ok := ps.subscriptions[topic]; ok {
        for _, sub := range subs {
            select {
            case <-sub.ctx.Done():
                continue
            case sub.messages <- msg:
            }
        }
    }

    return nil
}

func (ps *pubSub) Close(ctx context.Context) error {
    ps.mu.Lock()
    defer ps.mu.Unlock()

    if ps.closed {
        return errors.New("pubsub is already closed")
    }

    ps.closed = true // Устанавливаем флаг закрытия

    for _, subs := range ps.subscriptions {
        for _, sub := range subs {
            go sub.cancelFunc() // Отменяем все подписки
        }
    }

    // Ожидаем завершения всех горутин
    ps.wg.Wait()

    return nil
}


// NewPubSub creates a new instance of PubSub.
func NewPubSub() PubSub {
    return &pubSub{
        subscriptions: make(map[string][]*subscription),
    }
}
