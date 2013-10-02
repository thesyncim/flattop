//todo rename flattop

package main

import (
	"github.com/codegangsta/cli"
	"os"
)

var configDir = "/dockerimages/"

func main() {
	app := cli.NewApp()
	app.Name = "dockercm"
	app.Usage = "create, run and manage docker instances"
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.StringFlag{"c", configDir, "config dir to create and run container"},
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
				cli.StringSliceFlag{"-lxc-conf", &cli.StringSlice{}, "Add custom lxc options -lxc-conf=\"lxc.cgroup.cpuset.cpus = 0,1\""},
				cli.StringSliceFlag{"-volumes-from", &cli.StringSlice{}, "Mount volumes from the specified container"},
				cli.IntSliceFlag{"p", &cli.IntSlice{}, "Expose a container's port to the host (use 'docker port' to see the actual mapping)"},
				cli.StringFlag{"host", "", "Container host name"},
				cli.StringFlag{"u", "", "Username or UID"},
				cli.StringFlag{"w", "", "Working directory inside the container"},
				cli.StringFlag{"c", "0", "CPU shares (relative weight)"},
				cli.StringFlag{"image", "", "Image to run"},
				cli.StringFlag{"entrypoint", "", "Overwrite the default entrypoint of the image"},
				cli.StringFlag{"publicip", "", "Container public static ip"},
				cli.StringFlag{"privateip", "", "Container private static ip"},
				cli.BoolFlag{"privileged", "Give extended privileges to this container"},
				cli.BoolFlag{"r", "Replace Config File"},
				cli.IntFlag{"m", 0, "Memory limit (in bytes)"},
			},
			Action: func(c *cli.Context) {

				var container *Container

				container, err := ParseCreateContext(c)
				if err != nil {
					exit(err)
				}

				err = container.Create()
				if err != nil {
					exit(err)
				}

			},
		},
		{
			Name:      "start",
			ShortName: "s",
			Usage:     "run a container based on configfile",
			Flags:     []cli.Flag{},
			Action: func(c *cli.Context) {

				var container *Container

				container, err := ParseStartContext(c)
				if err != nil {
					exit(err)
				}
				err = container.Start()
				if err != nil {
					exit(err)
				}

			},
		},
	},
	{
			Name:      "reload",
			ShortName: "rl",
			Usage:     "run a container based on configfile",
			Flags:     []cli.Flag{},
			Action: func(c *cli.Context) {

				var container *Container

				container, err := ParseStartContext(c)
				if err != nil {
					exit(err)
				}
				err = container.Start()
				if err != nil {
					exit(err)
				}

			},
		},
	}

	app.Action = func(c *cli.Context) {
		cli.ShowAppHelp(c)
	}
	app.Run(os.Args)
}
