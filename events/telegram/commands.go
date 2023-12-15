package telegram

import (
	"context"
	"errors"
	"log"
	"net/url"
	"strings"
	"urlshortener/lib/e"
	"urlshortener/storage"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(ctx context.Context, text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command `%s` from user `%s` in chat `%d`", text, username, chatID)

	if isAddCmd(text) {
		return p.savePage(ctx, chatID, text, username)

	}

	switch text {
	case RndCmd:
		return p.SendRandom(ctx, chatID, username)
	case HelpCmd:
		return p.SendHelp(ctx, chatID)
	case StartCmd:
		return p.SendHello(ctx, chatID)
	default:
		return p.tg.SendMessage(ctx, chatID, msgUnknownCommand)
	}
}

func (p *Processor) savePage(ctx context.Context, chatID int, pageURL string, username string) (err error) {
	defer func() { err = e.Wrap("cant do cmd: save page", err) }()

	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	isExists, err := p.storage.IsExists(ctx, page)
	if err != nil {
		return err
	}
	if isExists {
		return p.tg.SendMessage(ctx, chatID, msgAlreadyExists) //
	}

	if err := p.storage.Save(ctx, page); err != nil {
		return err
	}
	if err := p.tg.SendMessage(ctx, chatID, msgSaved); err != nil {
		return err
	}
	return nil
}

func (p *Processor) SendRandom(ctx context.Context, chatID int, username string) (err error) {
	defer func() { err = e.Wrap("cant do command", err) }()

	page, err := p.storage.PickRandom(ctx, username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}
	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMessage(ctx, chatID, msgNoSavedPages)
	}
	if err := p.tg.SendMessage(ctx, chatID, page.URL); err != nil {
		return err
	}

	return p.storage.Remove(ctx, page)
}

func (p *Processor) SendHelp(ctx context.Context, chatID int) error {
	return p.tg.SendMessage(ctx, chatID, msgHelp)
}

func (p *Processor) SendHello(ctx context.Context, chatID int) error {
	return p.tg.SendMessage(ctx, chatID, msgHello)
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
