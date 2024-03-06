package stack_list

import "testing"

func TestNewStackList (t *testing.T) {
    stack := NewStackList[rune]()

    if (stack == nil) {
        t.Error("StackList creation failed.")
    }
}

func TestStackListPush (t *testing.T) {
    stack := NewStackList[rune]()

    if (stack == nil) {
        t.Error("StackList creation failed.")
    }

    val := 'A'
    stack.Push(val)

    if (stack.Head == nil) {
        t.Error("Push has failed to link node.")
    }

    if (stack.Head.Data != val) {
        t.Errorf("Head data = %c != %c", stack.Head.Data, val)
    }

    prev := val
    val = 'B'
    stack.Push(val)

    if (stack.Head.Data != val) {
        t.Errorf("Head data = %c != %c", stack.Head.Data, val)
    }
    if (stack.Head.Next.Data != prev) {
        t.Errorf("Head data = %c != %c", stack.Head.Next.Data, prev)
    }
}

func TestStackListPop (t *testing.T) {
    stack := NewStackList[rune]()
    stack_runes := []rune("ABCD")
    reverse_runes := []rune("DCBA")

    for _, r := range stack_runes {
        stack.Push(r)
    }

    for _, r := range reverse_runes {
        val, err := stack.Pop(); if (err != nil) {
            t.Errorf("Failed to pop %c", r)
        }
        if (val != r) {
            t.Errorf("Popped val %c != expected %c", val, r)
        }
    }

    _, err := stack.Pop()

    if (err == nil) {
        t.Error("Popping an empty stack should return an error")
    }
}

func TestStackListTop (t *testing.T) {
    stack := NewStackList[rune]()

    _, err := stack.Top()

    if (err == nil) {
        t.Error("Top checking an empty stack should return an error")
    }

    stack_runes := []rune("ABCD")

    for _, r := range stack_runes {
        stack.Push(r)
    }

    val, err := stack.Top()
    expected := 'D'

    if (val != expected) {
        t.Errorf("Top value %c != expected %c", val, expected)
    }
}

func TestStackListLen (t *testing.T) {
    stack := NewStackList[rune]()
    var want int = 0
    got := stack.Len()

    if (got != want) {
        t.Errorf("Length %d != expected %d", got, want)
    }

    stack_runes := []rune("ABCD")

    for _, r := range stack_runes {
        stack.Push(r)
    }

    _, err := stack.Pop(); if (err != nil) {
        t.Error("Pop failed while testing Len")
    }

    want = 3
    got = stack.Len()

    if (got != want) {
        t.Errorf("Length %d != expected %d", got, want)
    }
}
