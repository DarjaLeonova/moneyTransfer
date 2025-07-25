package queue

var JobsChan = make(chan TransferJob, 100)

func Enqueue(job TransferJob) {
	JobsChan <- job
}
