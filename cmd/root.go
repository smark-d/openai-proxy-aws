package cmd

import (
	"github.com/go-redis/redis"
	"github.com/smark-d/openai-proxy-aws/server"
	"github.com/smark-d/openai-proxy-aws/server/api"
	"github.com/spf13/cobra"
	"log"
	"strconv"
)

var RootCommand = &cobra.Command{
	Use:   "OpenAI Proxy",
	Short: "OpenAI Proxy",
	Long:  `OpenAI Proxy is a proxy server for OpenAI's API, which limits the number of requests per user`,
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the server",
	Long:  `Start the server with the specified port`,
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetString("port")
		server.Start(port)
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure the application",
	Long:  `Configure the application with various options`,
}

var addUserCmd = &cobra.Command{
	Use:   "addUser",
	Short: "Add a user",
	Long:  `Add a user with the specified limit count`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		username := args[0]
		limitCount := args[1]
		i, err := strconv.ParseInt(limitCount, 10, 64)
		if err != nil {
			i = 100
		}
		api.SetTotalCount(username, i)
		log.Printf("Add user %s with limit count %d", username, i)
	},
}

var listUserCmd = &cobra.Command{
	Use:   "listUser",
	Short: "List all users",
	Long:  `List all users with their limit counts`,
	Run: func(cmd *cobra.Command, args []string) {
		users := api.ListTotalUser()
		for _, user := range users {
			user = user[len(api.TOTAL_COUNT_PREFIX):]
			totalCount, err := api.GetTotalCount(user)
			if err != nil && err != redis.Nil {
				panic(err)
			}
			currCount, err := api.GetCurrCount(user)
			if err != nil && err != redis.Nil {
				panic(err)
			}
			log.Printf("User: %s, Limit: %d, Current: %d \n", user, totalCount, currCount)
		}
	},
}

var setLimitCmd = &cobra.Command{
	Use:   "setLimit",
	Short: "Set the limit for a user",
	Long:  `Set the limit for a user with the specified user and limit`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		user := args[0]
		limit := args[1]
		i, err := strconv.ParseInt(limit, 10, 64)
		if err != nil {
			i = 100
		}
		api.SetTotalCount(user, i)
	},
}

var removeUserCmd = &cobra.Command{
	Use:   "removeUser",
	Short: "Remove a user",
	Long:  `Remove a user with the specified user`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		user := args[0]
		log.Println(user)
		api.RemoveTotalCount(user)
	},
}

var addApiKeyCmd = &cobra.Command{
	Use:   "addApiKey",
	Short: "Add an API key",
	Long:  "Add an API key with the specified key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		key := args[0]
		api.AddOpenAIkey(key)
	},
}

func init() {
	serverCmd.Flags().String("port", "8080", "server port")
	configCmd.AddCommand(addUserCmd)
	configCmd.AddCommand(listUserCmd)
	configCmd.AddCommand(setLimitCmd)
	configCmd.AddCommand(removeUserCmd)
	configCmd.AddCommand(addApiKeyCmd)
	RootCommand.AddCommand(serverCmd)
	RootCommand.AddCommand(configCmd)
}
