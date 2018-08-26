/*

Copyright (C) 2017-2018  Ettore Di Giacinto <mudler@gentoo.org>
                         Daniele Rondina <geaaru@sabayonlinux.org>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.

*/

package namespace

import (
	"fmt"
	"log"

	tools "github.com/MottainaiCI/mottainai-cli/common"
	client "github.com/MottainaiCI/mottainai-server/pkg/client"
	setting "github.com/MottainaiCI/mottainai-server/pkg/settings"
	cobra "github.com/spf13/cobra"
	viper "github.com/spf13/viper"
)

func newNamespaceRemoveCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "remove <namespace> <absolute_path> [OPTIONS]",
		Short: "Remove a given path from a namespace",
		Args:  cobra.RangeArgs(2, 2),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var fetcher *client.Fetcher
			var v *viper.Viper = setting.Configuration.Viper

			fetcher = client.NewTokenClient(v.GetString("master"), v.GetString("apikey"))

			ns := args[0]
			path := args[1]
			if len(ns) == 0 || len(path) == 0 {
				log.Fatalln("You need to define a namespace and a path to delete")
			}

			res, err := fetcher.GetOptions("/api/namespace/"+ns+"/remove/"+path, map[string]string{})
			tools.CheckError(err)
			fmt.Println(string(res))
		},
	}

	return cmd
}
