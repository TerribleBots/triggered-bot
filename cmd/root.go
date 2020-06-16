package cmd

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/mitchellh/go-homedir"
	. "github.com/spf13/cobra"
	. "github.com/spf13/viper"
	"math/rand"
	"os"
	"time"
	. "triggered-bot/bot"
	. "triggered-bot/log"
)

var (
	cfgFile                 string
	sourceFile, includeFile string
	sampleRatio             float64
)

var rootCmd = &Command{
	Use:   "triggered-bot",
	Short: "A very stupid discord bot that is easily triggered.",
	Long: `
Provided with a source set of words, triggered-bot will randomly select a configurable percentage as trigger words.
If anyone makes a post that contains a trigger word, triggered-bot will kick the creator from the server, but send
them a dm with an apology and an invitation to rejoin. The messages posted to the channel and text in the apology are
also configurable and chosen at random from a set of templates.

Note: Triggered bot will not attempt to kick the server owner, nor does it have the ability to do so.`,
	Run: func(cmd *Command, args []string) {
		Log.Info("Creating bot...")
		b := newBot()
		Log.Info("Running bot...")
		b.Run()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "config.toml", "config file")
	rootCmd.Flags().StringVar(&sourceFile, "source", "words.txt", "source of words to be sampled by triggered bot")
	rootCmd.Flags().StringVar(&includeFile, "include", "", "words to be unconditionally included by triggered bot")
	rootCmd.Flags().Float64Var(&sampleRatio, "sample-ratio", .10, "percentage of words to be randomly sampled by source")
	SetDefault("reason-templates", ReasonTemplates)
	SetDefault("apology-templates", ApologyTemplates)
}

func initConfig() {
	if cfgFile != "" {
		SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		AddConfigPath(home)
		SetConfigName(".triggered-bot")
	}

	AutomaticEnv()
	WatchConfig()
	InitLogger()
	rand.Seed(time.Now().Unix())

	OnConfigChange(func(e fsnotify.Event) { InitLogger() })
}

func newBot() *Bot {
	token := GetString("token")

	if token == "" {
		Log.Fatal("No auth token provided")
	}

	return &Bot{
		Token:            token,
		ApologyTemplates: GetStringSlice("apology-templates"),
		ReasonTemplates:  GetStringSlice("reason-templates"),
		Matcher:          NewSimpleMatcher(words),
		Sampler:          NewSampler(sourceFile, includeFile, sampleRatio),
	}
}
