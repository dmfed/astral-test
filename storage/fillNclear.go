package storage

import (
	"math/rand"
	"time"
)

func GenerateRandomPayload(st Storage, n int) error {
	rand.Seed(time.Now().UnixNano())
	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	for i := 0; i < n; i++ {
		rand.Shuffle(len(letters), func(i, j int) {
			letters[i], letters[j] = letters[j], letters[i]
		})
		var e Element
		e.Payload = string(letters)
		if _, err := st.Put(e); err != nil {
			return err
		}
	}
	return nil
}

func DeleteAllFromStorage(st Storage) error {
	elements, err := st.Get()
	if err != nil {
		return err
	}
	for _, e := range elements {
		if err := st.Del(e.ID); err != nil {
			return err
		}
	}
	return nil
}
