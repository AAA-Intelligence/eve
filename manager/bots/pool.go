package bots

import (
	"log"
	"sync"
)

// A task holds message as input and a channel to send the answer to
// Every request to a bot needs all the information stored in message data.
type botTask struct {
	data   MessageData
	result chan BotAnswer
}

// BotPool manages a pool of running instances of the bot python script.
// In order to generate the bot answers faster for multiple users, any incoming request is forwarded to a free bot instance.
// For every incoming request a task is generated and executed by the next free instance.
// The Number of bot instances can be resized to allow dynamic load balancing.
type BotPool struct {
	mu    sync.Mutex
	size  int
	tasks chan botTask
	kill  chan struct{}
	wg    sync.WaitGroup
}

// NewBotPool creates a new pool with a given amount of bot instances
// The task buffer size is 128. This means the pool can have up to 128 tasks in the que.
func NewBotPool(size int) *BotPool {
	pool := &BotPool{
		tasks: make(chan botTask, 128),
		kill:  make(chan struct{}),
	}
	pool.Resize(size)
	return pool
}

// HandleRequest takes incoming requests and makes the task work off by the next free bot instance.
// The function always returns a answer with text. If an error occures errorBotAnswer(...) is returned.
func (p *BotPool) HandleRequest(data MessageData) BotAnswer {

	task := botTask{
		data:   data,
		result: make(chan BotAnswer),
	}
	p.tasks <- task

	answer := *errorBotAnswer(data.Mood, data.Affection)
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

// Wait locks thread until all bots are done
func (p *BotPool) Wait() {
	p.wg.Wait()
}

// worker starts a new bot instance that is running until it is killed by the pool it belongs to or an error occures.
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

// Resize changes the count of bot instances
// The size can be any positiv number including zero.
// Decreasing the size can take a while, because only bot instances that are not currently working can be destroyed.
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
