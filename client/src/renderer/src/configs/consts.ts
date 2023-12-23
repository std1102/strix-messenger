export const API_URL = 'http://192.168.1.77:7777'
export const IMAGE_URL = API_URL + '/api/v1/file/get?fileId='
export const SOCKET_URL = 'ws://192.168.1.77:7777'
export const WS_CALL_URL =
  'ws://192.168.1.77:7777/voip?voipSession={{voipToken}}&connType={{connType}}'

export const ACCESS_TOKEN_KEY = 'access_token'

export const CHAT_PREFIX = 'chat_'
export const CHAT_FILE = 'chat.js on'
export const RATCHET_FILE = 'ratchet.json'
export const AUTH_FILE = 'auth.json'
// export const AUTH_FILE = 'auth2.json'
export const AVATAR_FILE = 'avatar'

export const TEXT_TYPE = 'CHAT_TEXT' as const
export const IMAGE_TYPE = 'CHAT_IMAGE' as const
export const VIDEO_TYPE = 'CHAT_VIDEO' as const
export const FILE_TYPE = 'CHAT_FILE' as const

export const MESSAGE_EVENT = 'CHAT_TEXT'
export const CHAT_NEW_EVENT = 'CHAT_NEW'
export const CHAT_CLOSE = 'CHAT_CLOSE'
export const CHAT_VIDEO_EVENT = 'CHAT_VIDEO'
export const CALL_VIDEO_EVENT = 'CALL_VIDEO'
export const CHAT_AUDIO_EVENT = 'CHAT_AUDIO'
export const ACCEPT_CALL_EVENT = 'CHAT_ACCEPT'

export const AVATAR_DEFAULT = 'https://source.unsplash.com/RZrIJ8C0860'

export const VIDEO_MIME_TYPE = 'video/webm; codecs="opus,vp8"'
export const AUDIO_MIME_TYPE = 'audio/webm;codecs=opus'
