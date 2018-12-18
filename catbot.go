package gobot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	owm "github.com/briandowns/openweathermap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)
var (

	Prefix = "!"
	Color = 0xFF304B
	Name = ""
	Token = "NTIxOTc4MjgzNDgzNTk0NzYy.DvaQLA.AJq4srX_WfS4SEokGEvBsPESPdA"
	Cat = "57b50341-c9c0-4f25-a079-8059556baa60"
)


func main() {

	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("Error, can't create a discord session")
		return
	}

	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		fmt.Println("Error, can't open session")
		return
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	tr := &http.Transport{DisableKeepAlives: true}
	client := &http.Client{Transport: tr}

	if m.Author.ID == s.State.User.ID {
		return
	}


	help := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{Name: "Commands"},
		Color:  Color,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   Prefix + "help",
				Value:  "List of all commands",
				Inline: true,
			},
			{
				Name:   Prefix + "cat",
				Value:  "Random cat",
				Inline: true,
			},
			{
				Name:   Prefix + "weather",
				Value:  "Show the weather of Paris",
				Inline: true,
			},

		},

	}
	if strings.HasPrefix(m.Content, Prefix+"help") {
		s.ChannelMessageSendEmbed(m.ChannelID, help)
	}

	if strings.HasPrefix(m.Content, Prefix+"cat") {
		resp, err := client.Get("http://thecatapi.com/api/images/get?api_key=" + Cat + "&format=src")

		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "Error, no cat found :(")
		} else {
			s.ChannelMessageSendEmbed(m.ChannelID, &discordgo.MessageEmbed{
				Author: &discordgo.MessageEmbedAuthor{Name: "Cat Picture"},
				Color:  Color,
				Image: &discordgo.MessageEmbedImage{
					URL: resp.Request.URL.String(),
				},
			})
		}
	}

	if strings.HasPrefix(m.Content, Prefix+"weather") {
		 apiKey := "727b553a8b3950583765288fbdacbde5"

			w, err := owm.NewCurrent("C", "FR", apiKey)
			if err != nil {
				log.Fatalln(err)
			}

			w.CurrentByName("Paris")

		tmp := fmt.Sprintf("Temperature in %s: %v°C \n", w.Name, w.Main.Temp, w.Rain, w.Clouds, w.Snow)
			fmt.Printf("\n")

		s.ChannelMessageSend(m.ChannelID, tmp)


	}

}


