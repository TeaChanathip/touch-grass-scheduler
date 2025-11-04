import {
    LoginPayload,
    RegisterPayload,
    ResetPwdPayload,
} from "../../interfaces/Auth.interface"
import { User } from "../../interfaces/User.interface"
import { apiService, ApiService } from "../api.service"

export class AuthService {
    constructor(private readonly apiService: ApiService) {}

    private url = `${process.env.NEXT_PUBLIC_API_URL}/api/v1/auth`

    async getRegistrationMail(email: string): Promise<null> {
        const encodedParam = encodeURIComponent(email)
        return await this.apiService.get<null>(
            `${this.url}/registration-mail/${encodedParam}`
        )
    }

    async register(
        registrationToken: string,
        payload: RegisterPayload
    ): Promise<{ user: User }> {
        const encodedParam = encodeURIComponent(registrationToken)
        return await this.apiService.post<RegisterPayload, { user: User }>(
            `${this.url}/register/${encodedParam}`,
            payload
        )
    }

    async login(payload: LoginPayload): Promise<{ user: User }> {
        return await this.apiService.post<LoginPayload, { user: User }>(
            `${this.url}/login`,
            payload
        )
    }

    async logout(): Promise<null> {
        return await this.apiService.post<null, null>(`${this.url}/logout`)
    }

    async getResetPwdMail(email: string): Promise<null> {
        const encodedParam = encodeURIComponent(email)
        return await this.apiService.get<null>(
            `${this.url}/reset-password-mail/${encodedParam}`
        )
    }

    async resetPwd(payload: ResetPwdPayload): Promise<null> {
        return await this.apiService.put<ResetPwdPayload, null>(
            `${this.url}/reset-password`,
            payload
        )
    }
}

export const authService = new AuthService(apiService)
