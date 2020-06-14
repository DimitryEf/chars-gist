package cmd

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"io/ioutil"
	"os"
	"sync"
	"time"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "chars-gist",
	Short: "Creates a histogram of character usage",
	Long: `Creates a histogram of ASCII character usage 
in text files from the specified folder.`,

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println("Please set a path to folder with files")
			return
		}
		fmt.Println("Starting the app...")
		files, err := ioutil.ReadDir(args[1])
		if err != nil {
			fmt.Println(err)
			return
		}

		m := make(map[int32]int64)
		mutex := &sync.Mutex{}
		wg := &sync.WaitGroup{}
		wg.Add(len(files))

		for _, file := range files {
			go func(filename string) {
				f, err := os.Open(args[1] + "/" + filename)
				if err != nil {
					fmt.Println(err)
					return
				}

				defer func() {
					err := f.Close()
					if err != nil {
						fmt.Println(err)
						return
					}
				}()

				reader := bufio.NewReader(f)
				for {
					char, _, err := reader.ReadRune()
					if err != nil {
						if err == io.EOF {
							break
						}
						fmt.Println(err)
						return
					}
					mutex.Lock()
					if _, ok := m[char]; ok {
						m[char]++
					} else {
						m[char] = 1
					}
					mutex.Unlock()
				}
				wg.Done()
			}(file.Name())
		}

		wg.Wait()

		datetime := fmt.Sprint(time.Now().Format("2006-01-02_15-04-05"))
		resultFileName := "gist_" + datetime + ".txt"
		resultFile, err := os.Create(resultFileName)
		if err != nil {
			fmt.Println(err)
			return
		}
		for key, value := range m {
			_, err := resultFile.WriteString(fmt.Sprintf("%q %d\n", key, value))
			if err != nil {
				fmt.Println(err)
				return
			}
		}

		fmt.Printf("The result is saved in %s/%s\n", args[1], resultFileName)

	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.chars-gist.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().String("path", "", "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".chars-gist" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".chars-gist")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
