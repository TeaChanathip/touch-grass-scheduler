import { LoginPayload } from "../../../interfaces/LoginPayload.interface"
import { RegisterPayload } from "../../../interfaces/RegisterPayload.interface"
import { User } from "../../../interfaces/User.interface"
import {
    ApiErrorResponse,
    ApiService,
    isApiError,
} from "../../../services/api.service"
import { AuthService } from "../../../services/auth/auth.service"
import { createAppSlice } from "../../createAppSlice"
import { UsersService } from "../../../services/users/users.service"

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
const usersService = new UsersService(apiService)

export const userSlice = createAppSlice({
    name: "user",
    initialState,
    reducers: (create) => ({
        userRegister: create.asyncThunk(
            async (
                payload: {
                    registrationToken: string
                    registerPayload: RegisterPayload
                },
                { rejectWithValue }
            ) => {
                try {
                    const { user } = await authService.register(
                        payload.registrationToken,
                        payload.registerPayload
                    )
                    return user
                } catch (err) {
                    if (isApiError(err)) {
                        return rejectWithValue({
                            status: err.status,
                            message: err.message,
                        })
                    }
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
                    state.errMsg = undefined
                },
                rejected: (state, action) => {
                    const payload = action.payload as
                        | ApiErrorResponse
                        | undefined

                    state.status = "error"
                    state.user = undefined

                    // Other errors
                    if (payload === undefined) {
                        state.errMsg = action.error.message
                        return
                    }

                    // ApiError
                    switch (payload.message) {
                        case "actionToken parsing failed":
                            state.errMsg = "Invalid registration URL"
                            break
                        case "email already exists":
                            state.errMsg = "User already registered"
                            break
                        default:
                            state.errMsg = "Something went wrong"
                    }
                },
            }
        ),
        userLogin: create.asyncThunk(
            async (payload: LoginPayload, { rejectWithValue }) => {
                try {
                    const { user } = await authService.login(payload)
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
                        | ApiErrorResponse
                        | undefined

                    state.user = undefined

                    // Other errors
                    if (payload === undefined) {
                        state.status = "error"
                        state.errMsg = action.error.message
                        return
                    }

                    // ApiError
                    if (payload?.status === 401) {
                        state.status = "unauthenticated"
                        state.errMsg = "Invalid email or password"
                    } else {
                        state.status = "error"
                        state.errMsg = "Something went wrong"
                    }
                },
            }
        ),
        userLogout: create.asyncThunk(
            async (_, { rejectWithValue }) => {
                try {
                    await authService.logout()
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
                fulfilled: (state) => {
                    state.user = undefined
                    state.status = "unauthenticated"
                    state.errMsg = undefined
                },
                rejected: (state, action) => {
                    const payload = action.payload as
                        | { status?: number; message?: string }
                        | undefined

                    state.status = "error"
                    state.errMsg = payload?.message ?? action.error.message
                    state.user = undefined
                },
            }
        ),
        userAutoLogin: create.asyncThunk(
            async (_, { rejectWithValue }) => {
                try {
                    const { user } = await usersService.getMe()
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
                pending: (state): void => {
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

                    state.user = undefined

                    // Other errors
                    if (payload === undefined) {
                        state.status = "error"
                        state.errMsg = action.error.message
                        return
                    }

                    // ApiError
                    if (payload?.status === 401) {
                        state.status = "unauthenticated"
                        state.errMsg = undefined
                    } else {
                        state.status = "error"
                        state.errMsg = payload?.message
                    }
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
export const { userLogin, userRegister, userLogout, userAutoLogin } =
    userSlice.actions

// Selectors
export const { selectUser, selectUserStatus, selectUserErrMsg } =
    userSlice.selectors
