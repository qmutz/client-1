package chat

import (
	"testing"

	"github.com/keybase/client/go/protocol/chat1"
	"github.com/stretchr/testify/require"
)

func TestDevConversationBackedStorage(t *testing.T) {
	useRemoteMock = false
	defer func() { useRemoteMock = true }()
	ctc := makeChatTestContext(t, "TestDevConversationBackedStorage", 1)
	defer ctc.cleanup()

	users := ctc.users()
	tc := ctc.world.Tcs[users[0].Username]
	ri := ctc.as(t, users[0]).ri
	ctx := ctc.as(t, users[0]).startCtx
	uid := users[0].User.GetUID().ToBytes()

	key0 := "storage0"
	key1 := "storage1"
	storage := NewDevConversationBackedStorage(tc.Context(), chat1.ConversationMembersType_IMPTEAMNATIVE, false /* adminOnly */, func() chat1.RemoteInterface { return ri })
	settings := chat1.UnfurlSettings{
		Mode:      chat1.UnfurlMode_WHITELISTED,
		Whitelist: make(map[string]bool),
	}
	require.NoError(t, storage.Put(ctx, uid, uid, key0, settings))
	var settingsRes chat1.UnfurlSettings
	found, err := storage.Get(ctx, uid, uid, key0, &settingsRes)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, chat1.UnfurlMode_WHITELISTED, settingsRes.Mode)
	require.Zero(t, len(settingsRes.Whitelist))

	settings.Mode = chat1.UnfurlMode_NEVER
	require.NoError(t, storage.Put(ctx, uid, uid, key0, settings))
	found, err = storage.Get(ctx, uid, uid, key0, &settingsRes)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, chat1.UnfurlMode_NEVER, settingsRes.Mode)
	require.Zero(t, len(settingsRes.Whitelist))

	settings.Mode = chat1.UnfurlMode_WHITELISTED
	settings.Whitelist["MIKE"] = true
	require.NoError(t, storage.Put(ctx, uid, uid, key1, settings))
	found, err = storage.Get(ctx, uid, uid, key0, &settingsRes)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, chat1.UnfurlMode_NEVER, settingsRes.Mode)
	require.Zero(t, len(settingsRes.Whitelist))
	found, err = storage.Get(ctx, uid, uid, key1, &settingsRes)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, chat1.UnfurlMode_WHITELISTED, settingsRes.Mode)
	require.Equal(t, 1, len(settingsRes.Whitelist))
	require.True(t, settingsRes.Whitelist["MIKE"])

	found, err = storage.Get(ctx, uid, uid, "AHHHHH CANT FIND ME", &settingsRes)
	require.NoError(t, err)
	require.False(t, found)
}

func TestDevConversationBackedStorageTeamAdminOnly(t *testing.T) {
	useRemoteMock = false
	defer func() { useRemoteMock = true }()
	ctc := makeChatTestContext(t, t.Name(), 2)
	defer ctc.cleanup()

	users := ctc.users()
	tc := ctc.world.Tcs[users[0].Username]
	ri := ctc.as(t, users[0]).ri
	ctx := ctc.as(t, users[0]).startCtx
	uid := users[0].User.GetUID().ToBytes()
	readeruid := users[1].User.GetUID().ToBytes()

	conv := mustCreateChannelForTest(t, ctc, users[0], chat1.TopicType_CHAT, nil, chat1.ConversationMembersType_TEAM, users[1:]...)
	tlfid := conv.Triple.Tlfid

	key0 := "mykey"
	storage := NewDevConversationBackedStorage(tc.Context(), chat1.ConversationMembersType_TEAM, true /* adminOnly */, func() chat1.RemoteInterface { return ri })

	var msg string
	found, err := storage.Get(ctx, uid, tlfid, key0, &msg)
	require.NoError(t, err)
	require.False(t, found)

	err = storage.Put(ctx, uid, tlfid, key0, "hello")
	require.NoError(t, err)

	found, err = storage.Get(ctx, readeruid, tlfid, key0, &msg)
	require.NoError(t, err)
	require.True(t, found)
	require.Equal(t, msg, "hello")
}

func TestDevConversationBackedStorageTeamAdminOnlyReaderMisbehavior(t *testing.T) {
	useRemoteMock = false
	defer func() { useRemoteMock = true }()
	ctc := makeChatTestContext(t, t.Name(), 2)
	defer ctc.cleanup()

	users := ctc.users()
	tc := ctc.world.Tcs[users[0].Username]
	ri := ctc.as(t, users[0]).ri
	ctx := ctc.as(t, users[0]).startCtx
	uid := users[0].User.GetUID().ToBytes()

	readertc := ctc.world.Tcs[users[1].Username]
	readerri := ctc.as(t, users[1]).ri
	readerctx := ctc.as(t, users[1]).startCtx
	readeruid := users[1].User.GetUID().ToBytes()

	conv := mustCreateChannelForTest(t, ctc, users[0], chat1.TopicType_CHAT, nil, chat1.ConversationMembersType_TEAM, users[1:]...)
	tlfid := conv.Triple.Tlfid

	key0 := "mykey"
	storage := NewDevConversationBackedStorage(tc.Context(), chat1.ConversationMembersType_TEAM, true /* adminOnly */, func() chat1.RemoteInterface { return ri })
	readerstorage := NewDevConversationBackedStorage(readertc.Context(), chat1.ConversationMembersType_TEAM, true /* adminOnly */, func() chat1.RemoteInterface { return readerri })
	var msg string
	var readermsg string
	found, err := storage.Get(ctx, uid, tlfid, key0, &msg)
	require.NoError(t, err)
	require.False(t, found)
	found, err = readerstorage.Get(readerctx, readeruid, tlfid, key0, &readermsg)
	require.NoError(t, err)
	require.False(t, found)

	err = readerstorage.Put(readerctx, readeruid, tlfid, key0, "evil")
	require.Error(t, err)
	require.IsType(t, &DevStoragePermissionDeniedError{}, err, "got right error")

	// hack around side-protection and make channel anyway
	devconv := mustCreateChannelForTest(t, ctc, users[1], chat1.TopicType_DEV, &key0, chat1.ConversationMembersType_TEAM, users[0])
	larg := chat1.PostLocalArg{
		ConversationID: devconv.Id,
		Msg: chat1.MessagePlaintext{
			ClientHeader: chat1.MessageClientHeader{
				Conv:        devconv.Triple,
				TlfName:     devconv.TlfName,
				MessageType: chat1.MessageType_TEXT,
			},
			MessageBody: chat1.NewMessageBodyWithText(chat1.MessageText{Body: "reallyevil"}),
		},
	}
	_, err = ctc.as(t, users[1]).chatLocalHandler().PostLocal(ctc.as(t, users[1]).startCtx, larg)
	require.NoError(t, err)

	found, err = storage.Get(ctx, uid, tlfid, key0, &msg)
	require.Error(t, err, "got an error after misbehavior")
	require.IsType(t, &DevStorageAdminOnlyError{}, err, "got a permission error")

	found, err = readerstorage.Get(readerctx, readeruid, tlfid, key0, &readermsg)
	require.Error(t, err, "got an error after misbehavior")
	require.IsType(t, &DevStorageAdminOnlyError{}, err, "got a permission error")
}
