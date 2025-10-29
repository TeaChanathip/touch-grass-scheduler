import { LoginPayload, RegisterPayload } from "../../interfaces/Auth.interface"
import { User } from "../../interfaces/User.interface"
import { ApiService } from "../api.service"

export class AuthService {
    constructor(private readonly apiService: ApiService) {}

    private url = `${process.env.NEXT_PUBLIC_API_URL}/api/v1/auth`

    async getRegistrationMail(email: string) {
        const encodedParam = encodeURIComponent(email)
        return await this.apiService.get<null>(
            `${this.url}/registration-mail/${encodedParam}`
        )
    }

    async register(registrationToken: string, payload: RegisterPayload) {
        const encodedParam = encodeURIComponent(registrationToken)
        return await this.apiService.post<RegisterPayload, { user: User }>(
            `${this.url}/register/${encodedParam}`,
            payload
        )
    }

    async login(payload: LoginPayload) {
        return await this.apiService.post<LoginPayload, { user: User }>(
            `${this.url}/login`,
            payload
        )
    }

    async logout() {
        return await this.apiService.post<null, null>(`${this.url}/logout`)
    }
}
