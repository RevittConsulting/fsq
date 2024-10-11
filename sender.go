package fsq

type Sender struct {
	sender ISender
}

func NewSender(sender ISender) *Sender {
	return &Sender{
		sender: sender,
	}
}

func (s *Sender) Send(to string, subject string, body string) error {
	return s.sender.SendMail(to, subject, body)
}

type QueueSender struct {
	queue IQueue
}

func NewSenderWithQueue(queue IQueue) *QueueSender {
	return &QueueSender{
		queue: queue,
	}
}

func (s *QueueSender) Send(to string, subject string, body string) error {
	return s.queue.SendToQueue(to, subject, body)
}
