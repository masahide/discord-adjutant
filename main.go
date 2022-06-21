package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/kelseyhightower/envconfig"
)

type Specification struct {
	Token string
}
type UserState struct {
	Name string
}

var ()

func main() {
	var s Specification
	err := envconfig.Process("", &s)
	if err != nil {
		log.Fatal(err.Error())
	}

	var discord *discordgo.Session
	discord, err = discordgo.New("Bot " + s.Token)

	discord.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		fmt.Println("Bot is ready")
	})
	discord.AddHandler(voiceStateUpdate)
	err = discord.Open()
	defer discord.Close()
	fmt.Println("Listening...")
	stopBot := make(chan os.Signal, 1)
	signal.Notify(stopBot, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-stopBot

}

func voiceStateUpdate(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {
	if v == nil {
		log.Print("VoiceStateUpdate is null")
		return
	}
	// チャンネルへの入室ステータスが変更されたとき（ミュートON、OFFに反応しないように分岐）
	if v.BeforeUpdate != nil && v.BeforeUpdate.ChannelID == v.ChannelID {
		log.Printf("change state")
		return
	}
	if v.BeforeUpdate != nil && v.ChannelID == "" {
		log.Printf("exit")
		return
	}
	c, err := s.State.Channel(v.ChannelID)
	if err != nil {
		log.Printf("Channel err:%s", err)
		return
	}
	u, err := s.User(v.UserID)
	if err != nil {
		log.Printf("User err:%s", err)
		return
	}
	msg := fmt.Sprintf("%sが[%s]に入りました", u.Username, c.Name)
	s.ChannelMessageSend(v.ChannelID, msg)
	log.Printf("%sが[%s]に入りました", u.Username, c.Name)
	//pp.Println(v)
}
