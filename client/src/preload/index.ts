import { contextBridge, ipcRenderer } from 'electron'
import { electronAPI } from '@electron-toolkit/preload'
import IRatchetDetail from '../renderer/src/interfaces/IRatchetDetail'

// Custom APIs for renderer
const api = {
  readAuthFile: () => ipcRenderer.invoke('r_readAuthFile'),
  createAuthFile: () => ipcRenderer.invoke('r_createAuthFile'),
  writeAuthFile: (content: string) => ipcRenderer.invoke('r_writeAuthFile', content),
  checkAuthFile: (content: string) => ipcRenderer.invoke('r_checkAuthFile', content),
  existAuthFile: () => ipcRenderer.invoke('r_existAuthFile'),
  getInternalKey: () => ipcRenderer.invoke('r_getInternalKey'),

  updateAvatar: (avatar: string, filename: string) =>
    ipcRenderer.invoke('r_updateAvatar', avatar, filename),
  getAvatar: () => ipcRenderer.invoke('r_getAvatar'),

  random32Bytes: () => ipcRenderer.invoke('r_random32Bytes'),
  encryptblob: (data: ArrayBuffer, randomKey: string, mimeType: string) =>
    ipcRenderer.invoke('r_encryptblob', data, randomKey, mimeType),
  decryptblob: (data: ArrayBuffer, randomKey: string, mimeType: string) => {
    return ipcRenderer.invoke('r_decryptblob', data, randomKey, mimeType)
  },

  createRatchetFile: (username: string, ratchetDetail: IRatchetDetail, ratchetId) =>
    ipcRenderer.invoke('r_createRatchetFile', username, ratchetDetail, ratchetId),
  changeRatchetDetail: (receiver: string, ratchetDetail: IRatchetDetail) =>
    ipcRenderer.invoke('r_changeRatchetDetail', receiver, ratchetDetail),
  getRatchetId: (username: string) => ipcRenderer.invoke('r_getRatchetId', username),
  getRatchetDetailList: (currentUsername: string) =>
    ipcRenderer.invoke('r_getRatchetDetailList', currentUsername),
  getOldChatSessions: (currentUsername: string) =>
    ipcRenderer.invoke('r_getOldChatSessions', currentUsername),
  addMessageToRatchet: (receiver: string, messages: string) =>
    ipcRenderer.invoke('r_addMessageToRatchet', receiver, messages),
  getMessagesByUsername: (username: string) =>
    ipcRenderer.invoke('r_getMessagesByUsername', username),
  deleteMessage: (receiver: string, index: number) =>
    ipcRenderer.invoke('r_deleteMessage', receiver, index),
  changeConversationReaded: (username: string, isReaded: boolean) =>
    ipcRenderer.invoke('r_changeConversationReaded', username, isReaded)
}

// Use `contextBridge` APIs to expose Electron APIs to
// renderer only if context isolation is enabled, otherwise
// just add to the DOM global.
if (process.contextIsolated) {
  try {
    contextBridge.exposeInMainWorld('electron', electronAPI)
    contextBridge.exposeInMainWorld('api', api)
  } catch (error) {
    console.error(error)
  }
} else {
  // @ts-ignore (define in dts)
  window.electron = electronAPI
  // @ts-ignore (define in dts)
  window.api = api
}
