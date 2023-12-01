package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"

	"github.com/MikMik1011/gommand"
	"github.com/sampgo/sampgo"
)

var (
	players = make(map[string]*Roleplayer)
)

func writeJSONToFile(data *map[string]*Roleplayer, filePath string) error {
	// Convert map to JSON
	jsonData, err := json.MarshalIndent(*data, "", "    ")
	if err != nil {
		return err
	}

	// Write JSON data to file
	err = os.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Println("Data written to", filePath)
	return nil
}

func readJSONFromFile(filePath string) (map[string]*Roleplayer, error) {
	// Read JSON data from file
	jsonData, err := os.ReadFile(filePath)
	if err != nil {
		data := make(map[string]*Roleplayer)
		return data, err
	}

	// Unmarshal JSON data into map
	var data map[string]*Roleplayer
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		data = make(map[string]*Roleplayer)
		return data, err
	}
	return data, nil
}

func init() {
	players, err := readJSONFromFile("data/players.json")
	if err != nil {
		fmt.Println(err)
	}

	json, _ := json.Marshal(players)
	sampgo.Print(string(json))

	sampgo.Print("go init() called")
	sampgo.On("goModeInit", func() bool {

		gommand.SetCommandNotFound(ColorWhite, "[TVOJA MAMA] Command not found!")

		gommand.NewCompleteCommand("vehicle", []string{"v", "veh"},
			func(ctx gommand.Context) (err error) {
				if len(ctx.Args) < 1 {
					ctx.Player.SendMessage(ColorRed, "Usage: /vehicle <id>")
					return
				}
				plX, plY, plZ, err := ctx.Player.GetPos()
				id, err := strconv.Atoi(ctx.Args[0])
				if err != nil {
					ctx.Player.SendMessage(ColorRed, "[ERROR] Invalid vehicle ID!")
					return
				}
				veh, err := sampgo.NewVehicle(id, plX, plY, plZ, 0, uint8(rand.Intn(256)), uint8(rand.Intn(256)), 30, true)
				if err != nil {
					ctx.Player.SendMessage(ColorRed, "Unable to spawn car!")
					return
				}
				veh.PutPlayer(&ctx.Player, 0)
				sampgo.SetVehicleNumberPlate(veh.ID, "URMOM")
				msg := fmt.Sprintf("Vehicle id %d created", veh.ID)
				sampgo.Print(msg)
				ctx.Player.SendMessage(ColorWhite, msg)

				return
			})

		gommand.NewCompleteCommand("deletevehicle", []string{"dv", "delveh"},
			func(ctx gommand.Context) (err error) {
				vehID := sampgo.GetPlayerVehicleID(ctx.Player.ID)
				sampgo.DestroyVehicle(vehID)
				msg := fmt.Sprintf("Vehicle id %d destroyed", vehID)
				sampgo.Print(msg)
				ctx.Player.SendMessage(ColorWhite, msg)
				return
			})

		gommand.NewCompleteCommand("weapon", []string{"w", "wep"},
			func(ctx gommand.Context) (err error) {
				if len(ctx.Args) < 1 {
					ctx.Player.SendMessage(ColorRed, "Usage: /weapon <id>")
					return
				}
				weaponID, err := strconv.Atoi(ctx.Args[0])
				if err != nil {
					ctx.Player.SendMessage(ColorRed, "[ERROR] Invalid weapon ID!")
					return
				}
				sampgo.GivePlayerWeapon(ctx.Player.ID, weaponID, 10000)
				msg := fmt.Sprintf("Weapon id %d given", weaponID)
				sampgo.Print(msg)
				ctx.Player.SendMessage(ColorWhite, msg)
				return
			})

		gommand.NewCompleteCommand("fix", []string{"vfix"},
			func(ctx gommand.Context) (err error) {
				sampgo.RepairVehicle(sampgo.GetPlayerVehicleID(ctx.Player.ID))
				return
			})

		gommand.NewCompleteCommand("flip", []string{},
			func(ctx gommand.Context) error {
				sampgo.SetVehicleZAngle(sampgo.GetPlayerVehicleID(ctx.Player.ID), 0)
				return nil
			})

		gommand.NewCompleteCommand("respawn", []string{},
			func(ctx gommand.Context) error {
				ctx.Player.SetPos(-3, 3, 5)
				return nil
			})

		gommand.NewCompleteCommand("jetpack", []string{"jp"},
			func(ctx gommand.Context) error {
				sampgo.SetPlayerSpecialAction(ctx.Player.ID, sampgo.SpecialActionUsejetpack)
				return nil
			})

		sampgo.Print("commands registered!")
		return true
	})

	sampgo.On("goModeExit", func() bool {
		sampgo.Print("goModeExit!")
		writeJSONToFile(&players, "data/players.json")
		return true
	})

	sampgo.On("playerConnect", func(p sampgo.Player) bool {
		sampgo.Print(fmt.Sprintf("Player %s connected!", p.GetName()))
		sampgo.Print(fmt.Sprintf("Player ID is %d", p.ID))
		p.SendMessage(ColorWhite, "izes mi kurac!")
		playerName := p.GetName()

		_, exists := players[playerName]
		if !exists {
			players[playerName] = &Roleplayer{Money: 1000, Exp: rand.Intn(10)}
		}
		players[playerName].SetID(p.ID)
		sampgo.GivePlayerMoney(p.ID, players[playerName].GetMoney())

		p.SendMessage(ColorWhite, "Welcome to the server!")
		p.SendMessage(ColorWhite, fmt.Sprintf("You have $%d!", players[playerName].GetMoney()))
		p.SendMessage(ColorWhite, fmt.Sprintf("Your exp is %d!", players[playerName].GetExp()))
		p.SendMessage(ColorWhite, fmt.Sprintf("Your level is %d!", players[playerName].GetLevel()))
		sampgo.SetPlayerScore(p.ID, players[playerName].GetLevel())
		return true
	})

	sampgo.On("playerSpawn", func(p sampgo.Player) bool {

		p.SetPos(-3, 3, 5)
		return true
	})

	sampgo.On("playerDisconnect", func(p sampgo.Player, reason int) bool {
		sampgo.Print(fmt.Sprintf("Player %s disconnected!", p.GetName()))
		writeJSONToFile(&players, "data/players.json")
		return true
	})

	sampgo.On("playerDeath", func(victim sampgo.Player, killer sampgo.Player, reason int) bool {
		sampgo.Print(fmt.Sprintf("Player %s died!", victim.GetName()))
		players[victim.GetName()].RemoveMoney(100)
		victim.SendMessage(ColorWhite, "you died get rekt noob")
		victim.SetPos(-3, 3, 5)
		sampgo.SetPlayerHealth(victim.ID, 100)
		victim.Spawn()

		killerRP, exists := players[killer.GetName()]
		if !exists {
			return true
		}

		killerRP.AddExp(5)
		killer.SendMessage(ColorWhite, "gj m8")
		return true
	})

}

func main() {}
