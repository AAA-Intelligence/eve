package bots

import (
	"log"
	"sync"
)

type botTask struct {
	data   MessageData
	result chan BotAnswer
}

//BotPool manages a pool of bots
type BotPool struct {
	mu    sync.Mutex
	size  int
	tasks chan botTask
	kill  chan struct{}
	wg    sync.WaitGroup
}

// NewBotPool creates a new pool with a given amount of bot instances
func NewBotPool(size int) *BotPool {
	pool := &BotPool{
		tasks: make(chan botTask, 128),
		kill:  make(chan struct{}),
	}
	pool.Resize(size)
	return pool
}

// HandleRequest handels a message request and redirects it to a free bot instance
func (p *BotPool) HandleRequest(data MessageData) BotAnswer {

	task := botTask{
		data:   data,
		result: make(chan BotAnswer),
	}
	p.tasks <- task

	answer := BotAnswer{
		Text: "Ok",
	}
	select {
	case a, ok := <-task.result:
		if !ok {
			return answer
		}
		answer = a
	}

	return answer
}

// Close kills all running bot instances
func (p *BotPool) Close() {
	close(p.tasks)
}

//Wait locks thread until all bots are done
func (p *BotPool) Wait() {
	p.wg.Wait()
}

func (p *BotPool) worker() {
	defer p.wg.Done()
	bot, err := newBotInstance()
	// close bot if worker is killed
	defer bot.Close()
	if err != nil {
		log.Fatalln("error creating bot instance:", err)
		return
	}
	for {
		select {
		case task, ok := <-p.tasks:
			if !ok {
				return
			}
			task.result <- *bot.sendRequest(task.data)
		case <-p.kill:
			return
		}
	}
}

//Resize changes the count of bot instances
func (p *BotPool) Resize(n int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for p.size < n {
		p.size++
		p.wg.Add(1)
		go p.worker()
	}
	for p.size > n {
		p.size--
		p.kill <- struct{}{}
	}
}
