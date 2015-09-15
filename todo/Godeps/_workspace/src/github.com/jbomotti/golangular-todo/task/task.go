package task

import "fmt"

// TASKS
type Task struct {
  ID int64
  Title string
  Completed bool
}

func NewTask(title string) (*Task, error) {
  if title == "" {
    return nil, fmt.Errorf("Empty title")
  }
  return &Task{0, title, false}, nil
}

// LIST
type List struct {
  tasks []*Task
  lastID int64
}

func NewList() *List {
  return &List{}
}

//TASK/LIST Functions
// Save tasks to a list
func (l *List) Save(task *Task) error {
  if task.ID == 0 {
    l.lastID++
    task.ID = l.lastID
    l.tasks = append(l.tasks, dupeTask(task))
    return nil
  }

  for i, t := range l.tasks {
    if t.ID == task.ID {
      l.tasks[i] = dupeTask(task)
      return nil
    }
  }
  return fmt.Errorf("Unknown task")
}

// Duplicate a task (makes a deep copy), helped with pointers
func dupeTask(t *Task) *Task {
  dupe := *t
  return &dupe
}

// Return all tasks in a List
func (l *List) All() []*Task {
  return l.tasks
}

// Find a Task by ID
func (l *List) Find(ID int64) (*Task, bool) {
  for _, t := range l.tasks {
    if t.ID == ID {
      return t, true
    }
  }
  return nil, false
}
