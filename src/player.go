package main

import (
	"errors"

	"github.com/sampgo/sampgo"
)

type Roleplayer struct {
	Money int `json:"money"`
	Exp   int `json:"exp"`
	ID    int `json:"id"`
}

func Test() {
	sampgo.Print("Test")
}

func (r *Roleplayer) GetMoney() int {
	return r.Money
}

func (r *Roleplayer) AddMoney(money int) error {

	r.Money += money
	sampgo.GivePlayerMoney(r.ID, money)

	return nil
}

func (r *Roleplayer) RemoveMoney(money int) error {

	r.AddMoney(-money)

	return nil
}

func (r *Roleplayer) SetMoney(money int) {
	r.AddMoney(money - r.Money)
}

func (r *Roleplayer) GetExp() int {
	return r.Exp
}

func (r *Roleplayer) GetLevel() int {
	return r.Exp/10 + 1
}

func (r *Roleplayer) syncLevelToGame() {
	sampgo.SetPlayerScore(r.ID, r.GetLevel())
}

func (r *Roleplayer) AddExp(exp int) error {

	if exp < 0 {
		return errors.New("Exp cannot be negative")
	}
	r.Exp += exp
	r.syncLevelToGame()

	return nil
}
