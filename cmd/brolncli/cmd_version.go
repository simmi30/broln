package main

import (
	"fmt"

	"github.com/brsuite/broln/build"
	"github.com/brsuite/broln/lnrpc/lnclipb"
	"github.com/brsuite/broln/lnrpc/verrpc"
	"github.com/urfave/cli"
)

var versionCommand = cli.Command{
	Name:  "version",
	Usage: "Display brolncli and broln version info.",
	Description: `
	Returns version information about both brolncli and broln. If brolncli is unable
	to connect to broln, the command fails but still prints the brolncli version.
	`,
	Action: actionDecorator(version),
}

func version(ctx *cli.Context) error {
	ctxc := getContext()
	conn := getClientConn(ctx, false)
	defer conn.Close()

	versions := &lnclipb.VersionResponse{
		Brolncli: &verrpc.Version{
			Commit:        build.Commit,
			CommitHash:    build.CommitHash,
			Version:       build.Version(),
			AppMajor:      uint32(build.AppMajor),
			AppMinor:      uint32(build.AppMinor),
			AppPatch:      uint32(build.AppPatch),
			AppPreRelease: build.AppPreRelease,
			BuildTags:     build.Tags(),
			GoVersion:     build.GoVersion,
		},
	}

	client := verrpc.NewVersionerClient(conn)

	brolnVersion, err := client.GetVersion(ctxc, &verrpc.VersionRequest{})
	if err != nil {
		printRespJSON(versions)
		return fmt.Errorf("unable fetch version from broln: %v", err)
	}
	versions.broln = brolnVersion

	printRespJSON(versions)

	return nil
}
