import { LoginPayload } from "../../../interfaces/LoginPayload.interface"
import { RegisterPayload } from "../../../interfaces/RegisterPayload.interface"
import { User } from "../../../interfaces/User.interface"
import { ApiService, isApiError } from "../../../services/api.service"
import { AuthService } from "../../../services/auth/auth.service"
import { createAppSlice } from "../../createAppSlice"

export interface UserSliceState {
    user?: User
    status: "idle" | "loading" | "authenticated" | "unauthenticated" | "error"
    errMsg?: string
}

const initialState: UserSliceState = {
    user: undefined,
    status: "idle",
    errMsg: undefined,
}

const apiService = new ApiService()
const authService = new AuthService(apiService)

export const userSlice = createAppSlice({
    name: "user",
    initialState,
    reducers: (create) => ({
        login: create.asyncThunk(
            async (payload: LoginPayload, { rejectWithValue }) => {
                try {
                    const { user, token } = await authService.login(payload)

                    localStorage.setItem("token", token)

                    return user
                } catch (err) {
                    if (isApiError(err)) {
                        return rejectWithValue({
                            status: err.status,
                            message: err.message,
                        })
                    }
                    // unexpected error
                    throw err
                }
            },
            {
                pending: (state) => {
                    state.status = "loading"
                },
                fulfilled: (state, action) => {
                    state.status = "authenticated"
                    state.user = action.payload
                    state.errMsg = ""
                },
                rejected: (state, action) => {
                    const payload = action.payload as
                        | { status?: number; message?: string }
                        | undefined

                    // 401 Unauthorized
                    if (payload?.status === 401) {
                        state.status = "unauthenticated"
                        state.errMsg = ""
                    } else {
                        state.status = "error"
                        state.errMsg = action.error.message
                    }
                    state.user = undefined
                },
            }
        ),
        register: create.asyncThunk(
            async (payload: RegisterPayload, { rejectWithValue }) => {
                try {
                    const { user, token } = await authService.register(payload)

                    localStorage.setItem("token", token)

                    return user
                } catch (err) {
                    if (isApiError(err)) {
                        return rejectWithValue({
                            status: err.status,
                            message: err.message,
                        })
                    }
                    // unexpected error
                    throw err
                }
            },
            {
                pending: (state) => {
                    state.status = "loading"
                },
                fulfilled: (state, action) => {
                    state.status = "authenticated"
                    state.user = action.payload
                    state.errMsg = ""
                },
                rejected: (state, action) => {
                    const payload = action.payload as
                        | { status?: number; message?: string }
                        | undefined

                    // 401 Unauthorized
                    if (payload?.status === 401) {
                        state.status = "unauthenticated"
                        state.errMsg = ""
                    } else {
                        state.status = "error"
                        state.errMsg = action.error.message
                    }
                    state.user = undefined
                },
            }
        ),
    }),
    selectors: {
        selectUser: (state) => state.user,
        selectUserStatus: (state) => state.status,
        selectUserErrMsg: (state) => state.errMsg,
    },
})

// Actions
export const { login, register } = userSlice.actions

// Selectors
export const { selectUser, selectUserStatus, selectUserErrMsg } =
    userSlice.selectors

//TODO: action for get user if the token is available in local storage
