package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/go-co-op/gocron"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
	. "triggered-bot/log"
	"unicode"
)

type Bot struct {
	Token                             string
	ReasonTemplates, ApologyTemplates []string
	Matcher                           Matcher
	Sampler                           Sampler
}

func (b *Bot) Run() {
	dg, err := discordgo.New("Bot " + b.Token)
	if err != nil {
		Log.Fatal(err)
	}

	dg.AddHandler(b.messageHandler)
	err = dg.Open()
	if err != nil {
		Log.Fatal(err)
	}

	defer dg.Close()
	sched := gocron.NewScheduler(time.UTC)
	_, err = sched.Every(1).Day().Do(b.refreshWords)

	if err != nil {
		Log.Fatal(err)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

func (b *Bot) refreshWords() {
	Log.Info("Refreshing trigger words")
	words := b.Sampler.SampleWords()
	b.Matcher.SetWords(words)
}

func (b *Bot) messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if err := b.listen(s, m); err != nil {
		Log.Error(err)
	}
}

func (b *Bot) listen(s *discordgo.Session, m *discordgo.MessageCreate) error {
	if m.Author.ID == s.State.User.ID {
		return nil
	}

	c := preProcess(m.Content)
	if match := b.Matcher.Match(c); match != "" {
		return b.handleMatch(s, m, match)
	}

	return nil
}

func preProcess(candidate string) string {
	candidate = strings.ToLower(candidate)
	candidate = strings.TrimSpace(candidate)
	candidate = strings.TrimFunc(candidate, unicode.IsPunct)
	return candidate
}

func (b *Bot) handleMatch(s *discordgo.Session, m *discordgo.MessageCreate, match string) error {
	a := m.Author
	r := b.makeReason(match)
	msg := fmt.Sprintf("%s %s", a.Mention(), r)
	if _, err := s.ChannelMessageSend(m.ChannelID, msg); err != nil {
		return fmt.Errorf("unable to create message: %s", err)
	}

	g, err := s.Guild(m.GuildID)
	if err != nil {
		return fmt.Errorf("unable to get guild info: %s", err)
	}

	if a.ID != g.OwnerID {
		return b.kickMember(s, m, r, a.ID)
	}

	return nil
}

func (b *Bot) kickMember(s *discordgo.Session, m *discordgo.MessageCreate, r, id string) error {
	if err := b.sendApology(s, m, id); err != nil {
		return err
	}

	if err := s.GuildMemberDeleteWithReason(m.GuildID, id, r); err != nil {
		return fmt.Errorf("unable to kick user %s: %s", id, err)
	}

	return nil
}

func (b *Bot) sendApology(s *discordgo.Session, m *discordgo.MessageCreate, id string) error {
	uc, err := s.UserChannelCreate(id)
	if err != nil {
		return fmt.Errorf("unable to dm user %s: %s", id, err)
	}

	i, err := s.ChannelInviteCreate(m.ChannelID, discordgo.Invite{})
	if err != nil {
		return fmt.Errorf("unable to create invite: %s", err)
	}

	if _, err = s.ChannelMessageSend(uc.ID, b.makeApology(i.Code)); err != nil {
		return fmt.Errorf("unable to send apology invite: %s", err)
	}

	return nil
}

func (b *Bot) makeReason(match string) string {
	return fmt.Sprintf(sample(b.ReasonTemplates), match)
}

func (b *Bot) makeApology(code string) string {
	url := fmt.Sprintf("https://discord.gg/%s", code)
	return fmt.Sprintf(sample(b.ApologyTemplates), url)
}

func sample(s []string) string {
	return s[rand.Intn(len(s))]
}
