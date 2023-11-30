package main

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/sampgo/command"
	"github.com/sampgo/sampgo"
)

func init() {
	sampgo.Print("go init() called")

	sampgo.On("goModeInit", func() bool {
		sampgo.Print("Hello from Go!")
		cmd := command.NewCommand(command.Command{Name: "foo", Prefix: "!"})
		cmd.Handle(func(ctx command.Context) (err error) {
			ctx.Player.SendMessage(0xFFFFFF, "Hello from Go!")
			return
		})

		vehCMD := command.NewCommand(command.Command{Name: "vehicle", Prefix: ""})
		//vehCMD.SetAlias("v").SetAlias("veh")
		vehCMD.Handle(func(ctx command.Context) (err error) {
			if len(ctx.Args) < 1 {
				ctx.Player.SendMessage(0xFF0000, "Usage: /vehicle <id>")
				return
			}
			plX, plY, plZ, err := ctx.Player.GetPos()
			id, err := strconv.Atoi(ctx.Args[0])
			if err != nil {
				ctx.Player.SendMessage(0xFF0000, "[ERROR] Invalid vehicle ID!")
				return
			}
			veh := sampgo.CreateVehicle(id, plX, plY, plZ, 0, rand.Intn(256), rand.Intn(256), 30, true)
			sampgo.PutPlayerInVehicle(ctx.Player.ID, veh, 0)
			msg := fmt.Sprintf("Vehicle id %d created", veh)
			sampgo.Print(msg)
			ctx.Player.SendMessage(0xFFFFFF, msg)
			return
		})

		delVehCMD := command.NewCommand(command.Command{Name: "deletevehicle", Prefix: ""})
		//delVehCMD.SetAlias("dv").SetAlias("delveh")
		delVehCMD.Handle(func(ctx command.Context) (err error) {
			vehID := sampgo.GetPlayerVehicleID(ctx.Player.ID)
			sampgo.DestroyVehicle(vehID)
			msg := fmt.Sprintf("Vehicle id %d destroyed", vehID)
			sampgo.Print(msg)
			ctx.Player.SendMessage(0xFFFFFF, msg)
			return
		})

		weaponCMD := command.NewCommand(command.Command{Name: "weapon", Prefix: ""})
		//vehCMD.SetAlias("w").SetAlias("wep")
		weaponCMD.Handle(func(ctx command.Context) (err error) {
			if len(ctx.Args) < 1 {
				ctx.Player.SendMessage(0xFF0000, "Usage: /weapon <id>")
				return
			}
			weaponID, err := strconv.Atoi(ctx.Args[0])
			if err != nil {
				ctx.Player.SendMessage(0xFF0000, "[ERROR] Invalid weapon ID!")
				return
			}
			sampgo.GivePlayerWeapon(ctx.Player.ID, weaponID, 9999)
			msg := fmt.Sprintf("Weapon id %d given", weaponID)
			sampgo.Print(msg)
			ctx.Player.SendMessage(0xFFFFFF, msg)
			return
		})

		return true
	})

	sampgo.On("goModeExit", func() bool {
		sampgo.Print("goModeExit!")
		return true
	})

	sampgo.On("playerSpawn", func(p sampgo.Player) bool {
		sampgo.Print(fmt.Sprintf("Player ID is %d", p.ID))
		p.SendMessage(0xFFFFFF, "izes mi kurac!")
		p.SetPos(-3, 3, 5)
		return true
	})

}

func main() {}
