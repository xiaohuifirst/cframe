package main

import (
	"os"
	"strings"

	"github.com/ICKelin/cframe/pkg/etcdstorage"
	cli "github.com/urfave/cli/v2"
)

func main() {
	endpoints := []string{"127.0.0.1:2379"}
	envendpoints := os.Getenv("ETCD_ENDPOINTS")
	if len(envendpoints) > 0 {
		endpoints = strings.Split(envendpoints, ",")
	}

	userName := os.Getenv("ETCD_USER")
	passWord := os.Getenv("ETCD_PASSWORD")

	store := etcdstorage.NewEtcd(endpoints, userName, passWord)

	app := cli.NewApp()
	app.Usage = "cfctl manage namespace/edge of cframe"
	app.Commands = []*cli.Command{
		{
			Name:    "namespace",
			Aliases: []string{"ns"},
			Usage:   "manage namespace",
			Subcommands: []*cli.Command{
				{
					Name:  "add",
					Usage: "add a new namespace",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "name",
							Required: true,
						},
					},
					Action: func(ctx *cli.Context) error {
						name := ctx.String("name")
						addNamespace(name, store)
						return nil
					},
				},
				{
					Name:  "del",
					Usage: "delete a namespace",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:     "name",
							Required: true,
						},
					},
					Action: func(ctx *cli.Context) error {
						delNamespace(ctx.String("name"), store)
						return nil
					},
				},
				{
					Name:  "list",
					Usage: "list all namespaces",
					Action: func(ctx *cli.Context) error {
						listNamespace(store)
						return nil
					},
				},
			},
		},
		{
			Name:  "edge",
			Usage: "manage edge",
			Subcommands: []*cli.Command{
				{
					Name:  "add",
					Usage: "add a new edge",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:    "namespace",
							Aliases: []string{"ns"},
							Value:   "default",
						},
						&cli.StringFlag{
							Name:     "name",
							Required: true,
						},
						&cli.StringFlag{
							Name:     "listener",
							Required: true,
							Usage:    "edge listener, eg: 1.2.3.4:58423",
						},
						&cli.StringFlag{
							Name:     "cidr",
							Required: true,
							Usage:    "eg: 172.18.0.0/16",
						},
					},
					Action: func(ctx *cli.Context) error {
						ns := ctx.String("ns")
						edgeName := ctx.String("name")
						listen := ctx.String("listener")
						cidr := ctx.String("cidr")

						addEdge(ns, edgeName, listen, cidr, store)
						return nil
					},
				},
				{
					Name:  "del",
					Usage: "delete a edge",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:    "namespace",
							Aliases: []string{"ns"},
							Value:   "default",
						},
						&cli.StringFlag{
							Name: "name",
						},
					},
					Action: func(ctx *cli.Context) error {
						ns := ctx.String("ns")
						edgeName := ctx.String("name")
						delEdge(ns, edgeName, store)
						return nil
					},
				},
				{
					Name:  "list",
					Usage: "list all edges",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:    "namespace",
							Aliases: []string{"ns"},
							Value:   "default",
						},
					},
					Action: func(ctx *cli.Context) error {
						listEdges(ctx.String("ns"), store)
						return nil
					},
				},
			},
		},
		{
			Name:    "route",
			Aliases: []string{"ro"},
			Usage:   "route manager",
			Subcommands: []*cli.Command{
				{
					Name:  "add",
					Usage: "add a new route",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:    "namespace",
							Aliases: []string{"ns"},
							Usage:   "namespace",
							Value:   "default",
						},
						&cli.StringFlag{
							Name:     "listener",
							Usage:    "edge listener",
							Required: true,
						},
						&cli.StringFlag{
							Name:     "name",
							Usage:    "route name",
							Required: true,
						},
						&cli.StringFlag{
							Name:     "cidr",
							Usage:    "dst cidr block",
							Required: true,
						},
					},
					Action: func(ctx *cli.Context) error {
						ns := ctx.String("namespace")
						name := ctx.String("name")
						listener := ctx.String("listener")
						cidr := ctx.String("cidr")
						addRoute(ns, name, listener, cidr, store)
						return nil
					},
				},
				{
					Name:  "del",
					Usage: "del a route",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:    "namespace",
							Aliases: []string{"ns"},
							Usage:   "namespace",
							Value:   "default",
						},
						&cli.StringFlag{
							Name:     "name",
							Usage:    "route name",
							Required: true,
						},
					},
					Action: func(ctx *cli.Context) error {
						delRoute(ctx.String("namespace"), ctx.String("name"), store)
						return nil
					},
				},
				{
					Name:  "list",
					Usage: "list namespace routes",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:    "namespace",
							Aliases: []string{"ns"},
							Usage:   "namespace",
							Value:   "default",
						},
					},
					Action: func(ctx *cli.Context) error {
						listRoutes(ctx.String("namespace"), store)
						return nil
					},
				},
			},
		},
	}

	app.Run(os.Args)
}
