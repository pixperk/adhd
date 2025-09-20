package adhd

import "sync"

type SelectResult struct {
	Context ADHD
	Index   int
	Error   error
}

type RaceResult struct {
	Context ADHD
	Error   error
}

func Select(contexts ...ADHD) <-chan SelectResult {
	result := make(chan SelectResult, 1)

	if len(contexts) == 0 {
		close(result)
		return result
	}

	var once sync.Once

	for i, ctx := range contexts {
		go func(index int, context ADHD) {
			<-context.Done()
			once.Do(func() {
				result <- SelectResult{
					Context: context,
					Index:   index,
					Error:   context.Err(),
				}
				close(result)
			})
		}(i, ctx)
	}

	return result
}

func Race(contexts ...ADHD) <-chan RaceResult {
	result := make(chan RaceResult, 1)

	if len(contexts) == 0 {
		close(result)
		return result
	}

	var once sync.Once

	for _, ctx := range contexts {
		go func(context ADHD) {
			<-context.Done()
			once.Do(func() {
				result <- RaceResult{
					Context: context,
					Error:   context.Err(),
				}
				close(result)
			})
		}(ctx)
	}

	return result
}

func IsDone(ctx ADHD) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}

func WaitFor(ctx ADHD) error {
	<-ctx.Done()
	return ctx.Err()
}
