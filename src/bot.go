package src

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"

	"gopkg.in/mgo.v2"

	"github.com/tattsun/mood-linebot/src/model"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/pkg/errors"
)

type Bot struct {
	Config  *Config
	lineBot *linebot.Client

	userRepository *model.UserRepository
	moodRepository *model.MoodRepository
}

func NewBot(config *Config) (*Bot, error) {
	// Initialize Line Bot
	client := &http.Client{}
	bot, err := linebot.New(
		config.LINE.ChannelSecret,
		config.LINE.ChannelToken,
		linebot.WithHTTPClient(client))
	if err != nil {
		return nil, errors.Wrap(err, "failed to init line bot client")
	}

	// Initialize DB
	session, err := mgo.Dial(config.DB.MongoAddr)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect mongoDB")
	}
	db := session.DB(config.DB.MongoDatabase)
	userRepository, err := model.NewUserRepository(db)
	if err != nil {
		return nil, err
	}
	moodRepository, err := model.NewMoodRepository(db)
	if err != nil {
		return nil, err
	}

	// Bot creation
	return &Bot{
		Config:         config,
		lineBot:        bot,
		moodRepository: moodRepository,
		userRepository: userRepository,
	}, nil
}

func (b *Bot) RunServer() error {
	http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
		events, err := b.lineBot.ParseRequest(req)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}

		for _, event := range events {
			b.handleEvent(event)
		}
	})

	return http.ListenAndServe(":"+b.Config.Port, nil)
}

func (b *Bot) SendFeelingCheck() error {
	users, err := b.userRepository.FindAll()
	if err != nil {
		return errors.Wrap(err, "failed to find users")
	}

	messages := []string{
		"0: Miserable,nervous",
		"1: Sad,unhappy",
		"2: down,worried",
		"3: good,alright",
		"4: happy,excited",
		"5: pumped,energized",
	}
	buttons := make([]*linebot.QuickReplyButton, len(messages))
	for i, message := range messages {
		action := linebot.NewMessageAction(message, message)
		buttons[i] = linebot.NewQuickReplyButton("", action)
	}
	items := linebot.NewQuickReplyItems(buttons...)
	msg := linebot.NewTextMessage("How are you feeling today?").WithQuickReplies(items)

	for _, user := range users {
		_, err := b.lineBot.PushMessage(user.UserID, msg).Do()
		if err != nil {
			log.Print(err)
		}
	}

	return nil
}

func (b *Bot) handleEvent(event *linebot.Event) {
	defer func() {
		err := recover()
		if err != nil {
			log.Print(err)
		}
	}()

	if event.Type != linebot.EventTypeMessage {
		return
	}

	switch message := event.Message.(type) {
	case *linebot.TextMessage:
		if message.Text == "register" {
			user := &model.User{
				ID:     bson.NewObjectId(),
				UserID: event.Source.UserID,
			}
			err := b.userRepository.Create(user)
			if err != nil {
				b.sendError(event, err)
				return
			}
			b.sendTextMsg(event, "OK")
			return
		}

		splited := strings.Split(message.Text, ":")
		if len(splited) == 0 {
			b.sendError(event, errors.New("invalid msg"))
			return
		}

		mood, err := strconv.Atoi(splited[0])
		if err != nil {
			b.sendError(event, errors.New("invalid msg"))
			return
		}

		err = b.moodRepository.Create(&model.Mood{
			ID:        bson.NewObjectId(),
			UserID:    event.Source.UserID,
			Mood:      mood,
			Timestamp: time.Now(),
		})
		if err != nil {
			b.sendError(event, err)
			return
		}
		b.sendTextMsg(event, "Got it! It's marked in the books!")
	}
}

func (b *Bot) sendTextMsg(event *linebot.Event, text string) {
	_, err := b.lineBot.ReplyMessage(event.ReplyToken,
		linebot.NewTextMessage(text)).Do()
	if err != nil {
		log.Print(err)
	}
}

func (b *Bot) sendError(event *linebot.Event, srcErr error) {
	_, err := b.lineBot.ReplyMessage(event.ReplyToken,
		linebot.NewTextMessage(fmt.Sprintf("err: %s", srcErr))).Do()
	if err != nil {
		log.Print(err)
	}
}
