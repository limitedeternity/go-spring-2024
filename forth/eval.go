//go:build !solution

package main

import (
	"container/list"
	"fmt"
	"strconv"
	"strings"
)

type Stack struct {
	internal *list.List
}

func (c *Stack) Push(v any) {
	c.internal.PushFront(v)
}

func (c *Stack) Pop() {
	if c.internal.Len() == 0 {
		panic("stack is empty")
	}

	c.internal.Remove(c.internal.Front())
}

func (c *Stack) Front() *list.Element {
	return c.internal.Front()
}

func (c *Stack) Back() *list.Element {
	return c.internal.Back()
}

func (c *Stack) Len() int {
	return c.internal.Len()
}

const (
	EOC = ";"
)

type CommandFunc func() error

type Command struct {
	ID         int
	Func       CommandFunc
	Def        []int
	PushNum    int
	hasPushNum bool
}

type Evaluator struct {
	stack        Stack
	commandsData map[string]Command
	commandsList map[int]Command
	commandsNum  int
}

// NewEvaluator creates evaluator.
func NewEvaluator() *Evaluator {
	e := &Evaluator{
		stack:        Stack{internal: list.New()},
		commandsData: make(map[string]Command),
		commandsList: make(map[int]Command),
	}

	commands := map[string]CommandFunc{
		"+":    e.add,
		"-":    e.sub,
		"*":    e.mul,
		"/":    e.div,
		"over": e.over,
		"dup":  e.dup,
		"drop": e.drop,
		"swap": e.swap,
	}

	for name, fn := range commands {
		e.addNewCommand(name, Command{ID: e.commandsNum, Func: fn})
	}

	return e
}

// Process evaluates sequence of words or definition.
//
// Returns resulting stack state and an error.
func (e *Evaluator) Process(row string) ([]int, error) {
	var err error

	func() {
		if len(row) == 0 {
			err = fmt.Errorf("empty command")
			return
		}

		row = strings.ToLower(row)
		if row[0] == ':' {
			err = e.reserveNewCommand(row)
			return
		}

		data := strings.Split(row, " ")
		err = e.parseSentence(data)
	}()

	stackRepr := make([]int, e.stack.Len())
	for i, n := 0, e.stack.Back(); n != nil; i, n = i+1, n.Prev() {
		stackRepr[i] = n.Value.(int)
	}

	return stackRepr, err
}

func (e *Evaluator) addNewCommand(commandName string, command Command) {
	e.commandsList[command.ID] = command
	e.commandsData[commandName] = command
	e.commandsNum += 1
}

func (e *Evaluator) reserveNewCommand(row string) error {
	data := strings.SplitN(row[1:], " ", 3)
	if len(data) < 3 {
		return fmt.Errorf("invalid command definition")
	}

	name, cmd := data[1], data[2]
	if _, err := strconv.Atoi(name); err == nil {
		return fmt.Errorf("invalid command name")
	}

	definition, err := e.makeDefinition(strings.Split(cmd, " "))
	if err != nil {
		return err
	}

	delete(e.commandsData, name)
	e.addNewCommand(name, Command{ID: e.commandsNum, Def: definition})

	return nil
}

func (e *Evaluator) makeDefinition(commands []string) ([]int, error) {
	definition := make([]int, 0)

	for _, command := range commands {
		if command == EOC {
			break
		}

		if c, ok := e.commandsData[command]; ok {
			if c.Func != nil {
				definition = append(definition, c.ID)
			} else {
				definition = append(definition, c.Def...)
			}
		} else if number, err := strconv.Atoi(command); err == nil {
			definition = append(definition, e.commandsNum)
			e.addNewCommand(command, Command{ID: e.commandsNum, PushNum: number, hasPushNum: true})
		} else {
			return nil, fmt.Errorf("invalid command")
		}
	}

	return definition, nil
}

func (e *Evaluator) parseSentence(data []string) error {
	for _, word := range data {
		if number, err := strconv.Atoi(word); err == nil {
			e.stack.Push(number)
			continue
		}

		c, ok := e.commandsData[word]
		if !ok {
			return fmt.Errorf("invalid command")
		}

		ok, err := e.evaluate(c)
		if err != nil {
			return err
		}

		if ok {
			continue
		}

		for _, command := range c.Def {
			commandData, ok := e.commandsList[command]
			if !ok {
				return fmt.Errorf("invalid command")
			}

			_, err := e.evaluate(commandData)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (e *Evaluator) evaluate(c Command) (bool, error) {
	if c.Func != nil {
		if err := c.Func(); err != nil {
			return false, err
		}

		return true, nil
	}

	if c.hasPushNum {
		e.stack.Push(c.PushNum)
		return true, nil
	}

	return false, nil
}

func (e *Evaluator) add() error {
	first := e.stack.Front()
	if first == nil {
		return fmt.Errorf("insufficient stack size")
	}

	e.stack.Pop()

	second := e.stack.Front()
	if second == nil {
		return fmt.Errorf("insufficient stack size")
	}

	e.stack.Pop()
	e.stack.Push(first.Value.(int) + second.Value.(int))
	return nil
}

func (e *Evaluator) sub() error {
	first := e.stack.Front()
	if first == nil {
		return fmt.Errorf("insufficient stack size")
	}

	e.stack.Pop()

	second := e.stack.Front()
	if second == nil {
		return fmt.Errorf("insufficient stack size")
	}

	e.stack.Pop()
	e.stack.Push(second.Value.(int) - first.Value.(int))
	return nil
}

func (e *Evaluator) mul() error {
	first := e.stack.Front()
	if first == nil {
		return fmt.Errorf("insufficient stack size")
	}

	e.stack.Pop()

	second := e.stack.Front()
	if second == nil {
		return fmt.Errorf("insufficient stack size")
	}

	e.stack.Pop()
	e.stack.Push(first.Value.(int) * second.Value.(int))
	return nil
}

func (e *Evaluator) div() error {
	first := e.stack.Front()
	if first == nil {
		return fmt.Errorf("insufficient stack size")
	}

	e.stack.Pop()

	second := e.stack.Front()
	if second == nil {
		return fmt.Errorf("insufficient stack size")
	}

	e.stack.Pop()

	if denominator := first.Value.(int); denominator == 0 {
		return fmt.Errorf("division by zero")
	}

	e.stack.Push(second.Value.(int) / first.Value.(int))
	return nil
}

func (e *Evaluator) dup() error {
	first := e.stack.Front()
	if first == nil {
		return fmt.Errorf("insufficient stack size")
	}

	e.stack.Push(first.Value.(int))
	return nil
}

func (e *Evaluator) over() error {
	first := e.stack.Front()
	if first == nil {
		return fmt.Errorf("insufficient stack size")
	}

	second := first.Next()
	if second == nil {
		return fmt.Errorf("insufficient stack size")
	}

	e.stack.Push(second.Value.(int))
	return nil
}

func (e *Evaluator) swap() error {
	first := e.stack.Front()
	if first == nil {
		return fmt.Errorf("insufficient stack size")
	}

	e.stack.Pop()

	second := e.stack.Front()
	if second == nil {
		return fmt.Errorf("insufficient stack size")
	}

	e.stack.Pop()
	e.stack.Push(first.Value.(int))
	e.stack.Push(second.Value.(int))
	return nil
}

func (e *Evaluator) drop() error {
	first := e.stack.Front()
	if first == nil {
		return fmt.Errorf("insufficient stack size")
	}

	e.stack.Pop()
	return nil
}
