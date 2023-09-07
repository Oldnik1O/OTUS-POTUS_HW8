// В Go нет концепции "исключений" в том виде, как в других языках
// Если команда вызывает панику - рекомендуется использовать обработку ошибок
 
// Интерфейс команды и очереди:

go
package main

import (
  "fmt"
  "sync"
  "time"
)

type Command interface {
  Execute()
}

type CommandQueue struct {
  commands    []Command
  mutex       sync.Mutex
  cond        *sync.Cond
  stopChannel chan bool
}

func NewCommandQueue() *CommandQueue {
  q := &CommandQueue{}
  q.cond = sync.NewCond(&q.mutex)
  q.stopChannel = make(chan bool, 1)
  return q
}

func (q *CommandQueue) Enqueue(c Command) {
  q.mutex.Lock()
  defer q.mutex.Unlock()
  q.commands = append(q.commands, c)
  q.cond.Signal()
}

func (q *CommandQueue) Dequeue() Command {
  q.mutex.Lock()
  defer q.mutex.Unlock()
  for len(q.commands) == 0 && !q.HasStopSignal() {
    q.cond.Wait()
  }
  if len(q.commands) == 0 {
    return nil
  }
  cmd := q.commands[0]
  q.commands = q.commands[1:]
  return cmd
}

func (q *CommandQueue) HasStopSignal() bool {
  select {
  case <-q.stopChannel:
    return true
  default:
    return false
  }
}

func (q *CommandQueue) SendStopSignal() {
  q.stopChannel <- true
}


// Реализация выполнения команд в отдельном потоке:

go
func ExecuteCommands(queue *CommandQueue) {
  for {
    cmd := queue.Dequeue()
    if cmd == nil {
      return
    }
    cmd.Execute()
  }
}

// Определение команды для старта и остановки выполнения:

go
type StartCommand struct {
  queue *CommandQueue
}

func (s *StartCommand) Execute() {
  go ExecuteCommands(s.queue)
}

type HardStopCommand struct {
  queue *CommandQueue
}

func (h *HardStopCommand) Execute() {
  h.queue.SendStopSignal()
}

// Пример использования:

go
type PrintCommand struct {
  message string
}

func (p *PrintCommand) Execute() {
  fmt.Println(p.message)
}

func main() {
  queue := NewCommandQueue()

  // Start execution in a separate goroutine
  startCmd := &StartCommand{queue: queue}
  startCmd.Execute()

  // Enqueue some commands
  for i := 0; i < 5; i++ {
    queue.Enqueue(&PrintCommand{message: fmt.Sprintf("Command %d", i)})
  }

  // Hard stop after 2 seconds
  time.Sleep(2 * time.Second)
  stopCmd := &HardStopCommand{queue: queue}
  stopCmd.Execute()

  time.Sleep(3 * time.Second) // Just to see the output before the program terminates
}