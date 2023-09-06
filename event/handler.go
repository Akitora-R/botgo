package event

import "github.com/tencent-connect/botgo/dto"

type Handler interface {
	GetIntent() dto.Intent
	ParseAndHandle(payload *dto.WSPayload) error
	GetReadyHandler() ReadyHandler
	GetErrorNotifyHandler() ErrorNotifyHandler
}

type HandlerImpl struct {
	Ready       ReadyHandler
	ErrorNotify ErrorNotifyHandler
	Plain       PlainEventHandler

	Guild       GuildEventHandler
	GuildMember GuildMemberEventHandler
	Channel     ChannelEventHandler

	Message             MessageEventHandler
	MessageReaction     MessageReactionEventHandler
	ATMessage           ATMessageEventHandler
	DirectMessage       DirectMessageEventHandler
	MessageAudit        MessageAuditEventHandler
	MessageDelete       MessageDeleteEventHandler
	PublicMessageDelete PublicMessageDeleteEventHandler
	DirectMessageDelete DirectMessageDeleteEventHandler

	Audio AudioEventHandler

	Thread     ThreadEventHandler
	Post       PostEventHandler
	Reply      ReplyEventHandler
	ForumAudit ForumAuditEventHandler

	Interaction InteractionEventHandler

	Intent dto.Intent
}

func NewHandler(handlers ...any) *HandlerImpl {
	var hi = HandlerImpl{}
	var i dto.Intent
	for _, handler := range handlers {
		switch handle := handler.(type) {
		case ReadyHandler:
			hi.Ready = handle
		case ErrorNotifyHandler:
			hi.ErrorNotify = handle
		case PlainEventHandler:
			hi.Plain = handle
		case AudioEventHandler:
			hi.Audio = handle
			i = i | dto.EventToIntent(
				dto.EventAudioStart, dto.EventAudioFinish,
				dto.EventAudioOnMic, dto.EventAudioOffMic,
			)
		case InteractionEventHandler:
			hi.Interaction = handle
			i = i | dto.EventToIntent(dto.EventInteractionCreate)
		default:
		}
	}
	i = i | hi.registerRelationHandlers(i, handlers...)
	i = i | hi.registerMessageHandlers(i, handlers...)
	i = i | hi.registerForumHandlers(i, handlers...)
	hi.Intent = i
	return &hi
}

func (hi *HandlerImpl) GetIntent() dto.Intent {
	return hi.Intent
}

func (hi *HandlerImpl) GetReadyHandler() ReadyHandler {
	return hi.Ready
}

func (hi *HandlerImpl) GetErrorNotifyHandler() ErrorNotifyHandler {
	return hi.ErrorNotify
}

func (hi *HandlerImpl) registerForumHandlers(i dto.Intent, handlers ...interface{}) dto.Intent {
	for _, h := range handlers {
		switch handle := h.(type) {
		case ThreadEventHandler:
			hi.Thread = handle
			i = i | dto.EventToIntent(
				dto.EventForumThreadCreate, dto.EventForumThreadUpdate, dto.EventForumThreadDelete,
			)
		case PostEventHandler:
			hi.Post = handle
			i = i | dto.EventToIntent(dto.EventForumPostCreate, dto.EventForumPostDelete)
		case ReplyEventHandler:
			hi.Reply = handle
			i = i | dto.EventToIntent(dto.EventForumReplyCreate, dto.EventForumReplyDelete)
		case ForumAuditEventHandler:
			hi.ForumAudit = handle
			i = i | dto.EventToIntent(dto.EventForumAuditResult)
		default:
		}
	}
	return i
}

// registerRelationHandlers 注册频道关系链相关handlers
func (hi *HandlerImpl) registerRelationHandlers(i dto.Intent, handlers ...interface{}) dto.Intent {
	for _, h := range handlers {
		switch handle := h.(type) {
		case GuildEventHandler:
			hi.Guild = handle
			i = i | dto.EventToIntent(dto.EventGuildCreate, dto.EventGuildDelete, dto.EventGuildUpdate)
		case GuildMemberEventHandler:
			hi.GuildMember = handle
			i = i | dto.EventToIntent(dto.EventGuildMemberAdd, dto.EventGuildMemberRemove, dto.EventGuildMemberUpdate)
		case ChannelEventHandler:
			hi.Channel = handle
			i = i | dto.EventToIntent(dto.EventChannelCreate, dto.EventChannelDelete, dto.EventChannelUpdate)
		default:
		}
	}
	return i
}

// registerMessageHandlers 注册消息相关的 handler
func (hi *HandlerImpl) registerMessageHandlers(i dto.Intent, handlers ...interface{}) dto.Intent {
	for _, h := range handlers {
		switch handle := h.(type) {
		case MessageEventHandler:
			hi.Message = handle
			i = i | dto.EventToIntent(dto.EventMessageCreate)
		case ATMessageEventHandler:
			hi.ATMessage = handle
			i = i | dto.EventToIntent(dto.EventAtMessageCreate)
		case DirectMessageEventHandler:
			hi.DirectMessage = handle
			i = i | dto.EventToIntent(dto.EventDirectMessageCreate)
		case MessageDeleteEventHandler:
			hi.MessageDelete = handle
			i = i | dto.EventToIntent(dto.EventMessageDelete)
		case PublicMessageDeleteEventHandler:
			hi.PublicMessageDelete = handle
			i = i | dto.EventToIntent(dto.EventPublicMessageDelete)
		case DirectMessageDeleteEventHandler:
			hi.DirectMessageDelete = handle
			i = i | dto.EventToIntent(dto.EventDirectMessageDelete)
		case MessageReactionEventHandler:
			hi.MessageReaction = handle
			i = i | dto.EventToIntent(dto.EventMessageReactionAdd, dto.EventMessageReactionRemove)
		case MessageAuditEventHandler:
			hi.MessageAudit = handle
			i = i | dto.EventToIntent(dto.EventMessageAuditPass, dto.EventMessageAuditReject)
		default:
		}
	}
	return i
}

func (hi *HandlerImpl) ParseAndHandle(payload *dto.WSPayload) error {
	if payload.OPCode == dto.WSDispatchEvent {
		switch payload.Type {
		case dto.EventGuildCreate:
			return hi.guildHandler(payload, payload.RawMessage)
		case dto.EventGuildUpdate:
			return hi.guildHandler(payload, payload.RawMessage)
		case dto.EventGuildDelete:
			return hi.guildHandler(payload, payload.RawMessage)
		case dto.EventChannelCreate:
			return hi.channelHandler(payload, payload.RawMessage)
		case dto.EventChannelUpdate:
			return hi.channelHandler(payload, payload.RawMessage)
		case dto.EventChannelDelete:
			return hi.channelHandler(payload, payload.RawMessage)
		case dto.EventGuildMemberAdd:
			return hi.guildMemberHandler(payload, payload.RawMessage)
		case dto.EventGuildMemberUpdate:
			return hi.guildMemberHandler(payload, payload.RawMessage)
		case dto.EventGuildMemberRemove:
			return hi.guildMemberHandler(payload, payload.RawMessage)
		case dto.EventMessageCreate:
			return hi.messageHandler(payload, payload.RawMessage)
		case dto.EventMessageDelete:
			return hi.messageDeleteHandler(payload, payload.RawMessage)
		case dto.EventMessageReactionAdd:
			return hi.messageReactionHandler(payload, payload.RawMessage)
		case dto.EventMessageReactionRemove:
			return hi.messageReactionHandler(payload, payload.RawMessage)
		case dto.EventAtMessageCreate:
			return hi.atMessageHandler(payload, payload.RawMessage)
		case dto.EventPublicMessageDelete:
			return hi.publicMessageDeleteHandler(payload, payload.RawMessage)
		case dto.EventDirectMessageCreate:
			return hi.directMessageHandler(payload, payload.RawMessage)
		case dto.EventDirectMessageDelete:
			return hi.directMessageDeleteHandler(payload, payload.RawMessage)
		case dto.EventAudioStart:
			return hi.audioHandler(payload, payload.RawMessage)
		case dto.EventAudioFinish:
			return hi.audioHandler(payload, payload.RawMessage)
		case dto.EventAudioOnMic:
			return hi.audioHandler(payload, payload.RawMessage)
		case dto.EventAudioOffMic:
			return hi.audioHandler(payload, payload.RawMessage)
		case dto.EventMessageAuditPass:
			return hi.messageAuditHandler(payload, payload.RawMessage)
		case dto.EventMessageAuditReject:
			return hi.messageAuditHandler(payload, payload.RawMessage)
		case dto.EventForumThreadCreate:
			return hi.threadHandler(payload, payload.RawMessage)
		case dto.EventForumThreadUpdate:
			return hi.threadHandler(payload, payload.RawMessage)
		case dto.EventForumThreadDelete:
			return hi.threadHandler(payload, payload.RawMessage)
		case dto.EventForumPostCreate:
			return hi.postHandler(payload, payload.RawMessage)
		case dto.EventForumPostDelete:
			return hi.postHandler(payload, payload.RawMessage)
		case dto.EventForumReplyCreate:
			return hi.replyHandler(payload, payload.RawMessage)
		case dto.EventForumReplyDelete:
			return hi.replyHandler(payload, payload.RawMessage)
		case dto.EventForumAuditResult:
			return hi.forumAuditHandler(payload, payload.RawMessage)
		case dto.EventInteractionCreate:
			return hi.interactionHandler(payload, payload.RawMessage)
		}
	} else {
		if hi.Plain != nil {
			return hi.Plain(payload, payload.RawMessage)
		}
	}
	return nil
}

func (hi *HandlerImpl) guildHandler(payload *dto.WSPayload, message []byte) error {
	data := &dto.WSGuildData{}
	if err := ParseData(message, data); err != nil {
		return err
	}
	if hi.Guild != nil {
		return hi.Guild(payload, data)
	}
	return nil
}

func (hi *HandlerImpl) channelHandler(payload *dto.WSPayload, message []byte) error {
	data := &dto.WSChannelData{}
	if err := ParseData(message, data); err != nil {
		return err
	}
	if hi.Channel != nil {
		return hi.Channel(payload, data)
	}
	return nil
}

func (hi *HandlerImpl) guildMemberHandler(payload *dto.WSPayload, message []byte) error {
	data := &dto.WSGuildMemberData{}
	if err := ParseData(message, data); err != nil {
		return err
	}
	if hi.GuildMember != nil {
		return hi.GuildMember(payload, data)
	}
	return nil
}

func (hi *HandlerImpl) messageHandler(payload *dto.WSPayload, message []byte) error {
	data := &dto.WSMessageData{}
	if err := ParseData(message, data); err != nil {
		return err
	}
	if hi.Message != nil {
		return hi.Message(payload, data)
	}
	return nil
}

func (hi *HandlerImpl) messageDeleteHandler(payload *dto.WSPayload, message []byte) error {
	data := &dto.WSMessageDeleteData{}
	if err := ParseData(message, data); err != nil {
		return err
	}
	if hi.MessageDelete != nil {
		return hi.MessageDelete(payload, data)
	}
	return nil
}

func (hi *HandlerImpl) messageReactionHandler(payload *dto.WSPayload, message []byte) error {
	data := &dto.WSMessageReactionData{}
	if err := ParseData(message, data); err != nil {
		return err
	}
	if hi.MessageReaction != nil {
		return hi.MessageReaction(payload, data)
	}
	return nil
}

func (hi *HandlerImpl) atMessageHandler(payload *dto.WSPayload, message []byte) error {
	data := &dto.WSATMessageData{}
	if err := ParseData(message, data); err != nil {
		return err
	}
	if hi.ATMessage != nil {
		return hi.ATMessage(payload, data)
	}
	return nil
}

func (hi *HandlerImpl) publicMessageDeleteHandler(payload *dto.WSPayload, message []byte) error {
	data := &dto.WSPublicMessageDeleteData{}
	if err := ParseData(message, data); err != nil {
		return err
	}
	if hi.PublicMessageDelete != nil {
		return hi.PublicMessageDelete(payload, data)
	}
	return nil
}

func (hi *HandlerImpl) directMessageHandler(payload *dto.WSPayload, message []byte) error {
	data := &dto.WSDirectMessageData{}
	if err := ParseData(message, data); err != nil {
		return err
	}
	if hi.DirectMessage != nil {
		return hi.DirectMessage(payload, data)
	}
	return nil
}

func (hi *HandlerImpl) directMessageDeleteHandler(payload *dto.WSPayload, message []byte) error {
	data := &dto.WSDirectMessageDeleteData{}
	if err := ParseData(message, data); err != nil {
		return err
	}
	if hi.DirectMessageDelete != nil {
		return hi.DirectMessageDelete(payload, data)
	}
	return nil
}

func (hi *HandlerImpl) audioHandler(payload *dto.WSPayload, message []byte) error {
	data := &dto.WSAudioData{}
	if err := ParseData(message, data); err != nil {
		return err
	}
	if hi.Audio != nil {
		return hi.Audio(payload, data)
	}
	return nil
}

func (hi *HandlerImpl) threadHandler(payload *dto.WSPayload, message []byte) error {
	data := &dto.WSThreadData{}
	if err := ParseData(message, data); err != nil {
		return err
	}
	if hi.Thread != nil {
		return hi.Thread(payload, data)
	}
	return nil
}

func (hi *HandlerImpl) postHandler(payload *dto.WSPayload, message []byte) error {
	data := &dto.WSPostData{}
	if err := ParseData(message, data); err != nil {
		return err
	}
	if hi.Post != nil {
		return hi.Post(payload, data)
	}
	return nil
}

func (hi *HandlerImpl) replyHandler(payload *dto.WSPayload, message []byte) error {
	data := &dto.WSReplyData{}
	if err := ParseData(message, data); err != nil {
		return err
	}
	if hi.Reply != nil {
		return hi.Reply(payload, data)
	}
	return nil
}

func (hi *HandlerImpl) forumAuditHandler(payload *dto.WSPayload, message []byte) error {
	data := &dto.WSForumAuditData{}
	if err := ParseData(message, data); err != nil {
		return err
	}
	if hi.ForumAudit != nil {
		return hi.ForumAudit(payload, data)
	}
	return nil
}

func (hi *HandlerImpl) messageAuditHandler(payload *dto.WSPayload, message []byte) error {
	data := &dto.WSMessageAuditData{}
	if err := ParseData(message, data); err != nil {
		return err
	}
	if hi.MessageAudit != nil {
		return hi.MessageAudit(payload, data)
	}
	return nil
}

func (hi *HandlerImpl) interactionHandler(payload *dto.WSPayload, message []byte) error {
	data := &dto.WSInteractionData{}
	if err := ParseData(message, data); err != nil {
		return err
	}
	if hi.Interaction != nil {
		return hi.Interaction(payload, data)
	}
	return nil
}
