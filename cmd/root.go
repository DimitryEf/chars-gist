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

var (
	cfgFile string
	path    string
	out     string
)

var rootCmd = &cobra.Command{
	Use:   "chars-gist",
	Short: "Creates a histogram of character usage",
	Long: `Creates a histogram of ASCII character usage 
in text files from the specified folder.
Example:
chars-gist --path "C:/example"
or
chars-gist --path "example"`,

	// Run function create histogram of character usage
	Run: func(cmd *cobra.Command, args []string) {
		// Check args
		if path == "" {
			fmt.Println("Please set a path to folder with files")
			return
		}
		fmt.Println("Starting the app...")

		// Getting a folder with files
		files, err := ioutil.ReadDir(path)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Making a map to collect the result and some sync things
		m := make(map[int32]int64)
		mutex := &sync.Mutex{}
		wg := &sync.WaitGroup{}
		wg.Add(len(files))

		for _, file := range files {
			// Reading each file in goroutine
			go func(filename string) {
				// Open the file
				f, err := os.Open(path + "/" + filename)
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

				// Reading the file
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
					// Adding a char in the map
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

		// Waiting all goroutines
		wg.Wait()

		// Creating a file for the result with time stamp
		datetime := fmt.Sprint(time.Now().Format("2006-01-02_15-04-05"))
		resultFileName := "gist_" + datetime + ".txt"
		if out != "" {
			resultFileName = out + "/" + resultFileName
		}
		resultFile, err := os.Create(resultFileName)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Writing the results line by line
		for key, value := range m {
			_, err := resultFile.WriteString(fmt.Sprintf("%q %d\n", key, value))
			if err != nil {
				fmt.Println(err)
				return
			}
		}

		fmt.Printf("The result is saved in %s/%s\n", path, resultFileName)

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
	//rootCmd.Flags().String("path", "", "Help message for toggle")
	rootCmd.Flags().StringVar(&path, "path", "", "example: --path \"C:/example\"")
	rootCmd.Flags().StringVar(&out, "out", "", "example: --out \"C:/out\"")
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
