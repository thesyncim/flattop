//todo rename flattop

package main

import (
	"github.com/codegangsta/cli"
	"log"
	"os"
	"strings"
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
				cli.StringSliceFlag{"dns", &cli.StringSlice{}, "Set custom dns servers"},
				cli.StringSliceFlag{"v", &cli.StringSlice{}, "Bind mount a volume (e.g. from the host: -v /host:/container, from docker: -v /container)"},
				cli.StringSliceFlag{"e", &cli.StringSlice{}, "Set environment variables"},
				cli.IntSliceFlag{"p", &cli.IntSlice{}, "Expose a container's port to the host (use 'docker port' to see the actual mapping)"},
				cli.StringFlag{"host", "", "Container host name"},
				cli.StringFlag{"u", "", " Username or UID"},
				cli.StringFlag{"w", "", "Working directory inside the container"},
				cli.StringFlag{"c", "0", "CPU shares (relative weight)"},
				cli.StringFlag{"image", "", "Image to run"},
				cli.StringFlag{"publicip", "", "Container public static ip"},
				cli.StringFlag{"privateip", "", "Container private static ip"},
				cli.BoolFlag{"privileged", " Give extended privileges to this container"},
				cli.BoolFlag{"r", "Replace Config File"},
				cli.IntFlag{"m", 0, "Memory limit (in bytes)"},
			},
			Action: func(c *cli.Context) {

				if len(c.Args()) < 2 {
					exit("you must provide a app name and command to run")
				}
				appname := c.Args()[0]
				containerConf := new(Config)
				containerConf.Hostname = c.String("host")
				containerConf.Dns = c.StringSlice("dns")
				containerConf.Image = c.String("image")
				containerConf.Cpu = c.String("c")
				containerConf.Environment = c.StringSlice("e")
				containerConf.User = c.String("u")
				containerConf.WorkingDirectory = c.String("w")
				containerConf.Port = c.IntSlice("p")
				containerConf.Volumes = c.StringSlice("v")
				containerConf.PrivateIp = c.String("privateip")
				containerConf.PublicIp = c.String("publicip")
				containerConf.Name = appname
				containerConf.Command = strings.Join(c.Args()[1:], " ")

				if err := containerConf.ValidateContainer(); err != nil {
					exit(err)
				}

				if err := containerConf.ValidateNetwork(); err != nil {
					exit(err)
				}

				containerConf.Save(c.GlobalString("c")+appname, c.Bool("r"))

			},
		},
		{
			Name:      "run",
			ShortName: "r",
			Usage:     "run a container based on configfile",
			Flags:     []cli.Flag{},
			Action: func(c *cli.Context) {

				if len(c.Args()) < 1 {
					log.Fatal("you must provide a name for app")
				}

				appname := c.Args()[0]

				containerConf := new(Config)

				containerConf.LoadConfigFile(configDir + appname)
				if err := containerConf.ValidateContainer(); err != nil {
					exit(err)
				}

				if err := containerConf.ValidateNetwork(); err != nil {
					exit(err)
				}

				containerId := containerConf.Run()

				containerConf.SetPrivateIp(containerId)
				containerConf.setPublicIp(containerId)

			},
		},
	}

	app.Action = func(c *cli.Context) {
		cli.ShowAppHelp(c)
	}
	app.Run(os.Args)
}
