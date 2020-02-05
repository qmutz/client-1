package chat

import (
	"context"
	"encoding/json"

	"github.com/keybase/client/go/chat/globals"
	"github.com/keybase/client/go/chat/types"
	"github.com/keybase/client/go/chat/utils"
	"github.com/keybase/client/go/protocol/chat1"
	"github.com/keybase/client/go/protocol/gregor1"
	"github.com/keybase/client/go/protocol/keybase1"
)

type DevConversationBackedStorage struct {
	globals.Contextified
	utils.DebugLabeler

	mt        chat1.ConversationMembersType
	adminOnly bool

	ri func() chat1.RemoteInterface
}

func NewDevConversationBackedStorage(g *globals.Context, mt chat1.ConversationMembersType, adminOnly bool, ri func() chat1.RemoteInterface) *DevConversationBackedStorage {
	return &DevConversationBackedStorage{
		Contextified: globals.NewContextified(g),
		mt:           mt,
		adminOnly:    adminOnly,
		DebugLabeler: utils.NewDebugLabeler(g.GetLog(), "DevConversationBackedStorage", false),
		ri:           ri,
	}
}

func (s *DevConversationBackedStorage) Put(ctx context.Context, uid gregor1.UID, tlfid chat1.TLFID, name string, src interface{}) (err error) {
	defer s.Trace(ctx, func() error { return err }, "Put(%s)", name)()

	info, err := CreateNameInfoSource(ctx, s.G(), s.mt).LookupName(ctx, tlfid, false /* public */, "")
	if err != nil {
		return err
	}
	tlfname := info.CanonicalName

	dat, err := json.Marshal(src)
	if err != nil {
		return err
	}

	conv, err := NewConversation(ctx, s.G(), uid, tlfname, &name, chat1.TopicType_DEV,
		s.mt, keybase1.TLFVisibility_PRIVATE, s.ri, NewConvFindExistingNormal)
	if err != nil {
		return err
	}
	if s.adminOnly && !conv.ReaderInfo.UntrustedTeamRole.IsAdminOrAbove() {
		return NewDevStoragePermissionDeniedError(conv.ReaderInfo.UntrustedTeamRole)
	}
	if _, _, err = NewBlockingSender(s.G(), NewBoxer(s.G()), s.ri).Send(ctx, conv.GetConvID(),
		chat1.MessagePlaintext{
			ClientHeader: chat1.MessageClientHeader{
				Conv:        conv.Info.Triple,
				TlfName:     tlfname,
				MessageType: chat1.MessageType_TEXT,
			},
			MessageBody: chat1.NewMessageBodyWithText(chat1.MessageText{
				Body: string(dat),
			}),
		}, 0, nil, nil, nil); err != nil {
		return err
	}
	if s.adminOnly && (conv.ConvSettings == nil || conv.ConvSettings.MinWriterRoleInfo == nil || conv.ConvSettings.MinWriterRoleInfo.Role != keybase1.TeamRole_ADMIN) {
		_, err := s.ri().SetConvMinWriterRole(ctx, chat1.SetConvMinWriterRoleArg{ConvID: conv.Info.Id, Role: keybase1.TeamRole_ADMIN})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *DevConversationBackedStorage) Get(ctx context.Context, uid gregor1.UID, tlfid chat1.TLFID, name string,
	dest interface{}) (found bool, err error) {
	defer s.Trace(ctx, func() error { return err }, "Get(%s)", name)()

	info, err := CreateNameInfoSource(ctx, s.G(), s.mt).LookupName(ctx, tlfid, false /* public */, "")
	if err != nil {
		return false, err
	}
	tlfname := info.CanonicalName

	convs, err := FindConversations(ctx, s.G(), s.DebugLabeler, types.InboxSourceDataSourceAll, s.ri, uid,
		tlfname, chat1.TopicType_DEV, s.mt,
		keybase1.TLFVisibility_PRIVATE, name, nil)
	if err != nil {
		return false, err
	}
	if len(convs) == 0 {
		return false, nil
	}
	conv := convs[0]
	tv, err := s.G().ConvSource.Pull(ctx, conv.GetConvID(), uid, chat1.GetThreadReason_GENERAL,
		&chat1.GetThreadQuery{
			MessageTypes: []chat1.MessageType{chat1.MessageType_TEXT},
		}, &chat1.Pagination{Num: 1})
	if err != nil {
		return false, err
	}
	if len(tv.Messages) == 0 {
		return false, nil
	}
	msg := tv.Messages[0]
	if !msg.IsValid() {
		return false, nil
	}
	body := msg.Valid().MessageBody
	if !body.IsType(chat1.MessageType_TEXT) {
		return false, nil
	}
	if s.adminOnly {
		if conv.ConvSettings == nil || conv.ConvSettings.MinWriterRoleInfo == nil {
			return false, NewDevStorageAdminOnlyError("no conversation settings")
		}
		if conv.ConvSettings.MinWriterRoleInfo.Role != keybase1.TeamRole_ADMIN {
			return false, NewDevStorageAdminOnlyError("minWriterRole was not admin")
		}
	}
	if err = json.Unmarshal([]byte(body.Text().Body), dest); err != nil {
		return false, err
	}
	return true, nil
}
