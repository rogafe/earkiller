package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

func init() {
	log.SetFlags(log.Llongfile | log.LstdFlags)
}

var token string
var buffer = make([][]byte, 0)

var guildID []*discordgo.Guild

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token = os.Getenv("DISCORD_TOKEN")

	if token == "" {
		fmt.Println("NeedToken pls")
		return
	}

	// Load the sound file.

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}

	// Register ready as a callback for the ready events.
	dg.AddHandler(ready)

	dg.AddHandler(guildCreate)

	// Register messageCreate as a callback for the messageCreate events.
	dg.AddHandler(messageCreate)

	// Register guildCreate as a callback for the guildCreate events.

	// We need information about guilds (which includes their channels),
	// messages and voice states.
	dg.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsGuildVoiceStates

	// Open the websocket and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
	}

	// Wait here until CTRL-C or other term signal is received.
	log.Println("earkiller is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()

}

// This function will be called (due to AddHandler above) when the bot receives
// the "ready" event from Discord.
func ready(s *discordgo.Session, event *discordgo.Ready) {
	config := GetConfig()
	log.Println("ready")
	// Set the playing status.
	s.UpdateGameStatus(0, config.Status)
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.

func ComesFromDM(s *discordgo.Session, m *discordgo.MessageCreate) (bool, error) {
	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		if channel, err = s.Channel(m.ChannelID); err != nil {
			return false, err
		}
	}

	return channel.Type == discordgo.ChannelTypeDM, nil
}
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	channel, err := s.State.Channel(m.ChannelID)
	if err != nil {
		log.Println(err)
	}
	_ = channel
	// config := GetConfig()
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	// config := GetConfig()
	msg := string(m.Message.Content)
	// _, found := Find(config.Admin, m.Author.ID)
	found := true
	if found {
		if msg[0:1] == "*" {
			splitString := strings.Split(m.Content, " ")
			// log.Println(splitString, len(splitString), m.Author.ID)
			switch splitString[0] {
			case "*TAGUEULE":
				filename := "tg.dca"
				log.Println("case *TAGUEULE")
				if len(splitString) > 1 {
					m1 := regexp.MustCompile(`[<@!>]`)
					id := splitString[1]
					id = m1.ReplaceAllString(id, "")
					for _, guild := range guildID {
						for _, voice := range guild.VoiceStates {
							log.Println(voice.UserID, id, voice.GuildID)
							if voice.UserID == id {
								log.Printf("found user %s in guild %s ", voice.UserID, guild.Name)
								err := playSound(s, voice.GuildID, voice.ChannelID, filename)
								if err != nil {
									fmt.Println("Error playing sound:", err)
								}

							}

						}
					}
				}
			case "*DECIDE":
				filename := "decide.dca"
				log.Println("case *DECIDE")
				if len(splitString) > 1 {
					m1 := regexp.MustCompile(`[<@!>]`)
					id := splitString[1]
					id = m1.ReplaceAllString(id, "")
					for _, guild := range guildID {
						for _, voice := range guild.VoiceStates {
							log.Println(voice.UserID, id, voice.GuildID)
							if voice.UserID == id {
								log.Printf("found user %s in guild %s ", voice.UserID, guild.Name)
								err := playSound(s, voice.GuildID, voice.ChannelID, filename)
								if err != nil {
									fmt.Println("Error playing sound:", err)
								}

							}

						}
					}
				}
			case "*PIRATE":
				filename := "pirate.dca"
				log.Println("case *PIRATE")
				if len(splitString) > 1 {
					m1 := regexp.MustCompile(`[<@!>]`)
					id := splitString[1]
					id = m1.ReplaceAllString(id, "")
					for _, guild := range guildID {
						for _, voice := range guild.VoiceStates {
							log.Println(voice.UserID, id, voice.GuildID)
							if voice.UserID == id {
								log.Printf("found user %s in guild %s ", voice.UserID, guild.Name)
								err := playSound(s, voice.GuildID, voice.ChannelID, filename)
								if err != nil {
									fmt.Println("Error playing sound:", err)
								}

							}

						}
					}
				}
			case "*help":
				// s.ChannelMessageSend(m.ChannelID, "Command are, *DECIDE and !TAGUEULE use it with @nickofuser")
				embed := NewEmbed().
					SetTitle("Help").
					SetDescription("Merci a fanta pour les sound byte").
					AddField("*DECIDE", "Nom mais c'est pas toi qui decide").
					AddField("*TAGGUELE", "TG").
					SetColor(0x00ff00).MessageEmbed

				s.ChannelMessageSendEmbed(m.ChannelID, embed)

			}

		}
	}
}

// This function will be called (due to AddHandler above) every time a new
// guild is joined.
func guildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {

	if event.Guild.Unavailable {
		return
	}
	guildID = append(guildID, event.Guild)

}

// loadSound attempts to load an encoded sound file from disk.
func loadSound(filename string) error {

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening dca file :", err)
		return err
	}

	var opuslen int16

	for {
		// Read opus frame length from dca file.
		err = binary.Read(file, binary.LittleEndian, &opuslen)

		// If this is the end of the file, just return.
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			err := file.Close()
			if err != nil {
				return err
			}
			return nil
		}

		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return err
		}

		// Read encoded pcm from dca file.
		InBuf := make([]byte, opuslen)
		err = binary.Read(file, binary.LittleEndian, &InBuf)

		// Should not be any end of file errors
		if err != nil {
			fmt.Println("Error reading from dca file :", err)
			return err
		}

		// Append encoded pcm data to the buffer.
		buffer = append(buffer, InBuf)

	}
}

// playSound plays the current buffer to the provided channel.
func playSound(s *discordgo.Session, guildID, channelID string, filename string) (err error) {

	err = loadSound(filename)
	if err != nil {
		fmt.Println("Error loading sound: ", err)
		fmt.Println("Please.")
	}
	// Join the provided voice channel.
	vc, err := s.ChannelVoiceJoin(guildID, channelID, false, true)
	if err != nil {
		return err
	}

	// Sleep for a specified amount of time before playing the sound
	time.Sleep(1 * time.Millisecond)

	// Start speaking.
	vc.Speaking(true)

	// Send the buffer data.
	for _, buff := range buffer {
		vc.OpusSend <- buff
	}

	// Stop speaking
	vc.Speaking(false)

	// Sleep for a specificed amount of time before ending.
	time.Sleep(1 * time.Millisecond)

	// Disconnect from the provided voice channel.
	vc.Disconnect()

	for i := range buffer {
		buffer[i] = make([]byte, 0)
	}
	return nil
}
