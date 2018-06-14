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

package profile

import (
	"fmt"
	"os"
	"os/user"
	path "path/filepath"

	common "github.com/MottainaiCI/mottainai-cli/common"
	tools "github.com/MottainaiCI/mottainai-cli/common"
	cobra "github.com/spf13/cobra"
	v "github.com/spf13/viper"
)

func newProfileCreateCommand() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "create <profile-name> <api-url> [OPTIONS]",
		Short: "Create a new profile",
		Args:  cobra.RangeArgs(2, 2),
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			var name, master, f string
			var conf common.ProfileConf

			name = args[0]
			master = args[1]

			if v.Get("profiles") == nil {
				// POST: No configuration file found

				conf = *common.NewProfileConf()
				err = conf.AddProfile(name, master)
				tools.CheckError(err)

			} else {
				// POST: A configuration file is already present.

				err = v.Unmarshal(&conf)
				tools.CheckError(err)

				p, _ := conf.GetProfile(name)
				if p != nil {
					fmt.Printf("Profile %s is already present.\n", name)
					return
				}

				err = conf.AddProfile(name, master)
				tools.CheckError(err)
			}

			// Create new viper configuration to avoid
			// write of command line arguments/settings
			viper := v.New()
			viper.SetConfigType("yaml")
			viper.Set("profiles", conf.Profiles)

			if v.ConfigFileUsed() != "" {
				f = v.ConfigFileUsed()
			} else {

				user, _ := user.Current()
				f = fmt.Sprintf("%s/%s/%s.yml",
					user.HomeDir,
					common.MCLI_HOME_PATH,
					common.MCLI_CONFIG_NAME)

				// Create directory where save file if doesn't exists
				if _, err := os.Stat(path.Dir(f)); os.IsNotExist(err) {
					err = os.Mkdir(path.Dir(f), 0760)
					tools.CheckError(err)
				}
			}

			err = viper.WriteConfigAs(f)
			tools.CheckError(err)

			fmt.Printf("Profile %s with url %s added on file %s.\n",
				name, master, f)
		},
	}

	return cmd
}