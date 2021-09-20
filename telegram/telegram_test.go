package telegram

import "testing"

func TestSendMessage(t *testing.T) {
	type args struct {
		chatID int64
		text   string
		token  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"send hello", args{
			chatID: 0,
			text:   "test",
			token:  "xxxx",
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tl := NewTelegramLog(tt.args.chatID, tt.args.token)
			if err := tl.SendMessage(tt.args.text); (err != nil) != tt.wantErr {
				t.Errorf("SendMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
