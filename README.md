# ppgen

## `ppgen nop -type Thing -src my_thing.go`

Generate an implementation of `Thing` that does nothing.

[Example](./example):
```go
//go:generate ppgen nop -type UserDB -src db.go

type UserDB interface {
	Create(name string) (*User, error)
	Get(id string) (u *User, ok bool, err error)
	Delete(*User) error
}
```

## `ppgen interceptor -type Thing -src my_thing.go`

Generates a type that allows intercepting calls to an underlying interface. The callbacks will tell you what function
is being invoked and give you a chance to override the error being returned.

The interceptor will be used like, for instance to add timeouts and an error counter to all function calls:

```go
aroundHook := func(ctx context.Context, funcName string, callback func(context.Context) error) error {
    ctx, cancel := context.WithTimeout(ctx, time.Second)
    defer cancel()
    err := callback(ctx)
    if err != nil {
        errorCounter.Inc(funcName)
    }
    return err
}
thing := InterceptThing(thing, aroundHook)
```

```go
type thingInterceptor struct {
    around func(context.Context, string, func(context.Context) error) error
    wrap Thing
}

func InterceptThing(thing Thing, around func(context.Context, string, func(context.Context) error) error) Thing {
    return &thingInterceptor{
        around: callback,
        wrap: thing,
    }
}

func (iceptor *thingInterceptor) FunctionName(ctx context.Context, valA *pb.TypeA, valB *pb.TypeB) (valC *pb.TypeC, err error) {
    actualErr := iceptor.around(
        ctx,
        "FunctionName",
        func(childCtx context.Context) error {
            valC, err = iceptor.wrap.FunctionName(childCtx, valA, valB)
            return err
        },
    )
    err = actualErr
    return valC, err
}
```

## ppgen grpc -type Thing -src my_thing.go

Generates a gRPC server and client bridge for an interface. The interface must respect the following
convention for methods on the interface:

```go
// my_thing.go
type Thing interface {
    FunctionName(context.Context, *pb.TypeA, *pb.TypeB) (*pb.TypeC, error)
}
```

* The first argument _MUST_ be a `context.Context` from the standard library.
* The other arguments _MUST_ protobuf types.
* The return values _MUST_ be protobuf types.
* The last return value _MUST_ be an error.
