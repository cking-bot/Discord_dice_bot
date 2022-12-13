package command

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/mitchellh/mapstructure"
	"log"
	"net/http"
	"os"
	"rps_bot/models"
)

func Roll(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Access options in the order provided by the user.

	// Or convert the slice into a map
	var response string
	var options = make(map[string]interface{})
	for _, option := range i.ApplicationCommandData().Options {
		options[option.Name] = option.Value
	} //given by discord

	var cmd models.Command
	err := mapstructure.Decode(options, &cmd) //throws map into struct using mapstructure
	if err != nil {
		log.Println(err)
	} //structs are easier then maps

	log.Println(cmd)

	body, err := json.Marshal(cmd)
	if err != nil {
		log.Print(err)
		return
	}

	req, err := http.NewRequest("POST", os.Getenv("URL")+"/roll", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Print(err)
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		log.Print(err)
		return
	}

	//Responds to the command
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		// Ignore type for now, they will be discussed in "responses"
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf(
				response,
			),
		},
	})
}
