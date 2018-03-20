package nap

func scatter(n int, fn func(i int) error) (err error) {
	errors := make(chan error, n)

	for i := 0; i < n; i++ {
		go func(i int) { errors <- fn(i) }(i)
	}

	for i := 0; i < n; i++ {
		if innerErr := <-errors; innerErr != nil {
			err = innerErr
		}
	}

	return err
}
