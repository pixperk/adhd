package adhd

import (
	"testing"
	"time"
)

func TestBackground(t *testing.T) {
	ctx := Background()

	if ctx.Done() != nil {
		t.Error("Background context should have nil Done channel")
	}

	if ctx.Err() != nil {
		t.Error("Background context should have nil error")
	}

	if ctx.Value("key") != nil {
		t.Error("Background context should return nil for any value")
	}

	_, ok := ctx.Deadline()
	if ok {
		t.Error("Background context should not have deadline")
	}
}

func TestTODO(t *testing.T) {
	ctx := TODO()

	if ctx.Done() != nil {
		t.Error("TODO context should have nil Done channel")
	}

	if ctx.Err() != nil {
		t.Error("TODO context should have nil error")
	}

	if ctx.Value("key") != nil {
		t.Error("TODO context should return nil for any value")
	}

	_, ok := ctx.Deadline()
	if ok {
		t.Error("TODO context should not have deadline")
	}
}

func TestWithCancel(t *testing.T) {
	parent := Background()
	ctx, cancel := WithCancel(parent)

	select {
	case <-ctx.Done():
		t.Error("Context should not be done initially")
	default:
	}

	if ctx.Err() != nil {
		t.Error("Context should not have error initially")
	}

	cancel()

	select {
	case <-ctx.Done():
		t.Log("Context done after cancel")
	case <-time.After(100 * time.Millisecond):
		t.Error("Context should be done after cancel")
	}

	if ctx.Err() != ErrCanceled {
		t.Errorf("Expected ErrCanceled, got %v", ctx.Err())
	}
}

func TestWithTimeout(t *testing.T) {
	parent := Background()
	ctx, cancel := WithTimeout(parent, 50*time.Millisecond)
	defer cancel()

	select {
	case <-ctx.Done():
		t.Error("Context should not be done initially")
	default:
	}

	select {
	case <-ctx.Done():
		t.Log("Context timed out as expected")
	case <-time.After(200 * time.Millisecond):
		t.Error("Context should have timed out")
	}

	if ctx.Err() != ErrDeadlineExceeded {
		t.Errorf("Expected ErrDeadlineExceeded, got %v", ctx.Err())
	}
}

func TestWithDeadline(t *testing.T) {
	parent := Background()
	deadline := time.Now().Add(50 * time.Millisecond)
	ctx, cancel := WithDeadline(parent, deadline)
	defer cancel()

	actualDeadline, ok := ctx.Deadline()
	if !ok {
		t.Error("Deadline context should have deadline")
	}

	if !actualDeadline.Equal(deadline) {
		t.Errorf("Expected deadline %v, got %v", deadline, actualDeadline)
	}

	select {
	case <-ctx.Done():
		t.Log("Context met deadline as expected")
	case <-time.After(200 * time.Millisecond):
		t.Error("Context should have met deadline")
	}

	if ctx.Err() != ErrDeadlineExceeded {
		t.Errorf("Expected ErrDeadlineExceeded, got %v", ctx.Err())
	}
}

func TestWithDeadlinePast(t *testing.T) {
	parent := Background()
	pastDeadline := time.Now().Add(-1 * time.Hour)
	ctx, cancel := WithDeadline(parent, pastDeadline)
	defer cancel()

	select {
	case <-ctx.Done():
		t.Log("Past deadline context immediately done")
	default:
		t.Error("Past deadline context should be immediately done")
	}

	if ctx.Err() != ErrDeadlineExceeded {
		t.Errorf("Expected ErrDeadlineExceeded, got %v", ctx.Err())
	}
}

func TestWithValue(t *testing.T) {
	parent := Background()
	key := "testKey"
	value := "testValue"

	ctx := WithValue(parent, key, value)

	if ctx.Value(key) != value {
		t.Errorf("Expected %v, got %v", value, ctx.Value(key))
	}

	if ctx.Value("nonexistent") != nil {
		t.Error("Should return nil for nonexistent key")
	}

	if ctx.Done() != parent.Done() {
		t.Error("Value context should delegate Done to parent")
	}

	if ctx.Err() != parent.Err() {
		t.Error("Value context should delegate Err to parent")
	}
}

func TestValueChaining(t *testing.T) {
	parent := Background()
	ctx1 := WithValue(parent, "key1", "value1")
	ctx2 := WithValue(ctx1, "key2", "value2")

	if ctx2.Value("key1") != "value1" {
		t.Error("Should find key1 in parent context")
	}

	if ctx2.Value("key2") != "value2" {
		t.Error("Should find key2 in current context")
	}

	if ctx2.Value("key3") != nil {
		t.Error("Should return nil for nonexistent key")
	}
}

func TestSelect(t *testing.T) {
	ctx1, cancel1 := WithTimeout(Background(), 100*time.Millisecond)
	ctx2, cancel2 := WithTimeout(Background(), 200*time.Millisecond)
	defer cancel1()
	defer cancel2()

	result := <-Select(ctx1, ctx2)

	if result.Index != 0 {
		t.Errorf("Expected index 0, got %d", result.Index)
	}

	if result.Error != ErrDeadlineExceeded {
		t.Errorf("Expected ErrDeadlineExceeded, got %v", result.Error)
	}

	t.Logf("Select returned context %d with error %v", result.Index, result.Error)
}

func TestSelectEmpty(t *testing.T) {
	result := <-Select()

	if result.Context != nil || result.Index != 0 || result.Error != nil {
		t.Error("Empty select should return zero result")
	}
}

func TestRace(t *testing.T) {
	ctx1, cancel1 := WithCancel(Background())
	ctx2, cancel2 := WithTimeout(Background(), 50*time.Millisecond)
	defer cancel1()
	defer cancel2()

	cancel1()

	result := <-Race(ctx1, ctx2)

	if result.Error != ErrCanceled {
		t.Errorf("Expected ErrCanceled, got %v", result.Error)
	}

	t.Logf("Race won by context with error %v", result.Error)
}

func TestRaceEmpty(t *testing.T) {
	result := <-Race()

	if result.Context != nil || result.Error != nil {
		t.Error("Empty race should return zero result")
	}
}

func TestIsDone(t *testing.T) {
	ctx := Background()
	if IsDone(ctx) {
		t.Error("Background context should not be done")
	}

	cancelCtx, cancel := WithCancel(Background())
	if IsDone(cancelCtx) {
		t.Error("Cancel context should not be done initially")
	}

	cancel()
	if !IsDone(cancelCtx) {
		t.Error("Cancel context should be done after cancel")
	}
}

func TestWaitFor(t *testing.T) {
	ctx, cancel := WithTimeout(Background(), 50*time.Millisecond)
	defer cancel()

	err := WaitFor(ctx)
	if err != ErrDeadlineExceeded {
		t.Errorf("Expected ErrDeadlineExceeded, got %v", err)
	}

	t.Logf("WaitFor returned error %v", err)
}

func TestCancelMultipleTimes(t *testing.T) {
	ctx, cancel := WithCancel(Background())

	cancel()
	cancel()
	cancel()

	if ctx.Err() != ErrCanceled {
		t.Errorf("Expected ErrCanceled, got %v", ctx.Err())
	}

	t.Log("Multiple cancels handled correctly")
}

func TestTimeoutCancelRace(t *testing.T) {
	ctx, cancel := WithTimeout(Background(), 100*time.Millisecond)

	time.Sleep(10 * time.Millisecond)
	cancel()

	<-ctx.Done()

	if ctx.Err() != ErrCanceled {
		t.Errorf("Expected ErrCanceled when canceled before timeout, got %v", ctx.Err())
	}

	t.Log("Cancel won race against timeout")
}
