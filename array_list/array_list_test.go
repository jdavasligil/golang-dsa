package arraylist

import "testing"

func TestNewArrayList(t *testing.T) {
	stack := NewArrayList[rune]()

	if stack == nil {
		t.Error("ArrayList creation failed.")
	}
}

func TestArrayListPush(t *testing.T) {
	stack := NewArrayList[rune]()

	if stack == nil {
		t.Error("ArrayList creation failed.")
	}

	val := 'A'
	stack.Push(val)

	if stack.Len() == 0 {
		t.Error("Push has failed to link node.")
	}

	if stack.Data[0] != val {
		t.Errorf("Data = %c != %c", stack.Data[0], val)
	}

	prev := val
	val = 'B'
	stack.Push(val)

	if stack.Data[stack.Len()-1] != val {
		t.Errorf("Head data = %c != %c", stack.Data[stack.Len()-1], val)
	}
	if stack.Data[stack.Len()-2] != prev {
		t.Errorf("Head data = %c != %c", stack.Data[stack.Len()-2], prev)
	}
}

func TestArrayListPop(t *testing.T) {
	stack := NewArrayList[rune]()
	stack_runes := []rune("ABCD")
	reverse_runes := []rune("DCBA")

	for _, r := range stack_runes {
		stack.Push(r)
	}

	for _, r := range reverse_runes {
		val, err := stack.Pop()
		if err != nil {
			t.Errorf("Failed to pop %c", r)
		}
		if val != r {
			t.Errorf("Popped val %c != expected %c", val, r)
		}
	}

	_, err := stack.Pop()

	if err == nil {
		t.Error("Popping an empty stack should return an error")
	}
}

func TestArrayListTop(t *testing.T) {
	stack := NewArrayList[rune]()

	_, err := stack.Top()

	if err == nil {
		t.Error("Top checking an empty stack should return an error")
	}

	stack_runes := []rune("ABCD")

	for _, r := range stack_runes {
		stack.Push(r)
	}

	val, err := stack.Top()
	expected := 'D'

	if val != expected {
		t.Errorf("Top value %c != expected %c", val, expected)
	}
}

func TestArrayListLen(t *testing.T) {
	stack := NewArrayList[rune]()
	var want int = 0
	got := stack.Len()

	if got != want {
		t.Errorf("Length %d != expected %d", got, want)
	}

	stack_runes := []rune("ABCD")

	for _, r := range stack_runes {
		stack.Push(r)
	}

	_, err := stack.Pop()
	if err != nil {
		t.Error("Pop failed while testing Len")
	}

	want = 3
	got = stack.Len()

	if got != want {
		t.Errorf("Length %d != expected %d", got, want)
	}
}
