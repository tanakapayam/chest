package main

import "sync"

func processEjsonDocs(public string) {
	var wg sync.WaitGroup

	for _, doc := range docs {
		wg.Add(1)
		go updateEjsonDoc(
			doc,
			public,
			&wg,
		)
	}
	wg.Wait()
}
