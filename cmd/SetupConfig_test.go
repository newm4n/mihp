package main

import (
	"github.com/newm4n/mihp/internal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEmailChecker(t *testing.T) {
	data := []struct {
		address string
		correct bool
		email   string
		name    string
	}{
		{
			address: "Ferdinand<ferdinand.neman@gmail.com>",
			correct: true,
			email:   "ferdinand.neman@gmail.com",
			name:    "Ferdinand",
		}, {
			address: "Ferdinand.Neman<ferdinand.neman@gmail.com>",
			correct: true,
			email:   "ferdinand.neman@gmail.com",
			name:    "Ferdinand.Neman",
		}, {
			address: "Ferdinand<ferdinand neman@gmail.com>",
			correct: false,
		}, {
			address: "<ferdinand.neman@gmail.com>",
			correct: true,
			email:   "ferdinand.neman@gmail.com",
			name:    "",
		}, {
			address: "Ferdinand Neman<ferdinand.neman@gmail.com>",
			correct: true,
			email:   "ferdinand.neman@gmail.com",
			name:    "Ferdinand Neman",
		}, {
			address: "Ferdinand<ferdinand.neman@thisshouldnotvalid.com>",
			correct: false,
		},
	}

	for _, d := range data {
		em, err := internal.NewMailbox(d.address)
		if d.correct {
			assert.NoError(t, err)
			assert.Equal(t, em.Email, d.email)
			assert.Equal(t, em.Name, d.name)
		} else {
			assert.Error(t, err)
		}
	}

}
