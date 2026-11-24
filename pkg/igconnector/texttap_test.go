package igconnector

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"go.mau.fi/util/dbutil"
	"go.mau.fi/util/exsync"
	"gopkg.in/yaml.v3"
	"maunium.net/go/mautrix/bridgev2"
	"maunium.net/go/mautrix/bridgev2/commands"
	"maunium.net/go/mautrix/bridgev2/database"
	"maunium.net/go/mautrix/bridgev2/networkid"
	"maunium.net/go/mautrix/event"

	"go.mau.fi/mautrix-meta/pkg/metaid"
)

type textTapFixture struct {
	ic     *IGClient
	portal *bridgev2.Portal
	card   *database.Message
	taps   []*metaid.PostbackButton
}

func newTextTapFixture(t *testing.T, textReplyTaps bool) *textTapFixture {
	t.Helper()
	ctx := context.Background()
	db, err := dbutil.NewWithDialect(filepath.Join(t.TempDir(), "bridge.db"), "sqlite3")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = db.Close() })
	connector := &IGConnector{Config: Config{TextReplyTaps: textReplyTaps}}
	bridge := bridgev2.NewBridge("test", db, zerolog.Nop(), nil, &pollMatrixStub{}, connector, commands.NewProcessor)
	bridge.BackgroundCtx = ctx
	connector.Bridge = bridge
	if err = bridge.DB.Upgrade(ctx); err != nil {
		t.Fatal(err)
	}
	key := networkid.PortalKey{ID: "thread"}
	if err = bridge.DB.Portal.Insert(ctx, &database.Portal{
		PortalKey: key, MXID: "!room:example.com",
		Metadata: &metaid.PortalMetadata{IGID: "send-thread", IGThreadID: "history-thread"},
	}); err != nil {
		t.Fatal(err)
	}
	portal, err := bridge.GetPortalByKey(ctx, key)
	if err != nil {
		t.Fatal(err)
	}
	f := &textTapFixture{portal: portal}
	login := &bridgev2.UserLogin{UserLogin: &database.UserLogin{ID: "login"}}
	f.ic = &IGClient{Main: connector, UserLogin: login, pendingPostbackOTIDs: exsync.NewSet[string]()}
	f.ic.sendTap = func(_ context.Context, threadID string, button *metaid.PostbackButton, _ string) (string, error) {
		if threadID != "send-thread" {
			t.Errorf("sent to thread %q", threadID)
		}
		f.taps = append(f.taps, button)
		return "", nil
	}
	f.card = &database.Message{
		ID: metaid.MakeFBMessageID("77"), MXID: "$card", Room: key, SenderID: "card-ghost", Timestamp: time.Now(),
		Metadata: &metaid.MessageMetadata{
			PostbackButtons: []*metaid.PostbackButton{{Title: "Guide", Payload: "payload-a", CTAType: "postback"}},
			QuickReplies:    []*metaid.PostbackButton{{Title: "Maybe", Payload: "payload-b", CTAType: metaid.CTATypeQuickReply}},
		},
	}
	return f
}

func (f *textTapFixture) reply(text string) (*bridgev2.MatrixMessageResponse, bool, error) {
	return f.ic.handleMatrixPostbackTap(context.Background(), &bridgev2.MatrixMessage{
		MatrixEventBase: bridgev2.MatrixEventBase[*event.MessageEventContent]{
			Event:   &event.Event{ID: "$reply", Sender: "@user:example.com", Timestamp: time.Now().UnixMilli()},
			Portal:  f.portal,
			Content: &event.MessageEventContent{MsgType: event.MsgText, Body: text},
		},
		ReplyTo: f.card,
	})
}

func TestTextReplyDoesNotTapByDefault(t *testing.T) {
	f := newTextTapFixture(t, false)
	for _, text := range []string{"1", "2", "Guide", "maybe"} {
		resp, handled, err := f.reply(text)
		if !handled || resp != nil {
			t.Fatalf("%q: handled=%v resp=%+v, the text must not reach Instagram", text, handled, resp)
		}
		var status bridgev2.MessageStatus
		if !errors.As(err, &status) || !status.IsCertain {
			t.Fatalf("%q: want a certain failed status, got %#v", text, err)
		}
	}
	if len(f.taps) != 0 {
		t.Fatalf("sent %d taps with text_reply_taps off", len(f.taps))
	}
}

func TestTextReplyWithoutMatchIsStillText(t *testing.T) {
	f := newTextTapFixture(t, false)
	for _, text := range []string{"hello", "3", ""} {
		resp, handled, err := f.reply(text)
		if handled || resp != nil || err != nil {
			t.Errorf("%q: handled=%v resp=%+v err=%v, want plain text", text, handled, resp, err)
		}
	}
	if len(f.taps) != 0 {
		t.Fatalf("sent %d taps for text without a match", len(f.taps))
	}
}

func TestTextReplyTapsWhenEnabled(t *testing.T) {
	f := newTextTapFixture(t, true)
	resp, handled, err := f.reply("1")
	if err != nil || !handled || resp == nil || !resp.Pending {
		t.Fatalf("handled=%v resp=%+v err=%v", handled, resp, err)
	}
	if len(f.taps) != 1 || f.taps[0].Payload != "payload-a" {
		t.Fatalf("want one tap on the first button, got %+v", f.taps)
	}
}

func TestTextReplyTapsConfigDefaultsOff(t *testing.T) {
	for name, yml := range map[string]string{
		"key missing":    "tap_listen: 127.0.0.1:31000\n",
		"example":        ExampleConfig,
		"explicit false": "text_reply_taps: false\n",
	} {
		var cfg Config
		if err := yaml.Unmarshal([]byte(yml), &cfg); err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if cfg.TextReplyTaps {
			t.Errorf("%s: text_reply_taps must be false", name)
		}
	}
	var cfg Config
	if err := yaml.Unmarshal([]byte("text_reply_taps: true\n"), &cfg); err != nil || !cfg.TextReplyTaps {
		t.Fatalf("explicit true: %v %v", cfg.TextReplyTaps, err)
	}
}
