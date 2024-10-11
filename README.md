# FSQ

***

F***ing Simple Queue. A simple queue for sending emails.

***

## Description

Create an email service that sends emails. As long as it implements the `ISender` interface, it can be used with the queue.

***

## Usage

##### Create a new email service. Here we will use the smtp.go service.

```go
cfg := &SmtpConfig{
    Host:        "localhost",
    Port:        1025,
    DefaultFrom: "mail@mail.com",
}

smtp := NewSmtp(cfg)
```
##### Create a new sender service that uses the queue.

```go
sender := NewSenderWithQueue(smtp)
```

##### Create a new queue service. Here we will use the rabbit.go service.

```go
cfg := &RabbitConfig{
    RabbitHost: "localhost",
    RabbitPort: 5672,
}

rabbit, err := NewRabbitQueue(cfg, sender)
if err != nil {
    log.Fatalf("error creating rabbit queue: %v", err)
}
```

##### Create a new queue and run it.

```go
var shutdown atomics.AtomicBool
shutdown.Set(false)

queue := New(shuttingDown, rabbit)

if err := queue.Run(ctx); err != nil {
    log.Fatalf("error running queue: %v", err)
}
```