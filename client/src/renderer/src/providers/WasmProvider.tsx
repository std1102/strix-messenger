import { ReactNode, useEffect } from 'react'
import '../wasm/wasm_exec'

const WasmProvider = ({ children }: { children: ReactNode }) => {
  useEffect(() => {
    const go = new window.Go()

    const initWasm = async () => {
      const { instance } = await WebAssembly.instantiateStreaming(
        fetch('main.wasm'),
        go.importObject
      )
      await go.run(instance)
    }

    initWasm().then()
  }, [])

  return children
}

export default WasmProvider
