import { app, BrowserWindow, ipcMain, session, shell } from 'electron'
import { join } from 'path'
import { electronApp, is, optimizer } from '@electron-toolkit/utils'
import icon from '../../resources/icon.png?asset'
import * as fs from 'fs'
import { AUTH_FILE, AVATAR_FILE, CHAT_PREFIX } from '../renderer/src/configs/consts'
import IRatchetDetail from '../renderer/src/interfaces/IRatchetDetail'
import IMessage from '../renderer/src/interfaces/IMessage'
import { randomBytes } from 'crypto'
import { decryptblob, encryptblob } from '../renderer/src/crypto/cryptoLib'

function createWindow(): void {
  // Create the browser window.
  const mainWindow = new BrowserWindow({
    width: 1200,
    height: 920,
    show: false,
    autoHideMenuBar: true,
    ...(process.platform === 'linux' ? { icon } : {}),
    webPreferences: {
      contextIsolation: true,
      preload: join(__dirname, '../preload/index.js'),
      sandbox: false,
      enableBlinkFeatures: 'OverlayScrollbars',
      nodeIntegration: true,
      webSecurity: false
    }
  })
  mainWindow.webContents.openDevTools({
    mode: 'detach'
  })

  mainWindow.on('ready-to-show', () => {
    session.defaultSession.webRequest.onBeforeSendHeaders((details, callback) => {
      details.requestHeaders['User-Agent'] = 'DesktopClient'
      callback({ cancel: false, requestHeaders: details.requestHeaders })
    })
    mainWindow.show()
  })

  mainWindow.webContents.setWindowOpenHandler((details) => {
    shell.openExternal(details.url)
    return { action: 'deny' }
  })

  // HMR for renderer base on electron-vite cli.
  // Load the remote URL for development or the local html file for production.
  if (is.dev && process.env['ELECTRON_RENDERER_URL']) {
    mainWindow.loadURL(process.env['ELECTRON_RENDERER_URL'])
  } else {
    mainWindow.loadFile(join(__dirname, '../renderer/index.html'))
  }
}

// This method will be called when Electron has finished
// initialization and is ready to create browser windows.
// Some APIs can only be used after this event occurs.
app.whenReady().then(() => {
  // Set app user model id for windows
  electronApp.setAppUserModelId('com.electron')

  ipcMain.handle('r_readAuthFile', handleReadAuthFile)
  ipcMain.handle('r_writeAuthFile', handleWriteAuthFile)
  ipcMain.handle('r_checkAuthFile', handleCheckAuthFile)
  ipcMain.handle('r_getInternalKey', handleGetInternalKey)
  ipcMain.handle('r_existAuthFile', handleExistAuthFile)
  ipcMain.handle('r_createAuthFile', handleCreateAuthFile)

  ipcMain.handle('r_updateAvatar', handleUpdateAvatar)
  ipcMain.handle('r_getAvatar', handleGetAvatar)

  ipcMain.handle('r_random32Bytes', handleGetRandom32Bytes)
  ipcMain.handle('r_encryptblob', handleEncryptblob)
  ipcMain.handle('r_decryptblob', handleDecryptblob)

  ipcMain.handle('r_createRatchetFile', handleCreateRatchetFile)
  ipcMain.handle('r_changeRatchetDetail', handleChangeRatchetDetail)
  ipcMain.handle('r_getRatchetId', handleGetRatchetId)
  ipcMain.handle('r_getRatchetDetailList', handleGetRatchetDetailList)
  ipcMain.handle('r_getOldChatSessions', handleGetOldChatSessions)
  ipcMain.handle('r_addMessageToRatchet', handleAddMessageToRatchet)
  ipcMain.handle('r_getMessagesByUsername', handleGetMessagesByUsername)
  ipcMain.handle('r_deleteMessage', handleDeleteMessage)
  ipcMain.handle('r_changeConversationReaded', handleChangeConversationReaded)

  // Default open or close DevTools by F12 in development
  // and ignore CommandOrControl + R in production.
  // see https://github.com/alex8088/electron-toolkit/tree/master/packages/utils
  app.on('browser-window-created', (_, window) => {
    optimizer.watchWindowShortcuts(window)
  })

  createWindow()

  app.on('activate', function () {
    // On macOS it's common to re-create a window in the app when the
    // dock icon is clicked and there are no other windows open.
    if (BrowserWindow.getAllWindows().length === 0) createWindow()
  })
})

// Quit when all windows are closed, except on macOS. There, it's common
// for applications and their menu bar to stay active until the user quits
// explicitly with Cmd + Q.
app.on('window-all-closed', () => {
  if (process.platform !== 'darwin') {
    app.quit()
  }
})

// In this file you can include the rest of your app"s specific main process
// code. You can also put them in separate files and require them here.

function handleReadAuthFile() {
  try {
    if (fs.existsSync(AUTH_FILE)) {
      return fs
        .readFileSync(AUTH_FILE, {
          encoding: 'utf8'
        })
        .toString()
    }
  } catch (error) {
    console.error('ERROR', error)
  }

  return ''
}

function handleWriteAuthFile(_e: Electron.IpcMainInvokeEvent, content: string) {
  try {
    fs.writeFileSync(AUTH_FILE, JSON.stringify(content), {
      encoding: 'utf8'
    })
  } catch (error) {
    console.error('ERROR', error)
  }
}

function handleCheckAuthFile(_e: Electron.IpcMainInvokeEvent, newContent: string) {
  try {
    if (
      !fs.existsSync(AUTH_FILE) ||
      fs.readFileSync(AUTH_FILE, {
        encoding: 'utf8'
      }) === ''
    ) {
      handleWriteAuthFile(null as never, newContent)
      return true
    }

    const authContent = fs.readFileSync(AUTH_FILE, {
      encoding: 'utf8'
    })
    return authContent === newContent
  } catch (error) {
    console.error('ERROR', error)
  }
  return false
}

function handleGetInternalKey() {
  try {
    if (fs.existsSync(AUTH_FILE)) {
      return fs.readFileSync(AUTH_FILE, {
        encoding: 'utf8'
      })
    }
  } catch (error) {
    console.error('ERROR', error)
  }

  return ''
}

function handleExistAuthFile() {
  try {
    return fs.existsSync(AUTH_FILE)
  } catch (error) {
    console.error('ERROR', error)
  }

  return false
}

function handleCreateAuthFile() {
  try {
    fs.writeFileSync(AUTH_FILE, '')
  } catch (error) {
    console.error('ERROR', error)
  }
}

function handleCreateRatchetFile(
  _e: Electron.IpcMainInvokeEvent,
  username: string,
  ratchetDetail: IRatchetDetail,
  ratchetId: string
) {
  try {
    const filename = CHAT_PREFIX + username + '.json'

    fs.writeFileSync(filename, JSON.stringify({ ratchetDetail, ratchetId }), {
      encoding: 'utf8'
    })
  } catch (error) {
    console.error('ERROR', error)
  }
}

function handleChangeRatchetDetail(
  _e: Electron.IpcMainInvokeEvent,
  username: string,
  ratchetDetail: IRatchetDetail
) {
  try {
    const filename = CHAT_PREFIX + username + '.json'

    const fileContent = fs.readFileSync(filename, {
      encoding: 'utf8'
    })
    const parsedContent = JSON.parse(fileContent)
    if (ratchetDetail) parsedContent.ratchetDetail = ratchetDetail

    fs.writeFileSync(filename, JSON.stringify(parsedContent), {
      encoding: 'utf8'
    })
  } catch (error) {
    console.error('ERROR', error)
  }
}

function handleGetRatchetId(_e: Electron.IpcMainInvokeEvent, username: string) {
  try {
    const filename = CHAT_PREFIX + username + '.json'
    if (fs.existsSync(filename)) {
      const content = fs.readFileSync(filename, {
        encoding: 'utf8'
      })

      if (content) {
        const parsedContent = JSON.parse(content)
        return parsedContent.ratchetId
      }
    }
  } catch (error) {
    console.error('ERROR', error)
  }

  return ''
}

function handleGetRatchetDetailList(_e: Electron.IpcMainInvokeEvent, currentUsername: string) {
  const result: IRatchetDetail[] = []
  try {
    const ratchetFilenames = fs
      .readdirSync(process.cwd())
      .filter(
        (filename) =>
          filename.startsWith(CHAT_PREFIX) && filename !== CHAT_PREFIX + currentUsername + '.json'
      )

    for (const ratchetFilename of ratchetFilenames) {
      const content = fs.readFileSync(ratchetFilename, {
        encoding: 'utf8'
      })
      result.push(JSON.parse(content).ratchetDetail)
    }
  } catch (error) {
    console.error('ERROR', error)
  }

  return result
}

function handleGetOldChatSessions(_e: Electron.IpcMainInvokeEvent, currentUsername: string) {
  const result: {
    receiver: string
    ratchetId: string
    lastMessage: string
    updatedAt: string
    isReaded: boolean
    messages: IMessage[]
  }[] = []
  try {
    const chatFilenames = fs
      .readdirSync(process.cwd())
      .filter(
        (filename) =>
          filename.startsWith(CHAT_PREFIX) && filename !== CHAT_PREFIX + currentUsername + '.json'
      )

    for (const chatFilename of chatFilenames) {
      const content = fs.readFileSync(chatFilename, {
        encoding: 'utf8'
      })
      const parsedContent = JSON.parse(content)

      result.push({
        lastMessage: parsedContent?.lastMessage || '',
        receiver: chatFilename.replace(CHAT_PREFIX, '').replace('.json', ''),
        ratchetId: JSON.parse(content).ratchetId,
        updatedAt: parsedContent.updatedAt,
        isReaded: parsedContent.isReaded || false,
        messages: parsedContent.messages || []
      })
    }
    result.sort((a, b) => new Date(b.updatedAt).getTime() - new Date(a.updatedAt).getTime())
  } catch (error) {
    console.error('ERROR', error)
  }

  return result
}

function handleAddMessageToRatchet(
  _e: Electron.IpcMainInvokeEvent,
  receiver: string,
  messages: IMessage[]
) {
  try {
    const filename = CHAT_PREFIX + receiver + '.json'
    if (!fs.existsSync(filename)) {
      throw new Error('Ratchet file not found')
    } else {
      const fileContent = fs.readFileSync(filename, {
        encoding: 'utf8'
      })
      const parsedContent = JSON.parse(fileContent)
      if (!parsedContent.messages) parsedContent.messages = []
      parsedContent.messages = [...parsedContent.messages, ...messages]
      parsedContent.updatedAt = new Date().toISOString()
      if (messages.length) parsedContent.lastMessage = messages[messages.length - 1].content

      fs.writeFileSync(filename, JSON.stringify(parsedContent), {
        encoding: 'utf8'
      })
    }
  } catch (error) {
    console.error('ERROR', error)
  }
}

function handleGetMessagesByUsername(_e: Electron.IpcMainInvokeEvent, username: string) {
  try {
    const filename = CHAT_PREFIX + username + '.json'
    if (!fs.existsSync(filename)) {
      return []
    } else {
      const fileContent = fs.readFileSync(filename, {
        encoding: 'utf8'
      })
      const parsedContent = JSON.parse(fileContent)
      if (!parsedContent.messages) parsedContent.messages = []
      return parsedContent.messages
    }
  } catch (error) {
    console.error('ERROR', error)
  }
}

function handleUpdateAvatar(_e: Electron.IpcMainInvokeEvent, avatar: string, filename) {
  try {
    fs.writeFileSync(filename, avatar, {
      encoding: 'utf8'
    })
  } catch (error) {
    console.error('ERROR', error)
  }
}

function handleGetAvatar() {
  try {
    if (fs.existsSync(AVATAR_FILE)) {
      return fs.readFileSync(AVATAR_FILE, {
        encoding: 'utf8'
      })
    }
  } catch (error) {
    console.error('ERROR', error)
  }

  return 'https://source.unsplash.com/RZrIJ8C0860'
}

function handleGetRandom32Bytes() {
  try {
    const bytes = randomBytes(32)
    return Buffer.from(bytes).toString('base64')
  } catch (error) {
    console.error('ERROR', error)
  }
}

async function handleEncryptblob(
  _e: Electron.IpcMainInvokeEvent,
  data: ArrayBuffer,
  randomKey: string,
  mimeType: string
) {
  try {
    const res = await encryptblob(new Blob([data], { type: mimeType }), randomKey)
    return Buffer.from(await res.arrayBuffer()).toString('base64')
  } catch (error) {
    console.error('ERROR', error)
  }
}

async function handleDecryptblob(
  _e: Electron.IpcMainInvokeEvent,
  data: ArrayBuffer,
  randomKey: string,
  mimeType: string
) {
  try {
    const res = await decryptblob(new Blob([data], { type: mimeType }), randomKey)

    return Buffer.from(await res.arrayBuffer()).toString('base64')
  } catch (error) {
    console.error('ERROR', error)
  }
}

async function handleDeleteMessage(
  _e: Electron.IpcMainInvokeEvent,
  receiver: string,
  index: number
) {
  try {
    const filename = CHAT_PREFIX + receiver + '.json'
    const messages = handleGetMessagesByUsername(null as never, receiver)
    messages[index].isDeleted = true
    const fileContent = fs.readFileSync(filename, {
      encoding: 'utf8'
    })
    const parsedContent = JSON.parse(fileContent)
    if (parsedContent.messages.length - 1 === index) {
      parsedContent.lastMessage = 'Tin nhắn đã bị xóa'
    }

    parsedContent.messages = messages
    fs.writeFileSync(filename, JSON.stringify(parsedContent), {
      encoding: 'utf8'
    })
  } catch (error) {
    console.error('ERROR', error)
  }
}

async function handleChangeConversationReaded(
  _e: Electron.IpcMainInvokeEvent,
  username: string,
  isReaded: boolean
) {
  try {
    const filename = CHAT_PREFIX + username + '.json'
    const fileContent = fs.readFileSync(filename, {
      encoding: 'utf8'
    })
    const parsedContent = JSON.parse(fileContent)
    parsedContent.isReaded = isReaded
    fs.writeFileSync(filename, JSON.stringify(parsedContent), {
      encoding: 'utf8'
    })
  } catch (error) {
    console.error('ERROR', error)
  }
}
