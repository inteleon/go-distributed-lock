# go-distributed-lock
Extremely simple redis-backed distributed lock / counter

##### Release history
2019-06-13: Initial v0.0.1 release with a counter backed by redis or in-mem.

## Building and testing

Use the supplied Makefile to build and test this library.

    > make all
    go fmt ./...
    go build ./...
    go test ./...
    ok  	github.com/inteleon/go-distributed-lock/counter	(cached)
    ?   	github.com/inteleon/go-distributed-lock/mock_redis	[no test files]

## Using this library

In your application where you want to use this library. Edit your _go.mod_, add:

    require (
        ... other dependencies ...
        github.com/inteleon/go-distributed-lock v0.0.1
    )

## Usage in code

### Counter
The counter comes in two flavors, in-memory (for testing and local development) and redis (for usage in integration tests, staging, prod etc where there's a Redis available).

Both uses an int64 as store for the counter, with the semantic that a counter > 0 is "locked", while a nil or counter <= 0 is "not locked".

The primary use case is when we have some asynchronous processing such as consuming a finite number of messages from a queue, where the enqueuer knows how many messages that was enqueued but won't know when all the enqueued message have been processed. 

The number of messages is _Set(num)_ on a counter. Each consumer then can decrease this counter using Decr() once an item has been processed. When the counter reaches 0, we know that the entire "batch" has been completed.

##### Initialization of Redis counter

    // Create an instance of a RediCounter with 30 minute expiry.
    counter, err := counter.NewRedisCounter("my-counter", "redis:6379", "some-password", 1800)
   
    // Check for errors. Stop the application if there is a problem.
    if err != nil {
        logrus.Fatalf(err.Error())
    }
    
    // Sets the created instance as a application specific singleton
    lock.SetSingleton(counter)
    
##### Sample usage (in-memory counter)

    counter := NewInMemoryCounter("my-key", time.Now().Add(3*time.Second))
    counter.Set(1337)
    fmt.Printf("Is locked: %v\n", counter.IsLocked()) // true
    counter.Decr()
    fmt.Printf("Value: %v\n", counter.Get()) // 136
    for ;counter.Get() > 0; {
        counter.Decr()
    }
    fmt.Printf("Is locked: %v\n", counter.IsLocked()) // false
    
Should output

    Is locked: true
    Value: 136
    Is locked: false
    
    
## Mocking this library
The project provides a mock for its interface using [gomock](https://github.com/gomock).

Gomock and its mockgen tool can be installed using the following commands:

    go get github.com/golang/mock/gomock
    go install github.com/golang/mock/mockgen

To re-generate the mock, run:

    make mock