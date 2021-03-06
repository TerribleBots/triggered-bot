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
	"triggered-bot/text"
)

const externalConj = "Also,"
const inSectionConj = "and"

type Bot struct {
	Token                                              string
	ReasonTemplates, ApproxTemplates, ApologyTemplates []string
	Matcher                                            Matcher
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

	sched.StartAsync()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

func (b *Bot) refreshWords() {
	Log.Info("Refreshing trigger words")
	words := b.Matcher.GetSampler().SampleWords()
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

	c := text.Normalize(m.Content)

	if text.Overriden(c) {
		return nil
	}

	if match := b.Matcher.Match(c); match.AnyMatch() {
		return b.handleMatch(s, m, match)
	}

	return nil
}

func (b *Bot) handleMatch(s *discordgo.Session, m *discordgo.MessageCreate, match MatchResult) error {
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

func (b *Bot) makeReason(match MatchResult) string {
	m, a := match.matches, match.approximates
	return joinSentence(template(b.ReasonTemplates, m), template(b.ApproxTemplates, a))
}

func (b *Bot) makeApology(code string) string {
	url := fmt.Sprintf("https://discord.gg/%s", code)
	return fmt.Sprintf(sample(b.ApologyTemplates), url)
}

func template(templates []string, s []string) string {
	if len(s) > 0 {
		return fmt.Sprintf(sample(templates), join(s))
	} else {
		return ""
	}
}

//noinspection GoNilness
func join(strs []string) string {
	var out []string

	for _, s := range strs {
		out = append(out, fmt.Sprintf("\"%s\"", s))
	}

	if n := len(out); n == 0 {
		Log.Error("attempted to join empty slice of strings")
	} else if n == 1 {
		return out[0]
	} else if n == 2 {
		return joinPair(out[0], out[1])
	} else {
		return joinPair(strings.Join(out[:n-1], ","), out[n])
	}
	return ""
}

func joinPair(x, y string) string {
	return fmt.Sprintf("%s %s %s", x, inSectionConj, y)
}

func joinSentence(matchReason, approxReason string) string {
	if matchReason == "" {
		return approxReason
	} else if approxReason == "" {
		return matchReason
	} else {
		return fmt.Sprintf("%s. %s %s", matchReason, externalConj, approxReason)
	}
}

func sample(s []string) string {
	return s[rand.Intn(len(s))]
}
