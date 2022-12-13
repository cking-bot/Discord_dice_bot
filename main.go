package main

import (
	"flag"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"rps_bot/command"
)

// Bot parameters
var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers command globally") //server ID that we are using it on (currently blank) you can add a value if you want it only on one server
	BotToken       = flag.String("token", os.Getenv("DISCORD_TOKEN"), "Bot access token")                      //given from discord that allows the bot to run. We are currently using railway for this
	RemoveCommands = flag.Bool("rmcmd", false, "Remove all command after shutdowning or not")                  //lets you remove the command
)

var s *discordgo.Session

func init() { flag.Parse() } //Runs prior to main.go, bot needs this to parse the flags

func init() {
	var err error
	s, err = discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	} //creates a new session using the bot token
}

var (
	integerOptionMinValue          = 1.0
	dmPermission                   = false
	defaultMemberPermissions int64 = discordgo.PermissionManageServer

	Commands = []*discordgo.ApplicationCommand{
		{
			Name:        "roll",
			Description: "lets you choose the sides, number of dice, and an option for having advantage or disadvantage",

			Options: []*discordgo.ApplicationCommandOption{

				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "dice",
					Description: "lets you pick which type of dice.",
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "D-4",
							Value: "4",
						},
						{
							Name:  "D-6",
							Value: "6",
						},
						{
							Name:  "D-8",
							Value: "8",
						},
						{
							Name:  "D-10",
							Value: "10",
						},
						{
							Name:  "D-12",
							Value: "12",
						},
						{
							Name:  "D-20",
							Value: "20",
						},
					},
					Required: true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "number",
					Description: "Set the number of dice to roll",
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "1",
							Value: 1,
						},
						{
							Name:  "2",
							Value: 2,
						},
						{
							Name:  "3",
							Value: 3,
						},
						{
							Name:  "4",
							Value: 4,
						},
						{
							Name:  "5",
							Value: 5,
						},
						{
							Name:  "10",
							Value: 10,
						},
					},
					Required: true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "advantage",
					Description: "please select if you have advantage, disadvantage or none.",
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "advantage",
							Value: "1",
						},
						{
							Name:  "disadvantage",
							Value: "2",
						},
						{
							Name:  "none",
							Value: "0",
						},
					},
					Required: true,
				},
			},
		},
	}
	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"roll": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			command.Roll(s, i)
		},
	}
)

func init() {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := CommandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

//even listener, waits for someone to do something and handles that request
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) { //handles any requests
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	
	err = s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("Adding command...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(Commands))
	for i, v := range Commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, *GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop //creates a channel that listens for ctrl+c to stop the code

	if *RemoveCommands {
		log.Println("Removing command...")
		// // We need to fetch the command, since deleting requires the command ID.
		// // We are doing this from the returned command on line 375, because using
		// // this will delete all the command, which might not be desirable, so we
		// // are deleting only the command that we added.
		// registeredCommands, err := s.ApplicationCommands(s.State.User.ID, *GuildID)
		// if err != nil {
		// 	log.Fatalf("Could not fetch registered command: %v", err)
		// }

		for _, v := range registeredCommands {
			err := s.ApplicationCommandDelete(s.State.User.ID, *GuildID, v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}

	log.Println("Gracefully shutting down.")
}
