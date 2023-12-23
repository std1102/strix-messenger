import { ElectronAPI } from '@electron-toolkit/preload'
import IRatchetDetail from "../renderer/src/interfaces/IRatchetDetail";
import IMessage from "../renderer/src/interfaces/IMessage";
import { decryptblob } from "../renderer/src/crypto/cryptoLib";

declare global {
  interface Window {
    Go: any
    startUp: (pinValue: string) => Promise<void>
    generateInternalKeyBundle: () => Promise<string>
    populateExternalKeyBundle: () => Promise<{keyId: string, keyBundle: string}>
    regeneratePreKey: () => Promise<{keyId: string, keyBundle: string}>
    saveInternalKey: () => Promise<any>
    populateExternalKeyBundle: () => Promise<void>
    initRatchetFromInternal: (keyBundle: string) => Promise<any>
    initRatchetFromExternal: (externalKey: string,ephemeralKey: string, ratchetId: string ) => Promise<{ratchetId: string}>
    loadRatchet: (ratchetDetail: string) => Promise<string>
    saveRatchet: (ratchetId: string) => Promise<IRatchetDetail>
    loadInternalKey: (internalKey: string) => Promise<string>
    sendMessage: (ratchetId: string, isBinary: boolean, message: string) => Promise<any>
    isRatchetExist: (ratchetId: string) => Promise<boolean>
    receiveMessage: (data: string) => Promise<string>
    electron: ElectronAPI
    api: {
      readAuthFile(): Promise<string>
      createAuthFile(): Promise<void>
      writeAuthFile(content: string): Promise<void>
      checkAuthFile(content: string): Promise<boolean>
      existAuthFile(): Promise<boolean>
      getInternalKey(): Promise<string>

      random32Bytes(): Promise<string>
      encryptblob(data: ArrayBuffer, randomKey: string, mimeType: string): Promise<string>
      decryptblob(data: ArrayBuffer, randomKey: string, mimeType: string): Promise<string>

      updateAvatar(avatar: string, filename): Promise<void>
      getAvatar(): Promise<string>

      createRatchetFile(username: string, ratchetDetail: IRatchetDetail, ratchetId: string): Promise<void>
      changeRatchetDetail(receiver: string, ratchetDetail: IRatchetDetail): Promise<void>
      writeRatchetFile(username: string, ratchetData: any): Promise<void>
      changeConversationReaded(username: string, isReaded: boolean): Promise<void>
      getRatchetId(username: string): Promise<string>
      getRatchetDetailList(currentUsername: string): Promise<IRatchetDetail[]>
      getOldChatSessions(currentUsername: string): Promise<{receiver:string, ratchetId: string, lastMessage: string, updatedAt: string, isReaded: boolean, messages: IMessage[]}[]>
      addMessageToRatchet(receiver: string, messages: IMessage[]): Promise<void>
      getMessagesByUsername(username: string): Promise<IMessage[]>
      deleteMessage(receiver: string, index: number): Promise<void>
    }
  }
}
