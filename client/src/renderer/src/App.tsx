import HomeScreen from './screens/Home'
import { BrowserRouter, Route, Routes } from 'react-router-dom'
import { HOME_PAGE, SIGN_IN_PAGE, SIGN_UP_PAGE } from './configs/routes'
import SignInScreen from './screens/SignIn'
import SignUpScreen from './screens/SignUp'
import PinAuthentication from './components/PinAuthentication'
import WasmProvider from './providers/WasmProvider'
import { ToastContainer } from 'react-toastify'
import 'react-toastify/dist/ReactToastify.css'

const App = () => {
  return (
    <WasmProvider>
      <BrowserRouter>
        <Routes>
          <Route
            path={HOME_PAGE}
            element={
              <PinAuthentication>
                <HomeScreen />
              </PinAuthentication>
            }
          />
          <Route path={SIGN_IN_PAGE} element={<SignInScreen />} />
          <Route path={SIGN_UP_PAGE} element={<SignUpScreen />} />
        </Routes>
      </BrowserRouter>
      <ToastContainer />
    </WasmProvider>
  )
}

export default App
