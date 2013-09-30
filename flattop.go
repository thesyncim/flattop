//todo rename flattop

package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"log"
	"os"
)

var configDir = "/dockerimages/"

func main() {
	app := cli.NewApp()
	app.Name = "dockercm"
	app.Usage = "create, run and manage docker instances"
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{"c", configDir, "config file to create and run container"},
	}
	app.Commands = []cli.Command{
		{
			Name:      "create",
			ShortName: "c",
			Usage:     "create a container configuration",
			Flags: []cli.Flag{
				cli.StringSliceFlag{"dns", &cli.StringSlice{}, "set one or more dns servers"},
				cli.StringSliceFlag{"v", &cli.StringSlice{}, "Bind mount a volume (e.g. from the host: -v /host:/container, from docker: -v /container)"},
				cli.IntSliceFlag{"p", &cli.IntSlice{}, "set one or more ports to expose to host"},
				cli.StringFlag{"host", "", "set container hostname"},
				cli.StringFlag{"image", "", "set container Image"},
				cli.StringFlag{"publicip", "", "set public ip"},
				cli.StringFlag{"privateip", "", "set private static ip"},
				cli.BoolFlag{"privileged", "privileged mode default false"},
				cli.BoolFlag{"d", "detach from terminal"},
				cli.BoolFlag{"r", "replace existing settings"},
				cli.IntFlag{"m", 0, "set memory limit default unlimited"},
			},
			Action: func(c *cli.Context) {

				if len(c.Args()) < 1 {
					log.Fatalln("you must provide a app name")
				}
				containerConf := new(Config)
				containerConf.Hostname = c.String("host")
				containerConf.Dns = c.StringSlice("dns")
				containerConf.Image = c.String("image")
				containerConf.Port = c.IntSlice("p")
				containerConf.Detach = c.Bool("d")
				containerConf.Volumes = c.StringSlice("v")
				containerConf.PrivateIp = c.String("privateip")
				containerConf.PublicIp = c.String("publicip")
				containerConf.Name = c.Args()[0]

				if err := containerConf.ValidateContainer(); err != nil {
					log.Fatalln(err)
				}

				if err := containerConf.ValidateNetwork(); err != nil {
					log.Fatalln(err)
				}

				containerConf.Save(c.GlobalString("c")+c.Args()[0], c.Bool("r"))
				info("Configuration for container " + c.Args()[0] + "saved")

			},
		},
		{
			Name:      "run",
			ShortName: "r",
			Usage:     "run a container based on configfile",
			Flags:     []cli.Flag{},
			Action: func(c *cli.Context) {

				if c.Args()[0] == "" {
					log.Fatal("you must provide a name for app")
				}

			},
		},
	}

	app.Action = func(c *cli.Context) {
		cli.ShowAppHelp(c)
	}
	app.Run(os.Args)
}
