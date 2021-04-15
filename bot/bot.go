package bot

import (
	v1 "github.com/ebladrocher/voicebot/app/converter/v1"
	search "github.com/ebladrocher/voicebot/app/search"
	model "github.com/ebladrocher/voicebot/system/model"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// Bot ...
type Bot struct {
	Token      string
	APIKey     string
	IDChannel  int64
	TassAPIKey string
	IDUser     int
	Debug      bool
}

// Start ...
func (cfg *Bot) Start() {
	port := os.Getenv("PORT")
	log.Print(port)

	url := os.Getenv("URL")

	bot, err := tgbotapi.NewBotAPI(cfg.Token)
	if err != nil {
		log.Fatal("NewBotAPI: ", err)
	}

	bot.Debug = cfg.Debug
	log.Printf("Authorized on account %s", bot.Self.UserName)
	if cfg.Debug == true {
		_, err = bot.SetWebhook(tgbotapi.NewWebhook("" + bot.Token))
		if err != nil {
			log.Fatal("SetWebhook", err)
		}
		//defer log.Fatal(bot.RemoveWebhook())
	} else {
		_, err = bot.SetWebhook(tgbotapi.NewWebhook(url + bot.Token))
		if err != nil {
			log.Fatal("SetWebhook", err)
		}
	}

	updates := bot.ListenForWebhook("/" + bot.Token)

	time.Sleep(time.Millisecond * 500)
	updates.Clear()
	if cfg.Debug == true {
		go func() {
			log.Fatal("ListenAndServe:", http.ListenAndServe(":3000", nil))
		}()
	} else {
		go func() {
			log.Fatal("ListenAndServe:", http.ListenAndServe(":"+port, nil))
		}()
	}

	for update := range updates {
		if update.Message != nil {
			if update.Message.From.ID == cfg.IDUser {
				if update.Message.Voice != nil {
					cfgVoice := tgbotapi.NewVoiceShare(cfg.IDChannel, update.Message.Voice.FileID)
					bot.Send(cfgVoice)

				} else if update.Message.Audio != nil {
					if update.Message.Caption == "" {
						text := "дебил напиши название"
						bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, text))
						continue
					}

					caption := update.Message.Caption

					id := update.Message.Audio.FileID
					cfgFile := tgbotapi.FileConfig{
						FileID: id,
					}

					file, err := bot.GetFile(cfgFile)
					if err != nil {
						log.Print(err)
					}

					url := "https://api.telegram.org/file/bot" + cfg.Token + "/" + file.FilePath
					log.Print(url)

					idReq, tries := v1.GetID(url, cfg.APIKey)
					if tries <= 0 {
						text := "бесплатные преобразования закончились, попробуй завтра"
						bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, text))
						continue
					}
					text := "осталось попыток: " + strconv.FormatInt(tries, 10)
					bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, text))
					state := v1.GetStatus(idReq)
					requestBytes := v1.GetFile(state)

					// fileReader := tgbotapi.FileReader{
					// 	Name:   update.Message.Audio.FileID,
					// 	Reader: file2,
					// 	Size:   int64(update.Message.Audio.FileSize),
					// }

					fileBytes := tgbotapi.FileBytes{
						Name:  "voice_message",
						Bytes: requestBytes,
					}

					cfgVoice := tgbotapi.NewVoiceUpload(cfg.IDChannel, fileBytes)
					cfgVoice.Caption = caption
					bot.Send(cfgVoice)

				}
			} else {
				text := "это Inline бот"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
				bot.Send(msg)
			}

		} else if update.InlineQuery != nil {
			inlineQuery := update.InlineQuery
			queryID := inlineQuery.ID
			queryText := inlineQuery.Query

			var req search.RequestInline
			req.SetRequestInline(cfg.TassAPIKey, "getChatHistory", strconv.FormatInt(cfg.IDChannel, 10), "100", "-1", "0")

			posts, count := search.GetRequestInline(req)
			log.Println(count)

			cfgInline := tgbotapi.InlineConfig{
				InlineQueryID: queryID,
				CacheTime:     0,
			}

			result := make([]interface{}, 0)
			m := search.GetMap(posts)
			i := 0
			for id, text := range m {
				if strings.Contains(text, queryText) {
					voice := tgbotapi.InlineQueryResultCachedVoice{
						Type:    "voice",
						ID:      strconv.Itoa(i),
						VoiceID: id,
						Title:   text,
					}
					result = append(result, voice)
					i++
				}
			}

			cfgInline.Results = result
			if _, err := bot.AnswerInlineQuery(cfgInline); err != nil {
				log.Println(err)
			}

		}
	}
}

// NewBot ...
func NewBot(cfg *model.Config) (*Bot, error) {
	if cfg.Debug == true {
		return &Bot{
			Token:      cfg.Token,
			APIKey:     cfg.APIKey,
			IDChannel:  cfg.IDChannel,
			TassAPIKey: cfg.TassAPIKey,
			IDUser:     cfg.IDUser,
			Debug:      cfg.Debug,
		}, nil
	}

	token := os.Getenv("TOKEN")
	log.Println(token)
	apikey := os.Getenv("APIKEY")
	log.Println(apikey)
	idchannel, _ := strconv.ParseInt(os.Getenv("IDCHANNEL"), 10, 64)
	log.Println(idchannel)
	tassapikey := os.Getenv("TASSAPIKEY")
	log.Println(tassapikey)
	iduser, _ := strconv.Atoi(os.Getenv("IDUSER"))
	log.Println(iduser)
	return &Bot{
		Token:      token,
		APIKey:     apikey,
		IDChannel:  idchannel,
		TassAPIKey: tassapikey,
		IDUser:     iduser,
		Debug:      false,
	}, nil
}
