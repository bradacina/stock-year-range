package main

func feeder(list []string) <-chan string {
	return bufferedFeeder(list, 10)
}

func bufferedFeeder(list []string, bufSize uint) <-chan string {
	if bufSize == 0 {
		bufSize = 10
	}
	result := make(chan string, bufSize)

	// feed the values one by one to whoever is listening on the channel
	go func() {
		for _, v := range list {
			result <- v
		}

		close(result)
	}()

	return result
}
